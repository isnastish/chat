FROM golang:1.21 as build-image

WORKDIR /build

# copy go.mod and go.sum into build folder
COPY go.mod go.sum ./

# download all the dependencies
RUN go mod download

# copy all your source code into a build directory.
COPY *.go ./ 

COPY Makefile ./

# compile application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./instance

EXPOSE 8500/tcp
EXPOSE 8600/udp

# run
CMD ["/instance"]
