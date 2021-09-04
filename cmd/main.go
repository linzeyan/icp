package main

import (
	"fmt"
	"os"

	"github.com/linzeyan/icp"
	"github.com/spf13/cobra"
)

var (
	command = icpFlags()
)

func icpFlags() *cobra.Command {
	var (
		rootCmd = &cobra.Command{
			Use:   "icp",
			Short: `Check ICP status using WEST api`,
			Long:  `Check ICP status using WEST api`,
			Run: func(cmd *cobra.Command, args []string) {
				if len(os.Args) == 1 && len(args) == 0 {
					cmd.Help()
					os.Exit(0)
				}
				fmt.Println(icp.Domain+":", icp.Check())
			},
		}
	)
	cobra.OnInitialize(icp.ReadConf)
	rootCmd.PersistentFlags().StringVarP(&icp.ConfigFile, "config", "c", "", "Config File (default: $HOME/.env)")
	rootCmd.PersistentFlags().StringVarP(&icp.Domain, "domain", "d", "", "Domain Name (required)")
	rootCmd.MarkPersistentFlagRequired("domain")
	rootCmd.PersistentFlags().SortFlags = true
	// rootCmd.AddCommand(initCmd)
	return rootCmd
}

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
