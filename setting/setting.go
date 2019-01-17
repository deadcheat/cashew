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

	// Organization is who to authenticate user
	Organization string `yaml:"organization"`

	// ExpirationCheckInterval is a value of interval minute, used only in cashew-isolator
	// cashew-isolator is checking expiration of each tickets and delete all of them
	ExpirationCheckInterval int `yaml:"expiration_check_interval"`

	// GrantingDefaultExpire is seconds from last used for deletion ticket
	// This value will be set as second.
	// (T|P)GT will expires 7,200s (2 hour) as default.
	// If you want to use infinite expire policy, set 0
	GrantingDefaultExpire int `yaml:"granting_default_expire"`

	// Also, TGT will be timed-up in 28,800s (8 hour) from creation
	// This parameter also can be infinity when set 0
	GrantingHardTimeout int `yaml:"granting_hard_timeout"`

	// TicketNumberOfEachUsers is a threshold of tickets that users are allowed to create
	// This parameter is enabled only on ST/PT
	// Default value is 20
	TicketNumberOfEachUsers int `yaml:"ticket_number_of_each_users"`

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
	Driver     string            `yaml:"driver"`
	Name       string            `yaml:"name"`
	User       string            `yaml:"user"`
	Pass       string            `yaml:"pass"`
	Host       string            `yaml:"host"`
	Port       int               `yaml:"port"`
	Socket     string            `yaml:"socket"`
	Parameters map[string]string `yaml:"parameters"`
}

// Authenticator struct of authenticator setting
type Authenticator struct {
	Driver        string `yaml:"driver"`
	*AuthDatabase `yaml:"dbauth"`
	// *LDAP         `yaml:"ldapauth"`
}

// AuthDatabase struct of database setting
type AuthDatabase struct {
	*Database         `yaml:"database"`
	Table             string `yaml:"table"`
	UserNameColumn    string `yaml:"user_name_column"`
	PasswordColumn    string `yaml:"password_column"`
	SaltColumn        string `yaml:"salt_column"`
	*CustomAttributes `yaml:"custom_attributes"`
}

// CustomAttributes type alias to CustomAttribute slice
type CustomAttributes []CustomAttribute

// ToMap return map[string]string from []CustomAttribute
func (c *CustomAttributes) ToMap() (attr map[string]string) {
	attr = make(map[string]string)
	if c == nil {
		return
	}
	cas := *c
	for i := range cas {
		attr[cas[i].ColumnName] = cas[i].Name
	}
	return
}

// LDAP struct of LDAP authenticator
// FIXME when someone implement LDAP mode
// type LDAP struct {
// 	Host         string `yaml:"host"`
// 	Port         int    `yaml:"port"`
// 	BindDN       string `yaml:"bind_dn"`
// 	BindPassword string `yaml:"bind_password"`
// 	BaseDN       string `yaml:"base_dn"`
// 	Filter       string `yaml:"filter"`
// }

// Logging struct for logging setting
type Logging struct {
	Driver   string `yaml:"driver"`
	FileName string `yaml:"file"`
	LogLevel string `yaml:"level"`
}

// CustomAttribute struct for custom-attributes
type CustomAttribute struct {
	Name       string `yaml:"name"`
	ColumnName string `yaml:"column_name"`
}

var (
	// DefaultSetting default values for this
	DefaultSetting = App{
		UseSSL:                  false,
		SSLCertFile:             "",
		SSLCertKey:              "",
		Host:                    "127.0.0.1",
		Port:                    3000,
		URIPath:                 "",
		ExpirationCheckInterval: 30,
		GrantingDefaultExpire:   7200,
		GrantingHardTimeout:     28800,
		TicketNumberOfEachUsers: 20,
		Database:                nil,
		Authenticator:           nil,
		Logging: &Logging{
			Driver:   "stdout",
			LogLevel: "debug",
		},
	}
)
