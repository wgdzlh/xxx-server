package client

import (
	"io"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

	json "github.com/json-iterator/go"
	"go.uber.org/zap"
)

const (
	JSON_ACCEPT       = "application/json"
	JSON_CONTENT_TYPE = "application/json; charset=utf-8"
)

/*
Client 标准库http.Client，优化参数。建议不要直接用net/http库的默认Client
主要解决： 1. 无超时限制 2. 有时需要跳过证书验证
var DefaultHttpClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	},
	Timeout: time.Second * 10,
}
*/

type HttpClient struct {
	name string
	cli  *http.Client
}

func NewHttpClient(name string, timeoutSecs int) *HttpClient {
	c := &HttpClient{
		name: name,
		cli:  &http.Client{},
	}
	c.SetTimeout(timeoutSecs)
	return c
}

func (c *HttpClient) SetTimeout(timeoutSecs int) {
	c.cli.Timeout = time.Second * time.Duration(timeoutSecs)
}

// PostJson implements ClientRepo
func (c *HttpClient) PostJson(url string, req, resp any) (err error) {
	body, err := utils.GetJsonBody(req)
	if err != nil {
		log.Error(c.name+":Post:GetJsonBody", zap.Error(err))
		return
	}
	response, err := c.cli.Post(url, JSON_CONTENT_TYPE, body)
	if err != nil {
		log.Error(c.name+":Post", zap.Error(err))
		return
	}
	if resp == nil {
		return
	}
	defer response.Body.Close()
	respBs, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(c.name+":Post:ReadAll", zap.Error(err))
		return
	}
	if err = json.Unmarshal(respBs, resp); err != nil {
		log.Error(c.name+":Post:Unmarshal", zap.Error(err))
	}
	return
}

func (c *HttpClient) Get(urlPath string, query map[string]string, resp any) (err error) {
	if len(query) > 0 {
		queries := url.Values{}
		for k, v := range query {
			queries.Set(k, v)
		}
		if strings.ContainsRune(urlPath, '?') {
			urlPath += "&" + queries.Encode()
		} else {
			urlPath += "?" + queries.Encode()
		}
	}
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		log.Error(c.name+":Get:NewRequest", zap.Error(err))
		return
	}
	// req.Header.Set("accept", JSON_ACCEPT) // 默认不需要
	response, err := c.cli.Do(req)
	if err != nil {
		log.Error(c.name+":Get", zap.Error(err))
		return
	}
	if resp == nil {
		return
	}
	defer response.Body.Close()
	respBs, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(c.name+":Get:ReadAll", zap.Error(err))
		return
	}
	if err = json.Unmarshal(respBs, resp); err != nil {
		log.Error(c.name+":Get:Unmarshal", zap.Error(err))
	}
	return
}

func (c *HttpClient) PostFile(url, contentType string, body io.Reader, resp any) (err error) {
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Accept", "application/json")
	response, err := c.cli.Do(r)
	if err != nil {
		log.Error(c.name+":PostFile", zap.Error(err))
		return
	}
	if resp == nil {
		return
	}
	defer response.Body.Close()
	respBs, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(c.name+":PostFile:ReadAll", zap.Error(err))
		return
	}
	if err = json.Unmarshal(respBs, resp); err != nil {
		log.Error(c.name+":PostFile:Unmarshal", zap.Error(err))
	}
	return
}

func (c *HttpClient) Delete(url string, req, resp any) (err error) {
	var body io.Reader
	if req != nil {
		if body, err = utils.GetJsonBody(req); err != nil {
			log.Error(c.name+":Delete:GetJsonBody", zap.Error(err))
			return
		}
	}
	r, err := http.NewRequest("DELETE", url, body)
	if err != nil {
		return
	}
	r.Header.Set("Accept", "application/json")
	response, err := c.cli.Do(r)
	if err != nil {
		log.Error(c.name+":Delete", zap.Error(err))
		return
	}
	if resp == nil {
		return
	}
	defer response.Body.Close()
	respBs, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(c.name+":Delete:ReadAll", zap.Error(err))
		return
	}
	if len(respBs) == 0 {
		return
	}
	if err = json.Unmarshal(respBs, resp); err != nil {
		log.Error(c.name+":Delete:Unmarshal", zap.Error(err))
	}
	return
}
