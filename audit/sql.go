package audit

const (
	// QueryCustomsSubmittedBetweenDate Query the customs IDs submitted within the month
	QueryCustomsSubmittedBetweenDate = `SELECT DISTINCT customs_id
FROM log_customs_state
WHERE state = 'SUBMITTED'
  AND DATE_FORMAT(gmt_create, '%Y-%m') = ?;`

	QueryCustomsAuditData = `SELECT bb.bill_no,
       sca.customs_id,
       bc.mrn,
       sca.item_number,
       IFNULL(scvp.hs_code, sca.hs_code)                                     AS hs_code,
       CONCAT(FORMAT(IFNULL(scvp.eu_duty_rate, sca.duty_amount / sca.final_declared_value) * 100, 2),
               '%')                                                          AS eu_duty_rate,
       sca.product_no,
       bd.web_link,
       bd.description,
       scvp.price_screenshot
FROM service_customs_article sca
         INNER JOIN base_description bd ON sca.product_no = bd.product_no AND sca.country = bd.country
         INNER JOIN base_customs bc ON sca.customs_id = bc.customs_id
         LEFT JOIN service_bill_customs sbc ON sca.customs_id = sbc.customs_id
         LEFT JOIN base_bill bb ON sbc.bill_id = bb.bill_id
         LEFT JOIN service_customs_value_process scvp ON sca.customs_value_process_id = scvp.id
WHERE sca.customs_id =?`
)
