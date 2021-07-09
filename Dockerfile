FROM golang:alpine as builder
RUN mkdir -p /src/github.com/connyay/phlog
ADD . /src/github.com/connyay/phlog
WORKDIR /src/github.com/connyay/phlog
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /main ./cmd/server

FROM scratch
COPY --from=builder /main /app/
EXPOSE 8080
WORKDIR /app
CMD ["./main"]
