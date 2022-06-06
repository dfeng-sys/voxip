FROM golang:1.18 as build-env
WORKDIR /go/src/app
COPY . .
# download dependencies
RUN go get -d -v ./...
# install package
RUN go install -v ./...
RUN go build -o /go/bin/app
CMD ["app"]

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
COPY db/GeoLite2-Country.mmdb /db/
CMD ["/app"]
