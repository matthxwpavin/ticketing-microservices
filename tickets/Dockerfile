# consider use --build-arg for production building
FROM golang:1.22.4-alpine3.20

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.toml"]