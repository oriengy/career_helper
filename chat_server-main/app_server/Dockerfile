FROM debian:12-slim

WORKDIR /app

# Expect a prebuilt binary at ./dist/server in the build context.
COPY dist/server /app/server
COPY config.yaml /app/config.yaml

EXPOSE 8082

CMD ["/app/server", "-c", "/app/config.yaml"]
