package icp

import (
	"database/sql"
)

// CustomsICPBase Base info of customs icp
type CustomsICPBase struct {
	CustomsId      string         `db:"customs_id"`
	DeclareCountry string         `db:"declare_country"`
	Mrn            string         `db:"mrn"`
	DutyParty      sql.NullString `db:"duty_party"`
	PartnerName    sql.NullString `db:"partnerName"`
	BillNo         string         `db:"bill_no"`
	Mode           string         `db:"mode"`
}

// CustomsICPTax Tax info of customs
type CustomsICPTax struct {
	TaxType            string         `db:"tax_type"`
	ItemNumber         string         `db:"itemnr"`
	Destined           string         `db:"destined"`
	LocalCurrencyValue float64        `db:"declared_amount"`
	ImportDuty         float64        `db:"importDuty"`
	DutchCost          string         `db:"dutchCost"`
	DutchVat           string         `db:"dutchVat"`
	CountryPreFix      string         `db:"countryPreFix"`
	ProcessCode        string         `db:"process_code"`
	InvoiceDate        string         `db:"invoiceDate"`
	ProductNo          string         `db:"product_no"`
	HsCode             sql.NullString `db:"hs_code"`
	NetWeight          float64        `db:"net_weight"`
	Quantity           int            `db:"quantity"`
	Description        sql.NullString `db:"description"`
	Currency           string         `db:"currency"`
}

// CustomsICPImporter The importer address info for customs
type CustomsICPImporter struct {
	VatNo             sql.NullString `db:"vat_no"`
	EoriNo            sql.NullString `db:"eori_no"`
	ImportAddressCode string         `db:"importerAddressCode"`
}

// CustomsICPDelivery The delivery address info for customs
type CustomsICPDelivery struct {
	AddressCode   string         `db:"address_code"`
	Country       string         `db:"country"`
	City          string         `db:"city"`
	AddressDetail sql.NullString `db:"addressDetail"`
	PostalCode    sql.NullString `db:"postal_code"`
}

// ServiceICP sysafari.service_icp
type ServiceICP struct {
	DutyParty string `db:"duty_part"`
	Name      string `db:"name"`
	Year      int    `db:"year"`
	Month     int    `db:"month"`
	IcpDate   string `db:"icp_date"`
	Total     int    `db:"total"`
	Status    bool   `db:"status"`
}

// ServiceICPCustoms sysafari.service_icp_customs
type ServiceICPCustoms struct {
	IcpName   string `db:"icp_name"`
	CustomsId string `db:"customs_id"`
	TaxType   int    `db:"tax_type"`
	InExcel   bool   `db:"in_excel"`
}

// TaxObject tax information object,
// db 字段关联的为老版本查询
type TaxObject struct {
	Sn                   int
	BillNo               string `db:"bill_no"`
	TaxType              string `db:"taxType"`
	ItemNumber           string `db:"itemnr"`
	Destined             string `db:"destinedNumber"`
	ProcessCode          string
	ProcessStatus        int            `db:"processingStatus"`
	CustomsId            string         `db:"customs_id"`
	InvoiceDate          string         `db:"invoiceDate"`
	Currency             string         `db:"currency"`
	LocalCurrencyValue   float64        `db:"localCurrencyValue"`
	ImportDuty           float64        `db:"importDuty"`
	DutchCost            string         `db:"dutchCost"`
	DutchVat             string         `db:"dutchVat"`
	HsCode               sql.NullString `db:"hsCode"`
	NetWeight            float64        `db:"netWeight"`
	Quantity             int            `db:"quantity"`
	CountryPreFix        string         `db:"countryPreFix"`
	DutyParty            sql.NullString `db:"dutyParty"`
	PartnerName          string         `db:"partnerName"`
	CountryOfDestination string         `db:"countryOfDestination"`
	VatNo                string         `db:"vatNo"`
	EoriNo               sql.NullString `db:"eoriNo"`
	ImportAddressCode    string         `db:"importAddressCode"`
	AddressCode          string         `db:"addressCode"`
	AddressDetail        sql.NullString `db:"addressDetail"`
	PostalCode           sql.NullString `db:"postalCode"`
	City                 string         `db:"city"`
	ProductNo            string         `db:"productNo"`
	Description          sql.NullString `db:"description"`
	Mrn                  string         `db:"mrn"`
	Mode                 string         `db:"mode"`
	CompanyName          string         `db:"companyName"`
	InICPFile            string         `db:"hasInIcp"`
}

// TaxFileObject The object of the tax file
type TaxFileObject struct {
	Sn        int
	Mrn       string
	CustomsId string
	// 4, 115
	TaxType     int
	TaxFileLink string
}

// PodFileObject The object of the pod file
type PodFileObject struct {
	Sn          int            `db:"sn"`
	BillNo      string         `db:"bill_no"`
	CustomsId   string         `db:"customs_id"`
	Mrn         string         `db:"mrn"`
	TrackingNo  string         `db:"tracking_no"`
	PodFileLink sql.NullString `db:"uri"`
}

type CustomsServiceKeyObject struct {
	CustomsId  string `db:"customs_id"`
	MinIndexNo int    `db:"min_index_no"`
	ServiceKey string `db:"service_key"`
}
