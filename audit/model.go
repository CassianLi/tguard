package audit

import "database/sql"

type CustomsAuditObject struct {
	BillNo          sql.NullString `db:"bill_no"`
	CustomsId       string         `db:"customs_id"`
	InvoiceDate     sql.NullString `db:"invoice_date"`
	Mrn             sql.NullString `db:"mrn"`
	ItemNumber      string         `db:"item_number"`
	HsCode          sql.NullString `db:"hs_code"`
	EuDutyRate      string         `db:"eu_duty_rate"`
	ProductNo       string         `db:"product_no"`
	WebLink         sql.NullString `db:"web_link"`
	Description     string         `db:"description"`
	PriceScreenshot sql.NullString `db:"price_screenshot"`
}
