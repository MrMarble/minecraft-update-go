ARG GO_VERSION=1.17

# Install certificates
FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group && \
    apk add --no-cache ca-certificates

FROM scratch as final

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY minecraft /
USER 65534:65534
ENTRYPOINT ["/minecraft"]