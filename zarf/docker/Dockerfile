# Build the Go Binary.
FROM golang:1.18 AS base
ENV CGO_ENABLED 0
ARG BUILD_REF

# Copy the source code in the container.
COPY . /semantick

# Build the service binary.
WORKDIR /semantick
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go Binary in Alpine.
FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=base /semantick /semantick/semantick
WORKDIR /semantick
CMD ["./semantick"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="semantick" \
      org.opencontainers.image.authors="Bogdan Bojan <bogdanbojan03@gmail.com>" \
      org.opencontainers.image.source="https://github.com/bogdanbojan/semantick" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Bogdan Bojan"
