package script

/*
用户生成2022-09-19 10：34 以前的报关单的ICP 文件的SQL 查询语句
*/

const (
	// QueryCustomsIDsByVatSql 通过VatNo作为Importer查询customsId
	QueryCustomsIDsByVatSql = `SELECT DISTINCT customs_id
FROM base_address a
         LEFT JOIN base_customs c ON a.address_code = c.importer
WHERE a.vat_no = ? and c.declare_version=1 and
      c.declare_status='NORMAL' and c.customs_id is not null and c.mrn is not null;`

	// QueryOldICPFillDataSql 老版本一次查询customs_id 所有数据
	QueryOldICPFillDataSql = `SELECT b.bill_no ,                                                
       ta.tax_type                                                  AS taxType,
       ta.itemnr,
       IF(ta.tax_type = 'A00', '4a', '3b')                          AS destinedNumber,
       ta.processing_status                                         AS processingStatus,
       c.customs_id,                                                 
       sct.status_time                                              AS invoiceDate,
       'EUR'                                                        AS currency,
       ta.declared_amount                                           AS localCurrencyValue,
       ta.tax_fee                                                   AS importDuty,
       IF(ta.tax_type = 'A00', '0.00', 't.b.d.')                    AS dutchCost,
       0.00                                                         AS dutchVat,
       bdc.commodity_code                                           AS hsCode,
       cii.net_weight                                  	 			AS netWeight,
       cd.product_qty                                               AS quantity,
       IF(ta.tax_type = 'A00', 'NL', '')                            AS countryPreFix,
       cid.domestic_duty_taxparty                                   AS dutyParty,
       cdd.name                                                     AS partnerName,
       IF(ta.tax_type = 'A00', 'NL', ca.country)                    AS countryOfDestination,
       vatA.vat_no                                                  AS vatNo,
       vatA.eori_no                                                 AS eoriNo,
       vatA.address_code                                            AS importAddressCode,
       ca.address_code                                              AS addressCode,
       CONCAT(IFNULL(ca.address_line1, ''), IFNULL(ca.address_line2, ''),
              IFNULL(ca.address_line3, ''))  AS addressDetail,
       ca.postal_code                                               AS postalCode,
       ca.city,
       cd.product_no                                                AS productNo,
       d.description,
       c.mrn,
       IFNULL(com.name,'')                                                    AS companyName,
       b.mode,
       ''                                       					AS hasInIcp
FROM base_customs c
         INNER JOIN base_declaration_log dl ON dl.declaration_id = c.declaration_id
         INNER JOIN base_company com ON com.id = dl.company_id
         INNER JOIN base_reference_tracking t
                    ON t.customs_id = c.customs_id AND IF(c.type = 0, t.index_no = 1, t.tracking_no = c.customs_id)
         INNER JOIN base_bill b ON b.bill_id = t.bill_id
         INNER JOIN base_reference r ON r.reference = t.reference
         INNER JOIN base_address ca ON ca.address_code = r.consignee_address_code
         INNER JOIN service_customs_tax sct ON c.customs_id = sct.customs_id
         INNER JOIN base_customs_tax ta ON ta.xml_id = sct.xml_id AND ta.processing_status = c.tax_type
         LEFT JOIN base_customs_import_declaration cid ON cid.dossiernr = c.customs_id
         LEFT JOIN base_config_domestic_duty cdd ON cdd.domestic_duty_part = cid.domestic_duty_taxparty
         LEFT JOIN base_address vatA ON cid.importer_code = vatA.address_code
         LEFT JOIN base_customs_description cd ON cd.customs_id = c.customs_id AND cd.index_no = ta.itemnr
         INNER JOIN base_customs_import_item cii ON cii.customs_id = cd.customs_id AND cii.itemnr = cd.index_no
         LEFT JOIN base_description_code bdc ON bdc.status = 0 AND bdc.country = IF(
                c.sales_channel = 'OTHER' OR c.sales_channel = 'B2B', 'NL', ca.country) AND bdc.asin_no = cd.product_no
         LEFT JOIN base_description d ON d.product_no = bdc.asin_no AND d.country = IF(
                c.sales_channel = 'OTHER' OR c.sales_channel = 'B2B', 'NL', ca.country)
WHERE c.customs_id = ?
ORDER BY c.id, cd.index_no, ta.tax_type;`
)
