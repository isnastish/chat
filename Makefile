.SILENT:

TARGET_DIR = build
TARGET_APP = server
SRC = main.go server.go helpers.go registry.go

.PHONY: all clean

all: clean build run

build: $(SRC)
	go build -o $(TARGET_DIR)/$(TARGET_APP) $(SRC)

run:
	./build/server

clean:
	rm -rf $(TARGET_DIR)
