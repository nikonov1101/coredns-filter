default:
	@go run main.go

build/arm:
	CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOARCH=arm GOARM=7 GOOS=linux go build -o fdns_arm7 main.go
