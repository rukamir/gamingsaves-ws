FROM golang:1.12.9-stretch

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3306
EXPOSE 3000

CMD ["app"]