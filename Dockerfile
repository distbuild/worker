FROM rust:latest AS builder
WORKDIR /usr/src/worker
COPY . .
RUN make install && \
    make build

FROM gcr.io/distroless/base-debian11
COPY --from=builder /usr/src/worker/target/release/worker /usr/local/bin/worker
USER nonroot:nonroot
ENTRYPOINT ["worker"]
