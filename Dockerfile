FROM golang:1.22 AS build

WORKDIR /build

COPY . .
RUN go mod download
RUN go mod verify

RUN make build

FROM gcr.io/distroless/static

COPY --from=build /build/server /usr/local/bin/server

ENTRYPOINT ["server"]
