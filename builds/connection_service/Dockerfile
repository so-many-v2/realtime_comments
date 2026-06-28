# stage 1
FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/connection_service ./cmd/connection_service

# stage 2
FROM alpine:3.20

COPY --from=build /bin/connection_service /bin/connection_service

ENTRYPOINT ["/bin/connection_service"]