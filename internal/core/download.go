package core

import (
	"io"
	"net/http"
	"time"
    "os"
    "fmt"
	"strings"
	"crypto/tls"
    "golang.org/x/net/proxy"
	"github.com/go-shiori/shiori/internal/model"
)

var httpClient = &http.Client{Timeout: time.Minute}

func GetHttpClient() *http.Client {
    s5proxy := os.Getenv("SOCKS5_PROXY")
    if s5proxy == "" {
        return httpClient
    }
    dialer, err := proxy.SOCKS5("tcp", s5proxy, nil, proxy.Direct)
    if err != nil {
        fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		return nil
    }
    client := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
            Dial: dialer.Dial,
        },
    }
    return client
}

// DownloadBookmark downloads bookmarked page from specified URL.
// Return response body, make sure to close it later.
func DownloadBookmark(url string) (io.ReadCloser, string, error) {
	// Prepare download request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	// Send download request
	req.Header.Set("User-Agent", userAgent)
    // resp, err := httpClient.Do(req)
	resp, err := GetHttpClient().Do(req)
	if err != nil {
		return nil, "", err
	}

	// Get content type
	contentType := resp.Header.Get("Content-Type")

	return resp.Body, contentType, nil
}

func PreProcessBookmark(book model.Bookmark) (model.Bookmark, bool) {
    book.CreateArchive = true
    if book.Public == 1 {
        port := os.Getenv("SOCKS5_PROXY_PORT")
        if port == "" {
            port = "1881"
        }
        os.Setenv("SOCKS5_PROXY", "127.0.0.1:" + port)
    } else {
        os.Setenv("SOCKS5_PROXY", "")
    }
    is_local_page := strings.Contains(book.URL, "://theta")
    if is_local_page {
        book.URL = strings.Replace(book.URL, "://theta", "://127.0.0.1", 1)
    } else {
        if !strings.HasPrefix(book.URL, "http") {
            book.URL = "http://127.0.0.1/" + book.URL
            if !strings.HasSuffix(book.URL, ".html") {
                book.URL = book.URL + ".html"
            }
        }
    }
    fmt.Println("add bookmark url:", book.URL, "proxy:", os.Getenv("SOCKS5_PROXY"))
    return book, is_local_page
}
