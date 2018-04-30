package golinode

import (
	"fmt"
	"time"

	"github.com/go-resty/resty"
)

// Invoice represents a Invoice object
type Invoice struct {
	DateStr string `json:"date"`

	ID    int
	Label string
	Total float32
	Date  *time.Time `json:"-"`
}

type InvoiceItem struct {
	FromStr string `json:"from"`
	ToStr   string `json:"to"`

	Label     string
	Type      string
	UnitPrice int
	Quantity  int
	Amount    float32
	From      *time.Time `json:"-"`
	To        *time.Time `json:"-"`
}

// InvoicesPagedResponse represents a paginated Invoice API response
type InvoicesPagedResponse struct {
	*PageOptions
	Data []*Invoice
}

// Endpoint gets the endpoint URL for Invoice
func (InvoicesPagedResponse) Endpoint(c *Client) string {
	endpoint, err := c.Invoices.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends Invoices when processing paginated Invoice responses
func (resp *InvoicesPagedResponse) AppendData(r *InvoicesPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of Invoice
func (InvoicesPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InvoicesPagedResponse{})
}

// ListInvoices lists Invoices
func (c *Client) ListInvoices(opts *ListOptions) ([]*Invoice, error) {
	response := InvoicesPagedResponse{}
	err := c.ListHelper(response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Invoice) fixDates() *Invoice {
	v.Date, _ = parseDates(v.DateStr)
	return v
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *InvoiceItem) fixDates() *InvoiceItem {
	v.From, _ = parseDates(v.FromStr)
	v.To, _ = parseDates(v.ToStr)
	return v
}

// GetInvoice gets the template with the provided ID
func (c *Client) GetInvoice(id string) (*Invoice, error) {
	e, err := c.Invoices.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%s", e, id)
	r, err := c.R().SetResult(&Invoice{}).Get(e)
	if err != nil {
		return nil, err
	}
	return r.Result().(*Invoice).fixDates(), nil
}

// InvoiceItemsPagedResponse represents a paginated Invoice Item API response
type InvoiceItemsPagedResponse struct {
	*PageOptions
	Data []*InvoiceItem
}

// Endpoint gets the endpoint URL for InvoiceItems
func (InvoiceItemsPagedResponse) EndpointWithID(c *Client, id int) string {
	endpoint, err := c.InvoiceItems.EndpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// AppendData appends InvoiceItems when processing paginated Invoice Item responses
func (resp *InvoiceItemsPagedResponse) AppendData(r *InvoiceItemsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// SetResult sets the Resty response type of InvoiceItems
func (InvoiceItemsPagedResponse) SetResult(r *resty.Request) {
	r.SetResult(InvoiceItemsPagedResponse{})
}

// ListInvoiceItems lists Invoice Items
func (c *Client) ListInvoiceItems(id int, opts *ListOptions) ([]*InvoiceItem, error) {
	response := InvoiceItemsPagedResponse{}
	err := c.ListHelperWithID(response, id, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	return response.Data, err
}
