FROM golang:1.17

WORKDIR /go/src/app
COPY ./cmd/api .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["hill_user_crud"]