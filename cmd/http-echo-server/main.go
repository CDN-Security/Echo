package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/CDN-Security/Echo/pkg/config"
	"github.com/CDN-Security/Echo/pkg/controller"
	"github.com/gin-gonic/gin"
)

func startHTTPServer(r *gin.Engine, addr string) {
	slog.Info("starting HTTP server", slog.String("addr", addr))
	err := r.Run(addr)
	if err != nil {
		slog.Error("error occured while running HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func startHTTPSServer(r *gin.Engine, addr string, certificatePath string, privateKeyPath string) {
	slog.Info("starting HTTPS server", slog.String("addr", addr), slog.String("certificate_path", certificatePath), slog.String("private_key_path", privateKeyPath))
	err := r.RunTLS(
		addr,
		certificatePath,
		privateKeyPath,
	)
	if err != nil {
		slog.Error("error occured while running HTTPS server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		c.Next()
	})
	r.NoRoute(controller.Handler)
	for _, server := range config.DefaultConfig.ServerConfigs {
		if !server.Enable {
			slog.Warn("server is disabled", slog.String("host", server.Host), slog.Int("port", server.Port))
			continue
		}
		addr := fmt.Sprintf("%s:%d", server.Host, server.Port)
		if server.CertificatePath != "" && server.PrivateKeyPath != "" {
			go startHTTPSServer(r, addr, server.CertificatePath, server.PrivateKeyPath)
		} else {
			go startHTTPServer(r, addr)
		}
	}
	select {}
}
