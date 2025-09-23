package ecommon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// makeRequest 使用gin处理HTTP和JSON binding，callback处理解析后的Request
func makeRequest[TReq any](data interface{}, callback func(TReq)) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/test", func(c *gin.Context) {
		var parsedReq TReq
		if err := c.ShouldBindJSON(&parsedReq); err != nil {
			panic(err)
		}
		callback(parsedReq)
		c.Status(http.StatusOK)
	})

	var reqBody *bytes.Buffer
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	httpReq, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		panic(err)
	}

	if data != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		panic("HTTP request failed")
	}
}

// makeFormRequest 使用gin处理HTTP和Form binding，callback处理解析后的Request
func makeFormRequest[TReq any](formData map[string]string, callback func(TReq)) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/test", func(c *gin.Context) {
		var parsedReq TReq
		if err := c.ShouldBind(&parsedReq); err != nil {
			panic(err)
		}
		callback(parsedReq)
		c.Status(http.StatusOK)
	})

	var reqBody *strings.Reader
	if formData != nil {
		values := url.Values{}
		for k, v := range formData {
			values.Set(k, v)
		}
		reqBody = strings.NewReader(values.Encode())
	} else {
		reqBody = strings.NewReader("")
	}

	httpReq, err := http.NewRequest("POST", "/test", reqBody)
	if err != nil {
		panic(err)
	}

	if formData != nil {
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		panic("HTTP request failed")
	}
}

// makeGetRequest 使用gin处理HTTP GET和Query binding，callback处理解析后的Request
func makeGetRequest[TReq any](queryParams map[string]string, callback func(TReq)) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/test", func(c *gin.Context) {
		var parsedReq TReq
		if err := c.ShouldBindQuery(&parsedReq); err != nil {
			panic(err)
		}
		callback(parsedReq)
		c.Status(http.StatusOK)
	})

	reqURL := "/test"
	if queryParams != nil && len(queryParams) > 0 {
		values := url.Values{}
		for k, v := range queryParams {
			values.Set(k, v)
		}
		reqURL += "?" + values.Encode()
	}

	httpReq, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	if w.Code != http.StatusOK {
		panic("HTTP request failed")
	}
}
