FROM golang:1.24.1 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api_server

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=build /app/server /app/server
COPY --from=build /app/config.yaml /app/config.yaml

EXPOSE 8082

CMD ["/app/server", "-c", "/app/config.yaml"]
