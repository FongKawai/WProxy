package forward

import (
	"net/http"
	"testing"
)

func TestHandleHost(t *testing.T) {
	tests := []struct {
		name     string
		scheme   string
		host     string
		expected string
	}{
		{
			name:     "HTTPS without port",
			scheme:   "https",
			host:     "example.com",
			expected: "example.com:443",
		},
		{
			name:     "HTTP without port",
			scheme:   "http",
			host:     "example.com",
			expected: "example.com:80",
		},
		{
			name:     "HTTPS with port",
			scheme:   "https",
			host:     "example.com:8443",
			expected: "example.com:8443",
		},
		{
			name:     "HTTP with port",
			scheme:   "http",
			host:     "example.com:8080",
			expected: "example.com:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handleHost(tt.scheme, tt.host)
			if result != tt.expected {
				t.Errorf("handleHost(%q, %q) = %q, want %q", tt.scheme, tt.host, result, tt.expected)
			}
		})
	}
}

func TestHandleForwardLoopDetection(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// First request should succeed
	err = HandleForward(req)
	if err != nil {
		t.Errorf("First HandleForward failed: %v", err)
	}

	// Check that loop detection header was set
	if req.Header.Get("x-proxy-loop") != "1" {
		t.Error("Loop detection header not set")
	}

	// Second request with loop header should fail
	req2, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set("x-proxy-loop", "1")

	err = HandleForward(req2)
	if err != LoopDetected {
		t.Errorf("Expected LoopDetected error, got: %v", err)
	}
}

func TestHandleForwardWithCustomHeaders(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		proxyHost      string
		proxyScheme    string
		expectedHost   string
		expectedURLHost string
		expectedScheme string
	}{
		{
			name:            "Custom HTTPS host",
			url:             "http://original.com",
			proxyHost:       "target.com",
			proxyScheme:     "https",
			expectedHost:    "target.com:443",
			expectedURLHost: "target.com",
			expectedScheme:  "https",
		},
		{
			name:            "Custom HTTP host",
			url:             "https://original.com",
			proxyHost:       "target.com",
			proxyScheme:     "http",
			expectedHost:    "target.com:80",
			expectedURLHost: "target.com",
			expectedScheme:  "http",
		},
		{
			name:            "Custom host with port",
			url:             "http://original.com",
			proxyHost:       "target.com:8080",
			proxyScheme:     "http",
			expectedHost:    "target.com:8080",
			expectedURLHost: "target.com:8080",
			expectedScheme:  "http",
		},
		{
			name:            "No custom headers",
			url:             "https://original.com",
			proxyHost:       "",
			proxyScheme:     "",
			expectedHost:    "original.com:443",
			expectedURLHost: "original.com:443",
			expectedScheme:  "https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.proxyHost != "" {
				req.Header.Set("x-proxy-host", tt.proxyHost)
			}
			if tt.proxyScheme != "" {
				req.Header.Set("x-proxy-scheme", tt.proxyScheme)
			}

			err = HandleForward(req)
			if err != nil {
				t.Errorf("HandleForward failed: %v", err)
			}

			if req.URL.Host != tt.expectedURLHost {
				t.Errorf("Expected URL host %q, got %q", tt.expectedURLHost, req.URL.Host)
			}

			if req.URL.Scheme != tt.expectedScheme {
				t.Errorf("Expected URL scheme %q, got %q", tt.expectedScheme, req.URL.Scheme)
			}

			if req.Host != tt.expectedHost {
				t.Errorf("Expected req.Host %q, got %q", tt.expectedHost, req.Host)
			}
		})
	}
}
