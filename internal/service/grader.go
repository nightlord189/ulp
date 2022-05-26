package service

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/nightlord189/ulp/internal/model"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (s *Service) CreateAttempt(req model.AttemptRequest, file *multipart.FileHeader, fileSrc *multipart.File) error {
	var task model.TaskDB
	err := s.DB.GetEntityByField("id", strconv.Itoa(req.TaskID), &task)
	if err != nil {
		return fmt.Errorf("err get task from db: %w", err)
	}
	attempt := model.AttemptDB{
		TaskID:      req.TaskID,
		State:       "fail",
		Log:         "creating entity in db\n",
		RunningTime: 0,
		CreatorID:   req.CreatorID,
	}
	err = s.DB.CreateEntity(&attempt)
	if err != nil {
		return fmt.Errorf("err create attempt in db: %w", err)
	}
	defer func() {
		err = s.DB.UpdateEntity(&attempt)
		if err != nil {
			fmt.Println("error update attempt")
		}
	}()
	attemptDir, err := s.createAttemptFiles(&attempt, file, fileSrc)
	if err != nil {
		return fmt.Errorf("err write files to server's filesystem: %w", err)
	}
	err = createDockerfile(attemptDir, task.Dockerfile)
	if err != nil {
		return fmt.Errorf("err write dockerfile: %w", err)
	}
	err = buildAndRun(task, &attempt, attemptDir)
	if err != nil {
		return fmt.Errorf("err build and run: %w", err)
	}
	// TODO: create container
	// TODO: exec or make web-request to container
	return nil
}

func buildAndRun(task model.TaskDB, attempt *model.AttemptDB, attemptDir string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("error create docker-client: %w", err)
	}
	defer func() {
		err := cli.Close()
		if err != nil {
			fmt.Println("err close docker-client:", err)
		}
	}()
	err = build(ctx, cli, attemptDir, attempt.ID)
	if err != nil {
		attempt.Log += "docker build failed\n"
		return fmt.Errorf("error build image: %w", err)
	}
	defer removeImage(ctx, cli, attempt)
	attempt.Log += "docker build succeed\n"
	err = run(ctx, cli, attempt)
	if err != nil {
		attempt.Log += "run failed\n"
		return fmt.Errorf("error run image: %w", err)
	}
	return nil
}

func removeImage(ctx context.Context, cli *client.Client, attempt *model.AttemptDB) {
	imageName := fmt.Sprintf("ulp/grader/%d:latest", attempt.ID)
	_, err := cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{})
	if err != nil {
		fmt.Printf("error on removing image %s: %v", imageName, err)
	}
	attempt.Log += "image removed\n"
}

func run(ctx context.Context, cli *client.Client, attempt *model.AttemptDB) error {
	containerName := fmt.Sprintf("grader%d", attempt.ID)
	imageName := fmt.Sprintf("ulp/grader/%d:latest", attempt.ID)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		//Cmd:   []string{"cat", "/etc/hosts"},
		Tty: true,
	}, nil, nil, &specs.Platform{
		OS:           "linux",
		Architecture: "arm64",
	}, containerName)
	if err != nil {
		return fmt.Errorf("error create container: %w", err)
	}

	defer func() {
		fmt.Println("removing container")
		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			fmt.Println("error remove container:", err)
			return
		}
		fmt.Printf("container %s removed\n", resp.ID)
		attempt.Log += "docker container removed\n"
	}()

	if err := cli.ContainerStart(ctx, resp.ID,
		types.ContainerStartOptions{
			CheckpointID:  "",
			CheckpointDir: "",
		}); err != nil {
		return fmt.Errorf("error start container: %w", err)
	}

	defer func() {
		fmt.Println("stopping container")
		stopTimeout := 5 * time.Second
		if err := cli.ContainerStop(ctx, resp.ID, &stopTimeout); err != nil {
			fmt.Println("error stop container:", err)
			return
		}
		fmt.Printf("container %s stopped\n", resp.ID)
		attempt.Log += "docker container stopped\n"
	}()

	fmt.Println("container started")
	time.Sleep(3 * time.Second)
	fmt.Println("after 3 seconds")
	attempt.Log += "docker container started\n"

	return nil
}

func build(ctx context.Context, cli *client.Client, attemptDir string, attemptID int) error {
	dirPath := fmt.Sprintf("%s/", attemptDir)
	tar, err := archive.TarWithOptions(dirPath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("error on create tar-archive: %w", err)
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{fmt.Sprintf("ulp/grader/%d", attemptID)},
		Remove:     true,
	}
	res, err := cli.ImageBuild(ctx, tar, opts)
	if err != nil {
		return fmt.Errorf("error on build image: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Println("err on close body of build response:", err)
		}
	}()
	err = print(res.Body)
	if err != nil {
		return fmt.Errorf("error on reading response from docker build: %w", err)
	}
	return nil
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		fmt.Println(scanner.Text())
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func createDockerfile(dirPath, content string) error {
	fpath := path.Join(dirPath, "Dockerfile")
	return ioutil.WriteFile(fpath, []byte(content), 0644)
}

func (s *Service) createAttemptFiles(attempt *model.AttemptDB, file *multipart.FileHeader, fileSrc *multipart.File) (string, error) {
	attemptDirPath := path.Join(s.Config.AttemptsPath, fmt.Sprintf("%d", attempt.ID))
	fmt.Println(attemptDirPath)
	if err := os.Mkdir(attemptDirPath, os.ModePerm); err != nil {
		attempt.Log += "error on creating attempt directory\n"
		return attemptDirPath, fmt.Errorf("error create attempt directory: %w", err)
	}
	attemptFilePath := path.Join(attemptDirPath, file.Filename)
	dst, err := os.Create(attemptFilePath)
	if err != nil {
		attempt.Log += "error on creating attempt file\n"
		return attemptDirPath, fmt.Errorf("error on create file path: %w", err)
	}
	if _, err = io.Copy(dst, *fileSrc); err != nil {
		attempt.Log += "error on copy file from src\n"
		return attemptDirPath, fmt.Errorf("error on copy file from src: %w", err)
	}
	extension := filepath.Ext(file.Filename)
	if extension == ".zip" {
		err = unzipFile(attemptFilePath, attemptDirPath)
		if err != nil {
			return attemptDirPath, fmt.Errorf("error on unzipping archive: %w", err)
		}
		err = os.Remove(attemptFilePath)
		if err != nil {
			return attemptDirPath, fmt.Errorf("error on delete unzipped archive: %w", err)
		}
	}
	attempt.Log += "files created\n"
	return attemptDirPath, nil
}

func unzipFile(src, dst string) error {
	archive, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("error on open zip reader: %w", err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if strings.Contains(filePath, "__MACOSX") {
			fmt.Println("skip", filePath)
			continue
		}

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("error on mkdir: %w", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("error on openfile: %w", err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("error on open: %w", err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("error on copy: %w", err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
