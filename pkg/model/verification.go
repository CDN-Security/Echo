package model

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"

	"github.com/CDN-Security/Echo/pkg/config"
	"github.com/gin-gonic/gin"
)

func ExtractPathChallenge(c *gin.Context) string {
	basename := filepath.Base(c.Request.URL.Path)
	filename := basename[:len(basename)-len(filepath.Ext(basename))]
	return filename
}

func ExtractQueryChallenge(c *gin.Context) string {
	return c.Request.URL.Query().Get(config.DefaultConfig.ChallengeConfig.QueryName)
}

func ExtractHeaderChallenge(c *gin.Context) string {
	return c.Request.Header.Get(config.DefaultConfig.ChallengeConfig.HeaderName)
}

func ExtractCookieChallenge(c *gin.Context) string {
	cookie, err := c.Request.Cookie(config.DefaultConfig.ChallengeConfig.CookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func AcceptChallenge(challenge string, secret string) string {
	return Sha256([]byte(challenge + secret))
}

func Sha256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

func Md5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
