# Seatsurfing Backend

Seat booking server software which enables your organisation's employees to book seats, desks and rooms.

## Quick reference
* **Maintained by:** [Seatsurfing.app](https://seatsurfing.app/)
* **Where to get help:** [Documentation](https://seatsurfing.app/docs/)
* **Supported architectures:** amd64, arm64
* **License:** [GPL 3.0](https://github.com/seatsurfing/backend/blob/master/LICENSE)

## Supported tags
* ```latest``` refers to Seatsurfing Backend {{version}} as of {{date}}
* ```{{version}}``` as of {{date}}

## How to use this image
### Start using Docker Compose
```
version: '3.7'

services:
  server:
    image: seatsurfing/backend
    restart: always
    networks:
      sql:
      http:
    ports:
      - 8080:8080
    depends_on:
      - db
    environment:
      POSTGRES_URL: 'postgres://seatsurfing:DB_PASSWORD@db/seatsurfing?sslmode=disable'
      JWT_SIGNING_KEY: 'some_random_string'
      BOOKING_UI_BACKEND: 'booking-ui:3001'
      ADMIN_UI_BACKEND: 'admin-ui:3000'
      PUBLIC_URL: 'https://seatsurfing.your-domain.com'
      FRONTEND_URL: 'https://seatsurfing.your-domain.com'
  booking-ui:
    image: seatsurfing/booking-ui
    restart: always
    networks:
      http:
    environment:
      FRONTEND_URL: 'https://seatsurfing.your-domain.com'
  admin-ui:
    image: seatsurfing/admin-ui:dev
    restart: always
    networks:
      http:
    environment:
      FRONTEND_URL: 'https://seatsurfing.your-domain.com'
  db:
    image: postgres:15-alpine
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
  http:
```

This starts...
* a PostgreSQL database with data stored on Docker volume "db"
* a Seatsurfing Backend instance with port 8080 exposed.
* a Seatsurfing Booking UI instance which is accessible through the Backend instance at: :8080/ui/
* a Seatsurfing Admin UI instance which is accessible through the Backend instance at: :8080/admin/

### Running on Kubernetes
Please refer to our [Kubernetes documentation](https://docs.seatsurfing.app/kubernetes/).

## Environment variables
Please check out the [documentation](https://seatsurfing.app/docs/config) for information on available environment variables and further guidance.

**Hint**: When running in an IPV6-only Docker/Podman environment with multiple network interfaces bound to the Frontend containers, setting the ```LISTEN_ADDR``` environment variable can be necessary as NextJS binds to only one network interface by default. Set it to ```::``` to bind to any address.
