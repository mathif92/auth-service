# Build with standard golang image
FROM golang:1.17-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/app

# Package the app into a distroless image
FROM gcr.io/distroless/static

COPY --from=build /go/bin/app /

CMD ["/app"]