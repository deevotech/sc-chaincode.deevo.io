package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/deevotech/sc-chaincode.deevo.io/deevo-wallet/lib"
	"github.com/deevotech/sc-chaincode.deevo.io/deevo-wallet/util"
	"github.com/hyperledger/fabric/bccsp"
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
	toAddress string
	// value of transfer
	value float64
	// file key
	fileKey string
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
		w, err := s.getWallet().KeyGen(&bccsp.ECDSAP256KeyGenOpts{})
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
		wallet, err := s.getWallet().KeyImport(s.homeDirectory + "/" + s.fileKey)
		if err != nil {
			return err
		}
		//address, err := wallet.SaveAddress("newaddress.txt")
		//fmt.Println("newAddress: ", address)
		r, s, data, err := wallet.Transfer()
		fmt.Println("r: ", r)
		fmt.Println("s: ", s)
		fmt.Println("data: ", data)
		fmt.Println("Verify")
		redata, check := wallet.Receive([]byte(data), r, s, data)
		fmt.Println("check: ", check)
		fmt.Println("redata: ", redata)
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
	pflags.IntVar(&s.length, "l", 256, "Length of key")
	pflags.StringVar(&s.typeKey, "t", "rsa", "type of key")
	pflags.StringVar(&s.options, "p", "options", "options")
	pflags.StringVar(&s.toAddress, "to", "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykiBJ7X3fFG9CMsMCXkr4JksWG2oRy7rpWLkGTM48HhHKLPyDNv8jXoh7jjSYy9zLS9sJw1X2vE2P4Pc66hJtoirwxN8j", "address of account")
	pflags.IntVar(&s.value, "value", 10, "Value of transfering")
	pflags.IntVar(&s.fileKey, "fileKey", "private.key", "file key")
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
func (s *WalletCmd) getWallet() *lib.DeevoWallet {
	return &lib.DeevoWallet{
		HomeDir:   s.homeDirectory,
		TypeKey:   s.typeKey,
		KeyLen:    s.length,
		Ops:       s.options,
		ToAddress: s.toAddress,
		Value:     s.value,
	}
}
