FROM golang:1.21-bookworm AS server-builder
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD server/ server/
ADD go.mod .
ADD go.sum .
WORKDIR /go/src/app/server
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o main .

FROM gcr.io/distroless/static-debian12
COPY --from=server-builder /go/src/app/server/main /app/
COPY server/res/ /app/res
WORKDIR /app
EXPOSE 8080
USER 65532:65532
CMD ["./main"]