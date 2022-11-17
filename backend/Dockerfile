FROM golang:alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -v -o ./app .

EXPOSE 8080
CMD ["./app"]
