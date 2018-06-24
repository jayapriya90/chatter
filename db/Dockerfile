FROM postgres:9.6

ENV DEBIAN_FRONTEND=noninteractive

COPY create-db.sh /docker-entrypoint-initdb.d/
COPY init-db.sql /docker-entrypoint-initdb.d/

