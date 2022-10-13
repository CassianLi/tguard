/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"sysafari.com/customs/tguard/global"
	icp2 "sysafari.com/customs/tguard/icp"
	"time"

	"github.com/spf13/cobra"
)

const (
	MonthFormatLayout = "2006-01"
)

var month string
var offset int

// monthlyCmd represents the monthly command
var monthlyCmd = &cobra.Command{
	Use:   "monthly",
	Short: "生成一个月的ICP文件，将为该月每个有报关单的税代生成一个ICP文件",
	Long: `默认生成命令执行时当前月份的ICP文件，当前月有多个税代则生成多个税代的ICP文件。
命令可指定生成某个月份的ICP，也可指定生成命令执行时前几个月的ICP文件。
For example:

1. 生成当月ICP文件：		tguard monthly 
2. 生成上个月ICP文件：		tguard monthly -f 0
3. 生成2022-01的ICP文件： 	tguard monthly -m 2022-01
...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("monthly called")
		// Init database connection
		global.InitGlobalDatabaseConnection()

		makeICPForOneMonth()
	},
}

func init() {
	rootCmd.AddCommand(monthlyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monthlyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	monthlyCmd.Flags().IntVar(&offset, "offset", 0, "指定日期往前偏移的月份数，默认为0（表示不偏移月份，生成指定日期的ICP）")
	monthlyCmd.Flags().StringVar(&month, "month", time.Now().Format(MonthFormatLayout), "指定生成某月的ICP文件，默认为命令执行时当前月份ICP(2006-01)")
}

// makeICPForOneMonth To generate an ICP file for a month,
// one ICP file will be generated for each tax agent with a customs declaration for that month
func makeICPForOneMonth() {
	monthT, err := time.Parse(MonthFormatLayout, month)
	if err != nil {
		log.Panic("Date format error", err)
	}
	monthly := monthT.AddDate(0, -offset, 0)
	monthlyStr := monthly.Format(MonthFormatLayout)

	start := time.Now().UnixMilli()
	// make ICP for month
	icp2.MakeICPForOneMonth(monthlyStr)

	end := time.Now().UnixMilli()

	log.Printf("**** Generat ICP time costs: %d ms****\n", end-start)

}
