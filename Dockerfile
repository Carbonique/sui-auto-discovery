FROM golang:1.17-alpine AS build

WORKDIR /build
COPY ./go/go.mod .
COPY ./go/go.sum .

RUN go mod download

COPY ./go/*.go .

RUN CGO_ENABLED=0 go build -o auto-updater

FROM alpine:3.15

WORKDIR /app

COPY --from=build /build/auto-updater auto-updater

ENTRYPOINT ["/app/auto-updater"]
CMD ["--apps-config=/config/apps.json", "--check-interval=30"]
