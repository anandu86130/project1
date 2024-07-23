FROM golang:1.21.1

WORKDIR /golib

COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o golib

EXPOSE 8080

CMD ["./golib"]