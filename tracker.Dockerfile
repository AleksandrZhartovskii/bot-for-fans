# build stage
FROM golang as builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN cd cmd/release-tracker/ && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /app/cmd/release-tracker/main /app/
COPY cmd/release-tracker/.env .
ENTRYPOINT ["./main"]