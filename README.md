# smtp-service

Tiny HTTP service that sends plain text email over SMTP.

## Run

Fill in `.env`, then:

```sh
docker compose up -d --build
```

## Use

```sh
curl -X POST localhost:8080/send \
  -d '{"to":"someone@example.com","subject":"Hi","body":"Hello"}'
```

Returns `204` on success.
