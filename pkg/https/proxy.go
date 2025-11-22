package https

import "sync/atomic"

// GoProxyProvider is an interface for a Go proxy provider
type GoProxyProvider interface {
	// GetProxy returns the host and secret of a Go proxy
	GetProxy() (host string, secret string)
}

// WithGoProxyProvider sets the request to use a Go proxy provider
func WithGoProxyProvider(provider GoProxyProvider) func(cfg *Options) {
	return func(cfg *Options) {
		cfg.proxyProvider = provider
	}
}

// rrProxyProvider is a round-robin proxy provider
type rrProxyProvider struct {
	hosts  []string
	secret string
	index  int32
}

// GetProxy returns the host and secret of a Go proxy
func (p *rrProxyProvider) GetProxy() (host string, secret string) {
	id := atomic.AddInt32(&p.index, 1)
	return p.hosts[id%int32(len(p.hosts))], p.secret
}

// NewRRProxyProvider creates a new round-robin proxy provider
func NewRRProxyProvider(hosts []string, secret string) GoProxyProvider {
	return &rrProxyProvider{hosts, secret, -1}
}