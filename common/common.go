/*
Copyright 2020 Gravitational, Inc.
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

package common

import (
	"os"

	"github.com/gravitational/magnet"
	"github.com/olekukonko/tablewriter"
)

// HelpEnvs lists environment variables that can override build options
func HelpEnvs() error {
	var result [][]string

	for key, value := range magnet.EnvVars {
		if value.Secret {
			result = append(result, []string{key, "<redacted>", value.Default, value.Short})
		} else {
			result = append(result, []string{key, value.Value, value.Default, value.Short})
		}

	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value", "Default", "Short Description"})
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetReflowDuringAutoWrap(false)

	table.AppendBulk(result)
	table.Render()

	return nil
}
