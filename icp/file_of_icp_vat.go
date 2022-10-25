package icp

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/utils"
	"time"
)

type FileOfICPForVAT struct {
	// CustomsIDs The customs IDs that needs to be placed in the ICP file
	CustomsIDs []string `json:"customs_ids"`
	// The duty party
	VatNo string `json:"vat_no"`
	// TaxData Population data for tax information form
	TaxData []TaxObject `json:"tax_data"`
	// TaxFileData Population data for tax file information form
	TaxFileData []TaxFileObject `json:"tax_file_data"`
	// PodFileData The data used to fill the pod file table
	PodFileData []PodFileObject `json:"pod_file_data"`
	// FilePath The full path of ICP file.
	FilePath string `json:"file_path"`
	// FileName The ICP file name
	FileName string `json:"file_name"`
	// Errors The ICP errors
	Errors []string `json:"errors"`
}

// QueryCustomsIDs Query the customs IDs that has been used as VAT by the IMPORTER for declaration.
func (f *FileOfICPForVAT) QueryCustomsIDs() {
	var ids []string
	log.Printf("vat:%s", f.VatNo)
	err := global.Db.Select(&ids, QueryCustomsIDsByVatSql, f.VatNo)
	if err != nil || len(ids) == 0 {
		fmt.Println("Query customs ids failed, err: ", err)
		f.Errors = append(f.Errors, fmt.Sprintf("Can not query customs vat no: %s", f.VatNo))
	}
	log.Printf("Total cusotms: %d", len(ids))
	f.CustomsIDs = ids
}

// GenerateICP Begin to generate ICP file
func (f *FileOfICPForVAT) GenerateICP() string {
	if len(f.Errors) == 0 {
		f.readyICPFileInfo()
		if len(f.Errors) > 0 {
			log.Printf("Generating ICP ready ICP info error: %v \n", f.Errors)
		}

		f.generateFillData()
		if len(f.Errors) > 0 {
			log.Printf("Generating ICP query fill data error: %v \n", f.Errors)
		}

		f.createICPFile()
		if len(f.Errors) > 0 {
			log.Printf("Generating ICP create ICP excel file error: %v \n", f.Errors)
		}

		f.saveICPInfoIntoDB(true)
		// 不保存 ICP 与Customs 关系
		//f.saveCustomsInfoWithinICP()
		if len(f.Errors) > 0 {
			log.Printf("Save ICP and customs info failed, error: %v \n", f.Errors)
		}

		return f.FileName
	}
	return ""
}

// saveCustomsInfoWithinICP Save relations information for customs and ICP
func (f *FileOfICPForVAT) saveCustomsInfoWithinICP() {
	var customsICPs []ServiceICPCustoms

	for _, i2 := range f.TaxFileData {
		customsId := i2.CustomsId
		ci := ServiceICPCustoms{
			IcpName:   f.FileName,
			CustomsId: customsId,
			TaxType:   i2.TaxType,
			InExcel:   utils.In(customsId, f.CustomsIDs),
		}
		customsICPs = append(customsICPs, ci)
	}

	_, err := global.Db.NamedExec(InsertServiceICPCustoms, customsICPs)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Save ICP(%s)'s customs information failed: %v", f.FileName, err))
	}

}

// generateFillData Generate fill data for ICP file.
func (f *FileOfICPForVAT) generateFillData() {
	log.Printf("**** Begin to generate ICP file ****")
	for i, d := range f.CustomsIDs {
		log.Printf("**** %d cusotms ID: %s ****", i, d)
		icp := &CustomsICPOld{
			CustomsId: d,
		}
		icp.QueryICPFillData()
		if len(icp.Errors) == 0 {
			f.TaxData = append(f.TaxData, icp.TaxData...)
			f.TaxFileData = append(f.TaxFileData, icp.TaxFileData...)
			f.PodFileData = append(f.PodFileData, icp.PodFileData...)
		} else {
			f.Errors = append(f.Errors, icp.Errors...)
		}
	}
}

// readyICPFileInfo Get ready for icp file info
func (f *FileOfICPForVAT) readyICPFileInfo() {
	saveRoot := viper.GetString("icp.save-dir")
	if saveRoot == "" {
		log.Panic("ICP root save directory not set ..")
	}

	dt := time.Now()

	if f.FileName == "" {
		if f.VatNo == "" {
			f.Errors = append(f.Errors, fmt.Sprintf("Vat No. is required to generate ICP file, but is empty."))
			return
		}
		date, t := dt.Format(FileNameDateLayout), dt.Format(FileNameTimeLayout)
		f.FileName = fmt.Sprintf("VAT%s_%s_%s.xlsx", f.VatNo, date, t)
	}

	year, month := utils.GetCurrentYearMonth(dt)
	saveDir := filepath.Join(saveRoot, year, month)

	log.Println("ICP save dir: ", saveDir)
	if !utils.IsDir(saveDir) && !utils.CreateDir(saveDir) {
		f.Errors = append(f.Errors, fmt.Sprintf("Create save dir: %s failed.", saveDir))
		return
	}
	f.FilePath = filepath.Join(saveDir, f.FileName)
}

// createICPFile creates a ICP excel file
func (f *FileOfICPForVAT) createICPFile() {
	log.Println("**** Creating ICP excel ****")
	file := excelize.NewFile()
	icpDate := time.Now().Format(FileNameDateLayout)

	taxSheetName := fmt.Sprintf("%s_%s_%s", "ICP", f.VatNo, icpDate)
	err := FillTaxSheet(file, taxSheetName, f.TaxData)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Fill ICP sheet failed: %v", err))
	}

	taxFileSheetName := fmt.Sprintf("%s_%s_%s", "TAX", f.VatNo, icpDate)
	err = FillTaxFileSheet(file, taxFileSheetName, f.TaxFileData)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Fill ICP sheet failed: %v", err))
	}

	podSheetName := fmt.Sprintf("%s_%s_%s", "POD", f.VatNo, icpDate)
	err = FillPodSheet(file, podSheetName, f.PodFileData)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Fill POD sheet failed: %v", err))
	}

	log.Printf("**** Save ICP excel: %s ****\n", f.FilePath)
	if err := file.SaveAs(f.FilePath); err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Save ICP file on disk failed: %v", err))
	}
}

// saveICPInfoIntoDB Save ICP info to database
func (f *FileOfICPForVAT) saveICPInfoIntoDB(status bool) {
	dt := time.Now()

	serviceIcp := &ServiceICP{
		DutyParty: f.VatNo,
		Name:      f.FileName,
		Year:      dt.Year(),
		Month:     int(dt.Month()),
		IcpDate:   time.Now().UTC().Format("2006-01-02 15:04:05"),
		Total:     len(f.CustomsIDs),
		Status:    status,
	}
	_, err := global.Db.NamedExec(InsertServiceICP, serviceIcp)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Save ICP(%s) information failed: %v", f.FileName, err))
	}
}
