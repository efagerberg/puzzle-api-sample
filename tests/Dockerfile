FROM golang:1.9.0-alpine

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

RUN apk --no-cache add ca-certificates postgresql git openssh

ADD app.go main.go puzzle.go main_test.go postgres/wait-for-postgres.sh /gopath/src/github.com/efagerberg/puzzle-api-sample/
WORKDIR /gopath/src/github.com/efagerberg/puzzle-api-sample/

RUN go get github.com/lib/pq \
           github.com/gorilla/mux
ENTRYPOINT ["sh", "wait-for-postgres.sh", "database"]
CMD ["go", "test", "-v"]
