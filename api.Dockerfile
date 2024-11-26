# FROM golang:1.23.3-alpine as builder
# # RUN mkdir /app 
# # ADD . /app 
# WORKDIR /app 

# RUN apk update && apk add python3=3.11.5-r0 py3-pip && ln -sf python3 /usr/bin/python && apk add libc-dev && apk add gcc && apk add make  && apk add py3-pip && apk add libx11 && apk add libxcomposite && apk add libxdamage && apk add libxrandr && apk add libxi && apk add libxtst && apk add libxshmfence && apk add mesa-gl && apk add alsa-lib && apk add ttf-freefont
# # Tạo virtual environment
# RUN python3 -m venv /app/venv

# # Kích hoạt venv và cài đặt các package Python cần thiết
# RUN . /app/venv/bin/activate && pip install --upgrade pip && pip install tiktok_captcha_solver
# RUN . /app/venv/bin/activate && pip install --upgrade pip && pip install playwright
# RUN . /app/venv/bin/activate && pip install --upgrade pip && pip install pytest-playwright

# # # Cài đặt các trình duyệt cho Playwright (nếu cần)
# RUN . /app/venv/bin/activate && playwright install

# # Thêm venv vào PATH để các lệnh Python sử dụng venv
# ENV PATH="/app/venv/bin:$PATH"

# COPY go.mod go.sum ./
# RUN go mod download && go mod verify

# # Install compile daemon 
# # RUN go install github.com/githubnemo/CompileDaemon@latest
# RUN go install github.com/air-verse/air@latest

# COPY . .
# # CMD [ "air" ]
# COPY ./entrypoint.sh /entrypoint.sh 

# ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
# RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

# ENTRYPOINT [ "sh", "/entrypoint.sh" ]

FROM golang:1.23.3-alpine as builder
# RUN mkdir /app 
# ADD . /app 
WORKDIR /app 

RUN apk update && apk add libc-dev && apk add gcc && apk add make 
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install compile daemon 
# RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go install github.com/air-verse/air@latest

COPY . .
# CMD [ "air" ]
COPY ./entrypoint.sh /entrypoint.sh 

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

ENTRYPOINT [ "sh", "/entrypoint.sh" ]

# RUN go clean --modcache 
# RUN go mod download 
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .  

# FROM alpine:latest

# WORKDIR /app
# # Install PostgreSQL client tools
# RUN apk --no-cache add ca-certificates 
# RUN apk add --no-cache git make musl-dev go 

# # Copy the built binary from the builder stage
# COPY --from=builder /app/main .

# ENV GOROOT /usr/lib/go 
# ENV GOPATH /go 
# ENV PATH /go/bin:$PATH

# # # Build the application from source
# # FROM golang:1.22.2 AS build-stage

# # WORKDIR /app

# COPY go.mod go.sum ./
# # RUN go mod download

# # COPY *.go ./


# # # Run the tests in the container
# # FROM build-stage AS run-test-stage
# # RUN go test -v ./...

# # # Deploy the application binary into a lean image
# # FROM gcr.io/distroless/base-debian11 AS build-release-stage

# # WORKDIR /

# # COPY --from=build-stage /docker-gs-ping /docker-gs-ping

# # EXPOSE 8080

# # USER nonroot:nonroot

# # ENTRYPOINT ["/docker-gs-ping"]

# # Stage 1: Build Stage
# FROM golang:1.22.2-alpine as builder

# WORKDIR /app

# # Install build dependencies (use musl to avoid glibc dependency issues)
# RUN apk add --no-cache build-base

# # Set goarch to ensure the compability of the OS 
# ENV GOOS=linux GOARCH=arm64

# # Copy dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy application source code
# COPY . .

# # Build the application binary
# # RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
# RUN go build -o /app/main /app/main.go

# # Stage 2: Runtime Stage
# FROM alpine:latest

# WORKDIR /app

# # Install PostgreSQL client tools
# RUN apk add --no-cache postgresql-client

# # Copy the built binary from the builder stage
# COPY --from=builder /app/main .

# # Copy the wait-for.sh script
# COPY wait-for.sh /app/wait-for.sh
# RUN chmod +x /app/wait-for.sh

# # Copy the config.yaml file
# COPY config.yaml ./config.yaml

# # Expose the application's port
# EXPOSE 8080

# # Run the application with the wait-for script
# CMD ["/app/wait-for.sh", "postgres", "./main"]

# # Start from the official Go image
# FROM golang:1.22.2-alpine 
# # Set the Current Working Directory inside the container
# WORKDIR /app
# # Copy go.mod and go.sum files
# COPY go.mod go.sum ./
# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download
# # Copy the source from the current directory to the Working Directory inside the container
# COPY . .
# # Set environment variable for the views directory
# ENV VIEWS_DIR=/app/internal/views
# # Build the Go app
# RUN go build -o /app/main /app/main.go
# # Expose port 8080 to the outside world
# EXPOSE 8080
# # Set environment variable for Gin mode
# ENV GIN_MODE=release
# # Run the executable
# CMD ["/app/main"]