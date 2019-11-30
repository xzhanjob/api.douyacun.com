main:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1  go build -ldflags '-extldflags "-static"'  -a -o bin/douyacun main.go
