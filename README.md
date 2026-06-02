# UaaS — UUID as a Service

A microservice that generates and validates UUIDs. Because why call `uuid.New()` in your own process when you can pay network latency for it?

> This project is a meme. Inspired by [this post on programmerhumor.io](https://programmerhumor.io/backend-memes/micro-service-for-uuid-12nq). Not for serious use. Probably.

## What it does

- `POST /generateUUID` — generates a new UUID, stores it in SQLite, returns it.
- `GET /validateUUID?id=<uuid>` — checks whether a given UUID has been issued by this service.

Built with Go, `net/http`, `google/uuid`, and `modernc.org/sqlite` (pure-Go, no CGO).

## Setup

```bash
cp .env.example .env
go mod download
go run main.go
```

Server listens on the port from `.env` (default `8080`).

## Usage

Generate one:

```bash
curl -X POST http://localhost:8080/generateUUID
# {"id":"f47ac10b-58cc-4372-a567-0e02b2c3d479"}
```

Validate one:

```bash
curl "http://localhost:8080/validateUUID?id=f47ac10b-58cc-4372-a567-0e02b2c3d479"
# {"valid":false}  -> already issued (409)
# {"valid":true}   -> not in DB     (200)
```

## Status codes

| Endpoint            | Code | Meaning                        |
|---------------------|------|--------------------------------|
| `POST /generateUUID`| 201  | UUID created                   |
| `POST /generateUUID`| 409  | Collision (cosmic ray, panic)  |
| `GET /validateUUID` | 200  | Not in DB                      |
| `GET /validateUUID` | 409  | Already in DB                  |
| `GET /validateUUID` | 422  | Missing or malformed UUID      |

## Disclaimer

Do not deploy. Do not scale. Do not put this behind a load balancer. If your architecture diagram has a box labeled "UUID Service", please reconsider your life choices.
