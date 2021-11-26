# Seatsurfing Backend

Seat booking server software which enables your organisation's employees to book seats, desks and rooms.

This repository contains the Backend, which consists of:
* The Server (REST API Backend) written in Go
* Admin Web Interface ("Admin UI")
* User Self-Service Booking Web Interface ("Booking UI")
* Common TypeScript files for the two TypeScript/React web frontends

The mobile app is maintained in a [separate repository](https://github.com/seatsurfing/app).

## Quick reference
* **Maintained by:** [Seatsurfing.de](https://seatsurfing.de/)
* **Where to get help:** [Documentation](https://docs.seatsurfing.de/)
* **Docker architectures:** [amd64](http://hub.docker.com/r/seatsurfing/backend)
* **License:** [GPL 3.0](https://github.com/seatsurfing/backend/blob/master/LICENSE)
* **Mobile apps:** [Apple App Store](https://apps.apple.com/app/seatsurfing/id1579071273) and [Google Play](https://play.google.com/store/apps/details?id=de.seatsurfing.app)

## How to use this Docker image

### Start using ```docker run```
```
docker run \
    --rm \
    -p 8080:8080 \
    -e "POSTGRES_URL=postgres://seatsurfing:DB_PASSWORD@db/seatsurfing?sslmode=disable" \
    seatsurfing/backend
```

This starts the Seatsurfing Backend, connects to a PostgreSQL Database on host "db" and exposes port 8080.

### Start using Docker Compose
```
version: '3.7'

services:
  server:
    image: seatsurfing/backend
    restart: always
    networks:
      sql:
    ports:
      - 8080:8080
    environment:
      POSTGRES_URL: 'postgres://seatsurfing:DB_PASSWORD@db/seatsurfing?sslmode=disable'
      JWT_SIGNING_KEY: 'some_random_string'
  db:
    image: postgres:12
    restart: always
    networks:
      sql:
    volumes:
      - db:/var/lib/postgresql/data
    environment:
      POSTGRES_PASSWORD: DB_PASSWORD
      POSTGRES_USER: seatsurfing
      POSTGRES_DB: seatsurfing

volumes:
  db:

networks:
  sql:
```

This starts...
* a PostgreSQL database with data stored on Docker volume "db"
* a Seatsurfing Backend instance with port 8080 exposed.

## Environment variables
Please check out the [documentation](https://docs.seatsurfing.de/) for the latest information on available environment variables and further guidance.

| Environment Variable | Type | Default | Description |
| --- | --- | --- | --- |
| DEV | bool | 0 | Development Mode, set to 1 to enable  |
| PUBLIC_LISTEN_ADDR | string | 0.0.0.0:8080 | TCP/IP listen address and port |
| PUBLIC_URL | string | http://localhost:8080 | Public URL |
| FRONTEND_URL | string | http://localhost:8080 | Frontend URL (usually matches the Public URL) |
| APP_URL | string | seatsurfing:/// | App URL (should not be changed) |
| STATIC_ADMIN_UI_PATH | string | /app/adminui | Path to compiled Admin UI files |
| STATIC_BOOKING_UI_PATH | string | /app/bookingui | Path to compiled Booking UI files |
| POSTGRES_URL | string | postgres://postgres:root @ localhost/seatsurfing?sslmode=disable | PostgreSQL Connection |
| JWT_SIGNING_KEY | string | random string | JWT Signing Key |
| SMTP_HOST | string | 127.0.0.1:25 | SMTP server address and port (authentication is not supported at this time) |
| SMTP_SENDER_ADDRESS | string | no-reply@seatsurfing.local | SMTP sender address |
| MOCK_SENDMAIL | bool | 0 | SMTP mocking, set to 1 to enable |
| PRINT_CONFIG | bool | 0 | Print configuration on startup, set to 1 to enable |
| INIT_ORG_NAME | string | Sample Company | Your organization's name |
| INIT_ORG_DOMAIN | string | seatsurfing.local | Your organization's domain |
| INIT_ORG_USER | string | admin | Your organization's admin username |
| INIT_ORG_PASS | string | 12345678 | Your organization's admin password |
| INIT_ORG_COUNTRY | string | DE | Your organization's ISO country code |
| INIT_ORG_LANGUAGE | string | de | Your organization's ISO language code |
| ORG_SIGNUP_ENABLED | bool | 0 | Allow signup of new organizations, set to 1 to enable |
| ORG_SIGNUP_DOMAIN | string | .on.seatsurfing.local | Signup domain suffix |
| ORG_SIGNUP_ADMIN | string | admin | Admin username for new signups |
| ORG_SIGNUP_MAX_USERS | int | 50 | Maximum number of users for new organisations |