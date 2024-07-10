package config

type Config struct {
	// Mode is one among "debug", "prod", "test"
	Mode string

	// Secret is used to encrypt session cookies
	Secret string

	// BaseURL is used in emails to retrieve assets from the webserver (needs protocol)
	BaseURL string `yaml:"baseURL"`

	// ServerAddress is the address the web server will listen to
	ServerAddress string `yaml:"serverAddress"`

	Database struct {
		// Address is the hostname:port of the database
		Address string

		// Username is the username of the database
		Username string

		// Password is the password of the database
		Password string

		// Database is the name of the database
		Database string
	}

	Email struct {
		// Enabled indicates if the emails will be sent or not (set to false only for testing)
		Enabled bool

		// Address is the sender address
		Address string

		// Server is the address of the smtp server
		Server string

		// Port is the port of the smtp server
		Port string

		// Login is the login value of the smtp server
		Login string

		// Password is the password of the smtp server
		Password string
	}

	// SupportEmail is the email visible in the info page
	SupportEmail string `yaml:"supportEmail"`

	// RepoLink is the link to the repository visible in the info page
	RepoLink string `yaml:"repoLink"`

	// StoreLink is the link to the Play Store visible in the info page
	StoreLink string `yaml:"storeLink"`
}

var Runtime Config
