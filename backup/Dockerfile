FROM centos/mysql-57-centos7

ENV MYSQL_USER=user
ENV MYSQL_PASSWORD=pass
ENV MYSQL_DATABASE=db

EXPOSE 3310

ADD . /go/src/app
WORKDIR /go/src/app

VOLUME [ "./dump:/docker-entrypoint-initdb.d" ]
CMD ["go", "run", "main.go"]