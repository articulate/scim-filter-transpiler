FROM golang:1.13

ENV GO111MODULE on
ENV SERVICE_ROOT /service
ENV SERVICE_USER service

RUN mkdir -p $SERVICE_ROOT
RUN groupadd $SERVICE_USER && useradd --create-home --home $SERVICE_ROOT --gid $SERVICE_USER --shell /bin/bash $SERVICE_USER
WORKDIR $SERVICE_ROOT
