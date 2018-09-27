package foundation

import "fmt"

func init() {
	initializeApp()
}

// PrepareApp func loads config, starts database and create authenticator
func PrepareApp(configFile string) (err error) {

	// load config
	if err = Load(configFile); err != nil {
		return fmt.Errorf("failed to load config file %s \n", configFile)
	}

	// start database
	if err = StartDatabase(); err != nil {
		return fmt.Errorf("failed to start database %+v \n", err)
	}

	// prepare authenticator
	if err = PrepareAuthenticator(); err != nil {
		return fmt.Errorf("failed to prepare authenticator %+v \n", err)
	}

	return nil
}
