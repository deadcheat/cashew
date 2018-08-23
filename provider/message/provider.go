package message

// Provider hold messages
type Provider struct {
	info []string
	errs []string
}

// New return new provider
func New() *Provider {
	info := make([]string, 0)
	errs := make([]string, 0)
	return &Provider{
		info: info,
		errs: errs,
	}
}

// AddInfo add info message into inner slice
func (p *Provider) AddInfo(msg string) {
	p.info = append(p.info, msg)
}

// AddErr add info message into inner slice
func (p *Provider) AddErr(msg string) {
	p.errs = append(p.errs, msg)
}

// Info return string slice
func (p *Provider) Info() []string {
	return p.info
}

// Errors return string slice
func (p *Provider) Errors() []string {
	return p.errs
}
