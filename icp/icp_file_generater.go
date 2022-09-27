package icp

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"strings"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/utils"
	"time"
)

const (
	FileNameDateLayout = "200601"
	FileNameTimeLayout = "02150405"
	FloatDecimalPlaces = 6
)

type FileOfICP struct {
	// CustomsIDs The customs IDs that needs to be placed in the ICP file
	CustomsIDs []string `json:"customs_ids"`
	// The duty party
	DutyParty string `json:"duty_party"`
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

// QueryCustomsIDs Query customs IDs between the startDate and endDate
func (f *FileOfICP) QueryCustomsIDs(startDate string, endDate string) {
	var customsIds []string
	err := global.Db.Select(&customsIds, QueryCustomsIdForICPWithinOneMonthSql, f.DutyParty, startDate, endDate)
	if err != nil || len(customsIds) == 0 {
		f.Errors = append(f.Errors, fmt.Sprintf("Can not query customs for duty party %s between start date %s and end date %s", f.DutyParty, startDate, endDate))
	}
	log.Printf("Total cusotms: %d", len(customsIds))
	f.CustomsIDs = customsIds
}

// GenerateICP Begin to generate ICP file
func (f *FileOfICP) GenerateICP() string {
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

		return f.FileName
	}
	return ""
}

// generateFillData Generate fill data for ICP file.
func (f *FileOfICP) generateFillData() {
	log.Printf("**** Begin to generate ICP file ****")
	for i, d := range f.CustomsIDs {
		log.Printf("**** %d cusotms ID: %s ****", i, d)
		icp := &CustomsICP{
			CustomsId: d,
		}
		icp.QueryFillData()
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
func (f *FileOfICP) readyICPFileInfo() {
	saveRoot := viper.GetString("icp.save-dir")
	if saveRoot == "" {
		log.Panic("ICP root save directory not set ..")
	}

	var now time.Time
	if f.FileName == "" {
		if f.DutyParty == "" {
			f.Errors = append(f.Errors, fmt.Sprintf("Duty party is required to generate ICP file, but is empty."))
			return
		}
		now = time.Now()
		date, t := now.Format(FileNameDateLayout), now.Format(FileNameTimeLayout)
		f.FileName = fmt.Sprintf("%s_%s_%s.xlsx", f.DutyParty, date, t)
	} else {
		fp := strings.Split(f.FileName, "_")
		f.DutyParty = fp[0]
		dt := fp[1]
		d, err := time.Parse(FileNameDateLayout, dt)
		if err != nil {
			f.Errors = append(f.Errors, fmt.Sprintf("The ICP filename:%s invalid format(correct: BE0796544895_200601_02150405.xlsx)", f.FileName))
			return
		}
		now = d
	}

	year, month := utils.GetCurrentYearMonth(now)
	saveDir := filepath.Join(saveRoot, year, month)

	log.Println("ICP save dir: ", saveDir)
	if !utils.IsDir(saveDir) && !utils.CreateDir(saveDir) {
		f.Errors = append(f.Errors, fmt.Sprintf("Create save dir: %s failed.", saveDir))
		return
	}
	f.FilePath = filepath.Join(saveDir, f.FileName)
}

// createICPFile creates a ICP excel file
func (f *FileOfICP) createICPFile() {
	log.Println("**** Creating ICP excel ****")
	file := excelize.NewFile()
	icpDate := time.Now().Format(FileNameDateLayout)

	err := f.fillICPSheet(file, icpDate)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Fill ICP sheet failed: %v", err))
	}

	err = f.fillTaxFileSheet(file, icpDate)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Fill tax sheet failed: %v", err))
	}

	err = f.fillPodSheet(file, icpDate)
	if err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Fill POD sheet failed: %v", err))
	}

	log.Printf("**** Save ICP excel: %s ****\n", f.FilePath)
	if err := file.SaveAs(f.FilePath); err != nil {
		f.Errors = append(f.Errors, fmt.Sprintf("Save ICP file on disk failed: %v", err))
	}
}

// fillTaxSheet fill tax sheet
func (f *FileOfICP) fillICPSheet(file *excelize.File, date string) error {
	taxSheetName := fmt.Sprintf("%s_%s_%s", "ICP", f.DutyParty, date)
	log.Println("ICP sheet name: ", taxSheetName)
	file.SetSheetName("Sheet1", taxSheetName)

	TaxSheetHeaders := &[]interface{}{"SN", "BIll NO.", "Tax Type", "Itemnr", "Destined Number", "Processing Status",
		"Invoice Number", "Invoice Date", "Xml Id", "Currency Code", "LocalCurrency Value", "Import Duty",
		"Dutch Costs", "Ductch VAT", "Statistical Number", "Weight(KG)", "No. of Pieces", "Country Pre fix",
		"VAT Registration Number", "Partner Name", "Country of Destination", "VAT Number", "EORI Number",
		"Importer SS Code", "Address Code", "Address", "Postcode", "City", "Product No", "Description", "MRN",
		"Company Name", "Mode", "ICP/115"}

	err := file.SetSheetRow(taxSheetName, "A1", TaxSheetHeaders)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, datum := range f.TaxData {
		sn := i + 1
		idx := sn + 1
		err = file.SetCellInt(taxSheetName, fmt.Sprintf("A%d", idx), sn)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("B%d", idx), datum.BillNo)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("C%d", idx), datum.TaxType)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("D%d", idx), datum.ItemNumber)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("E%d", idx), datum.Destined)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("F%d", idx), datum.ProcessCode)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("G%d", idx), datum.CustomsId)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("H%d", idx), datum.InvoiceDate)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("I%d", idx), "")
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("J%d", idx), datum.Currency)
		err = file.SetCellFloat(taxSheetName, fmt.Sprintf("K%d", idx), datum.LocalCurrencyValue, FloatDecimalPlaces, 64)
		err = file.SetCellFloat(taxSheetName, fmt.Sprintf("L%d", idx), datum.ImportDuty, FloatDecimalPlaces, 64)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("M%d", idx), datum.DutchCost)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("N%d", idx), datum.DutchVat)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("O%d", idx), datum.HsCode)
		err = file.SetCellFloat(taxSheetName, fmt.Sprintf("P%d", idx), datum.NetWeight, FloatDecimalPlaces, 64)
		err = file.SetCellInt(taxSheetName, fmt.Sprintf("Q%d", idx), datum.Quantity)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("R%d", idx), datum.CountryPreFix)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("S%d", idx), datum.DutyParty)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("T%d", idx), datum.PartnerName)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("U%d", idx), datum.CountryOfDestination)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("V%d", idx), datum.VatNo)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("W%d", idx), datum.EoriNo)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("X%d", idx), datum.ImportAddressCode)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("Y%d", idx), datum.AddressCode)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("Z%d", idx), datum.AddressDetail)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AA%d", idx), datum.PostalCode)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AB%d", idx), datum.City)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AC%d", idx), datum.ProductNo)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AD%d", idx), datum.Description)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AE%d", idx), datum.Mrn)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AF%d", idx), datum.CompanyName)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AG%d", idx), datum.Mode)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("AH%d", idx), datum.InICPFile)

		if err != nil {
			return err
		}
	}

	return nil
}

// fillTaxFileSheet fill tax file sheet
func (f *FileOfICP) fillTaxFileSheet(file *excelize.File, date string) error {
	taxSheetName := fmt.Sprintf("%s_%s_%s", "TAX", f.DutyParty, date)
	log.Println("Tax file sheet name: ", taxSheetName)
	file.NewSheet(taxSheetName)

	TaxSheetHeaders := &[]interface{}{"SN", "MRN", "Tax receipt Link"}

	err := file.SetSheetRow(taxSheetName, "A1", TaxSheetHeaders)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, datum := range f.TaxFileData {
		sn := i + 1
		idx := sn + 1
		err = file.SetCellInt(taxSheetName, fmt.Sprintf("A%d", idx), sn)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("B%d", idx), datum.Mrn)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("C%d", idx), datum.TaxFileLink)

		if err != nil {
			return err
		}
	}

	return nil
}

// fillPodSheet fill pod file sheet
func (f *FileOfICP) fillPodSheet(file *excelize.File, date string) error {
	taxSheetName := fmt.Sprintf("%s_%s_%s", "POD", f.DutyParty, date)
	log.Println("POD sheet name: ", taxSheetName)
	file.NewSheet(taxSheetName)

	TaxSheetHeaders := &[]interface{}{"SN", "MRN No.", "POD Link"}

	err := file.SetSheetRow(taxSheetName, "A1", TaxSheetHeaders)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, datum := range f.PodFileData {
		sn := i + 1
		idx := sn + 1
		err = file.SetCellInt(taxSheetName, fmt.Sprintf("A%d", idx), sn)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("B%d", idx), datum.Mrn)
		err = file.SetCellStr(taxSheetName, fmt.Sprintf("C%d", idx), datum.PodFileLink.String)

		if err != nil {
			return err
		}
	}

	return nil
}
