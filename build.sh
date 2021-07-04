DIR=$PWD

echo "Building WASM..."
cd "$DIR/src"
GOOS=js GOARCH=wasm go build -o ../bin/spaghetti.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ../bin/wasm_exec.js

echo "Building Webpack..."
cd "$DIR"
npx webpack --mode=development && \
cp ./bin/spaghetti.js resources/spaghetti.js && \
rm -rf bin

echo "Done"
