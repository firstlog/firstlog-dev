// Licensed to Elasticsearch B.V under one or more agreements.
// Elasticsearch B.V. licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.
//
// Code generated from specification version 6.8.5: DO NOT EDIT

package esapi

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func newXPackWatcherAckWatchFunc(t Transport) XPackWatcherAckWatch {
	return func(watch_id string, o ...func(*XPackWatcherAckWatchRequest)) (*Response, error) {
		var r = XPackWatcherAckWatchRequest{WatchID: watch_id}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// XPackWatcherAckWatch - http://www.elastic.co/guide/en/elasticsearch/reference/current/watcher-api-ack-watch.html
//
type XPackWatcherAckWatch func(watch_id string, o ...func(*XPackWatcherAckWatchRequest)) (*Response, error)

// XPackWatcherAckWatchRequest configures the X Pack Watcher Ack Watch API request.
//
type XPackWatcherAckWatchRequest struct {
	ActionID []string
	WatchID  string

	MasterTimeout time.Duration

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
//
func (r XPackWatcherAckWatchRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "PUT"

	path.Grow(1 + len("_xpack") + 1 + len("watcher") + 1 + len("watch") + 1 + len(r.WatchID) + 1 + len("_ack") + 1 + len(strings.Join(r.ActionID, ",")))
	path.WriteString("/")
	path.WriteString("_xpack")
	path.WriteString("/")
	path.WriteString("watcher")
	path.WriteString("/")
	path.WriteString("watch")
	path.WriteString("/")
	path.WriteString(r.WatchID)
	path.WriteString("/")
	path.WriteString("_ack")
	if len(r.ActionID) > 0 {
		path.WriteString("/")
		path.WriteString(strings.Join(r.ActionID, ","))
	}

	params = make(map[string]string)

	if r.MasterTimeout != 0 {
		params["master_timeout"] = formatDuration(r.MasterTimeout)
	}

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
//
func (f XPackWatcherAckWatch) WithContext(v context.Context) func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.ctx = v
	}
}

// WithActionID - a list of the action ids to be acked.
//
func (f XPackWatcherAckWatch) WithActionID(v ...string) func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.ActionID = v
	}
}

// WithMasterTimeout - explicit operation timeout for connection to master node.
//
func (f XPackWatcherAckWatch) WithMasterTimeout(v time.Duration) func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.MasterTimeout = v
	}
}

// WithPretty makes the response body pretty-printed.
//
func (f XPackWatcherAckWatch) WithPretty() func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
//
func (f XPackWatcherAckWatch) WithHuman() func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
//
func (f XPackWatcherAckWatch) WithErrorTrace() func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
//
func (f XPackWatcherAckWatch) WithFilterPath(v ...string) func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
//
func (f XPackWatcherAckWatch) WithHeader(h map[string]string) func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
//
func (f XPackWatcherAckWatch) WithOpaqueID(s string) func(*XPackWatcherAckWatchRequest) {
	return func(r *XPackWatcherAckWatchRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
