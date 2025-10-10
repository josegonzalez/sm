package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var manifestFile string
var v bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "sm",
	Short:   "Simple secret management tool",
	Version: Version,
	Long: `Simple secret management tool used to protect secrets in the server apps. 
	
It relies on the key management system (KMS) provided by the server environment.
For example, Amazon Web Services KMS is used for servers running on EC2 virtual 
machines. A fake KMS implementation can be used when running on developer machines 
to avoid storing unencrypted secrets in version control systems.

sm can be used in a command-line to encrypt and decrypt configuration files.

For example:

	cat plaintext.config.yaml | sm encrypt > encrypted.config.yaml
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if v {
			fmt.Println(Version)
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var env string
var region string
var masterKeyID string
var out string

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sm/config.yaml)")
	RootCmd.PersistentFlags().StringVar(&manifestFile, "manifest", "./.sm/manifest", "path to the manifest file")
	RootCmd.Flags().BoolVarP(&v, "version", "v", false, "display version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")    // name of config file (without extension)
		viper.AddConfigPath("$HOME/.sm") // adding home directory as first search path
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}
}
