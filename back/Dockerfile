FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.* ./

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /app/wait

RUN chmod +x /app/wait

COPY . .

RUN go install

RUN go build -o main main.go

FROM alpine:latest

WORKDIR /srv

COPY --from=builder /app/main main

COPY --from=builder /app/wait wait

COPY /scripts/docker-run.sh docker-run.sh

EXPOSE 8080

CMD ["sh", "docker-run.sh"]