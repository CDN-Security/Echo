package controller

import (
	"log/slog"
	"net/http"

	"github.com/CDN-Security/Echo/pkg/config"
	"github.com/CDN-Security/Echo/pkg/model"
	http_grab_model "github.com/WangYihang/http-grab/pkg/model"
	"github.com/gin-gonic/gin"
)

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
	Response  string `json:"response"`
}

type EchoServerResponseBody struct {
	// Misc
	RemoteAddr string `json:"remote_addr"`
	ClientIp   string `json:"client_ip"`
	// TLS Related Fields
	TLS *http_grab_model.TLS `json:"tls,omitempty"`
	// Request Related Fields
	HTTP http_grab_model.HTTP `json:"http"`
	// Challenge and Response
	ChallengeResponse map[string]ChallengeResponse `json:"challenge_response"`
}

func Handler(c *gin.Context) {
	var httpRequest *http_grab_model.HTTPRequest
	var err error

	// GET /?cdn_challenge=${challenge} HTTP/1.1
	// Host: www.example.com
	queryChallenge := model.ExtractQueryChallenge(c)
	queryResponse := model.AcceptChallenge(queryChallenge, config.DefaultConfig.ChallengeConfig.SecretKey)

	// GET / HTTP/1.1
	// Host: www.example.com
	// CDN-Challenge: ${challenge}
	headerChallenge := model.ExtractHeaderChallenge(c)
	headerResponse := model.AcceptChallenge(headerChallenge, config.DefaultConfig.ChallengeConfig.SecretKey)

	// GET / HTTP/1.1
	// Host: www.example.com
	// Cookie: cdn_challenge=${challenge}
	cookieChallenge := model.ExtractCookieChallenge(c)
	cookieResponse := model.AcceptChallenge(cookieChallenge, config.DefaultConfig.ChallengeConfig.SecretKey)

	httpRequest, err = http_grab_model.NewHTTPRequest(c.Request)
	if err != nil {
		slog.Error("error occured while parsing HTTP request", slog.String("error", err.Error()))
	}

	responseBody := EchoServerResponseBody{
		// Misc
		RemoteAddr: c.Request.RemoteAddr,
		ClientIp:   c.ClientIP(),
		// Request Related Fields
		HTTP: http_grab_model.HTTP{
			Request:  httpRequest,
			Response: nil,
		},
		// Challenge and Response
		ChallengeResponse: map[string]ChallengeResponse{
			"query":  {Challenge: queryChallenge, Response: queryResponse},
			"header": {Challenge: headerChallenge, Response: headerResponse},
			"cookie": {Challenge: cookieChallenge, Response: cookieResponse},
		},
	}

	if c.Request.TLS != nil {
		responseBody.TLS = http_grab_model.NewTLS(c.Request.TLS)
	}

	c.JSON(http.StatusOK, responseBody)
}
