echo "Updating Noodle..."
GOOS=js GOARCH=wasm GOPROXY=direct go get -u github.com/lachee/noodle@matrix-rework
echo "Done"
