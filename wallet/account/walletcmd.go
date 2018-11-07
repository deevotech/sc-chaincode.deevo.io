package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/deevotech/sc-chaincode.deevo.io/wallet/lib"
	"github.com/deevotech/sc-chaincode.deevo.io/wallet/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version = "version"
)

type WalletCmd struct {
	// name of the wallet command (init, start, version)
	name string
	// rootCmd is the cobra command
	rootCmd *cobra.Command
	// My viper instance
	myViper *viper.Viper
	// homeDirectory is the location of the wallet's home directory
	homeDirectory string
	// length
	length int
	// type of keypairs
	typeKey string
	// options
	options string
	// to address
	to_address string
	// value of transfer
	value int
}

// NewCommand returns new WalletCmd ready for running
func NewCommand(name string) *WalletCmd {
	s := &WalletCmd{
		name:    name,
		myViper: viper.New(),
	}
	s.init()
	return s
}

// Execute runs this WalletCmd
func (s *WalletCmd) Execute() error {
	return s.rootCmd.Execute()
}

// init initializes the WalletCmd instance
// It intializes the cobra root and sub commands and
// registers command flgs with viper
func (s *WalletCmd) init() {
	// root command
	rootCmd := &cobra.Command{
		Use:   cmdName,
		Short: longName,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := s.configInit()
			if err != nil {
				return err
			}
			cmd.SilenceUsage = true
			util.CmdRunBegin(s.myViper)
			return nil
		},
	}
	s.rootCmd = rootCmd

	// create represents the server create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: fmt.Sprintf("Initialize the %s", shortName),
		Long:  "Generate the key pairs",
	}
	createCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.Errorf(extraArgsError, args, createCmd.UsageString())
		}
		err := s.getWallet().Create()
		if err != nil {
			util.Fatal("Creation failure: %s", err)
		}
		log.Info("Creation was successful")
		return nil
	}
	s.rootCmd.AddCommand(createCmd)

	// transferCmd represents the server transfer command
	transferCmd := &cobra.Command{
		Use:   "transfer",
		Short: fmt.Sprintf("Transfer the %s", shortName),
	}

	transferCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.Errorf(extraArgsError, args, transferCmd.UsageString())
		}
		err := s.getWallet().LoadPrivateKey()
		if err != nil {
			return err
		}

		signature, ciphertext, label, hash, newhash, hashed, opts := s.getWallet().Transfer()
		fmt.Println("signature: ", signature)
		fmt.Println("ciphertext: ", ciphertext)
		fmt.Println("label: ", label)
		fmt.Println("hash: ", hash)
		fmt.Println("newhash: ", newhash)
		fmt.Println("hashed: ", hashed)
		fmt.Println("opts: ", opts)
		return nil
	}
	s.rootCmd.AddCommand(transferCmd)
	s.registerFlags()
}

// registerFlags registers command flags with viper
func (s *WalletCmd) registerFlags() {
	// Get the default config file path
	cfg := util.GetDefaultConfigFile(cmdName)

	// All env variables must be prefixed
	s.myViper.SetEnvPrefix(envVarPrefix)
	s.myViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set specific global flags used by all commands
	pflags := s.rootCmd.PersistentFlags()
	// Don't want to use the default parameter for StringVarP. Need to be able to identify if home directory was explicitly set
	pflags.StringVarP(&s.homeDirectory, "home", "H", "", fmt.Sprintf("Directory to store wallet (default \"%s\")", filepath.Dir(cfg)))
	pflags.IntVar(&s.length, "l", 2048, "Length of key")
	pflags.StringVar(&s.typeKey, "t", "rsa", "type of key")
	pflags.StringVar(&s.options, "p", "options", "type of key")
	pflags.StringVar(&s.to_address, "to", "yyyy", "type of key")
	pflags.IntVar(&s.value, "value", 10, "type of key")
	//err := util.RegisterFlags(s.myViper, pflags, nil, nil)
	/*if err != nil {
		panic(err)
	}*/
}

// Configuration file is not required for some commands like version
func (s *WalletCmd) configRequired() bool {
	return s.name != version
}

// getServer returns a lib.Server for the init and start commands
func (s *WalletCmd) getWallet() *lib.Wallet {
	return &lib.Wallet{
		HomeDir:    s.homeDirectory,
		TypeKey:    s.typeKey,
		Length:     s.length,
		Ops:        s.options,
		To_address: s.to_address,
		Value:      s.value,
	}
}
