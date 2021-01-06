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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	lgr "github.com/jiuchen1986/cks/pkg/logger"
	etest "github.com/jiuchen1986/cks/test"
)

// dumyCmd represents the dumy command
var dumyCmd = &cobra.Command{
	Use:          "dumy",
	Short:        "A dummy subcommand for testing purpose.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// logger := lgr.GetGlobalStructuredLogger()
		logger := lgr.GetGlobalLogger()
		defer logger.Sync()
		logger.Debug("dumy called")

		pflag.Visit(func(f *pflag.Flag) {
			field := map[string]string{"name": f.Name, "value": f.Value.String()}
			logger.Info("Flags accessible in dumy", field)
		})

		// e := errors.New("error occured in dumy")
		// e := etest.ReturnError()
		// e := etest.ReturnWrappedError()
		e := etest.ReturnNestedError()
		// e := etest.ReturnFmtError()
		logger.UnwrappedStackErrorf(e, "get error in dumy %v", e)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
