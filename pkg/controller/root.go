package controller

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/CDN-Security/Echo/pkg/config"
	"github.com/CDN-Security/Echo/pkg/model"
	http_grab_model "github.com/WangYihang/http-grab/pkg/model"
	"github.com/gin-gonic/gin"
)

type Verification struct {
	Challenge string `json:"challenge"`
	Response  string `json:"response"`
}

type EchoServerResponseBody struct {
	// Misc
	RemoteAddr string `json:"remote_addr"`
	ClientIp   string `json:"client_ip"`
	Timestamp  int64  `json:"timestamp"`
	// TLS Related Fields
	TLS *http_grab_model.TLS `json:"tls,omitempty"`
	// Request Related Fields
	HTTP http_grab_model.HTTP `json:"http"`
	// Challenge and Response
	Verifications map[string]Verification `json:"verifications,omitempty"`
}

func Handler(c *gin.Context) {
	responseBody := EchoServerResponseBody{
		RemoteAddr:    c.Request.RemoteAddr,
		ClientIp:      c.ClientIP(),
		Timestamp:     time.Now().UnixMicro(),
		Verifications: make(map[string]Verification),
	}

	var httpRequest *http_grab_model.HTTPRequest
	var err error

	// GET /?cdn_challenge=${challenge} HTTP/1.1
	// Host: www.example.com
	queryChallenge := model.ExtractQueryChallenge(c)
	if queryChallenge != "" {
		queryResponse := model.AcceptChallenge(queryChallenge, config.DefaultConfig.VerificationConfig.SecretKey)
		responseBody.Verifications["query"] = Verification{Challenge: queryChallenge, Response: queryResponse}
	}

	// GET / HTTP/1.1
	// Host: www.example.com
	// CDN-Challenge: ${challenge}
	headerChallenge := model.ExtractHeaderChallenge(c)
	if headerChallenge != "" {
		headerResponse := model.AcceptChallenge(headerChallenge, config.DefaultConfig.VerificationConfig.SecretKey)
		responseBody.Verifications["header"] = Verification{Challenge: headerChallenge, Response: headerResponse}
		c.Header("Echo-Response", headerResponse)
	}

	// GET / HTTP/1.1
	// Host: www.example.com
	// Cookie: cdn_challenge=${challenge}
	cookieChallenge := model.ExtractCookieChallenge(c)
	if cookieChallenge != "" {
		cookieResponse := model.AcceptChallenge(cookieChallenge, config.DefaultConfig.VerificationConfig.SecretKey)
		responseBody.Verifications["cookie"] = Verification{Challenge: cookieChallenge, Response: cookieResponse}
		c.SetCookie("echo_response", cookieResponse, 0, "/", "", false, false)
	}

	httpRequest, err = http_grab_model.NewHTTPRequest(c.Request)
	if err != nil {
		slog.Error("error occured while parsing HTTP request", slog.String("error", err.Error()))
	}

	// Request Related Fields
	responseBody.HTTP = http_grab_model.HTTP{
		Request:  httpRequest,
		Response: nil,
	}

	if c.Request.TLS != nil {
		responseBody.TLS = http_grab_model.NewTLS(c.Request.TLS)
	}

	// GET /?status_code=200 HTTP/1.1
	// Host: www.example.com
	statusCode := c.Query("status_code")
	statusCodeInt, err := strconv.Atoi(statusCode)
	if err != nil {
		statusCodeInt = http.StatusOK
	}

	c.JSON(statusCodeInt, responseBody)
}
