package curl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// refer:
// http://www.ruanyifeng.com/blog/2011/09/curl.html
// http://www.ruanyifeng.com/blog/2019/09/curl-reference.html
type CURL struct {
	// -H 添加http请求标头
	URL     string
	Method  string
	Cookie  string
	Referer string

	Headers map[string]string // header
	Options map[string]bool

	// --data
	PostBytes        []byte               // binary data
	PostString       string               // string data
	PostFields       neturl.Values        // k-v format
	PostReader       io.Reader            // from reader
	PostFieldReaders map[string]io.Reader // from map reader
	PostFiles        neturl.Values        // key 文件名；value 文件路径

	Timeout time.Duration      // ctx 超时时间
	ctx     context.Context    // ctx
	cancel  context.CancelFunc // cancel
}

func NewCurl(method, url string) *CURL {
	return &CURL{
		Method:           method,
		URL:              url,
		Headers:          make(map[string]string),
		Options:          make(map[string]bool),
		PostBytes:        make([]byte, 0),
		PostFields:       neturl.Values{},
		PostFieldReaders: make(map[string]io.Reader),
		PostFiles:        neturl.Values{},
	}
}

func (c *CURL) SetUrl(url string) {
	if c.URL == "" {
		c.URL = url
	}
}

func (c *CURL) SetMethod(method string) error {
	Method := strings.ToUpper(method)
	switch method {
	case "GET", "POST", "DELETE", "OPTIONS", "HEAD", "PUT", "CONNECT", "TRACE":
		c.Method = Method
		return nil
	default:
		return fmt.Errorf("invalid Method %q", method)
	}
}

func (c *CURL) SetDefaultHeaders() {
	if c.Headers == nil {
		return
	}
	// 配置request default参数
	c.Headers["Content-Type"] = "application/json"
	c.Headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	c.Headers["Accept-Encoding"] = "gzip, deflate"
	c.Headers["Accept-Language"] = "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3"
	c.Headers["Connection"] = "keep-alive"
	c.Headers["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36"
}

func (c *CURL) SetHeader(key, value string) {
	if c.Headers == nil {
		return
	}
	key = http.CanonicalHeaderKey(key)
	c.Headers[key] = value
}

func (c *CURL) SetContext(ctx context.Context) {
	if c.ctx == nil {
		c.ctx, c.cancel = context.WithCancel(ctx)
	}
}

func (c *CURL) SetTimeout(sec int) {
	c.SetContext(context.Background())
	c.Timeout = time.Duration(sec) * time.Second
}

func (c *CURL) SetOption(key string, value bool) {
	if c.Options == nil {
		return
	}
	key = http.CanonicalHeaderKey(key)
	c.Options[key] = value
}

func (c *CURL) SetCookie(cookie string) {
	if c.Cookie == "" {
		c.Cookie = cookie
	}
}

// bytes for data content
func (c *CURL) SetPostBytes(p []byte) {
	if c.PostBytes == nil {
		return
	}
	c.PostBytes = p
}

// string for data content
func (c *CURL) SetPostString(p string) {
	if c.PostString == "" {
		c.PostString = p
	}
}

// map[string][]string for data content
func (c *CURL) SetPostFields(p neturl.Values) {
	if c.PostFields == nil {
		return
	}
	c.PostFields = p
}

// io.Reader for data content
func (c *CURL) SetPostReader(p io.Reader) {
	if c.PostReader == nil {
		return
	}
	c.PostReader = p
}

// field's io.Reader for data content
func (c *CURL) SetPostFieldReaders(p map[string]io.Reader) {
	if c.PostFieldReaders == nil {
		return
	}
	c.PostFieldReaders = p
}

// file data for content
func (c *CURL) SetPostFiles(p neturl.Values) {
	if c.PostFiles == nil {
		return
	}
	c.PostFiles = p
}

/*
 start Do http request

 if ctx is nil then set ctx equal context.Background()
 and meantime set timeout 5 second.
*/
func (c *CURL) Do(ctx context.Context) (resp *Response, err error) {
	if ctx != nil {
		c.SetContext(ctx)
	}

	if c.Timeout <= 0 || c.ctx == nil {
		// default Timeout 5 sec
		c.SetTimeout(5)
	}

	var httpRequest *http.Request
	var httpResponse *http.Response

	httpRequest, err = c.CreateRequest()
	if err != nil {
		return nil, errors.Wrap(err, "Curl.Do -> func CreateRequest failed.")
	}

	httpRequest = httpRequest.WithContext(c.ctx)

	ch := make(chan struct{}, 1)
	go func() {
		client := http.Client{}
		httpResponse, err = client.Do(httpRequest)
		ch <- struct{}{}
	}()

	select {
	case <-time.After(c.Timeout):
		err = errors.New("requestTimeout")
	case <-c.ctx.Done():
		err = c.ctx.Err()
	case <-ch:
		close(ch)
	}

	if err != nil {
		c.cancel()
		return
	}

	return c.HandleResponse(httpResponse)
}

/*
 return http's request and set headers.
*/
func (c *CURL) CreateRequest() (request *http.Request, err error) {
	if c.Method == "" || c.URL == "" {
		return nil, fmt.Errorf("method: [%+v] or URL: [%+v] is empty", c.Method, c.URL)
	}

	if c.PostString != "" || len(c.PostString) > 0 || len(c.PostFields) > 0 ||
		len(c.PostFieldReaders) > 0 || len(c.PostFiles) > 0 {
		request, err = c.createPostRequest()
	} else {
		request, err = http.NewRequest(c.Method, c.URL, nil)
	}

	if err != nil {
		return nil, errors.Wrap(err, "create request failed.")
	}

	if c.Headers != nil {
		for k, v := range c.Headers {
			request.Header.Add(k, v)
		}
	}

	if c.Cookie != "" {
		request.Header.Add("Cookie", c.Cookie)
	}

	if c.Referer != "" {
		request.Header.Add("Referer", c.Referer)
	}

	return
}

/*
 here for create post request.
 main is for data content.

 set request's body and header's 'Content-Type' params.
*/
func (c *CURL) createPostRequest() (request *http.Request, err error) {
	if c.Method != "POST" {
		return nil, fmt.Errorf("method:%v is not POST", c.Method)
	}

	var isSetHeader bool

	switch {
	case c.PostBytes != nil:
		b := bytes.NewReader(c.PostBytes)
		request, err = http.NewRequest(c.Method, c.URL, b)
		if err != nil {
			return nil, err
		}

	case c.PostString != "":
		b := strings.NewReader(c.PostString)
		request, err = http.NewRequest(c.Method, c.URL, b)
		if err != nil {
			return nil, err
		}

	case c.PostReader != nil:
		request, err = http.NewRequest(c.Method, c.URL, c.PostReader)
		if err != nil {
			return nil, err
		}

	case len(c.PostFields) > 0 || len(c.PostFieldReaders) > 0 || len(c.PostFiles) > 0:
		b := new(bytes.Buffer)
		bodyWriter := multipart.NewWriter(b)

		for key, val := range c.PostFields {
			for _, v := range val {
				_ = bodyWriter.WriteField(key, v)
			}
		}

		for key, val := range c.PostFieldReaders {
			fileWriter, _ := bodyWriter.CreateFormField(key)
			_, err = io.Copy(fileWriter, val)
			if err != nil {
				return nil, err
			}
		}

		// file
		for key, val := range c.PostFiles {
			for _, v := range val {
				_, err = os.Stat(v)
				if err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("PostFiles %s - %s not exist", key, v))
				}
				fileWriter, err := bodyWriter.CreateFormFile(key, v)
				if err != nil {
					return nil, errors.Wrap(err, "multipart CreateFormFile failed.")
				}

				fileInfo, err := os.Open(v)
				if err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("open file:%s failed", v))
				}

				_, err = io.Copy(fileWriter, fileInfo)
				_ = fileInfo.Close()
				if err != nil {
					return nil, errors.Wrap(err, "io.Copy fileInfo to fileWriter failed.")
				}
			}
		}
		_ = bodyWriter.Close()

		request, err = http.NewRequest(c.Method, c.URL, b)
		if err != nil {
			return nil, errors.Wrap(err, "create Request failed.")
		}

		request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		isSetHeader = true

	default:
		request, err = http.NewRequest(c.Method, c.URL, nil)
		if err != nil {
			return nil, err
		}
	}

	if !isSetHeader {
		if v, ok := c.Headers["Content-Type"]; ok {
			request.Header.Set("Content-Type", v)
		}
	}

	delete(c.Headers, "Content-Type")

	return
}
