FROM golang:1.22 AS build

WORKDIR /build

COPY . .
RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=${VERSION}" -a -o server ./cmd/server

FROM gcr.io/distroless/static

COPY --from=build /build/server /usr/local/bin/server

ENTRYPOINT ["server"]
