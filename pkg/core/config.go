package core

// GithubAppConfig represents Github App config
type GithubAppConfig struct {
	AppID      int64
	PrivateKey []byte
	InstallURL string
}
