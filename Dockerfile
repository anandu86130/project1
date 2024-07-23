FROM golang:1.21.1

WORKDIR /PROJECT

COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o PROJECT

EXPOSE 8080

CMD ["./PROJECT"]