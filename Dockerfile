########
FROM golang:1.17-buster as backendbuild

WORKDIR /go/src/github.com/analogj/tangle

COPY . /go/src/github.com/analogj/tangle

RUN go mod vendor && \
    go build -ldflags '-w -extldflags "-static"' -o /tangle webapp/backend/cmd/tangle/tangle.go && \
    chmod +x /tangle && \
    /tangle --version

#########
#FROM gcr.io/distroless/static-debian11 as runtime
#EXPOSE 8080
#
#COPY --from=backendbuild /tangle /tangle
##COPY --from=frontendbuild /tangle/dist /tangle/web
CMD ["/tangle", "start"]