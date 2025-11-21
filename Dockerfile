FROM gcr.io/distroless/static:nonroot

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/coredns /coredns

EXPOSE 53/udp
EXPOSE 53/tcp

ENTRYPOINT ["/coredns"]
