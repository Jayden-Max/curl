package curl

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCURL_CreateRequest(t *testing.T) {
	type fields struct {
		ctx              context.Context
		cancel           context.CancelFunc
		url              string
		method           string
		cookie           string
		referer          string
		Headers          map[string]string
		Options          map[string]bool
		PostBytes        []byte
		PostString       string
		PostFields       url.Values
		PostReader       io.Reader
		PostFieldReaders map[string]io.Reader
		PostFiles        url.Values
		timeout          time.Duration
	}
	tests := []struct {
		name        string
		fields      fields
		wantRequest *http.Request
		wantErr     bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{
				ctx:              tt.fields.ctx,
				cancel:           tt.fields.cancel,
				URL:              tt.fields.url,
				Method:           tt.fields.method,
				Cookie:           tt.fields.cookie,
				Referer:          tt.fields.referer,
				Headers:          tt.fields.Headers,
				Options:          tt.fields.Options,
				PostBytes:        tt.fields.PostBytes,
				PostString:       tt.fields.PostString,
				PostFields:       tt.fields.PostFields,
				PostReader:       tt.fields.PostReader,
				PostFieldReaders: tt.fields.PostFieldReaders,
				PostFiles:        tt.fields.PostFiles,
				Timeout:          tt.fields.timeout,
			}
			gotRequest, err := c.CreateRequest()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRequest, tt.wantRequest) {
				t.Errorf("CreateRequest() gotRequest = %v, want %v", gotRequest, tt.wantRequest)
			}
		})
	}
}

func TestCURL_Do(t *testing.T) {
	type fields struct {
		url        string
		method     string
		Headers    map[string]string
		Options    map[string]bool
		PostBytes  []byte
		PostString string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *Response
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				url:     "http://www.baidu.com",
				method:  "GET",
				Headers: make(map[string]string),
			},
			args:     args{context.TODO()},
			wantResp: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{
				URL:     tt.fields.url,
				Method:  tt.fields.method,
				Headers: tt.fields.Headers,
			}

			c.SetDefaultHeaders()

			gotResp, err := c.Do(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Do() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			strings.Trim(gotResp.Body, "\n")
			t.Logf("\nbody:%v", gotResp.Body)
		})
	}
}

func TestCURL_HandleResponse(t *testing.T) {
	type fields struct {
		ctx              context.Context
		cancel           context.CancelFunc
		url              string
		method           string
		cookie           string
		referer          string
		Headers          map[string]string
		Options          map[string]bool
		PostBytes        []byte
		PostString       string
		PostFields       url.Values
		PostReader       io.Reader
		PostFieldReaders map[string]io.Reader
		PostFiles        url.Values
		timeout          time.Duration
	}
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *Response
		wantErr      bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{
				ctx:              tt.fields.ctx,
				cancel:           tt.fields.cancel,
				URL:              tt.fields.url,
				Method:           tt.fields.method,
				Cookie:           tt.fields.cookie,
				Referer:          tt.fields.referer,
				Headers:          tt.fields.Headers,
				Options:          tt.fields.Options,
				PostBytes:        tt.fields.PostBytes,
				PostString:       tt.fields.PostString,
				PostFields:       tt.fields.PostFields,
				PostReader:       tt.fields.PostReader,
				PostFieldReaders: tt.fields.PostFieldReaders,
				PostFiles:        tt.fields.PostFiles,
				Timeout:          tt.fields.timeout,
			}
			gotResponse, err := c.HandleResponse(tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("HandleResponse() gotResponse = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestCURL_createPostRequest(t *testing.T) {
	type fields struct {
		ctx              context.Context
		cancel           context.CancelFunc
		url              string
		method           string
		cookie           string
		referer          string
		Headers          map[string]string
		Options          map[string]bool
		PostBytes        []byte
		PostString       string
		PostFields       url.Values
		PostReader       io.Reader
		PostFieldReaders map[string]io.Reader
		PostFiles        url.Values
		timeout          time.Duration
	}
	tests := []struct {
		name        string
		fields      fields
		wantRequest *http.Request
		wantErr     bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CURL{
				ctx:              tt.fields.ctx,
				cancel:           tt.fields.cancel,
				URL:              tt.fields.url,
				Method:           tt.fields.method,
				Cookie:           tt.fields.cookie,
				Referer:          tt.fields.referer,
				Headers:          tt.fields.Headers,
				Options:          tt.fields.Options,
				PostBytes:        tt.fields.PostBytes,
				PostString:       tt.fields.PostString,
				PostFields:       tt.fields.PostFields,
				PostReader:       tt.fields.PostReader,
				PostFieldReaders: tt.fields.PostFieldReaders,
				PostFiles:        tt.fields.PostFiles,
				Timeout:          tt.fields.timeout,
			}
			gotRequest, err := c.createPostRequest()
			if (err != nil) != tt.wantErr {
				t.Errorf("createPostRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRequest, tt.wantRequest) {
				t.Errorf("createPostRequest() gotRequest = %v, want %v", gotRequest, tt.wantRequest)
			}
		})
	}
}

func TestNewCurl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *CURL
	}{
		// TODO: Add test cases.
		{
			name: "curl1",
			args: args{url: "http://www.baidu.com"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCurl(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCurl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_ReadBody(t *testing.T) {
	type fields struct {
		Headers    map[string]string
		Cookie     string
		URL        string
		Body       []byte
		BodyReader io.ReadCloser
	}
	tests := []struct {
		name     string
		fields   fields
		wantBody []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				Headers:    tt.fields.Headers,
				Cookie:     tt.fields.Cookie,
				URL:        tt.fields.URL,
				BodyReader: tt.fields.BodyReader,
			}
			gotBody, err := r.ReadBody()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("ReadBody() gotBody = %v, want %v", gotBody, tt.wantBody)
			}

		})
	}
}
