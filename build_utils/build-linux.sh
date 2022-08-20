appName=$1

./git-info.sh

echo "build to linux."
go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64
go build -o $appName
