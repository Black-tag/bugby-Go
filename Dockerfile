FROM golang:1.24.0-alpine AS builder

WORKDIR /app 

COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bugby ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bugby .
COPY --from=builder /app/docs ./docs
COPY rbac_model.conf .
COPY rbac_policy.csv .
RUN chmod +x bugby

ENV PORT=8080

EXPOSE 8080

CMD [ "./bugby" ]
