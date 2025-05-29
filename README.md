# wedding_service

A simple yet elegant web service for hosting your wedding website — including invitations, planned activities, RSVP forms, and other delightful features to ensure a memorable celebration.

This project handles:

- Hosting your wedding landing page
- Registering invitations and guest responses
- Displaying event schedule and important updates

---

## Installation

You can run the service with Go or Docker. Here's how:

### Prerequisites

- [Go 1.24+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Running Locally (Go)

### Run tests

```bash
go test ./...
```

### Run go

```bash
go run .
```

### Build go

```bash
go build .
```

### Build docker container
```bash
docker-compose up --build
```

### Start docker container
```bash
docker-compose up -d
```

### Stop docker container
```bash
docker-compose down
```

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


