// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/agilebits/sm/secrets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt content using key management system",
	Long: `This command will decrypt content that was encrypted using encrypt command.

It requires access to the same key management system (KMS) that was used for encryption.

For example:

	cat encrypted-app-config.sm | sm decrypt > app-config.yml

`,
	Run: func(cmd *cobra.Command, args []string) {
		var message []byte
		var err error
		input := viper.GetString("input")
		fmt.Println(input)
		if input == "" {
			reader := bufio.NewReader(os.Stdin)
			message, err = ioutil.ReadAll(reader)
			if err != nil {
				log.Fatal("failed to read:", err)
			}
		} else {
			message, err = ioutil.ReadFile(input)
			if err != nil {
				log.Fatal("failed to read:", err)
			}
		}

		envelope := &secrets.Envelope{}
		if err = json.Unmarshal(message, &envelope); err != nil {
			log.Fatal("failed to Unmarshal:", err)
		}

		result, err := secrets.DecryptEnvelope(envelope)
		if err != nil {
			log.Fatal("failed to DecryptEnvelope:", err)
		}

		if out != "" {
			f, err := os.Create(out)
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to open %s for writing", out))
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			_, err = w.WriteString(string(result))
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to write output to %s", out))
			}
			w.Flush()
			fmt.Println(fmt.Sprintf("output written to %s", out))
		} else {
			fmt.Println(string(result))
		}
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringP("input", "i", "", "A file to get input from")
	decryptCmd.Flags().StringVarP(&out, "out", "o", "", "A file to write the output to")
	viper.BindPFlag("input", decryptCmd.Flags().Lookup("input"))
}
