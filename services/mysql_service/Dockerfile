# Minimal custom MySQL image for the project
# We keep it simple and rely on docker-compose to pass env vars from .env
# and to mount the persistent data volume.

FROM mysql:8.4

# Optional: set a sane default for character set/collation at server level.
# You may add a custom my.cnf later if needed by uncommenting the COPY below.
# COPY my.cnf /etc/mysql/conf.d/my.cnf

# Expose default MySQL port
EXPOSE 3306

# The official image entrypoint handles initialization via
# /docker-entrypoint-initdb.d when present. We provide schema + seed SQL files.
COPY services/mysql_service/initdb/ /docker-entrypoint-initdb.d/
