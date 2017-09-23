FROM golang:1.7.3 AS build-stage
WORKDIR /go/src/github.com/efagerberg/puzzle-api/
RUN go get github.com/lib/pq \
           github.com/gorilla/mux
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql
WORKDIR /root/
COPY --from=build-stage /go/src/github.com/efagerberg/puzzle-api/app .
COPY postgres/wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh app
ENTRYPOINT ["sh", "wait-for-postgres.sh", "database"]
CMD ["./app"]
