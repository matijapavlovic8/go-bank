FROM alpine:3.17.0

ARG GLIBC_VERSION=2.34-r0

RUN apk --no-cache add wget && \
    wget -q -nv -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget -q -nv https://github.com/sgerrand/alpine-pkg-glibc/releases/download/$GLIBC_VERSION/glibc-$GLIBC_VERSION.apk && \
    apk add --no-cache --force-overwrite glibc-$GLIBC_VERSION.apk && \
    rm glibc-$GLIBC_VERSION.apk && \
    apk del wget

WORKDIR /app

COPY . /app

ENTRYPOINT ["/app/myapp"]