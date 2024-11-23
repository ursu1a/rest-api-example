# Choose whatever you want, version >= 1.16
FROM golang:1.23-alpine

WORKDIR /app

# Install required packages for the build
RUN apk update && apk add --no-cache bash git curl
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
