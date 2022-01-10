FROM golang:1.16-alpine as builder
RUN apk --no-cache --no-progress add --virtual \
  build-deps \
  build-base
WORKDIR /src/booru
COPY . .
RUN GOOS=linux go build -o /booru .

FROM alpine:3.15
COPY --from=builder /booru /booru
COPY --from=builder /src/booru/templates /templates
RUN apk add --no-cache ca-certificates
COPY docker_config.json config.json
CMD ["/booru"]