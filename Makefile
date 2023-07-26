CAPROVER_TOKEN='45ef580af1b8d684afd536a5fed8de5516962c89a8c1a2a09ba3abb46d65dd5e'

start:
	go run main.go

lint:
	golangci-lint run --timeout=5m

.PHONY: test

test:
	go clean --testcache
	go test ./...

docker-build:
	docker build -t nightlord189/ulp .

docker-push:
	docker push nightlord189/ulp:latest

docker-run:
	docker run --rm -p 80:8080 -d --name ulp -e TEST_HOST='172.17.0.1' -e DB_HOST='172.17.0.1' -v /Users/aburavov/dev/volumes/ulp:/attempt -v /var/run/docker.sock:/var/run/docker.sock nightlord189/ulp:latest

ansible-install:
	ansible-galaxy collection install -r deployments/ansible/requirements.yml
	ansible-galaxy role install -r deployments/ansible/requirements.yml

deploy:
	ansible-playbook -i deployments/ansible/inventory.ini deployments/ansible/ansible.yml

deploy-ci:
	rm deploy.tar || true
	tar -cvf ./deploy.tar  ./*
	caprover deploy -t ./deploy.tar --host https://captain.app.tinygreencat.dev --appToken ${CAPROVER_TOKEN} --appName ulp
	rm deploy.tar

deploy:
	caprover deploy

