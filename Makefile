build:
	@go build -o "bin/simple_mail"

run: build
	@bin/simple_mail
