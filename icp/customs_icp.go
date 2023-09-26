package icp

import (
	"fmt"
	"strings"
	"sysafari.com/customs/tguard/global"
)

const (
	ProcessCodeTax        = "TAX"
	ProcessCodeTemTax     = "TMP_TAX"
	TaxFileUrlPrefixForNL = "https://board.sysafari.com/declarefile/-1/18-"
	TaxFileUrlPrefixForBE = "https://board.sysafari.com/declarefile-be?customsId=CUSTOMS_ID&statusCode=09"
)

type CustomsICP struct {
	CustomsId string `json:"customs_id"`
	// TAX , TAX_TMP
	ProcessCode    string          `json:"process_code"`
	DeclareCountry string          `json:"declare_country"`
	Mrn            string          `json:"mrn"`
	TaxData        []TaxObject     `json:"tax_data"`
	TaxFileData    []TaxFileObject `json:"tax_file_data"`
	PodFileData    []PodFileObject `json:"pod_file_data"`
	Errors         []string        `json:"errors"`
}

// QueryFillData Query fill data for the ICP file
func (icp *CustomsICP) QueryFillData() {
	icp.queryTaxData()
	icp.queryTaxFileData()
	icp.queryPodFileData()
}

// queryTaxData Query the fill data of the tax table
func (icp *CustomsICP) queryTaxData() {
	// base info
	var icpBase CustomsICPBase
	err := global.Db.Get(&icpBase, QueryCustomsICPBaseSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query icp base info failed. %v", icp.CustomsId, err))
	}
	icp.Mrn, icp.DeclareCountry = icpBase.Mrn, icpBase.DeclareCountry

	// tax info
	var taxInfo []CustomsICPTax
	err = global.Db.Select(&taxInfo, QueryCustomsICPTaxSql, icp.CustomsId, ProcessCodeTax)
	if err != nil || len(taxInfo) == 0 {
		// Query temporary tax information if official tax information is not available
		err = global.Db.Select(&taxInfo, QueryCustomsICPTaxSql, icp.CustomsId, ProcessCodeTemTax)
	}

	if err != nil || len(taxInfo) == 0 {
		// Query none-ec sql tax information is not available
		err = global.Db.Select(&taxInfo, QueryCustomsICPTaxSqlNoneEc, icp.CustomsId, ProcessCodeTax)
	}

	if err != nil || len(taxInfo) == 0 {
		// Query temporary tax information if official tax information is not available
		err = global.Db.Select(&taxInfo, QueryCustomsICPTaxSqlNoneEc, icp.CustomsId, ProcessCodeTemTax)
	}

	if err != nil || len(taxInfo) == 0 {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query tax info failed.%v", icp.CustomsId, err))
	} else {
		icp.ProcessCode = taxInfo[0].ProcessCode
	}

	// importer info
	var importerInfo CustomsICPImporter
	err = global.Db.Get(&importerInfo, QueryCustomsICPImporterSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query importer info failed.%v", icp.CustomsId, err))
	}

	// delivery info
	var deliveryInfo CustomsICPDelivery
	err = global.Db.Get(&deliveryInfo, QueryCustomsICPDeliverySql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query delivery address info failed.%v", icp.CustomsId, err))
	}

	// Company info
	var companyName string
	err = global.Db.Get(&companyName, QueryCustomsCompanySql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query company name failed.%v", icp.CustomsId, err))
	}

	// Query customs has inspection fine
	var inspectionFineCount int64
	err = global.Db.Get(&inspectionFineCount, QueryCustomsHasInspectionFineSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query inspection fine failed.%v", icp.CustomsId, err))
	}
	hasInspectionFine := ""
	if inspectionFineCount > 0 {
		hasInspectionFine = "Yes"
	}

	// into which icp file
	var icpFileNames string
	err = global.Db.Get(&icpFileNames, QueryCustomsHasInICPNameSql, icp.CustomsId)
	if err != nil {
		icpFileNames = ""
	}
	if len(icp.Errors) == 0 {
		taxData := combineToTaxData(icpBase, taxInfo, importerInfo, deliveryInfo, companyName, icpFileNames, hasInspectionFine)
		icp.TaxData = taxData
	}
}

// combineToTaxData The one that merges multiple objects is TaxObject
func combineToTaxData(baseInfo CustomsICPBase, taxInfo []CustomsICPTax, importerInfo CustomsICPImporter, deliveryInfo CustomsICPDelivery, companyName, hasInspectionFine, hasInIcpName string) []TaxObject {
	var taxData []TaxObject
	for _, tax := range taxInfo {
		t := TaxObject{
			BillNo:               baseInfo.BillNo,
			CustomsId:            baseInfo.CustomsId,
			Mrn:                  baseInfo.Mrn,
			DutyParty:            baseInfo.DutyParty,
			PartnerName:          baseInfo.PartnerName.String,
			Mode:                 baseInfo.Mode,
			TaxType:              tax.TaxType,
			ItemNumber:           tax.ItemNumber,
			Destined:             tax.Destined,
			LocalCurrencyValue:   tax.LocalCurrencyValue,
			ImportDuty:           tax.ImportDuty,
			DutchCost:            tax.DutchCost,
			DutchVat:             tax.DutchVat,
			CountryPreFix:        tax.CountryPreFix,
			ProcessCode:          tax.ProcessCode,
			InvoiceDate:          tax.InvoiceDate,
			ProductNo:            tax.ProductNo,
			HsCode:               tax.HsCode,
			NetWeight:            tax.NetWeight,
			Quantity:             tax.Quantity,
			Description:          tax.Description,
			Currency:             tax.Currency,
			VatNo:                importerInfo.VatNo.String,
			EoriNo:               importerInfo.EoriNo,
			ImportAddressCode:    importerInfo.ImportAddressCode,
			AddressCode:          deliveryInfo.AddressCode,
			CountryOfDestination: deliveryInfo.Country,
			AddressDetail:        deliveryInfo.AddressDetail,
			PostalCode:           deliveryInfo.PostalCode,
			City:                 deliveryInfo.City,
			CompanyName:          companyName,
			HasInspectionFine:    hasInspectionFine,
			InICPFile:            hasInIcpName,
		}
		taxData = append(taxData, t)
	}

	return taxData
}

// queryTaxFileData Query the fill data of the tax file table
func (icp *CustomsICP) queryTaxFileData() {
	taxType := 114
	if "TAX" == icp.ProcessCode {
		taxType = 4
	}
	if "NL" == icp.DeclareCountry {
		tf := TaxFileObject{
			Mrn:         icp.Mrn,
			CustomsId:   icp.CustomsId,
			TaxType:     taxType,
			TaxFileLink: TaxFileUrlPrefixForNL + icp.CustomsId,
		}
		icp.TaxFileData = append(icp.TaxFileData, tf)
	}

	if "BE" == icp.DeclareCountry {
		tf := TaxFileObject{
			Mrn:         icp.Mrn,
			CustomsId:   icp.CustomsId,
			TaxType:     taxType,
			TaxFileLink: strings.ReplaceAll(TaxFileUrlPrefixForBE, "CUSTOMS_ID", icp.CustomsId),
		}
		icp.TaxFileData = append(icp.TaxFileData, tf)
	}

}

// queryPodFileData Query the fill data of the pod file table
func (icp *CustomsICP) queryPodFileData() {
	var customsServiceKey CustomsServiceKeyObject
	err := global.Db.Get(&customsServiceKey, QueryCustomsServiceKeySql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query service_key  failed.%v", icp.CustomsId, err))
	}

	var podFiles []PodFileObject
	if "DECLARATION ONLY" == customsServiceKey.ServiceKey {
		err = global.Db.Select(&podFiles, QueryCustomsTrackingPodDeclareOnlySql, icp.CustomsId)
	} else {
		err = global.Db.Select(&podFiles, QueryCustomsTrackingPodSql, icp.CustomsId)
	}

	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query tracking pod  failed.%v", icp.CustomsId, err))
	}

	icp.PodFileData = append(icp.PodFileData, podFiles...)
}
