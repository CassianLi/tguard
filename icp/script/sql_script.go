package script

const (
	// QueryDutyPartiesForMonth SQL is used to query duty parties for a month
	QueryDutyPartiesForMonth = `SELECT DISTINCT c.duty_party
FROM log_clearance_process lcp
         INNER JOIN base_customs c ON lcp.customs_id = c.customs_id
WHERE LENGTH(c.duty_party) > 5
    AND DATE_FORMAT(lcp.gmt_create, '%Y-%m') = ?
  AND (lcp.process_code = 'TAX'
    OR lcp.process_code = 'TMP_TAX');`

	// QueryCustomsIdForICPWithinOneMonthSql SQL is used to query the CustomsId of tax receipts within a month
	QueryCustomsIdForICPWithinOneMonthSql = `SELECT distinct lcp.customs_id
FROM log_clearance_process lcp
         INNER JOIN base_customs c ON lcp.customs_id = c.customs_id
WHERE c.declare_version = 0 
  AND  c.duty_party = ?
  AND DATE_FORMAT(lcp.gmt_create, '%Y-%m') = ?
  AND (lcp.process_code = 'TAX'
    OR lcp.process_code = 'TMP_TAX');`

	// QueryCustomsICPBaseSql The SQL used to query base info of customs icp
	QueryCustomsICPBaseSql = `SELECT bc.customs_id,
       bc.declare_country,
       bc.mrn,
       bc.duty_party,
       bb.bill_no,
       bb.mode,
       cta.name AS partnerName
FROM base_customs bc
         INNER JOIN service_bill_customs sbc ON sbc.is_removed = 0 AND bc.customs_id = sbc.customs_id
         INNER JOIN base_bill bb ON sbc.bill_id = bb.bill_id
         LEFT JOIN config_tax_agency cta ON bc.duty_party = cta.vat_number
WHERE bc.customs_id = ?`

	// QueryDutyNeedVatNote Query duty weather need vat note
	QueryDutyNeedVatNote = `select is_need_vat_note from config_tax_agency where vat_number =?`

	// QueryCustomsICPTaxSql The SQL used to query tax info of customs
	QueryCustomsICPTaxSql = `SELECT bct.tax_type,
       bct.itemnr,
       IF(bct.tax_type = 'B00', '3b', '4a')       AS destined,
       bct.declared_amount,
       bct.tax_fee                                AS importDuty,
       IF(bct.tax_type = 'A00', '0.00', 't.b.d.') AS dutchCost,
       '0.00'                                       AS dutchVat,
       IF(bct.tax_type = 'A00', 'NL', '')         AS countryPreFix,
       lcp.process_code,
       DATE_FORMAT(lcp.gmt_create, '%Y/%m/%d')    AS invoiceDate,
       sca.product_no,
       IFNULL(scvp.hs_code, sca.hs_code) AS hs_code,
       sca.net_weight,
       sca.quantity,
       bd.description,
       'EUR'                                      AS currency
FROM log_clearance_process lcp
         INNER JOIN base_customs_tax bct ON bct.customs_id = lcp.customs_id AND
                                            IF(lcp.process_code = 'TAX', bct.processing_status = 4,
                                               bct.processing_status = 115)
         INNER JOIN service_customs_article sca ON bct.customs_id = sca.customs_id AND bct.itemnr = sca.item_number
         INNER JOIN service_customs_value_process scvp ON sca.customs_value_process_id = scvp.id
         INNER JOIN base_description bd ON scvp.description_id = bd.id
WHERE lcp.customs_id = ? AND lcp.process_code = ?
ORDER BY bct.itemnr, bct.tax_type;`

	// QueryCustomsICPTaxSqlNoneEc The SQL used to query tax info of none-ec customs
	QueryCustomsICPTaxSqlNoneEc = `SELECT bct.tax_type,
       bct.itemnr,
      IF(bct.tax_type = 'B00', '3b', '4a')       AS destined,
       bct.declared_amount,
       bct.tax_fee                                AS importDuty,
       IF(bct.tax_type = 'A00', '0.00', 't.b.d.') AS dutchCost,
       '0.00'                                     AS dutchVat,
       IF(bct.tax_type = 'A00', 'NL', '')         AS countryPreFix,
       lcp.process_code,
       DATE_FORMAT(lcp.gmt_create, '%Y/%m/%d')    AS invoiceDate,
       sca.product_no,
       sca.hs_code,
       sca.net_weight,
       sca.quantity,
       bd.description,
       'EUR'                                      AS currency
FROM log_clearance_process lcp
         INNER JOIN base_customs_tax bct ON bct.customs_id = lcp.customs_id AND
                                            IF(lcp.process_code = 'TAX', bct.processing_status = 4,
                                               bct.processing_status = 115)
         INNER JOIN service_customs_article sca ON bct.customs_id = sca.customs_id AND bct.itemnr = sca.item_number
         INNER JOIN base_description bd ON sca.product_no = bd.product_no AND bd.country = sca.country
WHERE lcp.customs_id = ? AND lcp.process_code = ?
ORDER BY bct.itemnr, bct.tax_type;`

	// QueryCustomsICPImporterSql The SQL used to query importer info for customs
	QueryCustomsICPImporterSql = `SELECT bc.vat_no,
       a.eori_no,
       a.address_code AS importerAddressCode
FROM   service_customs_address sca
    INNER JOIN base_customs bc ON sca.customs_id = bc.customs_id
         INNER JOIN base_address a ON sca.address_code = a.address_code
WHERE sca.customs_id = ?
  AND sca.type = 'IMPORTER';`

	// QueryCustomsICPDeliverySql The SQL used to query delivery address info for customs
	QueryCustomsICPDeliverySql = `SELECT a.address_code,
       a.country,
       a.city,
       CONCAT(IFNULL(a.address_line1, ''), IFNULL(a.address_line2, ''), IFNULL(a.address_line3, '')) AS addressDetail,
       a.postal_code
FROM service_customs_address sca
         INNER JOIN base_address a ON sca.address_code = a.address_code
WHERE sca.customs_id = ?
  AND sca.type = 'DELIVERY';`

	// QueryCustomsCompanySql  The SQL used to query company name of customs
	QueryCustomsCompanySql = `SELECT bc.name
FROM base_customs c
         INNER JOIN base_declaration_log bdl ON c.declaration_id = bdl.declaration_id
         INNER JOIN base_company bc ON bdl.company_id = bc.id
WHERE c.customs_id = ? ;`

	// QueryCustomsHasInICPNameSql Query the ICP file name that already contains the Customs
	QueryCustomsHasInICPNameSql = `	SELECT GROUP_CONCAT(distinct sic.icp_name) FROM service_icp_customs sic WHERE sic.customs_id = ? GROUP BY customs_id;`

	// QueryCustomsTrackingPodSql Query the customs' tracking pod
	QueryCustomsTrackingPodSql = `SELECT b.bill_no,c.customs_id,
       c.mrn AS mrn, 
       t.tracking_no, bf.uri
FROM base_reference_tracking t
    	 INNER JOIN base_bill b ON t.bill_id = b.bill_id
    	 INNER JOIN base_customs c ON t.customs_id = c.customs_id
         LEFT JOIN base_track_logistics_info btli ON t.tracking_no = btli.tracking_no AND btli.index_no = 0
         LEFT JOIN base_file bf ON bf.id = btli.file_id
WHERE t.customs_id = ? ;`

	// QueryCustomsServiceKeySql Query the Customs service key
	QueryCustomsServiceKeySql = `SELECT t.customs_id, MIN(index_no) AS min_index_no, br.service_key
FROM base_reference_tracking t
         INNER JOIN base_reference br ON t.reference = br.reference
WHERE t.customs_id =?;`

	// QueryCustomsTrackingPodDeclareOnlySql Query the pod if the customs use DECLARE ONLY
	QueryCustomsTrackingPodDeclareOnlySql = `SELECT b.bill_no,
       c.customs_id,
       c.mrn AS mrn,
       t.tracking_no,
       MIN(t.index_no) as min_index_no,
       bf.uri
FROM base_reference_tracking t
         INNER JOIN base_bill b ON t.bill_id = b.bill_id
         INNER JOIN base_customs c ON t.customs_id = c.customs_id
         LEFT JOIN base_track_logistics_info btli
                   ON t.tracking_no = btli.tracking_no AND btli.index_no = 0
         LEFT JOIN base_file bf ON bf.id = btli.file_id
WHERE t.customs_id = ? ;`

	// QueryCustomsHasInspectionFineSql 查询报关单是否有过查验罚款
	QueryCustomsHasInspectionFineSql = `SELECT COUNT(1) FROM log_clearance_process WHERE customs_id = ? and process_code='INSPECTION_FINE';`

	// QueryIcpHasExistTotalSql 查询ICP是否已经存在
	QueryIcpHasExistTotalSql = `SELECT COUNT(*) FROM service_icp WHERE duty_part = ? AND year = ? AND month = ?;`

	// UpdateIcpIsNewestSql 更新ICP为非最新。便于新生成的ICP文件成为最新
	UpdateIcpIsNewestSql = `UPDATE service_icp SET is_newest = 0 WHERE duty_part = ? AND year = ? AND month = ?;`

	// InsertServiceICP Insert row into service_icp
	InsertServiceICP = `INSERT INTO service_icp (duty_part, name, year, month, icp_date,total, status, vat_note, is_newest) 
values (:duty_part, :name, :year, :month, :icp_date,:total,:status,:vat_note,:is_newest);`

	// InsertServiceICPCustoms Insert row into service_icp_customs
	InsertServiceICPCustoms = `INSERT INTO service_icp_customs (icp_name, xml_id, customs_id, tax_type,  in_excel) 
values (:icp_name, '', :customs_id, :tax_type, :in_excel);`
)
