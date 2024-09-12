package icp

import (
	"fmt"
	"strings"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/icp/script"
)

const (
	ProcessCodeTax        = "TAX"
	ProcessCodeTemTax     = "TMP_TAX"
	TaxFileUrlPrefixForNL = "https://board.sysafari.com/declarefile/-1/18-"
	TaxFileUrlPrefixForBE = "https://board.sysafari.com/declarefile-be?customsId=CUSTOMS_ID&statusCode=09"
)

// CustomsICP 生成ICP表格文件，主要是ICP Excel文件的制作，不包含后续压缩包和存储路径等的操作
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
	err := global.Db.Get(&icpBase, script.QueryCustomsICPBaseSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query icp base info failed. %v", icp.CustomsId, err))
	}
	icp.Mrn, icp.DeclareCountry = icpBase.Mrn, icpBase.DeclareCountry

	// 查询当前customs 是否是拆分报关 has_split
	var hasSplit bool
	err = global.Db.Get(&hasSplit, script.QueryCustomsHasSplitSql, icp.CustomsId)
	if err != nil {
		fmt.Println("Query customs has split failed, continue to make as no-split", err, icp.CustomsId)
	}

	// 如果是拆分报关，查询拆分报关的税金信息
	queryCustomsTaxSql := script.QueryCustomsICPTaxSql
	if hasSplit {
		fmt.Printf("The customs_id:%s is split customs, query split customs tax info.\n", icp.CustomsId)
		queryCustomsTaxSql = script.QuerySplitCustomsTaxSql
	}

	// tax info, 如果没有正式税则信息（TAX），则查询临时税则信息（TMP_TAX）
	var taxInfo []CustomsICPTax
	err = global.Db.Select(&taxInfo, queryCustomsTaxSql, icp.CustomsId, ProcessCodeTax)
	if err != nil || len(taxInfo) == 0 {
		// 查询临时税金信息
		fmt.Printf("The customs_id:%s query TAX info failed, try to query TMP_TAX info.\n", icp.CustomsId)
		err = global.Db.Select(&taxInfo, queryCustomsTaxSql, icp.CustomsId, ProcessCodeTemTax)
	}

	if err != nil || len(taxInfo) == 0 {
		// Query none-ec sql tax information is not available
		err = global.Db.Select(&taxInfo, script.QueryCustomsICPTaxSqlNoneEc, icp.CustomsId, ProcessCodeTax)
	}

	if err != nil || len(taxInfo) == 0 {
		// Query temporary tax information if official tax information is not available
		err = global.Db.Select(&taxInfo, script.QueryCustomsICPTaxSqlNoneEc, icp.CustomsId, ProcessCodeTemTax)
	}

	if err != nil || len(taxInfo) == 0 {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query tax info failed.%v", icp.CustomsId, err))
	} else {
		icp.ProcessCode = taxInfo[0].ProcessCode
	}

	// importer info
	var importerInfo CustomsICPImporter
	err = global.Db.Get(&importerInfo, script.QueryCustomsICPImporterSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query importer info failed.%v", icp.CustomsId, err))
	}

	// delivery info
	var deliveryInfo CustomsICPDelivery
	err = global.Db.Get(&deliveryInfo, script.QueryCustomsICPDeliverySql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query delivery address info failed.%v", icp.CustomsId, err))
	}

	// Company info
	var companyName string
	err = global.Db.Get(&companyName, script.QueryCustomsCompanySql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query company name failed.%v", icp.CustomsId, err))
	}

	// Query customs has inspection fine
	var inspectionFineCount int64
	err = global.Db.Get(&inspectionFineCount, script.QueryCustomsHasInspectionFineSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query inspection fine failed.%v", icp.CustomsId, err))
	}
	hasInspectionFine := ""
	if inspectionFineCount > 0 {
		hasInspectionFine = "Yes"
	}

	// into which icp file
	var icpFileNames string
	err = global.Db.Get(&icpFileNames, script.QueryCustomsHasInICPNameSql, icp.CustomsId)
	if err != nil {
		icpFileNames = ""
	}
	if len(icp.Errors) == 0 {
		taxData := combineToTaxData(icpBase, taxInfo, importerInfo, deliveryInfo, companyName, hasInspectionFine, icpFileNames)
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
	err := global.Db.Get(&customsServiceKey, script.QueryCustomsServiceKeySql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query service_key  failed.%v", icp.CustomsId, err))
	}

	var podFiles []PodFileObject
	if "DECLARATION ONLY" == customsServiceKey.ServiceKey {
		err = global.Db.Select(&podFiles, script.QueryCustomsTrackingPodDeclareOnlySql, icp.CustomsId)
	} else {
		err = global.Db.Select(&podFiles, script.QueryCustomsTrackingPodSql, icp.CustomsId)
	}

	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query tracking pod  failed.%v", icp.CustomsId, err))
	}

	icp.PodFileData = append(icp.PodFileData, podFiles...)
}
