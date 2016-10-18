FROM golang:1.7.1

ENV USER root

RUN go get github.com/mitchellh/gox \
	    && go get golang.org/x/tools/cmd/cover

EXPOSE 8200

WORKDIR $GOPATH/src/github.com/hashicorp/vault/
CMD ["/bin/bash"]
