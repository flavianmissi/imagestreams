FROM docker.io/library/golang:1.18

WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 8080
CMD ["app"]
