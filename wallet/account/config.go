package main

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/deevotech/sc-chaincode.deevo.io/wallet/util"
)

const (
	longName     = "Deevo Wallet"
	shortName    = "deevo-wallet"
	cmdName      = "deevo-wallet"
	envVarPrefix = "DEEVO_WALLET"
	homeEnvVar   = "DEEVO_WALLET_HOME"
)
const (
	defaultCfgTemplate = ``
)

var (
	extraArgsError = "Unrecognized arguments found: %v\n\n%s"
)

// Initialize config
func (s *WalletCmd) configInit() (err error) {
	if !s.configRequired() {
		return nil
	}

	s.homeDirectory, err = util.ValidateAndReturnAbsConf(s.homeDirectory, cmdName)
	if err != nil {
		return err
	}

	log.Debugf("Home directory: %s", s.homeDirectory)

	// Read the config
	s.myViper.AutomaticEnv() // read in environment variables that match
	if err != nil {
		return err
	}

	return nil
}
func (s *WalletCmd) createDefaultConfigFile() error {
	//cfg := ""
	// Now write the file
	// cfgDir := filepath.Dir(s.cfgFileName)
	// err = os.MkdirAll(cfgDir, 0755)
	//if err != nil {
	//	return err
	//}

	// Now write the file""
	//return ioutil.WriteFile(s.cfgFileName, []byte(cfg), 0644)
	return nil
}
