FROM gcr.io/distroless/static:nonroot

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/coredns /coredns

EXPOSE 53 53/udp

WORKDIR /
USER nonroot:nonroot
ENTRYPOINT ["/coredns"]
