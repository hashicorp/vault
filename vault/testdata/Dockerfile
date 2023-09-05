FROM docker.mirror.hashicorp.services/ubuntu:22.04

ARG plugin

COPY ${plugin} /bin/plugin

ENTRYPOINT [ "/bin/plugin" ]