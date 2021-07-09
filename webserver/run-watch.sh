go build
./webserver.exe \
    -dir ../resources/shader/ \
    -dir ../src/ \
    -dir ../src/spaghetti/ \
    -dir ../src/spaghetti/ui/ \
    -dir ../src/js/ \
    -cmd "build.bat" \
    -filter **/*.** \
    -resources ../resources/