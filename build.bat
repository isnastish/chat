@echo off

set BUILD_OPTIONS=-v

::format source files
go fmt ./...

::build
if not exist build (
    mkdir build
)
pushd build

:: build main application
go build -o ./server.exe  %BUILD_OPTIONS% ../ 

:: build study examples
go build -o ./ex8_3.exe %BUILD_OPTIONS% ../misc/ex8_3

popd 

echo Build successfull.