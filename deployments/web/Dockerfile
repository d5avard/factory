FROM golang:1.23 AS build

ARG USER="webuser"
ENV USER="${USER}"

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    "${USER}"

RUN apt-get update && apt-get install -y ca-certificates
RUN git clone https://github.com/d5avard/factory.git

WORKDIR /go/factory
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/web ./cmd/web/web.go

### Final stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

ARG USER

USER ${USER}:${USER}

ARG APPLICATION="web"
ARG DESCRIPTION="A simple web server that serves static files and a health check endpoint"
ARG VERSION="0.1.0"
ARG VCS_REF="unknown"
ARG BUILD_DATE="unknown"    
ARG PACKAGE="d5avard/factory"

LABEL org.opencontainers.image.ref.name="${PACKAGE}" \
    org.opencontainers.image.authors="Dany Savard>" \
    org.opencontainers.image.documentation="https://github.com/${PACKAGE}/README.md" \
    org.opencontainers.image.description="${DESCRIPTION}" \
    org.opencontainers.image.licenses="Apache 2.0" \
    org.opencontainers.image.source="https://github.com/${PACKAGE}"

COPY --from=build /go/bin/${APPLICATION} /web
CMD ["/web", "--port", "80"]