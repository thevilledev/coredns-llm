ARG TARGETOS
ARG TARGETARCH

FROM gcr.io/distroless/static:nonroot

COPY ./linux/${TARGETARCH}/coredns /coredns

EXPOSE 53/udp
EXPOSE 53/tcp

ENTRYPOINT ["/coredns"]
