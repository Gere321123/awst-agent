package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type AgentProxy struct {
	BackendURL string
	Proxy      *httputil.ReverseProxy
	Token      string
}

func NewAgentProxy(backendURL, token string) (*AgentProxy, error) {
	remote, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Backend server is not available (502). Start the backend on port 3000!"))
	}

	return &AgentProxy{
		BackendURL: backendURL,
		Proxy:      proxy,
		Token:      token,
	}, nil
}

func (p *AgentProxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Redirect(w, r, "http://localhost:5000/buy", http.StatusTemporaryRedirect)
			return
		}

		// Token validation
		if token != "Bearer "+p.Token && token != p.Token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		p.Proxy.ServeHTTP(w, r)
	}
}
