FROM gcr.io/distroless/static-debian11

ARG plugin

COPY ${plugin} /bin/plugin

ENTRYPOINT [ "/bin/plugin" ]