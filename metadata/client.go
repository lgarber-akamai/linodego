package metadata

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

const APIHost = "169.254.169.254"
const APIProto = "http"
const APIVersion = "v1"

type ClientCreateOptions struct {
	HTTPClient *http.Client

	BaseURLOverride string
	VersionOverride string
	SchemeOverride  string

	DisableTokenInit bool
}

type Client struct {
	resty *resty.Client

	apiBaseURL  string
	apiProtocol string
	apiVersion  string
}

func NewClient(ctx context.Context, opts *ClientCreateOptions) (*Client, error) {
	var result Client

	shouldUseHTTPClient := false
	shouldSkipTokenGeneration := false

	if opts != nil {
		shouldUseHTTPClient = opts.HTTPClient != nil
		shouldSkipTokenGeneration = opts.DisableTokenInit

		if opts.BaseURLOverride != "" {
			result.SetBaseURL(opts.BaseURLOverride)
		}

		if opts.VersionOverride != "" {
			result.SetVersion(opts.VersionOverride)
		}
	}

	if shouldUseHTTPClient {
		result.resty = resty.NewWithClient(opts.HTTPClient)
	} else {
		result.resty = resty.New()
	}

	if debugEnv, ok := os.LookupEnv("LINODE_DEBUG"); ok {
		debugBool, err := strconv.ParseBool(debugEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse debug bool: %s", err)
		}
		result.resty.SetDebug(debugBool)
	}

	result.updateHostURL()

	if !shouldSkipTokenGeneration {
		token, err := result.GenerateToken(ctx, GenerateTokenOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to generate metadata token: %s", err)
		}

		result.UseToken(token)
	}

	return &result, nil
}

func (c *Client) UseToken(token string) *Client {
	c.resty.SetHeader("X-Metadata-Token", token)
	return c
}

func (c *Client) SetBaseURL(baseURL string) *Client {
	baseURLPath, _ := url.Parse(baseURL)

	c.apiBaseURL = path.Join(baseURLPath.Host, baseURLPath.Path)
	c.apiProtocol = baseURLPath.Scheme

	c.updateHostURL()

	return c
}

func (c *Client) SetVersion(version string) *Client {
	c.apiVersion = version

	c.updateHostURL()

	return c
}

func (c *Client) updateHostURL() {
	apiProto := APIProto
	baseURL := APIHost
	apiVersion := APIVersion

	if c.apiBaseURL != "" {
		baseURL = c.apiBaseURL
	}

	if c.apiVersion != "" {
		apiVersion = c.apiVersion
	}

	if c.apiProtocol != "" {
		apiProto = c.apiProtocol
	}

	c.resty.SetHostURL(fmt.Sprintf("%s://%s/%s", apiProto, baseURL, apiVersion))
}

// R wraps resty's R method
func (c *Client) R(ctx context.Context) *resty.Request {
	return c.resty.R().
		ExpectContentType("application/json").
		SetHeader("Content-Type", "application/json").
		SetContext(ctx)
}
