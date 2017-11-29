# Adapted from tcnksm/dockerfile-gox -- thanks!

FROM debian:stable

RUN apt-get update -y && apt-get install --no-install-recommends -y -q \
                         curl \
                         zip \
                         build-essential \
                         ca-certificates \
                         git mercurial bzr \
               && rm -rf /var/lib/apt/lists/*

ENV GOVERSION 1.9.2
RUN mkdir /goroot && mkdir /gopath
RUN curl https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz \
           | tar xvzf - -C /goroot --strip-components=1

ENV GOPATH /gopath
ENV GOROOT /goroot
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

RUN go get github.com/mitchellh/gox

RUN mkdir -p /gopath/src/github.com/hashicorp/vault
WORKDIR /gopath/src/github.com/hashicorp/vault
CMD make bin
