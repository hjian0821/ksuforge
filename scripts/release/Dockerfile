FROM golang:1.25-bookworm

ARG WAILS_VERSION=v2.15.0

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        gcc-mingw-w64-x86-64 \
        git \
        libgtk-3-dev \
        libwebkit2gtk-4.1-dev \
        make \
        pkg-config \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/wailsapp/wails/v2/cmd/wails@${WAILS_VERSION}

# The source tree is bind-mounted from the host and intentionally owned by the
# invoking user rather than root. Allow Git (and Go VCS stamping) to inspect it.
RUN git config --system --add safe.directory /src

ENV PATH="/go/bin:${PATH}"
WORKDIR /src
