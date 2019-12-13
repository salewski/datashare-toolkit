/*
Copyright 2019 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/GoogleCloudPlatform/bq-datashare-toolkit/client/internal/injestion"
	"github.com/spf13/cobra"
)

var (
	fileName string
	rawData  string
)

// inputCmd represents the input command
var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "BQDS client input service",
	Long: `The BQDS client input service will parse the raw data or a file and iterate over
the data if new line. For example:
k
-d '{"a": 123}\n{"b": "xyz"}'
-f "$PWD/input/file.txt"
`,
	Args: func(cmd *cobra.Command, args []string) error {
		fileName, _ := cmd.Flags().GetString("fileName")
		rawData, _ := cmd.Flags().GetString("rawData")
		//fmt.Println(args)
		//if len(args) < 1 && fileName == "" {
		if rawData == "" && fileName == "" {
			log.Fatalf("Input requires either rawData or fileName")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		rawData, _ := cmd.Flags().GetString("rawData")
		// These PersistantFlags are required and validated in parent command
		projectID, _ := cmd.Flags().GetString("projectID")
		topicID, _ := cmd.Flags().GetString("topicID")
		if rawData != "" {
			log.Debugln("rawData:", rawData)
			err := injestion.Run(projectID, topicID, rawData)
			if err != nil {
				log.Fatal(err)
			}
		}
		fileName, _ := cmd.Flags().GetString("fileName")
		if fileName != "" {
			log.Debugln("FileName:", fileName)
			log.Fatalf("Publishing messages from a fileName is still a WIP...")
		}
	},
}

func init() {
	clientCmd.AddCommand(inputCmd)
	inputCmd.Flags().StringVarP(&rawData, "rawData", "d", "", "raw data to parse")
	inputCmd.Flags().StringVarP(&fileName, "fileName", "f", "", "file name to read")
}
