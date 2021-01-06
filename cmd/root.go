/*
Copyright © 2020 Xin Chen <devops.chen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	erh "github.com/jiuchen1986/cks/pkg/error"
	lgr "github.com/jiuchen1986/cks/pkg/logger"
)

var (
	// prefer this naming pattern for variables binding to flags
	// use cmd name + "Flag" + variable name
	// this makes more readable when those varabiles are used across multiple files
	rootCmdFlagCfgFile           string
	rootCmdFlagLogLevel          string
	rootCmdFlagErrHandleWithExit string
	// undo is usually called when the whole program finishes
	// this is used with zap logger
	undo func()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cks",
	Short: "A test Kubernetes distribution",
	Long: `Example usage:

xxxxxxxx`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		erh.ExitOnErr(err, undo)
	}

	// as error.ExitOnErr calls os.Exit(1)
	// that doesn't respect defer, defer is able to called here
	// even if function is already called in error.ExitOnErr
	defer func() {
		if undo != nil {
			undo()
		}
	}()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	var desc string

	rootCmd.PersistentFlags().StringVar(&rootCmdFlagCfgFile, "config", "/var/lib/eke.yaml", "config file")
	desc = fmt.Sprintf("log level (support %s)", lgr.PrintAvailLogLevel())
	rootCmd.PersistentFlags().StringVar(&rootCmdFlagLogLevel, "log-level", "info", desc)
	desc = fmt.Sprintf("how error information is given when handling error by exiting (support %s)",
		erh.PrintAvailExitOnErr())
	rootCmd.PersistentFlags().StringVar(&rootCmdFlagErrHandleWithExit, "err-handling", "simple", desc)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootCmdFlagCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(rootCmdFlagCfgFile)
	} else {
		// Search config in home directory with name ".eke" (without extension).
		viper.AddConfigPath("/var/lib")
		viper.SetConfigName("eke")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// always first setup log system and then error handling
	initLogger()
	initErrHandling()
}

func initErrHandling() {
	if er := erh.UpdateErrHandling(rootCmdFlagErrHandleWithExit); er != nil {
		erh.ExitOnErr(er)
	}
}

func initLogger() {
	opts := []lgr.LogOption{}

	// configure log level
	opt, er := lgr.NewLogLevelOption(rootCmdFlagLogLevel)
	if er != nil {
		erh.ExitOnErr(er)
	}

	opts = append(opts, opt)

	// initialize global logger
	var err error
	undo, _, err = lgr.InitLogger(opts...)
	if err != nil {
		erh.ExitOnErr(err, undo)
	}
}
