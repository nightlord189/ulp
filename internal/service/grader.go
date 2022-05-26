package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/nightlord189/ulp/internal/model"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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
	err = s.createAttemptFiles(&attempt, file, fileSrc)
	if err != nil {
		return fmt.Errorf("err write files to server's filesystem: %w", err)
	}
	return nil
}

func (s *Service) createAttemptFiles(attempt *model.AttemptDB, file *multipart.FileHeader, fileSrc *multipart.File) error {
	attemptDirPath := path.Join(s.Config.AttemptsPath, fmt.Sprintf("%d", attempt.ID))
	fmt.Println(attemptDirPath)
	if err := os.Mkdir(attemptDirPath, os.ModePerm); err != nil {
		attempt.Log += "error on creating attempt directory\n"
		return fmt.Errorf("error create attempt directory: %w", err)
	}
	attemptFilePath := path.Join(attemptDirPath, file.Filename)
	dst, err := os.Create(attemptFilePath)
	if err != nil {
		attempt.Log += "error on creating attempt file\n"
		return fmt.Errorf("error on create file path: %w", err)
	}
	if _, err = io.Copy(dst, *fileSrc); err != nil {
		attempt.Log += "error on copy file from src\n"
		return fmt.Errorf("error on copy file from src: %w", err)
	}
	// TODO: zip files
	extension := filepath.Ext(file.Filename)
	if extension == ".zip" {
		err = unzipFile(attemptFilePath, attemptDirPath)
		if err != nil {
			return fmt.Errorf("error on unzipping archive: %w", err)
		}
		err = os.Remove(attemptFilePath)
		if err != nil {
			return fmt.Errorf("error on delete unzipped archive: %w", err)
		}
	}
	attempt.Log += "files created\n"
	return nil
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
