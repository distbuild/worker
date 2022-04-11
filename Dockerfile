FROM rust:latest AS builder
WORKDIR /usr/src/worker
COPY . .
RUN make install && \
    make build

FROM scratch
COPY --from=builder /usr/src/worker/target/release/worker /usr/local/bin/worker
ENTRYPOINT ["worker"]
