FROM golang:1.15.8-alpine AS build
WORKDIR /src
RUN apk --no-cache add build-base git mercurial gcc
RUN go get github.com/bshuster-repo/logrus-logstash-hook
RUN go get github.com/sirupsen/logrus
COPY . .
RUN go build -o /bin/service .
FROM alpine
WORKDIR /app
COPY --from=build /bin/service /app/
CMD ["./service"]