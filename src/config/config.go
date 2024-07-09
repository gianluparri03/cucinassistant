package config

type Config struct {
	Debug bool

	BaseURL       string `yaml:"baseURL"`
	ServerAddress string `yaml:"serverAddress"`

	Database struct {
		Address  string
		Username string
		Password string
		Database string
	}

	Email struct {
		Enabled  bool
		Address  string
		Server   string
		Port     string
		Login    string
		Password string
	}

	SupportEmail string `yaml:"supportEmail"`

	RepoLink  string `yaml:"repoLink"`
	StoreLink string `yaml:"storeLink"`
}

var Runtime Config
