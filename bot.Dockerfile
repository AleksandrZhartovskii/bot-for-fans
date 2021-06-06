# build stage
FROM golang as builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN cd cmd/bot/ && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/cmd/bot/main /app/
COPY cmd/bot/.env .
ADD internal/repository/pgdb/migrations migrations
EXPOSE 8080
ENTRYPOINT ["./main"]