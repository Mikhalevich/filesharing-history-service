FROM golang:1.16-alpine as builder

WORKDIR /app

RUN GOOS=linux GOARCH=amd64 go get -u -tags="no_mysql no_sqlite3 no_mssql no_redshift no_clickhouse" github.com/pressly/goose/cmd/goose

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/filesharing-history-service

FROM alpine:latest

WORKDIR /app/
COPY --from=builder /go/bin/filesharing-history-service /app/filesharing-history-service
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/docker/ /app/
COPY --from=builder /app/run.sh /app/run.sh

ENTRYPOINT ["./run.sh"]
