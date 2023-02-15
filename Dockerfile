FROM    golang:1.19.5-alpine
RUN     mkdir -p /app
ENTRYPOINT go run main.go >> ./logs/run.log 2>&1