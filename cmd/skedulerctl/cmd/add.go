/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"

	"github.com/filariow/gardenia/pkg/skeduler"
	"github.com/filariow/gardenia/pkg/valvedprotos"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := skeduler.NewClientFromEnv()
		if err != nil {
			return err
		}

		d, err := cmd.Flags().GetInt64("duration")
		if err != nil {
			return err
		}

		cr, err := cmd.Flags().GetString("cron")
		if err != nil {
			return err
		}

		req := valvedprotos.AddSkeduleRequest{
			Skedule: &valvedprotos.Skedule{
				DurationSec: d,
				CronSkedule: cr,
			},
		}

		rep, err := c.AddSkedule(cmd.Context(), &req)
		if err != nil {
			return err
		}

		m, err := json.Marshal(struct{ JobName string }{JobName: rep.GetJobName()})
		if err != nil {
			return err
		}

		fmt.Println(m)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().Int64("duration", 'd', "The duration of the job")
	addCmd.PersistentFlags().String("cron", "c", "The cron spec for the job")
}
