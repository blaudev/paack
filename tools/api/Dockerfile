FROM golang as go

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

EXPOSE 5010

ENTRYPOINT ["./run"]