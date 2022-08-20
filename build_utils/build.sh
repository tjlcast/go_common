osType=$1
appName="xxxx"

if [ -z $osType ]; then
    osType="linux"
fi

if [ $osType = mac ]; then
    echo "Choose: mac"
    ./build-mac.sh $appName
else
    echo "Choose: linux"
    ./build-linux.sh $appName
fi

# back to mac.
go env -w CGO_ENABLED=1 GOOS=darwin GOARCH=amd64
echo "Back to mac."

echo "Finish build."
