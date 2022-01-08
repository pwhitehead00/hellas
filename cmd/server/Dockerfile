FROM golang:alpine AS build

RUN apk update && \
    apk add git make && \
    rm -rf /var/cache/apk/*

WORKDIR /build

COPY . .
RUN go mod download
RUN go mod verify

RUN make build

FROM scratch
WORKDIR /app

COPY --from=build /build/server /app/server
COPY --from=build /build/server.key /app/server.key
COPY --from=build /build/server.crt /app/server.crt

ENTRYPOINT ["/app/server"]
