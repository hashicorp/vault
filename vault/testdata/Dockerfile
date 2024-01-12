# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

FROM docker.mirror.hashicorp.services/ubuntu:22.04

ARG plugin

RUN groupadd nonroot && useradd -g nonroot nonroot

USER nonroot

COPY ${plugin} /bin/plugin

ENTRYPOINT [ "/bin/plugin" ]