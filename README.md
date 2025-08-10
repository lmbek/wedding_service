# wedding_service
<!--![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)-->

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
    make go-generate
```

This command does two important things:
1. Generates self-signed certificates for HTTPS
2. Generates Swagger API documentation

### Certificate Generation and Usage

The application uses self-signed certificates for HTTPS. These certificates are generated using the `go-generate` command and are stored in the `webserver/certificate/self_sign_cert` directory.

#### Certificate Files

The following certificate files are generated:
- `localhost_wedding_service.crt`: The certificate file
- `localhost_wedding_service.key`: The private key file
- `localhost_wedding_service.pem`: A combined PEM file containing both the certificate and key

#### Docker Integration

When running the application in Docker, the certificates are:
1. Generated during the Docker build process
2. Copied to the appropriate locations in the container
3. Stored in a named volume (`wedding-certificates`) for persistence between container restarts

This ensures that the certificates are always available to the application, even after rebuilding or restarting the container.

#### Importing Certificates for Browser Use

To use the application with HTTPS in your browser without security warnings, you need to import the certificate into your browser or operating system's certificate store:

**For Windows:**
1. Open the certificate file (`webserver/certificate/self_sign_cert/localhost_wedding_service.crt`)
2. Click "Install Certificate"
3. Select "Current User" or "Local Machine"
4. Select "Place all certificates in the following store"
5. Click "Browse" and select "Trusted Root Certification Authorities"
6. Click "Next" and then "Finish"

**For macOS:**
1. Open Keychain Access
2. Import the certificate file (`webserver/certificate/self_sign_cert/localhost_wedding_service.crt`)
3. Find the imported certificate in the list
4. Double-click on it
5. Expand the "Trust" section
6. Set "When using this certificate" to "Always Trust"

**For Linux (Ubuntu/Debian):**
```bash
sudo cp webserver/certificate/self_sign_cert/localhost_wedding_service.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

**For Firefox (all platforms):**
Firefox uses its own certificate store, so you need to import the certificate separately:
1. Open Firefox
2. Go to Settings/Preferences
3. Search for "certificates"
4. Click "View Certificates"
5. Go to the "Authorities" tab
6. Click "Import" and select the certificate file
7. Check "Trust this CA to identify websites" and click "OK"

## Running Locally

### Run tests

Run `make test` to execute all tests (`go test ./...`).

Please note that testing only involves testing the go service (without integration to sql etc.)

### Run test coverage (and show in browser)

Run `make test-coverage-html` to execute all tests and show test coverage.


### Run benchmarks

Run `make bench` to execute benchmarks (`go test -bench=.`).

Please note that benching only involves testing the go service (without integration to sql etc.)

### Build and Run

Use the Makefile targets to handle everything smoothly:

- To build, clean, and start your service and Docker containers as daemons, simply run:

  `make` or `make all`

This starts the docker-compose so you dont have to run anything in docker other than

`make down`

When you want to stop the containers

### Docker

If you want to run Docker commands separately, you can also use:

- `make docker-build` to build and start Docker containers in the foreground

- `make run-as-daemon` to start Docker containers as daemons

- `make docker-stop` to stop all Docker containers

### Go build

If you want to build or run go without docker use
- `make go-build`  builds go with go build

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
    make docker-build
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

## Project Structure

The project follows a clean, standard Go project structure:

```
wedding_service/
├── services/                # Docker-related files and service configurations
│   ├── grafana/              # Grafana configuration
│   │   ├── dashboards/       # Grafana dashboards
│   │   └── provisioning/     # Grafana provisioning
│   ├── mysql/                # MySQL configuration
│   └── wedding-service/      # Wedding service configuration and Dockerfile
├── bin/                      # Compiled binaries (for local development)
├── webserver/                # Webserver code
│   ├── api/                  # API handlers
│   ├── certificate/          # Certificate management
│   ├── database/             # Database access
│   └── website/              # Website handlers
├── buildtag/                 # Build tags
├── config/                   # Configuration
├── logging/                  # Logging utilities
├── metrics/                  # Metrics collection
├── .env                      # Environment variables
├── docker-compose.yml        # Docker Compose configuration
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── main.go                   # Application entry point
├── Makefile                  # Build automation
└── README.md                 # This file
```

## Docker Configuration

The application now uses a multi-stage Docker build process defined in `services/wedding-service/Dockerfile`. This approach:

1. Builds the Go application in a builder stage
2. Copies only the necessary files to a minimal Alpine image
3. Configures the runtime environment

The Docker Compose configuration has been simplified to use this Dockerfile instead of mounting volumes directly.

## Docker Volume Configuration

The application uses named volumes for better management and portability:

- `wedding-mysql-data`: MySQL database data
- `wedding-grafana-data`: Grafana configuration and data
- `wedding-certificates`: SSL certificates for the wedding service
- `wedding-prometheus-data`: Prometheus time-series data

Benefits of this approach:
- Improved performance compared to bind mounts
- Better portability across different environments
- Easier backup and restore
- Cleaner separation of concerns
- Persistence of certificates between container restarts
- Retention of historical metrics data

## Metrics and Monitoring

The application uses a comprehensive metrics and monitoring stack:

1. **Wedding Service**: Exposes custom metrics at the `/metrics` endpoint
2. **Prometheus**: Collects and stores metrics from the wedding service
3. **Grafana**: Visualizes metrics from Prometheus

### Custom Metrics

The application exposes the following custom metrics:

- **HTTP Request Metrics**:
  - `wedding_http_requests_total`: Counter of HTTP requests (labels: method, path, status)
  - `wedding_http_request_duration_seconds`: Histogram of request durations (labels: method, path, status)
  - `wedding_http_response_size_bytes`: Histogram of response sizes (labels: method, path, status)
  - `wedding_http_active_connections`: Gauge of active HTTP connections

- **Database Metrics**:
  - `wedding_database_queries_total`: Counter of database queries (labels: operation, table)
  - `wedding_database_query_duration_seconds`: Histogram of query durations (labels: operation, table)

- **Error Metrics**:
  - `wedding_errors_total`: Counter of errors (labels: type, code)

### Prometheus

Prometheus is a time-series database that collects and stores metrics. It's configured to scrape metrics from the wedding service every 15 seconds. The Prometheus server is accessible at:

- **URL**: http://localhost:9090
- **Features**:
  - Query metrics using PromQL
  - View metric graphs
  - Set up alerts (not configured by default)
  - Explore metric targets and their health

### Grafana Dashboards

The project includes pre-configured Grafana dashboards for visualizing metrics:

1. **HTTP Requests Dashboard**: Visualizes HTTP request metrics
2. **Database Dashboard**: Visualizes database query metrics
3. **System Dashboard**: Visualizes system metrics (CPU, memory, disk, network)

Grafana is configured to use Prometheus as its data source, allowing it to visualize all metrics collected by Prometheus.

For detailed setup instructions, see the [Grafana Setup Guide](services/grafana/README.md).

## Logging

This project uses structured logging with Go's `slog` package (introduced in Go 1.22) for a simple yet powerful logging solution.

### Logging Features

- **Structured JSON Logs**: All logs are output in JSON format for easy parsing
- **Log Levels**: Support for Debug, Info, Warning, and Error levels
- **Configurable Output**: Logs can be directed to stdout (default) or a file
- **Source Information**: Optional inclusion of source file and line numbers

## Environment Configuration

The application uses a `.env` file for configuration. The file is organized into logical sections:

### Application Settings
- `APP_NAME`: Name of the application (default: wedding-go-service)
- `DEBUG`: Sets the log level (0=None, 1=Error, 2=Warning, 3=Info, 4=Debug)

### Server Configuration
- `WEDDING_SERVICE_HTTP_PORT`: HTTP port for the service (default: 8080)
- `WEDDING_SERVICE_HTTPS_PORT`: HTTPS port for the service (default: 8443)
- `WEDDING_SERVICE_HOSTNAMES`: Hostnames for the service (default: localhost)

### Certificate Paths
- `SELF_SIGNED_CERT_PATH`: Path to the self-signed certificate
- `SELF_SIGNED_KEY_PATH`: Path to the self-signed key

## Host Protection

The application includes a middleware that protects against unauthorized hosts by checking the `Host` header of each request. This helps prevent DNS rebinding attacks and ensures that only authorized domains can access your service.

### Configuring Protected Hosts

Protected hosts are configured using the `WEDDING_SERVICE_HOSTNAMES` environment variable in the `.env` file. This variable supports a flexible format that allows specifying multiple domains and their aliases:

```
WEDDING_SERVICE_HOSTNAMES=domain1:alias1,alias2|domain2:alias3,alias4
```

Where:
- `domain1`, `domain2` are the primary domain names
- `alias1`, `alias2`, `alias3`, `alias4` are aliases for those domains

For example:
```
WEDDING_SERVICE_HOSTNAMES=example.com:www.example.com,api.example.com|example.org:www.example.org
```

This configuration would allow requests to:
- example.com
- www.example.com
- api.example.com
- example.org
- www.example.org

For local development, you can simply use:
```
WEDDING_SERVICE_HOSTNAMES=localhost
```

### How Host Protection Works

1. When the application starts, it parses the `WEDDING_SERVICE_HOSTNAMES` variable into a map of allowed hosts
2. For each incoming request, the middleware checks if the `Host` header matches any of the allowed hosts
3. If the host is allowed, the request proceeds normally
4. If the host is not allowed, the middleware returns a 403 Forbidden response

WebSocket connections are exempt from host validation to allow for hot reloading during development.

### MySQL Configuration
- `MYSQL_SERVICE_HOSTNAME`: Hostname for the MySQL service
- `MYSQL_SERVICE_HTTP_PORT`: Port for the MySQL service
- `MYSQL_SERVICE_USER`: Username for the MySQL service
- `MYSQL_SERVICE_PASSWORD`: Password for the MySQL service
- `MYSQL_SERVICE_ROOT_PASSWORD`: Root password for the MySQL service
- `MYSQL_SERVICE_DATABASE`: Database name for the MySQL service

### Docker Build Configuration
- `BUILD_DIR`: Directory for Docker build files

### Logging Configuration
- `LOG_FORMAT`: Set to "text" for human-readable logs (default is JSON)
- `LOG_SOURCE`: Set to "true" to include source file information
- `LOG_FILE`: Path to a log file (if not set, logs go to stdout)

### Viewing Logs

Logs can be viewed in several ways:

1. **Console Output**: When running locally, logs appear in your terminal
2. **Docker Logs**: Use `docker logs wedding-go-service` to view container logs
3. **Grafana Dashboard**: Access the Grafana UI at http://localhost:3000 to view logs
   - Navigate to Explore and select the Docker Logs data source
   - Use the container filter to select `wedding-go-service`
   - The logs are in JSON format and can be parsed using LogQL

### Grafana Setup

The project includes a pre-configured Grafana setup for metrics and log visualization:

1. Start the services with `make all`
2. Open Grafana at http://localhost:3000 (default credentials: admin/admin)
3. Navigate to Dashboards to view the pre-configured dashboards
4. For logs, navigate to Explore and select the Docker Logs data source

#### Troubleshooting Grafana Connection Issues

If Grafana cannot connect to Prometheus:

1. Check that all services are running with `docker ps`
2. Verify that the Prometheus service is healthy with `docker ps` (look for "healthy" status)
3. Check the Grafana datasource configuration in `services/grafana/provisioning/datasources/prometheus.yml`
4. Ensure the URL is set to `http://wedding-prometheus:9090`
5. Check that Prometheus can scrape metrics from the wedding-go-service:
   - Open Prometheus at http://localhost:9090
   - Go to Status > Targets to see if the wedding-go-service target is "UP"
   - If not, check the Prometheus configuration in `services/prometheus/prometheus.yml`
6. Restart the services with `docker-compose down && docker-compose up -d`

If Prometheus cannot scrape metrics from the wedding-go-service:

1. Verify that the wedding-go-service is exposing metrics at `/metrics`
2. Check that the wedding-go-service is accessible from Prometheus:
   - From the Prometheus container: `docker exec -it wedding-prometheus wget -q --no-check-certificate --spider https://wedding-go-service:8443/metrics`
3. Check the Prometheus configuration in `services/prometheus/prometheus.yml`:
   - Ensure the target is set to `wedding-go-service:8443`
   - Ensure the scheme is set to `https`
   - Ensure TLS verification is disabled with `insecure_skip_verify: true`

Note: The metrics endpoint is available on the HTTPS port (8443) because the HTTP server redirects to HTTPS. The `insecure_skip_verify` option is needed because we're using a self-signed certificate.
