FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN go mod tidy && CGO_ENABLED=0 go build -o akasha_whisper

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/akasha_whisper .

RUN chmod +x akasha_whisper

EXPOSE 8882

ENTRYPOINT ["/app/akasha_whisper"]