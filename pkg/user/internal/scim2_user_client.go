/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"net/http"
	"net/url"
	"strings"
)

type AddUserJSONBody struct {
	Username string  `json:"userName"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     Name    `json:"name"`
	Emails   []Email `json:"emails"`
}

type Name struct {
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

type Email struct {
	Primary bool   `json:"primary"`
	Value   string `json:"value"`
}

type PatchMeJSONBody struct {
	Schemas    []string         `json:"schemas"`
	Operations []PatchOperation `json:"Operations"`
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// The interface specification for the client above.
type ClientInterface interface {
	AddUser(ctx context.Context, body AddUserJSONBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	GetMe(ctx context.Context, attributes string, reqEditors ...RequestEditorFn) (*http.Response, error)
	PatchMe(ctx context.Context, body PatchMeJSONBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	GetSchemas(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) AddUser(ctx context.Context, body AddUserJSONBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddUserRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func NewAddUserRequest(server string, body AddUserJSONBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddUserRequestWithBody(server, "application/json", bodyReader)
}

func NewAddUserRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/Users"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) GetMe(ctx context.Context, attributes string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMeRequest(c.Server, attributes)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, reqEditors)
}

func (c *Client) PatchMe(ctx context.Context, body PatchMeJSONBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchMeRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, reqEditors)
}

func (c *Client) GetSchemas(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := newRequest(c.Server, "GET", "/Schemas", nil)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, reqEditors)
}

func (c *Client) do(ctx context.Context, req *http.Request, reqEditors []RequestEditorFn) (*http.Response, error) {
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func NewGetMeRequest(server string, attributes string) (*http.Request, error) {
	req, err := newRequest(server, "GET", "/Me", nil)
	if err != nil {
		return nil, err
	}
	if attributes != "" {
		req.URL.RawQuery = url.Values{"attributes": {attributes}}.Encode()
	}
	return req, nil
}

func NewPatchMeRequest(server string, body PatchMeJSONBody) (*http.Request, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := newRequest(server, "PATCH", "/Me", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/scim+json")
	return req, nil
}

// newRequest creates a request for a path relative to the server URL
func newRequest(server, method, operationPath string, body io.Reader) (*http.Request, error) {
	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	queryURL, err := serverURL.Parse("." + operationPath)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, queryURL.String(), body)
}
