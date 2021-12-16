########
FROM golang:1.17-buster as backendbuild

WORKDIR /go/src/github.com/analogj/tangle

COPY . /go/src/github.com/analogj/tangle

RUN go mod vendor && \
    go build -ldflags '-w -extldflags "-static"' -o tangle webapp/backend/cmd/tangle/tangle.go

########
FROM gcr.io/distroless/static-debian11 as runtime
EXPOSE 8080

COPY --from=backendbuild /go/src/github.com/analogj/tangle/tangle /tangle
#COPY --from=frontendbuild /scrutiny/dist /scrutiny/web
CMD ["/tangle", "start"]