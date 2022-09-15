FROM golang:1.18-alpine AS build
WORKDIR /src
RUN apk --no-cache add build-base git mercurial gcc
ADD go.* ./
RUN go mod download
COPY . .
RUN go build -o /bin/service .
FROM alpine
WORKDIR /app
COPY --from=build /bin/service /app/
CMD ["./service"]