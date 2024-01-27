@echo off

::format
go fmt ./...

::build
if not exist build (mkdir build)
pushd build
go build -o . ../ 
popd 