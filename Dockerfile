FROM golang:1.25-alpine AS builder
WORKDIR /workspace

RUN apk add --no-cache git bash

COPY . .

RUN hack/build-coredns-with-llm.sh

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /workspace/_build/coredns/coredns /coredns

EXPOSE 53/udp
EXPOSE 53/tcp

ENTRYPOINT ["/coredns"]
