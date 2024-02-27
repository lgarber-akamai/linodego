package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// paginatedResponse represents a single response from a paginated
// endpoint.
type paginatedResponse[T any] struct {
	Page    int `json:"page"    url:"page,omitempty"`
	Pages   int `json:"pages"   url:"pages,omitempty"`
	Results int `json:"results" url:"results,omitempty"`
	Data    []T `json:"data"`
}

// getPaginatedResults aggregates results from the given
// paginated endpoint using the provided ListOptions.
// nolint:funlen
func getPaginatedResults[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
	opts *ListOptions,
) ([]T, error) {
	var resultType paginatedResponse[T]

	result := make([]T, 0)

	req := client.R(ctx).SetResult(resultType)

	if opts == nil {
		opts = &ListOptions{PageOptions: &PageOptions{Page: 0}}
	}

	if opts.PageOptions == nil {
		opts.PageOptions = &PageOptions{Page: 0}
	}

	// Apply all user-provided list options to the base request
	if err := applyListOptionsToRequest(opts, req); err != nil {
		return nil, err
	}

	// Makes a request to a particular page and
	// appends the response to the result
	handlePage := func(page int) error {
		req.SetQueryParam("page", strconv.Itoa(page))

		res, err := coupleAPIErrors(req.Get(endpoint))
		if err != nil {
			return err
		}

		response := res.Result().(*paginatedResponse[T])

		opts.Page = page
		opts.Pages = response.Pages
		opts.Results = response.Results

		result = append(result, response.Data...)
		return nil
	}

	// This helps simplify the logic below
	startingPage := 1
	pageDefined := opts.Page > 0

	if pageDefined {
		startingPage = opts.Page
	}

	// Get the first page
	if err := handlePage(startingPage); err != nil {
		return nil, err
	}

	// If the user has explicitly specified a page, we don't
	// need to get any other pages.
	if pageDefined {
		return result, nil
	}

	// Get the rest of the pages
	for page := 2; page <= opts.Pages; page++ {
		if err := handlePage(page); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// doGETRequest runs a GET request using the given client and API endpoint,
// and returns the result
func doGETRequest[T any](
	ctx context.Context,
	client *Client,
	endpoint string,
) (*T, error) {
	var resultType T

	req := client.R(ctx).SetResult(&resultType)
	r, err := coupleAPIErrors(req.Get(endpoint))
	if err != nil {
		return nil, err
	}

	return r.Result().(*T), nil
}

// doPOSTRequest runs a PUT request using the given client, API endpoint,
// and options/body.
func doPOSTRequest[T, O any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options ...O,
) (*T, error) {
	var resultType T

	req := client.R(ctx).SetResult(&resultType)

	// `null` is not accepted by the API
	if len(options) > 1 && options[0] != nil {
		body, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}

		req = req.SetBody(body)
	}

	r, err := coupleAPIErrors(req.Post(endpoint))
	if err != nil {
		return nil, err
	}

	return r.Result().(*T), nil
}

// doPUTRequest runs a PUT request using the given client, API endpoint,
// and options/body.
func doPUTRequest[T, O any](
	ctx context.Context,
	client *Client,
	endpoint string,
	options ...O,
) (*T, error) {
	var resultType T

	req := client.R(ctx).SetResult(&resultType)

	// `null` is not accepted by the API
	if len(options) > 1 && options[0] != nil {
		body, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}

		req = req.SetBody(body)
	}

	r, err := coupleAPIErrors(req.Put(endpoint))
	if err != nil {
		return nil, err
	}

	return r.Result().(*T), nil
}

// doDELETERequest runs a DELETE request using the given client
// and API endpoint.
func doDELETERequest(
	ctx context.Context,
	client *Client,
	endpoint string,
) error {
	req := client.R(ctx)
	_, err := coupleAPIErrors(req.Delete(endpoint))
	return err
}

// formatAPIPath allows us to safely build an API request with path escaping
func formatAPIPath(format string, args ...any) string {
	escapedArgs := make([]any, len(args))
	for i, arg := range args {
		if typeStr, ok := arg.(string); ok {
			arg = url.PathEscape(typeStr)
		}

		escapedArgs[i] = arg
	}

	return fmt.Sprintf(format, escapedArgs...)
}
