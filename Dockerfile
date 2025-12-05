FROM golang:1.25.5-trixie

RUN apt-get update && apt-get install -y \
    build-essential \
    wget \
    unzip \
    && rm -rf /var/lib/apt/lists/*

ENV Z3_VERSION=4.15.4
RUN wget https://github.com/Z3Prover/z3/releases/download/z3-${Z3_VERSION}/z3-${Z3_VERSION}-x64-glibc-2.39.zip \
    && unzip z3-${Z3_VERSION}-x64-glibc-2.39.zip \
    && mv z3-${Z3_VERSION}-x64-glibc-2.39 /opt/z3 \
    && rm z3-${Z3_VERSION}-x64-glibc-2.39.zip

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/opt/z3/include"
ENV CGO_LDFLAGS="-L/opt/z3/bin"
ENV LD_LIBRARY_PATH="/opt/z3/bin"

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN /opt/z3/bin/z3 --version

CMD ["go", "test", "-v", "./..."]
