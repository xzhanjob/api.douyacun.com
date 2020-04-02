main:
	go build -a -o douyacun main.go
test:
	go build -gcflags "all=-N -l" -o douyacun main.go
