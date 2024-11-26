# start with base image
FROM postgres:latest

# import data into container
# All scripts in docker-entrypoint-initdb.d/ are automatically executed during container startup
COPY ./database/migration.sql /docker-entrypoint-initdb.d/
RUN echo "Copy migration.sql into /docker-entrypoint-initdb.d/ completed"
# docker exec -it test_db psql -U postgres -d cockroachdb -f /docker-entrypoint-initdb.d/migration.sql