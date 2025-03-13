package smtp

type SMTPConfig struct {
	Server   string
	Port     string
	Sender   string
	Password string
}

func NewSMTPClient(cfg SMTPConfig) *SMTPConfig {
	return &SMTPConfig{
		Server:   cfg.Server,
		Port:     cfg.Port,
		Sender:   cfg.Sender,
		Password: cfg.Password,
	}
}
