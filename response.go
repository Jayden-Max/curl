package curl

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Response struct {
	Headers      map[string]string `json:"headers"`
	Cookie       string            `json:"Cookie"`
	URL          string            `json:"URL"`
	Body         string            `json:"body"`
	HttpResponse *http.Response
	BodyReader   io.ReadCloser
}

func (r *Response) ReadBody() (string, error) {
	body, err := ioutil.ReadAll(r.BodyReader)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(body)), nil
}

func (r *Response) Close() {
	if r.BodyReader != nil {
		_ = r.BodyReader.Close()
	}
}

func (c *CURL) HandleResponse(resp *http.Response) (response *Response, err error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	response = new(Response)
	response.HttpResponse = resp
	response.Headers = rcHeader(resp.Header)
	// Referer info
	location, _ := resp.Location()
	if location != nil {
		locationUrl := location.String()
		response.Headers["Location"] = locationUrl

		if c.Options["Redirect"] {
			response, err = rcReferer(c.URL, locationUrl)
			return
		}
	}

	response.Headers["Status"] = resp.Status
	response.Headers["Status-Code"] = strconv.Itoa(resp.StatusCode)
	response.Headers["Proto"] = resp.Proto
	response.Cookie = rcCookie(resp)
	response.URL = c.URL

	response.BodyReader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		response.BodyReader, err = gzip.NewReader(resp.Body)
	}
	if strings.Contains(resp.Header.Get("Content-Encoding"), "deflate") {
		response.BodyReader = flate.NewReader(resp.Body)
	}

	response.Body, err = response.ReadBody()
	if err != nil {
		return nil, err
	}

	return
}

func rcCookie(resp *http.Response) string {
	cookie := make([]string, 0, len(resp.Cookies()))
	for _, v := range resp.Cookies() {
		cookie = append(cookie, v.Name+"="+v.Value)
	}

	sort.Strings(cookie)

	return strings.TrimSpace(strings.Join(cookie, "; "))
}

func rcReferer(url, locationUrl string) (*Response, error) {
	curl := NewCurl(url, "GET")
	curl.Referer = locationUrl

	resp, err := curl.Do(context.TODO())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func rcHeader(header map[string][]string) map[string]string {
	headers := make(map[string]string)
	for k, v := range header {
		headers[k] = strings.Join(v, " ")
	}

	return headers
}
