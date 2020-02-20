############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add --no-cache git
# Create appuser.
ENV USER=appuser
ENV UID=10001
# Enable the go modules feature
ENV GO111MODULE=on
# Set the GOPROXY environment variable
ENV GOPROXY=https://goproxy.io
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
WORKDIR /everest
COPY go.mod .
COPY go.sum .
# Fetch dependencies.
# Using go get.
# RUN go get -d -v
# Using go mod.
RUN go mod download
COPY . .
RUN go mod verify
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /everest/bin/server ./cmd/server/main.go
############################
# STEP 2 build a small image
############################
FROM alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk update && apk add --no-cache bash
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable.
COPY --from=builder /everest/bin/server /everest/bin/server
COPY --from=builder /everest/configs/viper.yaml /etc/everest/viper.yaml
# Use an unprivileged user.
USER appuser:appuser
CMD ["./everest/bin/server"]