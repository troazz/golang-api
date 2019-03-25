FROM mysql

ENV MYSQL_DATABASE=mydb \
    MYSQL_ROOT_PASSWORD=password \
    MYSQL_ALLOW_EMPTY_PASSWORD=no

ADD schema.sql /docker-entrypoint-initdb.d

EXPOSE 3306