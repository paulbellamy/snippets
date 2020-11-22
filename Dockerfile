# BASE
FROM golang:1.15-alpine as build
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main ./

# PROD
FROM alpine as production
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true
RUN apk --no-cache add tzdata
WORKDIR /app
COPY --from=build /app/main /app
ENTRYPOINT ["/app/main"]
