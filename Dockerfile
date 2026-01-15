FROM golang:alpine AS builder

# Build arguments for version info
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

LABEL maintainer="cylonchau"

WORKDIR /apps
COPY ./ /apps

RUN apk add --no-cache git upx bash make ca-certificates tzdata && \
    # If VERSION is 'dev', try to get git info
    if [ "$VERSION" = "dev" ]; then \
        VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
        GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
    fi && \
    # Build with version info
    go build -ldflags="-s -w \
        -X 'github.com/cylonchau/prism/pkg/cli.Version=${VERSION}' \
        -X 'github.com/cylonchau/prism/pkg/cli.GitCommit=${GIT_COMMIT}' \
        -X 'github.com/cylonchau/prism/pkg/cli.BuildDate=${BUILD_DATE}'" \
        -o _output/prism ./cmd/prism && \
    upx -1 _output/prism && \
    chmod +x _output/prism

FROM alpine:latest AS runner

# Read Terraform version from config file
ARG TERRAFORM_VERSION
COPY .terraform-version /tmp/.terraform-version

# Install runtime dependencies and Terraform
RUN apk add --no-cache ca-certificates tzdata curl unzip && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    # Read Terraform version from file if not provided as build arg
    if [ -z "$TERRAFORM_VERSION" ]; then \
        TERRAFORM_VERSION=$(grep terraform_version /tmp/.terraform-version | cut -d'=' -f2); \
    fi && \
    echo "Installing Terraform ${TERRAFORM_VERSION}..." && \
    # Detect architecture
    ARCH=$(uname -m) && \
    case ${ARCH} in \
        x86_64) TF_ARCH="amd64" ;; \
        aarch64) TF_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: ${ARCH}" && exit 1 ;; \
    esac && \
    # Download and install Terraform
    curl -fsSL "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TF_ARCH}.zip" -o /tmp/terraform.zip && \
    unzip /tmp/terraform.zip -d /usr/local/bin && \
    chmod +x /usr/local/bin/terraform && \
    # Verify installation
    terraform version && \
    # Cleanup
    rm -rf /tmp/* && \
    apk del curl unzip

WORKDIR /apps

# Copy binary from builder
COPY --from=builder /apps/_output/prism /usr/bin/prism

# Create non-root user
RUN addgroup -g 1000 prism && \
    adduser -D -u 1000 -G prism prism && \
    mkdir -p /apps/.terraform.d /apps/workdir && \
    chown -R prism:prism /apps

USER prism

VOLUME ["/apps/workdir"]

ENTRYPOINT ["prism"]
CMD ["--help"]

EXPOSE 8080/tcp