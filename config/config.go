package config

// Config contains all the fields read from a config file
type Config struct {
	// Debugging, if enabled, will add a shutdown prompt
	Debugging bool

	// Testing, if enabled, will disable emails
	Testing bool

	// SessionSecret is used to encrypt session cookies
	SessionSecret string

	// BaseURL is an URL from which CucinAssistant will be accessed,
	// like https://ca.gianlucaparri.me
	BaseURL string `yaml:"baseURL"`

	// Port is the port the server will listen to
	Port string

	// Database is the postgresql's connection string
	Database string

	Email struct {
		// Address is the sender address
		Address string

		// Server is the address of the smtp server
		Server string

		// Port is the port of the smtp server
		Port string

		// Password is the password of the smtp server
		Password string
	}
}

// Runtime contains the config read from the file
var Runtime Config
