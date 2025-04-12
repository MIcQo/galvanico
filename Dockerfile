# use the official Bun image
# see all versions at https://hub.docker.com/r/oven/bun/tags
FROM oven/bun:1@sha256:7eb9c0438a42438d884891f5460d6f5b89c20797cb58062b6d28ccba725a8c42 AS base
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

FROM golang:1.24@sha256:1ecc479bc712a6bdb56df3e346e33edcc141f469f82840bab9f4bc2bc41bf91d AS builder
WORKDIR /go/src/app
COPY . .
COPY --from=prerelease /temp/prod/dist public
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM busybox:uclibc@sha256:5ea5116fb0855f9c513dac58fabc55509a8da6df6945587cbf9bc206cb9af0eb AS busybox

FROM gcr.io/distroless/base@sha256:27769871031f67460f1545a52dfacead6d18a9f197db77110cfc649ca2a91f44
WORKDIR /go
COPY --from=busybox /bin/ls /bin/ls
COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/stat /bin/stat
COPY --from=builder /go/bin/app /go/bin/app
CMD ["/go/bin/app", "serve"]
