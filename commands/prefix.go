package commands

// Register struct class
type Register struct{}

func handle(args string) error {
	main.Config.Token = args
	return nil
}

func (r *Register) description() string { return "Changes command prefix." }
func (r *Register) usage() string       { return "<newprefix>" }
