package app

// Setting struct of application setting
type Setting struct {
	// Server Setting
	UseSSL      bool   `yaml:"ssl"`
	SSLCertFile string `yaml:"ssl_cert"`
	SSLCertKey  string `yaml:"ssl_key"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`

	// Route Setting
	URIPath string `yaml:"uri_path"`

	// Database Setting
	Database `yaml:"database"`

	// Authenticate Setting
	Authenticator `yaml:"authenticator"`

	// Logging Setting
	Logging `yaml:"logging"`
}

// Database struct of database setting
type Database struct {
	Driver string `yaml:"driver"`
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

// Authenticator struct of authenticator setting
type Authenticator struct {
	Driver        string `yaml:"driver"`
	*AuthDatabase `yaml:"database"`
	*LDAP         `yaml:"ldap"`
}

// AuthDatabase struct of database setting
type AuthDatabase struct {
	Driver      string `yaml:"driver"`
	Name        string `yaml:"name"`
	User        string `yaml:"user"`
	Pass        string `yaml:"pass"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Table       string `yaml:"table"`
	UserNameKey string `yaml:"user_name_key"`
	PasswordKey string `yaml:"password_key"`
}

// LDAP struct of LDAP authenticator
// FIXME when someone implement LDAP mode
type LDAP struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	BindDN       string `yaml:"bind_dn"`
	BindPassword string `yaml:"bind_password"`
	BaseDN       string `yaml:"base_dn"`
	Filter       string `yaml:"filter"`
}

// Logging struct for logging setting
type Logging struct {
	FileName string `yaml:"file"`
	LogLevel string `yaml:"level"`
}
