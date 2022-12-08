FROM golang:1.18

RUN mkdir /app

COPY go.mod /app
COPY go.sum /app
COPY *.go /app

WORKDIR /app

RUN go build -o repo-watcher .

CMD ["/app/repo-watcher"]