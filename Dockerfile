FROM golang:1.21

WORKDIR /ewallet

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD [ "ewallet" ]