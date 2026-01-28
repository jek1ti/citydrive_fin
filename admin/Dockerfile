FROM golang:1.24.5 AS build
WORKDIR /src
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
COPY . .
WORKDIR /src/admin
RUN /bin/sh -c 'set -e; for i in 1 2 3 4 5; do go mod download && exit 0; echo "go mod download failed, retry ${i}/5" >&2; sleep 2; done; exit 1'
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/admin ./cmd

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates busybox-extras
COPY --from=build /out/admin /app/admin
EXPOSE 50053
ENTRYPOINT ["/app/admin"]
