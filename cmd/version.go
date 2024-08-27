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
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	commit    string
	buildDate string
	version   = "v0.0.0-unknown"
)

var versionTemplate = `argo-bw-secrets
  Version:		{{ .Version }}
  Go Version:		{{ .GoVersion }}
  Git Commit:		{{ .GitCommit }}
  Build Date:		{{ .BuildDate }}
  Build Settings:	{{ .BuildSettings }}
`

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print argo-bw-secrets version",
	RunE: func(cmd *cobra.Command, args []string) error {
		build, _ := debug.ReadBuildInfo()

		settings := []string{}
		for _, i := range build.Settings {
			settings = append(settings, fmt.Sprintf("%s %s", i.Key, i.Value))
		}
		data := struct {
			BuildDate     string
			Version       string
			GoVersion     string
			GitCommit     string
			BuildSettings string
		}{
			buildDate,
			version,
			build.GoVersion,
			commit,
			strings.Join(settings, " "),
		}

		tmpl := template.Must(template.New("").Parse(versionTemplate))
		return tmpl.Execute(cmd.OutOrStdout(), data)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
