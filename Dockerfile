# Build Stage
ARG BUILDPLATFORM
FROM --platform=${BUILDPLATFORM} golang:1.24.4 AS build-stage

LABEL app="gutenberg-ingester"
LABEL REPO="https://github.com/astr0n8t/gutenberg-ingester"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN mkdir -p /var/gutenberg-ingester /data

COPY *.go ./
# Copy all internal modules
COPY cmd/*.go ./cmd/
COPY pkg/ ./pkg/
COPY internal/*.go ./internal/
COPY version/*.go ./version/

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /gutenberg-ingester

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /gutenberg-ingester /gutenberg-ingester
COPY --from=build-stage --chown=nonroot:nonroot /data /data
COPY --from=build-stage --chown=nonroot:nonroot /var/gutenberg-ingester /var/gutenberg-ingester

USER nonroot:nonroot

WORKDIR /data

ENTRYPOINT ["/gutenberg-ingester"]
