FROM golang:1.23 AS build

ARG USER="webuser"
ENV USER="${USER}"

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    "${USER}" \
    && groupadd tlsaccess \
    && usermod -aG tlsaccess ${USER}

RUN mkdir -p /etc/letsencrypt/live/ \
    && mkdir -p /etc/letsencrypt/archive/ \
    && chgrp -R tlsaccess /etc/letsencrypt/live/ \
    && chgrp -R tlsaccess /etc/letsencrypt/archive/ \
    && chmod -R 750 /etc/letsencrypt/live/ \
    && chmod -R 750 /etc/letsencrypt/archive/

RUN apt-get update && apt-get install -y ca-certificates

COPY / /factory/
WORKDIR /factory

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/web ./cmd/web/web.go

### Final stage
FROM scratch

ARG APPLICATION="web"
ARG DESCRIPTION="A simple web server that serves static files and a health check endpoint"    
ARG PACKAGE="d5avard/factory"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /factory/web/web/templates /app/templates/
COPY --from=build /go/bin/${APPLICATION} /app/${APPLICATION}
COPY --from=build /etc/letsencrypt/live/ /etc/letsencrypt/live/
COPY --from=build /etc/letsencrypt/archive/ /etc/letsencrypt/archive/

ENV PATH="${PATH}:/app"

ARG USER
USER ${USER}:${USER}

LABEL org.opencontainers.image.ref.name="${PACKAGE}" \
    org.opencontainers.image.authors="Dany Savard>" \
    org.opencontainers.image.documentation="https://github.com/${PACKAGE}/README.md" \
    org.opencontainers.image.description="${DESCRIPTION}" \
    org.opencontainers.image.licenses="Apache 2.0" \
    org.opencontainers.image.source="https://github.com/${PACKAGE}"

WORKDIR /app

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["web"]