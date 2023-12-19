build:
	@go build -o bin/gosoc

run:
	@./bin/gosoc

postgres:
 	docker run --name gosoc -e DB_USER=postgres -e DB_NAME=postgres -e DB_PASS=goS0c@pgcontainer -d postgres