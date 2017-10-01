FROM golang:1.9.0-alpine AS build-stage

RUN apk --no-cache add ca-certificates postgresql git

ENV GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PROJPATH=/gopath/src/github.com/efagerberg/puzzle-api-sample

WORKDIR /gopath/src/github.com/efagerberg/puzzle-api-sample/
RUN go get github.com/lib/pq \
           github.com/gorilla/mux
ADD . /gopath/src/github.com/efagerberg/puzzle-api-sample/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o puzzle-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql
WORKDIR /root/
COPY --from=build-stage /gopath/src/github.com/efagerberg/puzzle-api-sample/puzzle-api .
COPY postgres/wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh puzzle-api
ENTRYPOINT ["sh", "wait-for-postgres.sh", "database"]
CMD ["./puzzle-api"]
