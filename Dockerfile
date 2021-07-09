FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY go.mod /src/
RUN go mod download
COPY cmd/ /src/cmd/
COPY pkg/ /src/pkg/
RUN ls -lah
RUN CGO_ENABLED=0 go build -o /bin/server /src/cmd/server

FROM scratch
COPY --from=build /bin/server /bin/server
ENTRYPOINT ["/bin/server"]