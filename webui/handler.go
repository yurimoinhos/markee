package webui

import (
	"embed"
	"errors"
	"io/fs"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/aggi-tech/aggipay/platform/config"
	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var embeddedFiles embed.FS

var distFS fs.FS

func init() {
	sub, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		panic(err)
	}
	distFS = sub
}

func RegisterRoutes(router *gin.Engine, cfg *config.Config) {
	h := newHandler(cfg)
	router.GET("/app", h)
	router.GET("/app/*filepath", h)
	router.POST("/app/*filepath", h)
	router.PUT("/app/*filepath", h)
	router.PATCH("/app/*filepath", h)
	router.DELETE("/app/*filepath", h)
	router.HEAD("/app/*filepath", h)
	router.OPTIONS("/app/*filepath", h)
}

func newHandler(cfg *config.Config) gin.HandlerFunc {
	nextURL := strings.TrimSpace(cfg.Frontend.NextInternalURL)
	proxy := buildProxy(nextURL)

	return func(c *gin.Context) {
		if proxy != nil {
			proxy.ServeHTTP(c.Writer, c.Request)
			if c.Writer.Written() {
				return
			}
		}

		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.JSON(http.StatusBadGateway, gin.H{"error": "frontend indisponível"})
			return
		}

		requestPath := strings.TrimPrefix(c.Param("filepath"), "/")
		if requestPath == "" {
			serveIndex(c)
			return
		}

		if hasAssetExtension(requestPath) && serveStaticFile(c, requestPath) {
			return
		}
		serveIndex(c)
	}
}

func buildProxy(target string) *httputil.ReverseProxy {
	if target == "" {
		return nil
	}
	u, err := url.Parse(target)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = u.Host
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		_ = err
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`{"error":"frontend indisponível"}`))
			return
		}

		fallbackPath := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/app"), "/")
		if fallbackPath == "" {
			fallbackPath = "index.html"
		}
		if hasAssetExtension(fallbackPath) {
			if b, ct, readErr := readStaticFile(fallbackPath); readErr == nil {
				w.Header().Set("Content-Type", ct)
				w.Header().Set("Cache-Control", "no-cache")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(b)
				return
			}
		}
		if b, ct, readErr := readStaticFile("index.html"); readErr == nil {
			w.Header().Set("Content-Type", ct)
			w.Header().Set("Cache-Control", "no-cache")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(b)
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("frontend indisponível"))
	}
	return proxy
}

func hasAssetExtension(p string) bool {
	ext := strings.ToLower(path.Ext(p))
	return ext != ""
}

func serveIndex(c *gin.Context) {
	if !serveStaticFile(c, "index.html") {
		c.String(http.StatusBadGateway, "frontend indisponível")
	}
}

func serveStaticFile(c *gin.Context, p string) bool {
	b, contentType, err := readStaticFile(p)
	if err != nil {
		return false
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, contentType, b)
	return true
}

func readStaticFile(p string) ([]byte, string, error) {
	p = strings.TrimPrefix(path.Clean("/"+p), "/")
	b, err := fs.ReadFile(distFS, p)
	if err != nil {
		return nil, "", err
	}
	ext := strings.ToLower(path.Ext(p))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		switch ext {
		case ".html":
			contentType = "text/html; charset=utf-8"
		case ".js":
			contentType = "application/javascript"
		case ".css":
			contentType = "text/css; charset=utf-8"
		default:
			contentType = "application/octet-stream"
		}
	}
	if len(b) == 0 {
		return nil, "", errors.New("empty file")
	}
	return b, contentType, nil
}
