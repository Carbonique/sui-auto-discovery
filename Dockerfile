FROM golang:1.17-alpine AS build

WORKDIR /build
COPY ./go/go.mod .
COPY ./go/go.sum .

RUN go mod download

COPY ./go/*.go .

RUN CGO_ENABLED=0 go build -o sui-auto-discovery

FROM alpine:3.15

WORKDIR /app

COPY --from=build /build/sui-auto-discovery sui-auto-discovery

ENTRYPOINT ["/app/sui-auto-discovery"]
CMD ["--apps-config=/config/apps.json", "--check-interval=30"]
