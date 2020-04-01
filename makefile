main:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -ldflags '-extldflags "-static"'  -a -o douyacun main.go
test:
	go build -gcflags "all=-N -l" -o douyacun main.go
