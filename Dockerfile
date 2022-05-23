FROM golang:alpine AS build-env
RUN apk --no-cache add git
ADD . /go/src/github.com/dubuqingfeng/bit-node-crawler
RUN cd /go/src/github.com/dubuqingfeng/bit-node-crawler && \
   go mod download && \
   go build -v -o /src/bin/bit-node-crawler main.go
   
FROM alpine
RUN apk --no-cache add openssl ca-certificates tzdata
WORKDIR /app
COPY --from=build-env /src/bin /app/
COPY --from=build-env /go/src/github.com/dubuqingfeng/bit-node-crawler/configs /app/configs
COPY --from=build-env /go/src/github.com/dubuqingfeng/bit-node-crawler/logs /app/logs
ENTRYPOINT ./bit-node-crawler