FROM golang as build

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -v

# ---

FROM postgres:alpine

USER postgres

COPY --from=build /go/src/app/container-angel /entrypoint

ENTRYPOINT ["/entrypoint"]
# ENTRYPOINT ["docker-entrypoint.sh"]
# CMD ["postgres"]
