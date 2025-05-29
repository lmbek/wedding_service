# wedding_service
![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)

A simple yet elegant web service for hosting your wedding website — including invitations, planned activities, RSVP forms, and other delightful features to ensure a memorable celebration.

This project handles:

- Hosting your wedding landing page
- Registering invitations and guest responses
- Displaying event schedule and important updates

---

## Installation

First install Go and Docker how you want them (preferably as CLI)
You can run the service with Go or Docker. Here's how:

### Prerequisites

- [Go 1.24+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Before running
First run
```bash
    go generate ./...
```

Please be aware:
This will generate the localhost certificates. (as this is the only go generate use in this project).

I have a script to put it into the windows certificate store automatically if you need it. (message me)

If you use windows then add it to cert store manually, if linux then put it into that cert store

## Running Locally

### Run tests

Run `make test` to execute all tests (`go test ./...`).

Please note that testing only involves testing the go service (without integration to sql etc.)

### Run benchmarks

Run `make bench` to execute benchmarks (`go test -bench=.`).

Please note that benching only involves testing the go service (without integration to sql etc.)

### Build and Run

Use the Makefile targets to handle everything smoothly:

- To build, clean, and start your service and Docker containers as daemons, simply run:

  `make` or `make all`

### Docker

If you want to run Docker commands separately, you can also use:

- `make build-docker` to build and start Docker containers in the foreground

- `make run-as-daemon` to start Docker containers as daemons

- `make stop-docker` to stop all Docker containers

### Go build

If you want to build or run go without docker use
- `make build-go`  builds go with go build

- `make rm-executable` removes the binary that was build

Please note that the name of the binary is managed from the .env file

---

## Deploy (Manual via Docker Save + SCP)

### Prerequisites

1. SSH key-based authentication
2. Mobile-based 2FA on the server (e.g., TOTP)
3. SSH hardening:
    - Disable root login
    - Disable password authentication (custom auth is best always)
4. Use a password-protected private key

---

### Step-by-step Deployment

**1. Build and Save the Docker Image**
```bash
    docker-compose up --build
    docker save wedding_service | gzip > wedding_service.tar.gz
```
**2. Transfer to Remote Server**
Note i use both windows and linux. On linux you would use id_rsa probably. But i write for windows users mostly.
Also here you would normally have private/public key, PAM (mobile auth) or a TOTP script and password.
Also use config file in ssh or some agent.
```bash
    scp -i ~/.ssh/yourprivatessh.key wedding_service.tar.gz user@example.com:/yourdestinationhere/
```
**3. Load and Run on the Server**
Also note, you can zip and unzip with other things than gzip
```bash
    ssh -i ~/.ssh/yourprivatessh.key user@example.com

    # On the server:
    cd /yourdestinationhere
    gunzip wedding_service.tar.gz
    docker load < wedding_service.tar

    # Start the container
    docker run -d -p 80:8080 wedding_service
```
---

## Bonus: Autodeploy (Recommended)

Build a small auto-deploy service on your server that:

- Accepts uploaded `.tar.gz` images
- Validates with a signed token or secret key
- Automatically:
    - Loads the image
    - Restarts the container
- Sends notifications via whatever

Optional improvements:

- GitHub webhook listener
- SHA256 checksum validation
- Logging and rollback on failure

---

## Security Tips

- Use a firewall (e.g., UFW); only expose ports 22, 80, and 443
- Force HTTPS (Caddy, Nginx, or Traefik recommended. Or develop your own, because using third party is always overkill if you ask me)
- Protect admin/RSVP routes:
    - Use secret URLs or tokens
    - Optionally require basic auth
- Sanitize all inputs to avoid injection vulnerabilities
- Lock down your database with least-privilege and IP rules

---


