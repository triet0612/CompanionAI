FROM golang:latest
COPY ./backend /app/
WORKDIR /app
RUN go build -o ./companion.bin ./cmd/main.go

ENTRYPOINT [ "/app/companion.bin" ]
