package script

// 当前文件用于编写SQL查询语句，以及其他的SQL语句相关的常量
// 主要用于“拆分报关” 查询ICP相关数据的SQL语句
// date: After 2024-08-29

const (

	// QueryCustomsByDutyPartyForMonthAfterSplitSql 查询指定月份的指定dutyParty的所有的customs_id。注意排除拆分报关时的子报关单
	QueryCustomsByDutyPartyForMonthAfterSplitSql = `
SELECT DISTINCT c.customs_id
FROM base_customs c
         INNER JOIN stats_customs_info sci ON c.customs_id = sci.customs_id
         INNER JOIN log_clearance_process lcp ON lcp.customs_id = c.customs_id
WHERE c.declare_version = 0
  AND c.duty_party = ?
  AND sci.is_master = 1
  AND DATE_FORMAT(lcp.gmt_create
          , '%Y-%m') = ?
  AND (lcp.process_code = 'TAX'
    OR lcp.process_code = 'TMP_TAX');`

	// QueryCustomsHasSplitSql 查询指定的customs_id是否是拆分报关
	QueryCustomsHasSplitSql = `
SELECT has_split
FROM stats_customs_info
WHERE  customs_id = ?;`

	// QuerySplitCustomsTaxSql 查询拆分报关的税金信息。
	// 注意拆分报关时需要通过service_customs_supply_article来获取子报关的品类信息，并与最终查验申报的Article对应从而获取税金信息
	QuerySplitCustomsTaxSql = `
SELECT bct.tax_type,
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
       IFNULL(scvp.hs_code, sca.hs_code)          AS hs_code,
       sca.net_weight,
       sca.quantity,
       bd.description,
       'EUR'                                      AS currency
FROM log_clearance_process lcp
         INNER JOIN base_customs_tax bct ON bct.customs_id = lcp.customs_id AND
                                            IF(lcp.process_code = 'TAX', bct.processing_status = 4,
                                               bct.processing_status = 115)
         INNER JOIN service_customs_supply_article scsa ON bct.customs_id = scsa.customs_id AND bct.itemnr = scsa.item
         INNER JOIN service_customs_article sca ON scsa.article_id = sca.id
         INNER JOIN service_customs_value_process scvp ON sca.customs_value_process_id = scvp.id
         INNER JOIN base_description bd ON scvp.description_id = bd.id
WHERE lcp.process_code = 'TAX'
  AND lcp.customs_id = ?
ORDER BY bct.itemnr, bct.tax_type;`
)
