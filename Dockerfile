FROM node:15 AS admin-ui-builder
RUN mkdir -p /usr/src/app /usr/src/commons/ts/
WORKDIR /usr/src/commons/ts/
ADD commons/ts/ .
RUN npm install
RUN npm run build
WORKDIR /usr/src/app
ADD admin-ui/ .
RUN npm install --legacy-peer-deps
RUN npm run build

FROM node:15 AS booking-ui-builder
RUN mkdir -p /usr/src/app /usr/src/commons/ts/
WORKDIR /usr/src/commons/ts/
ADD commons/ts/ .
RUN npm install
RUN npm run build
WORKDIR /usr/src/app
ADD booking-ui/ .
RUN npm install --legacy-peer-deps
RUN npm run build

FROM golang:1.15-alpine AS server-builder
RUN apk --update add --no-cache git gcc g++
RUN export GOBIN=$HOME/work/bin
WORKDIR /go/src/app
ADD server/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine
COPY --from=server-builder /go/src/app/main /app/
COPY --from=admin-ui-builder /usr/src/app/build/ /app/adminui/
COPY --from=booking-ui-builder /usr/src/app/build/ /app/bookingui/
ADD server/res/ /app/res
ADD docker-entrypoint.sh /app/
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["./main"]