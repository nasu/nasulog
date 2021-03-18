FROM golang:1.15-alpine as builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/nasulog -ldflags '-s -w'
###
FROM scratch
COPY --from=builder /go/bin/nasulog /go/bin/nasulog
EXPOSE 8080
ENTRYPOINT ["/go/bin/nasulog"]
