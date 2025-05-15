# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1@sha256:a02c6162266611419fd84c8f96dbdbf3029532e2491314dee5172a27223e5428 AS base
WORKDIR /usr/src/app

FROM base AS install
WORKDIR /temp/prod
COPY client/package.json client/bun.lock /temp/prod/
RUN bun install --frozen-lockfile --production

FROM base AS prerelease
WORKDIR /temp/prod
COPY --from=install /temp/prod/node_modules node_modules
COPY ./client .
RUN if [ "$(uname -m)" = "aarch64" ]; then \
      bun install @rollup/rollup-linux-arm64-musl --no-save; \
    elif [ "$(uname -m)" = "x86_64" ]; then \
      bun install @rollup/rollup-linux-x64-musl --no-save; \
    fi
ENV NODE_ENV=production
ENV VITE_BACKEND_URL=/
RUN bun run build-only

FROM golang:1.24@sha256:86b4cff66e04d41821a17cea30c1031ed53e2635e2be99ae0b4a7d69336b5063 AS builder
WORKDIR /go/src/app
COPY . .
COPY --from=prerelease /temp/prod/dist public
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM busybox:uclibc@sha256:ed0455a0f719c9f203d2ec745682ed78854395fe419bfac8b6ea2ee2e9413409 AS busybox

FROM gcr.io/distroless/base@sha256:27769871031f67460f1545a52dfacead6d18a9f197db77110cfc649ca2a91f44
WORKDIR /go
COPY --from=busybox /bin/ls /bin/ls
COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/stat /bin/stat
COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /go/src/app/public /go/public
CMD ["/go/bin/app", "serve"]
