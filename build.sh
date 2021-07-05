DIR=$PWD
if [ -z "$MODE" ]; then
    MODE="development"
fi

echo ">> Building $MODE"

echo "Clearing BIN"
rm -rf "$DIR/resources/bin"

echo "Building WASM..."
cd "$DIR/src"
GOOS=js GOARCH=wasm go build -tags $MODE -o ../resources/bin/spaghetti.wasm .
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ../resources/bin/wasm_exec.js

echo "Building Webpack..."
cd "$DIR"
npx webpack --mode=$MODE 

if [ "$MODE" == "production" ]; then
    echo "Clearing temporary files"
    rm "$DIR/resources/bin/spaghetti.wasm"
    rm "$DIR/resources/bin/wasm_exec.js"
fi

echo "Done"
