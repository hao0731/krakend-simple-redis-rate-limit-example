## Redis rate limit in KrakenD

This is an example of a Redis-based rate limit in KrakenD.

## Get Started

1. Build the HTTP Server Plugin:

```sh
sudo sh ./scripts/build-redis-rate-limit.sh --path ./src/plugins/server/redis-rate-limit
```

2. Serve KrakenD and Redis through Docker Compose:

```sh
docker-compose up -d
```