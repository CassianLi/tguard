package icp

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

// FillTaxSheet fill tax sheet
func FillTaxSheet(file *excelize.File, sheetName string, taxData []TaxObject) error {
	log.Println("ICP sheet name: ", sheetName)
	file.SetSheetName("Sheet1", sheetName)

	TaxSheetHeaders := &[]interface{}{"SN", "BIll NO.", "Tax Type", "Itemnr", "Destined Number", "Processing Status",
		"Invoice Number", "Invoice Date", "Xml Id", "Currency Code", "LocalCurrency Value", "Import Duty",
		"Dutch Costs", "Ductch VAT", "Statistical Number", "Weight(KG)", "No. of Pieces", "Country Pre fix",
		"VAT Registration Number", "Partner Name", "Country of Destination", "VAT Number", "EORI Number",
		"Importer SS Code", "Address Code", "Address", "Postcode", "City", "Product No", "Description", "MRN",
		"Company Name", "Mode", "ICP/115"}

	err := file.SetSheetRow(sheetName, "A1", TaxSheetHeaders)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, datum := range taxData {
		sn := i + 1
		idx := sn + 1
		err = file.SetCellInt(sheetName, fmt.Sprintf("A%d", idx), sn)
		err = file.SetCellStr(sheetName, fmt.Sprintf("B%d", idx), datum.BillNo)
		err = file.SetCellStr(sheetName, fmt.Sprintf("C%d", idx), datum.TaxType)
		err = file.SetCellStr(sheetName, fmt.Sprintf("D%d", idx), datum.ItemNumber)
		err = file.SetCellStr(sheetName, fmt.Sprintf("E%d", idx), datum.Destined)
		err = file.SetCellStr(sheetName, fmt.Sprintf("F%d", idx), datum.ProcessCode)
		err = file.SetCellStr(sheetName, fmt.Sprintf("G%d", idx), datum.CustomsId)
		err = file.SetCellStr(sheetName, fmt.Sprintf("H%d", idx), datum.InvoiceDate)
		err = file.SetCellStr(sheetName, fmt.Sprintf("I%d", idx), "")
		err = file.SetCellStr(sheetName, fmt.Sprintf("J%d", idx), datum.Currency)
		err = file.SetCellFloat(sheetName, fmt.Sprintf("K%d", idx), datum.LocalCurrencyValue, FloatDecimalPlaces, 64)
		err = file.SetCellFloat(sheetName, fmt.Sprintf("L%d", idx), datum.ImportDuty, FloatDecimalPlaces, 64)
		err = file.SetCellStr(sheetName, fmt.Sprintf("M%d", idx), datum.DutchCost)
		err = file.SetCellStr(sheetName, fmt.Sprintf("N%d", idx), datum.DutchVat)
		err = file.SetCellStr(sheetName, fmt.Sprintf("O%d", idx), datum.HsCode.String)
		err = file.SetCellFloat(sheetName, fmt.Sprintf("P%d", idx), datum.NetWeight, FloatDecimalPlaces, 64)
		err = file.SetCellInt(sheetName, fmt.Sprintf("Q%d", idx), datum.Quantity)
		err = file.SetCellStr(sheetName, fmt.Sprintf("R%d", idx), datum.CountryPreFix)
		err = file.SetCellStr(sheetName, fmt.Sprintf("S%d", idx), datum.DutyParty.String)
		err = file.SetCellStr(sheetName, fmt.Sprintf("T%d", idx), datum.PartnerName)
		err = file.SetCellStr(sheetName, fmt.Sprintf("U%d", idx), datum.CountryOfDestination)
		err = file.SetCellStr(sheetName, fmt.Sprintf("V%d", idx), datum.VatNo)
		err = file.SetCellStr(sheetName, fmt.Sprintf("W%d", idx), datum.EoriNo.String)
		err = file.SetCellStr(sheetName, fmt.Sprintf("X%d", idx), datum.ImportAddressCode)
		err = file.SetCellStr(sheetName, fmt.Sprintf("Y%d", idx), datum.AddressCode)
		err = file.SetCellStr(sheetName, fmt.Sprintf("Z%d", idx), datum.AddressDetail.String)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AA%d", idx), datum.PostalCode.String)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AB%d", idx), datum.City)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AC%d", idx), datum.ProductNo)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AD%d", idx), datum.Description.String)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AE%d", idx), datum.Mrn)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AF%d", idx), datum.CompanyName)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AG%d", idx), datum.Mode)
		err = file.SetCellStr(sheetName, fmt.Sprintf("AH%d", idx), datum.InICPFile)

		if err != nil {
			return err
		}
	}

	return nil
}

// FillTaxFileSheet fill tax file sheet
func FillTaxFileSheet(file *excelize.File, sheetName string, taxFileData []TaxFileObject) error {
	log.Println("Tax file sheet name: ", sheetName)
	file.NewSheet(sheetName)

	TaxSheetHeaders := &[]interface{}{"SN", "MRN", "Tax receipt Link"}

	err := file.SetSheetRow(sheetName, "A1", TaxSheetHeaders)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, datum := range taxFileData {
		sn := i + 1
		idx := sn + 1
		err = file.SetCellInt(sheetName, fmt.Sprintf("A%d", idx), sn)
		err = file.SetCellStr(sheetName, fmt.Sprintf("B%d", idx), datum.Mrn)
		err = file.SetCellStr(sheetName, fmt.Sprintf("C%d", idx), datum.TaxFileLink)

		if err != nil {
			return err
		}
	}

	return nil
}

// FillPodSheet fill pod file sheet
func FillPodSheet(file *excelize.File, sheetName string, podFileData []PodFileObject) error {
	log.Println("POD sheet name: ", sheetName)
	file.NewSheet(sheetName)

	TaxSheetHeaders := &[]interface{}{"SN", "Bill No.", "Customs ID", "MRN No.", "Tracing No.", "POD Link", "Invoice"}

	err := file.SetSheetRow(sheetName, "A1", TaxSheetHeaders)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for i, datum := range podFileData {
		sn := i + 1
		idx := sn + 1
		err = file.SetCellInt(sheetName, fmt.Sprintf("A%d", idx), sn)
		err = file.SetCellStr(sheetName, fmt.Sprintf("B%d", idx), datum.BillNo)
		err = file.SetCellStr(sheetName, fmt.Sprintf("C%d", idx), datum.CustomsId)
		err = file.SetCellStr(sheetName, fmt.Sprintf("D%d", idx), datum.Mrn)
		err = file.SetCellStr(sheetName, fmt.Sprintf("E%d", idx), datum.TrackingNo)
		err = file.SetCellStr(sheetName, fmt.Sprintf("F%d", idx), datum.PodFileLink.String)

		if err != nil {
			return err
		}
	}

	return nil
}
