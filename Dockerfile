FROM cflynnus/golang-1.12.0-packr2 as builder
ARG production_build="false"
WORKDIR /go/src/myblog
COPY . .
RUN \
  mkdir ./dist; \
  echo "production_build: $production_build"; \
  if [ "$production_build" = "true" ]; then \
    packr2; \
  else \
    cp -rf ./posts/ ./dist/posts/;  \
    cp -rf ./templates/ ./dist/templates/; \
    mkdir -p ./dist/web/static/; \
    cp -rf ./web/static/ ./dist/web/static/; \
  fi; 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./dist/main .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/myblog/dist/ ./
CMD ["./main"]
