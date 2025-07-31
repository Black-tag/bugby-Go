FROM golang:1.24.0-alpine AS builder

WORKDIR /app 

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bugby ./cmd

FROM alpine:latest

WORKDIR /app

RUN apk update && apk add --no-cache postgresql-client

COPY --from=builder /app/bugby .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /app/docs ./docs
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

COPY rbac_model.conf .
COPY rbac_policy.csv .
COPY entrypoint.sh /app/entrypoint.sh


RUN chmod +x bugby
RUN chmod +x /app/entrypoint.sh



ENV PORT=8080

EXPOSE 8080

# CMD [ "./bugby" ]
ENTRYPOINT ["/app/entrypoint.sh"]