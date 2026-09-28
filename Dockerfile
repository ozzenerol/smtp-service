FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /smtp-service .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /smtp-service /app/smtp-service
COPY config.yaml /app/config.yaml
ENV SERVER_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/smtp-service"]
