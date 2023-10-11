package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

type ReverseProxy struct {
	pathprefix string
	proxyMap   map[string]*httputil.ReverseProxy
	mutx       *sync.RWMutex
}

var (
	defaultRP *ReverseProxy
	once      sync.Once
)

func GetDefaultReverseProxy() *ReverseProxy {
	if defaultRP == nil {
		once.Do(func() {
			defaultRP = NewReverseProxy()
		})
	}

	return defaultRP
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{
		proxyMap: map[string]*httputil.ReverseProxy{},
		mutx:     &sync.RWMutex{},
	}
}

func (p *ReverseProxy) SetPathPrefix(prefix string) {
	p.pathprefix = prefix
}

func (p *ReverseProxy) RegisterProxy(name, upstream string) error {
	targetURL, err := url.Parse(upstream)
	if err != nil {
		return err
	}

	pp := httputil.NewSingleHostReverseProxy(targetURL)
	pp.Director = func(req *http.Request) {
		req.Host = targetURL.Host
		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme

		suffix := ""
		if strings.HasSuffix(req.URL.Path, "/") {
			suffix = "/"
		}

		req.URL.Path = path.Join(targetURL.Path, req.URL.Path) + suffix

		klog.Infof("req.Host=%s req.Path=%s", req.Host, req.URL.Path)

		targetQuery := targetURL.RawQuery
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}

	p.mutx.Lock()
	p.proxyMap[name] = pp
	p.mutx.Unlock()

	return nil
}

func (p *ReverseProxy) UnregsiterProxy(name string) {
	p.mutx.Lock()
	delete(p.proxyMap, name)
	p.mutx.Unlock()
}

func (p *ReverseProxy) Proxy(c *gin.Context) {
	urlpath := strings.TrimPrefix(c.Param("wildcard"), "/")

	klog.Infof("path=%s", urlpath)

	name := ""

	idx := strings.Index(urlpath, "/")
	if idx == -1 {
		name = urlpath
		c.Request.URL.Path = ""
	} else {
		name = urlpath[:idx]
		c.Request.URL.Path = urlpath[idx:]
	}

	pp := p.GetProxyHandler(name)
	if pp != nil {
		pp(c)
	} else {
		c.String(400, "invalid path=%s", c.Request.URL.Path)
	}
}

func (p *ReverseProxy) GetProxyHandler(name string) func(c *gin.Context) {
	p.mutx.RLock()
	pp := p.proxyMap[name]
	p.mutx.RUnlock()

	if pp == nil {
		return nil
	}

	return func(c *gin.Context) {
		klog.Infof("proxy=%s", c.Request.URL.Path)
		pp.ServeHTTP(c.Writer, c.Request)
	}
}
