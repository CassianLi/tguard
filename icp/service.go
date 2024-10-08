package icp

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/icp/script"
)

// MakeICPForOneMonth Make ICP for one month
func MakeICPForOneMonth(month string) {
	var dutyParties []string
	err := global.Db.Select(&dutyParties, script.QueryDutyPartiesForMonth, month)
	if err != nil {
		log.Panicf("Query duty parties in the month %s , error:%v \n", month, err)
	}
	if len(dutyParties) == 0 {
		log.Panicf("No duty party used to customs in the month %s \n", month)
	}

	log.Printf("There are %d duty party in this month %s \n", len(dutyParties), month)

	for _, dutyParty := range dutyParties {
		filename, errs := MakeICPForDutyPart(dutyParty, month)
		if len(errs) > 0 {
			log.Printf("Error creating ICP for duty party %s in the month %s, erros: %v\n", dutyParty, month, errs)
		} else {
			log.Printf("Generat ICP for duty party %s in the month %s success ,the filename: %s\n", dutyParty, month, filename)
		}
	}
}

// MakeICPForDutyPart Make ICP file for the duty party
func MakeICPForDutyPart(dutyParty string, month string) (string, []string) {
	log.Printf("Making ICP for duty party %s in the month %s \n", dutyParty, month)
	icp := &FileOfICP{
		DutyParty: dutyParty,
		Month:     month,
	}

	// 1. 查询ICP基础数据。本月内，指定的dutyParty的所有的customs_id
	// 需要排除作为拆分报关的子报关单
	icp.QueryCustomsIDs()

	openVatNote := viper.GetBool("zip.vat-note-open")
	if openVatNote {
		fmt.Println("Need check whether duty need vat-note..")
		if icp.DutyNeedVatNote() {
			icp.GenerateVatNotesZip()
		}
	}

	filename := icp.GenerateICP()
	errs := icp.Errors
	return filename, errs
}

// MakeICPByVatNo Make ICP file by VAt No.
func MakeICPByVatNo(vatNo string) (string, []string) {
	log.Printf("Making ICP by vat no %s  \n", vatNo)
	icp := &FileOfICPForVAT{
		VatNo: vatNo,
	}

	icp.QueryCustomsIDs()
	filename := icp.GenerateICP()
	errs := icp.Errors
	fmt.Println("errors: ", errs)
	return filename, errs
}
