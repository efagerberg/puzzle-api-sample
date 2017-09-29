FROM golang:1.9.0-alpine AS build-stage

RUN apk --no-cache add ca-certificates postgresql git

WORKDIR /gopath/src/github.com/efagerberg/puzzle-api-sample/
RUN go get github.com/lib/pq \
           github.com/gorilla/mux
COPY web/ /gopath/src/github.com/efagerberg/puzzle-api-sample/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql
WORKDIR /root/
COPY --from=build-stage /gopath/src/github.com/efagerberg/puzzle-api-sample/app .
COPY postgres/wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh app
ENTRYPOINT ["sh", "wait-for-postgres.sh", "database"]
CMD ["./app"]
