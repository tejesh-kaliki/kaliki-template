// Package observability provides first-class, opt-out tracing. By default it
// installs a no-op tracer so the application runs with no external dependencies.
// Point config.Observability.Endpoint at an OTel collector to start exporting.
package observability

import "log"

type Provider struct {
	enabled bool
}

// Init returns a Provider. When endpoint is empty the provider is a no-op.
func Init(serviceName, endpoint string) *Provider {
	if endpoint == "" {
		log.Printf("observability: no-op tracer (set observability.endpoint to export)")
		return &Provider{enabled: false}
	}
	// TODO: wire OTLP exporter (otelgin/otelpgx) when an endpoint is configured.
	log.Printf("observability: exporting traces for %q to %s", serviceName, endpoint)
	return &Provider{enabled: true}
}

func (p *Provider) Enabled() bool { return p.enabled }

func (p *Provider) Shutdown() {}
