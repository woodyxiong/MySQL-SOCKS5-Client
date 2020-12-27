SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -a -ldflags '-extldflags="-static"' -o ./build ./
