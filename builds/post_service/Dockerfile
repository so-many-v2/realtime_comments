# stage 1
FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/post_service ./cmd/post_service

# stage 2
FROM alpine:3.20

COPY --from=build /bin/post_service /bin/post_service

ENTRYPOINT ["/bin/post_service"]