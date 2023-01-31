FROM golang:1.18.3 AS builder

COPY . /build_dir
WORKDIR /build_dir
RUN go build -o bin/app cmd/app/main.go
EXPOSE 8080


ENTRYPOINT ["/build_dir/bin/app"]
