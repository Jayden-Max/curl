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
