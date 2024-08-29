package icp

import (
	"fmt"
	"sysafari.com/customs/tguard/global"
	"sysafari.com/customs/tguard/icp/script"
)

/*
用户生成2022-09-19 10：34 以前的报关单的ICP 文件
*/

type CustomsICPOld struct {
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

// QueryICPFillData Query fill data for the ICP file
func (icp *CustomsICPOld) QueryICPFillData() {
	icp.queryTaxData()
	icp.queryTaxFileData()
	icp.queryPodFileData()
}

// queryTaxData Query ICP fill data
func (icp *CustomsICPOld) queryTaxData() {
	var taxData []TaxObject
	err := global.Db.Select(&taxData, script.QueryOldICPFillDataSql, icp.CustomsId)
	if err != nil {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s query ICP fill data failed, error:%v", icp.CustomsId, err))
	} else {
		icp.TaxData = taxData
	}
}

// queryTaxFileData Query the fill data of the tax file table
func (icp *CustomsICPOld) queryTaxFileData() {
	if len(icp.TaxData) > 0 {
		first := icp.TaxData[0]
		tf := TaxFileObject{
			Mrn:         first.Mrn,
			CustomsId:   icp.CustomsId,
			TaxType:     first.ProcessStatus,
			TaxFileLink: TaxFileUrlPrefixForNL + icp.CustomsId,
		}
		icp.TaxFileData = append(icp.TaxFileData, tf)
	} else {
		icp.Errors = append(icp.Errors, fmt.Sprintf("The customs_id:%s 's ICP fill data is empty.", icp.CustomsId))
	}

}

// queryPodFileData Query the fill data of the pod file table
func (icp *CustomsICPOld) queryPodFileData() {
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
