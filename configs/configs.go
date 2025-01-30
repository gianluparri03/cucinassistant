package configs

// Test (env `CA_TEST`) if set to true disables the requirement for other
// config fields to be set. The only one *always* required is Database.
// Default: false.
var Test bool

// Debug (env `CA_DEBUG`) indicates if the server has to be run in
// debug mode: this only implies a more verbose logging and
// an input to shut it down.
// Default: false.
var Debug bool

// BaseURL (env `CA_BASEURL`) is the URL from which the server will be accessible,
// like https://cucinassistant.com or http://localhost:8080.
// Does not need the trailing slash.
var BaseURL string

// Port (env `CA_PORT`) is the local port to which the server will listen.
var Port string

// SessionSecret (env `CA_SESSIONSECRET`) is used to encrypt session cookies.
var SessionSecret string

// Database (env `CA_DATABASE`) is the PostgreSQL's connection string.
var Database string

// EmailEnabled (env `CA_EMAIL_ENABLED`) indicates if the server should send emails
// or write their content in the logs.
var EmailEnabled bool

// EmailSender (env `CA_EMAIL_SENDER`) is the sender address of the emails.
var EmailSender string

// EmailServer (env `CA_EMAIL_SERVER`) is the hostname of the email server.
var EmailServer string

// EmailPort (env `CA_EMAIL_PORT`) is the port of the email server.
var EmailPort string

// EmailLogin (env `CA_EMAIL_LOGIN`) is used to login to the mail server.
var EmailLogin string

// EmailPassword (env `CA_EMAIL_PASSWORD`) is used to login to the mail server.
var EmailPassword string
