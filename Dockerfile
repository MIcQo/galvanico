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
RUN bun run build-only

FROM golang:1.24@sha256:991aa6a6e4431f2f01e869a812934bd60fbc87fb939e4a1ea54b8494ab9d2fc6 AS builder
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12@sha256:3d0f463de06b7ddff27684ec3bfd0b54a425149d0f8685308b1fdf297b0265e9
COPY --from=builder /go/bin/app /
COPY --from=prerelease /temp/prod/dist public/
CMD ["/app", "serve"]