# Envoy CrowdSec Bouncer (ext_proc)

Small gRPC service implementing Envoy's `ext_proc` API to bounce requests using CrowdSec decisions.

![coverage](https://raw.githubusercontent.com/FINTLabs/envoy-crowdsec-bouncer/refs/heads/badges/.badges/main/coverage.svg)

## Layout
- `cmd`: command.
- `internal/extproc`: Envoy ExternalProcessor implementation.
- `internal/crowdsec`: stubbed CrowdSec client to be replaced with LAPI calls.

## Running locally
```bash
# config can come from flags, env (prefix ENVOY_BOUNCER_), or --config (json/yaml)
ENVOY_BOUNCER_SERVER_GRPCPORT=9000 \
ENVOY_BOUNCER_SERVER_LOGLEVEL=info \
ENVOY_BOUNCER_BOUNCER_ENABLED=true \
ENVOY_BOUNCER_BOUNCER_LAPIURL="http://127.0.0.1:8080" \
ENVOY_BOUNCER_BOUNCER_APIKEY="your-api-key" \
ENVOY_BOUNCER_BOUNCER_METRICS=true \
ENVOY_BOUNCER_BOUNCER_TICKERINTERVAL="10s" \
ENVOY_BOUNCER_BOUNCER_METRICSINTERVAL="5m" \
ENVOY_BOUNCER_BOUNCER_BANSTATUSCODE=403 \
ENVOY_BOUNCER_WAF_ENABLED=false \
ENVOY_BOUNCER_WAF_APPSECURL="http://127.0.0.1:7422" \
ENVOY_BOUNCER_WAF_APIKEY="your-appsec-key" \
go run .

# equivalent via flags (see cmd/root.go for all options)
go run . \
  --server.grpcPort=9000 \
  --server.logLevel=info \
  --bouncer.enabled \
  --bouncer.lapiUrl=http://127.0.0.1:8080 \
  --bouncer.apiKey=your-api-key \
  --bouncer.metrics \
  --bouncer.tickerInterval=10s \
  --bouncer.metricsInterval=5m \
  --bouncer.banStatusCode=403 \
  --waf.enabled=false \
  --waf.appSecURL=http://127.0.0.1:7422 \
  --waf.apiKey=your-appsec-key

# use a config file instead
go run . --config config.yaml
```

The current client is a stub that always allows. Replace `internal/crowdsec` with a real LAPI implementation to enforce decisions.

### AppSec / partial body
- The ext_proc handler asks Envoy to stream request bodies (`ProcessingMode_STREAMED`), so you can plug CrowdSec AppSec to inspect partial chunks.
- Implement `EvaluateBody` in `internal/crowdsec` to call CrowdSec AppSec/LAPI with body fragments and decide mid-stream.

## Envoy configuration (example)
Add the ext_proc HTTP filter and point it at this service:
```yaml
http_filters:
  - name: envoy.filters.http.ext_proc
    typed_config:
      "@type": type.googleapis.com/envoy.extensions.filters.http.ext_proc.v3.ExternalProcessor
      grpc_service:
        envoy_grpc:
          cluster_name: crowdsec_ext_proc
      message_timeout: 0.5s
      failure_mode_allow: false
clusters:
  - name: crowdsec_ext_proc
    type: STRICT_DNS
    lb_policy: ROUND_ROBIN
    load_assignment:
      cluster_name: crowdsec_ext_proc
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: 127.0.0.1
                    port_value: 9000
```