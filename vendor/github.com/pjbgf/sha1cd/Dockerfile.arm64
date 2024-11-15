FROM golang:1.20@sha256:2edf6aab2d57644f3fe7407132a0d1770846867465a39c2083770cf62734b05d

ENV GOOS=linux
ENV GOARCH=arm64
ENV CGO_ENABLED=1
ENV CC=aarch64-linux-gnu-gcc
ENV PATH="/go/bin/${GOOS}_${GOARCH}:${PATH}"
ENV PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig

# install build & runtime dependencies
RUN dpkg --add-architecture arm64 \
    && apt update \
    && apt install -y --no-install-recommends \
        gcc-aarch64-linux-gnu \
        libc6-dev-arm64-cross \
        pkg-config \
    && rm -rf /var/lib/apt/lists/*

COPY . /src/workdir

WORKDIR /src/workdir

RUN go build ./...
