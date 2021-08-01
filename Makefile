#!make

PASSWORD:=lungo
PROJECT_NAME:=lungo
PROJECT_KEY:=com.github.felix-kaestner.lungo

PWD:=$(shell pwd)
TIME:=$(shell date +%s)
UID:=$(shell id -u)
GID:=$(shell id -g)

docs:
	@godoc -http=:6060

fmt:
	@gofmt -w .

test:
	@go test -v ./... -race

coverage:
	@go test ./... -cover -race -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html

sonar:
	@docker-compose up -d
	@echo -n "Waiting for Sonarqube to start..."
	@until curl -u admin:admin -s http://localhost:9000/system/info > /dev/null 2>&1; do echo -n "."; sleep 1; done
	@until curl -u admin:admin -s http://localhost:9000/api/authentication/validate | grep '"valid":true' > /dev/null 2>&1; do echo -n "."; sleep 1; done
	@curl -u admin:admin -X POST -s "http://localhost:9000/api/users/change_password?login=admin&previousPassword=admin&password=$(PASSWORD)"
	@curl -u admin:$(PASSWORD) -X POST -s "http://localhost:9000/api/projects/create?name=$(PROJECT_NAME)&project=$(PROJECT_KEY)" > /dev/null

lint: coverage
	@set -e; \
	SONAR_LOGIN=$$(curl -u admin:$(PASSWORD) -X POST -s "http://localhost:9000/api/user_tokens/generate?name=$(TIME)" | jq '.token'); \
	sh -c "docker run --rm --user="$(UID):$(GID)" --network sonarqube -v $(PWD):/usr/src sonarsource/sonar-scanner-cli:4.6 sonar-scanner -Dsonar.login=$$SONAR_LOGIN -Dsonar.projectName=$(PROJECT_NAME) -Dsonar.projectKey=$(PROJECT_KEY) -Dsonar.host.url=http://sonarqube:9000 -Dsonar.exclusions='**/*_test.go' -Dsonar.test.inclusions='**/*_test.go' -Dsonar.coverage.exclusions='assert.go,lungo.go,**/*config.go' -Dsonar.go.coverage.reportPaths=coverage.out"

clean:
	@docker-compose down -v
	@rm -rf coverage.html coverage.out cert.pem key.pem .scannerwork .cache
