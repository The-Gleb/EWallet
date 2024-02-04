FROM golang:1.21.3

WORKDIR /usr/src/app

COPY . .

# RUN go build -o ./bin/app cmd/main.go

# ENTRYPOINT [ "go", "run", "cmd/main.go" ]

# RUN go install github.com/cosmtrek/air@latest

# RUN air init




# EXPOSE 8080

# RUN go get -u github.com/mattes/migrate

# RUN make migrateup

# RUN go build -o ewallet cmd/main.go

# CMD [ "go run cmd/main.go" ]