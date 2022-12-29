/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"sysafari.com/customs/tguard/audit"
	"sysafari.com/customs/tguard/global"
	"time"

	"github.com/spf13/cobra"
)

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "生成指定月份的报关自检文件",
	Long:  `指定月份生成系统已报关的报关自检文件. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("audit called")

		// Init database connection
		global.InitGlobalDatabaseConnection()

		makeAudit()
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// auditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// auditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	auditCmd.Flags().IntVar(&offset, "offset", 0, "指定日期往前偏移的月份数，默认为0（表示不偏移月份，生成指定日期的ICP）")
	auditCmd.Flags().StringVar(&month, "month", time.Now().Format(MonthFormatLayout), "指定月份，默认(2006-01)")
}

func makeAudit() {
	monthT, err := time.Parse(MonthFormatLayout, month)
	if err != nil {
		log.Panic("Date format error", err)
	}
	monthly := monthT.AddDate(0, -offset, 0)
	monthlyStr := monthly.Format(MonthFormatLayout)

	start := time.Now().UnixMilli()
	// make ICP for month
	customsAudit := audit.CustomsAudit{
		Month: monthlyStr,
	}
	customsAudit.MakeAudit()

	end := time.Now().UnixMilli()

	log.Printf("**** Generat audit time costs: %d ms****\n", end-start)
}
