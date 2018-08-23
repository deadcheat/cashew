package setting

// Loader config load interface
type Loader interface {
	Load(id string) (*App, error)
}

// App struct of application setting
type App struct {
	// Server Setting
	UseSSL      bool   `yaml:"ssl"`
	SSLCertFile string `yaml:"ssl_cert"`
	SSLCertKey  string `yaml:"ssl_key"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`

	// Route Setting
	URIPath string `yaml:"uri_path"`

	// Database Setting
	*Database `yaml:"database"`

	// Authenticate Setting
	*Authenticator `yaml:"authenticator"`

	// Logging Setting
	*Logging `yaml:"logging"`
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
	*AuthDatabase `yaml:"dbauth"`
	*LDAP         `yaml:"ldapauth"`
}

// AuthDatabase struct of database setting
type AuthDatabase struct {
	*Database      `yaml:"database"`
	Table          string `yaml:"table"`
	UserNameColumn string `yaml:"user_name_column"`
	PasswordColumn string `yaml:"password_column"`
	SaltColumn     string `yaml:"salt_column"`
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
	Driver   string `yaml:"driver"`
	FileName string `yaml:"file"`
	LogLevel string `yaml:"level"`
}

var (
	// DefaultSetting default values for this
	DefaultSetting = App{
		UseSSL:        false,
		SSLCertFile:   "",
		SSLCertKey:    "",
		Host:          "127.0.0.1",
		Port:          3000,
		URIPath:       "",
		Database:      nil,
		Authenticator: nil,
		Logging: &Logging{
			Driver:   "stdout",
			LogLevel: "debug",
		},
	}
)
