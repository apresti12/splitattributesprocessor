FROM golang:1.23 as build
ARG OTEL_VERSION=0.113.0
WORKDIR /app
COPY . .
#RUN go install go.opentelemetry.io/collector/cmd/builder@v0.113.0

RUN curl --proto '=https' --tlsv1.2 -fL -o ocb \
    https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.113.0/ocb_0.113.0_linux_amd64
RUN chmod +x ocb


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ./ocb --config=builder-config.yaml

# 4317 - default OTLP receiver
# 55678 - op
# 55679 - zpages
EXPOSE 4317/tcp 55678/tcp 55679/tcp 13133/tcp 4318/tcp
CMD ["--config", "config.yaml"]
ENTRYPOINT ["./otelcol-dev/otelcol-dev"]



