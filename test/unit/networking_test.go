package unit

import (
	"context"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
)

func TestIPv6Range_ListMany(t *testing.T) {
	client := createMockClient(t)

	const totalResults = 4231
	const pageSize = 500
	totalPages := int(math.Ceil(float64(totalResults) / float64(pageSize)))

	expectedPages := make([]int, totalPages)
	for i := range expectedPages {
		expectedPages[i] = i + 1
	}

	buildResponse := func(page, pageSize int) map[string]any {
		pageResults := pageSize
		if page*pageSize > totalResults {
			pageResults = totalResults - pageSize*(page-1)
		}

		ranges := make([]linodego.IPv6Range, pageResults)
		for i := range ranges {
			ranges[i] = linodego.IPv6Range{
				Range:  "1234::5678",
				Region: "us-mia",
				Prefix: 64,
			}
		}

		return map[string]any{
			"page":    page,
			"pages":   totalPages,
			"results": totalResults,
			"data":    ranges,
		}
	}

	pagesRequested := make([]int, 0)

	httpmock.RegisterRegexpResponder(
		"GET", mockRequestURL(t, "networking/ipv6/ranges"),
		func(request *http.Request) (*http.Response, error) {
			page, _ := strconv.Atoi(request.URL.Query().Get("page"))
			currentPageSize, _ := strconv.Atoi(request.URL.Query().Get("page_size"))

			// Default to page 1
			if !request.URL.Query().Has("page") {
				page = 1
			}

			pagesRequested = append(pagesRequested, page)
			return httpmock.NewJsonResponse(200, buildResponse(page, currentPageSize))
		},
	)

	results, err := client.ListIPv6Ranges(context.Background(), &linodego.ListOptions{PageSize: pageSize})
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, reflect.DeepEqual(pagesRequested, expectedPages), cmp.Diff(pagesRequested, expectedPages))
	require.Equal(t, totalResults, len(results))
}
