@echo off

set BUILD_OPTIONS=-v

::format source files
go fmt ./...

::build
if not exist build (
    mkdir build
)
pushd build
go build -o ./server.exe  %BUILD_OPTIONS% ../ 
popd 

echo Build successfull.