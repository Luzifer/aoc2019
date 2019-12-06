default: fulltest

fulltest:
	mkdir -p cover
	go test -cover -coverprofile ./cover/cover.profile -v *.go
	go tool cover -html ./cover/cover.profile -o ./cover/index.html

day%:
	go test -cover -v \
		day$*.go day$*_test.go \
		helpers.go intcode.go
