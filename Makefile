build:
	mkdir -p .build
	GOOS=linux CGO_ENABLED=0 go build -o .build/main src/main.go
	zip -j function.zip .build/main