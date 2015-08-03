FROM golang:1.5
MAINTAINER sthysel <sthysel@gmail.com>
ENV REFRESHED_AT 2015-08-03

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
  apt-transport-https \
  build-essential \
  git \
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN env --unset=DEBIAN_FRONTEND

# WORKDIR is set to $GOPATH in master image
RUN git clone https://github.com/hashicorp/vault ${GOPATH}/src/github.com/hashicorp/vault
WORKDIR ${GOPATH}/src/github.com/hashicorp/vault
RUN make bootstrap
RUN make dev

CMD ["/bin/bash"]
