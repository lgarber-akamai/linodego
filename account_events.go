package linodego

import (
	"time"

	"github.com/go-resty/resty"
)

// Event represents an action taken on the Account.
type Event struct {
	CreatedStr string `json:"created"`
	UpdatedStr string `json:"updated"`

	ID              int
	Status          string
	Action          string
	PercentComplete int `json:"percent_complete"`
	Rate            string
	Read            bool
	Seen            bool
	TimeRemaining   int
	Username        string
	Entity          *EventEntity
	Created         *time.Time `json:"-"`
}

// EventEntity provides detailed information about the Event's
// associated entity, including ID, Type, Label, and a URL that
// can be used to access it.
type EventEntity struct {
	ID    int
	Label string
	Type  string
	URL   string
}

// EventsPagedResponse represents a paginated Events API response
type EventsPagedResponse struct {
	*PageOptions
	Data []*Event
}

// endpoint gets the endpoint URL for Event
func (EventsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Events.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// endpointWithID gets the endpoint URL for a specific Event
func (EventsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.Events.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends Events when processing paginated Event responses
func (resp *EventsPagedResponse) appendData(r *EventsPagedResponse) {
	(*resp).Data = append(resp.Data, r.Data...)
}

// setResult sets the Resty response type of Events
func (EventsPagedResponse) setResult(r *resty.Request) {
	r.SetResult(EventsPagedResponse{})
}

// ListEvents gets a collection of Event objects representing actions taken
// on the Account. The Events returned depend on the token grants and the grants
// of the associated user.
func (c *Client) ListEvents(opts *ListOptions) ([]*Event, error) {
	response := EventsPagedResponse{}
	err := c.listHelper(&response, opts)
	for _, el := range response.Data {
		el.fixDates()
	}
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// fixDates converts JSON timestamps to Go time.Time values
func (v *Event) fixDates() *Event {
	v.Created, _ = parseDates(v.CreatedStr)
	return v
}
