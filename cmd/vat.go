/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/icp"
)

var vatNo string

// vatCmd represents the vat command
var vatCmd = &cobra.Command{
	Use:   "vat",
	Short: "生成VAT No. 的ICP 文件",
	Long: `根据海关要求，生成系统建成以来，指定VAT No.报关的所有ICP文件. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vat called")
		// Init database connection
		global.InitGlobalDatabaseConnection()

		icp.MakeICPByVatNo(vatNo)
	},
}

func init() {
	rootCmd.AddCommand(vatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vatCmd.PersistentFlags().String("foo", "", "A help for foo")

	vatCmd.Flags().StringVar(&vatNo, "vat", "", "VAT number")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
