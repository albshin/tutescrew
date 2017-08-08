package commands

// Prefix struct class
type Prefix struct{}

func (p *Prefix) handle() error {
	return nil
}

func (p *Prefix) description() string { return "Changes command prefix." }
func (p *Prefix) usage() string       { return "<newprefix>" }
