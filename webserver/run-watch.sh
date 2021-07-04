go build
./webserver.exe \
    -dir ../src/ \
    -dir ../src/spaghetti/ \
    -dir ../src/js/ \
    -cmd "build.bat" \
    -filter **/*.** \
    -resources ../resources/