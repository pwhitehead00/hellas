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

COPY --from=build /build/server /usr/local/bin/server

ENTRYPOINT ["server"]
