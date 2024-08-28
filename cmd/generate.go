/*
Copyright © 2024 Richard Cox <code@swrm.io>

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
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/swrm-io/argo-bw-secrets/pkg/replacer"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:          "generate",
	Short:        "Replace placeholder strings with Bitwarden Secret",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {

		info, err := os.Stdin.Stat()
		if err != nil {
			log.Fatal(err)
		}

		if info.Mode()&os.ModeCharDevice != 0 {
			fmt.Println("The command is intended to work with pipes.")
			fmt.Println("cat file.json | argo-bw-secrets")
			os.Exit(1)
		}

		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		replacer, err := replacer.New(
			viper.GetString("API_URL"),
			viper.GetString("IDENTITY_URL"),
			viper.GetString("TOKEN"),
		)

		if err != nil {
			return err
		}

		result, err := replacer.Replace(string(stdin))
		if err != nil {
			return err
		}

		fmt.Println(result)
		return nil

	},
}

func init() {
	viper.SetDefault("API_URL", "https://api.bitwarden.com")
	viper.SetDefault("IDENTITY_URL", "https://identity.bitwarden.com")

	rootCmd.AddCommand(generateCmd)
}
