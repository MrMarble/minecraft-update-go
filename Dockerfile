ARG GO_VERSION=1.17

## Build container
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /minecraft /src/cmd/minecraft

## Final container
FROM scratch AS final

COPY --from=builder /minecraft /minecraft

ENTRYPOINT ["/minecraft"]