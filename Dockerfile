# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh libc6-compat

WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux

# Copy go mod and sum files
COPY go.mod ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go build -o /main ./main.go

EXPOSE 8001

ENTRYPOINT ["/main"]