# Seatsurfing Backend

Seat booking server software which enables your organisation's employees to book seats, desks and rooms.

## Quick reference
* **Maintained by:** [Seatsurfing.app](https://seatsurfing.app/)
* **Where to get help:** [Documentation](https://seatsurfing.app/docs/)
* **Supported architectures:** amd64, arm64, arm v7
* **License:** [GPL 3.0](https://github.com/seatsurfing/backend/blob/master/LICENSE)

## Supported tags
* ```latest``` refers to Seatsurfing Backend {{version}} as of {{date}}
* ```{{version}}``` as of {{date}}

## How to use this image
## Breaking change in version 1.13
Up to version 1.12, the ```backend``` Docker image included all the static web resources for the Booking and Admin interfaces. With version 1.13, we separated the web interfaces from the backend server. You therefore need to start the ```booking-ui``` and ```admin-ui``` Docker images separately. The backend has an integrated HTTP proxy which forwards incoming requests for ```/ui/``` and ```/admin/``` to the corresponding backends. If you prefer to handle request routing with a preceding reverse proxy (such as Traefik or nginx), you can disable proxy functionality by setting the environment variable ```DISABLE_UI_PROXY=1```.

### Start using Docker Compose
```
version: '3.7'

services:
  server:
    image: seatsurfing/backend
    #build:
    #  context: .
    #  dockerfile: Dockerfile
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
    #build:
    #  context: .
    #  dockerfile: Dockerfile.booking-ui
    restart: always
    networks:
      http:
    environment:
      FRONTEND_URL: 'https://seatsurfing.your-domain.com'
  admin-ui:
    image: seatsurfing/admin-ui:dev
    #build:
    #  context: .
    #  dockerfile: Dockerfile.admin-ui
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
Please check out the [documentation](https://seatsurfing.app/docs/) for the latest information on available environment variables and further guidance.

### Backend
| Environment Variable | Type | Default | Description |
| --- | --- | --- | --- |
| DEV | bool | 0 | Development Mode, set to 1 to enable  |
| PUBLIC_LISTEN_ADDR | string | 0.0.0.0:8080 | TCP/IP listen address and port |
| PUBLIC_URL | string | http://localhost:8080 | Public URL |
| FRONTEND_URL | string | http://localhost:8080 | Frontend URL (usually matches the Public URL) |
| ADMIN_UI_BACKEND | string | localhost:3000 | Host serving the Admin UI frontend |
| BOOKING_UI_BACKEND | string | localhost:3001 | Host serving the Booking UI frontend |
| DISABLE_UI_PROXY | bool | 0 | Disable proxy for admin and booking UI, set to 1 to disable the proxy |
| POSTGRES_URL | string | postgres://postgres:root @ localhost/seatsurfing?sslmode=disable | PostgreSQL Connection |
| JWT_SIGNING_KEY | string | random string | JWT Signing Key |
| SMTP_HOST | string | 127.0.0.1 | SMTP server address |
| SMTP_PORT | int | 25 | SMTP server port |
| SMTP_START_TLS | bool | 0 | Use SMTP STARTTLS extension, set to 1 to enable |
| SMTP_INSECURE_SKIP_VERIFY | bool | 0 | Disable SMTP TLS certificate validation |
| SMTP_AUTH | bool | 0 | SMTP authentication, set to 1 to enable |
| SMTP_AUTH_USER | string |  | SMTP auth username |
| SMTP_AUTH_PASS | string |  | SMTP auth password |
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
| ORG_SIGNUP_DELETE | bool | 0 | Allow admins to delete their own organisation |

### Frontend (Admin UI, Booking UI)
| Environment Variable | Type | Default | Description |
| --- | --- | --- | --- |
| FRONTEND_URL | string | ```req.url``` | Frontend URL |
| PORT | int | 3000 (Admin UI), 3001 (Booking UI) | The server's HTTP port |
| LISTEN_ADDR | string | | TCP/IP listen address (defaults to NextJS' ```hostname``` setting) |

**Hint**: When running in an IPV6-only Docker/Podman environment with multiple network interfaces bound to the Frontend containers, setting the ```LISTEN_ADDR``` environment variable can be necessary as NextJS binds to only one network interface by default. Set it to ```::``` to bind to any address.
