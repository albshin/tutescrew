package config

// CASConfig holds env var related to CAS
type CASConfig struct {
	AuthURL     string `required:"true"`
	RedirectURL string `required:"true"`
}
