FROM debian:stable-slim

WORKDIR /app

COPY config/local.yaml ./config/local.yaml
COPY storage/storage.db ./storage/storage.db

COPY ad-market ./ad-market

CMD ["./ad-market"]
