FROM golang:alpine AS build

RUN apk update && \
    apk add git make && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN make build

FROM scratch
WORKDIR /app

COPY --from=build /app/server /app/server

ENTRYPOINT ["/app/server"]
