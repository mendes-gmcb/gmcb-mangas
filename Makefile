install:
	@echo Validating dependecies...
	@go mod tidy
	@echo Generating vendor...
	@go mod vendor

run:
	go run ./main.go

dev-up:
	cd ./infra/local && docker-compose -p "q2_pagseguro_gateway" up

dev-stop:
	cd ./infra/local && docker-compose -p "q2_pagseguro_gateway" stop

dev-destroy:
	docker run -it --rm -v $(PWD)/infra/local:/data/dev -w /data/dev alpine ./destroy.sh

test:
	go test ./... -short -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

coverage:
	@echo "${COLOR_YELLOW}Running project coverage...${COLOR_WHITE}\n"
	@go test ./... -v -coverprofile=cover.out
	@go tool cover -html=cover.out -o cover.html
	@echo "${COLOR_GREEN}Coverage completed successfully.${COLOR_WHITE}"