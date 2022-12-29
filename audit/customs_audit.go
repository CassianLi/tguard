package audit

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/oss"
	"sysafari.com/customs/tguard/utils"
)

type CustomsAudit struct {
	// Month exp: 2006-01
	Month string `json:"month"`

	AuditData []CustomsAuditObject

	Errors []string
}

// queryCustomsAuditData  Query Customs Audit Data
func (ca *CustomsAudit) queryCustomsAuditData() {
	fmt.Println(ca.Month)
	var customsIds []string
	err := global.Db.Select(&customsIds, QueryCustomsSubmittedBetweenDate, ca.Month)
	if err != nil {
		log.Panic(err)
		ca.Errors = append(ca.Errors, fmt.Sprintf("Query customs list within month:%s, error:%v", ca.Month, err))
	}

	log.Infof("The customs total: %d within month: %s", len(customsIds), ca.Month)

	for idx, id := range customsIds {
		log.Infof("%d customs id: %s", idx, id)
		var audit CustomsAuditObject
		err := global.Db.Get(&audit, QueryCustomsAuditData, id)

		if err != nil {
			ca.Errors = append(ca.Errors, fmt.Sprintf("Query customs :%s audit info failed, error:%v", id, err))
		} else {
			ca.AuditData = append(ca.AuditData, audit)
		}
	}
}

func (ca *CustomsAudit) fileAuditExcel(fp string) error {
	sname := "Sheet1"
	file := excelize.NewFile()
	header := &[]interface{}{
		"Bill NO.", "Invoice No.", "Invoice Date", "Itemnr", "Statistical Number",
		"Duty(%)", "Product No.", "Link", "Description", "MRN", "Screenshot",
	}

	err := file.SetSheetRow(sname, "A1", header)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ossClient := oss.Client{
		Endpoint:        viper.GetString("oss.endpoint"),
		AccessKeyId:     viper.GetString("oss.access-key"),
		AccessKeySecret: viper.GetString("oss.access-secret"),
		BucketName:      viper.GetString("oss.bucket"),
	}

	tmpDir := viper.GetString("audit.tmp-dir")

	for i, datum := range ca.AuditData {
		idx := i + 2
		err = file.SetCellStr(sname, fmt.Sprintf("A%d", idx), datum.BillNo.String)
		err = file.SetCellStr(sname, fmt.Sprintf("B%d", idx), datum.CustomsId)
		err = file.SetCellStr(sname, fmt.Sprintf("C%d", idx), datum.InvoiceDate.String)
		err = file.SetCellStr(sname, fmt.Sprintf("D%d", idx), datum.ItemNumber)
		err = file.SetCellStr(sname, fmt.Sprintf("E%d", idx), datum.HsCode.String)
		err = file.SetCellStr(sname, fmt.Sprintf("F%d", idx), datum.EuDutyRate)
		err = file.SetCellStr(sname, fmt.Sprintf("G%d", idx), datum.ProductNo)
		err = file.SetCellStr(sname, fmt.Sprintf("H%d", idx), datum.WebLink.String)
		err = file.SetCellStr(sname, fmt.Sprintf("I%d", idx), datum.Description)
		err = file.SetCellStr(sname, fmt.Sprintf("J%d", idx), datum.Mrn.String)

		screenshotName := datum.PriceScreenshot.String
		if screenshotName != "" && !strings.Contains(screenshotName, "http") {
			tmpPath := filepath.Join(tmpDir, screenshotName)
			err = ossClient.DownloadOssFile(screenshotName, tmpPath)
			if err == nil {
				absPaht, _ := filepath.Abs(tmpPath)
				fmt.Println("abs path:", absPaht)
				err = file.SetRowHeight(sname, idx, 200)
				err = file.AddPicture(sname, fmt.Sprintf("K%d", idx), tmpPath, `{"autofit": true}`)
				if err != nil {
					log.Panic(err)
				}
			}
		}

		if err != nil {
			return err
		}
	}

	if err := file.SaveAs(fp); err != nil {
		return err
	}

	return nil
}

// MakeAudit make audit file
func (ca *CustomsAudit) MakeAudit() {
	ca.queryCustomsAuditData()
	if len(ca.Errors) > 0 {
		log.Errorf("Query audit data failed, err: %v", ca.Errors)
		return
	}

	auditSavePath := viper.GetString("audit.save-dir")
	if !utils.IsExists(auditSavePath) {
		utils.CreateDir(auditSavePath)
	}
	auditFilename := ca.Month + ".xlsx"

	err := ca.fileAuditExcel(filepath.Join(auditSavePath, auditFilename))
	if err != nil {
		log.Error("Generate audit file failed, err: ", err)
	}
}
