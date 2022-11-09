FROM docker.io/golang:1.19-alpine AS builder
ADD . /work/
WORKDIR /work
RUN CGO_ENABLED=0 go build -o /work/bin/quicklink

FROM docker.io/alpine
COPY --from=builder /work/bin/* /usr/bin/
ENV ADDR=:80
CMD /usr/bin/quicklink -pg "$PGDATABASE" -addr "$ADDR"
