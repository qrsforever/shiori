package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	fp "path/filepath"
	"strconv"
	"strings"

	"github.com/go-shiori/shiori/internal/core"
	"github.com/go-shiori/shiori/internal/model"
	"github.com/julienschmidt/httprouter"
)

// apiInsertViaExtension is handler for POST /api/bookmarks/ext
func (h *handler) apiInsertViaExtension(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Make sure session still valid
	err := h.validateSession(r)
	checkError(err)

	// Decode request
	request := model.Bookmark{}
	err = json.NewDecoder(r.Body).Decode(&request)
	checkError(err)

	// Clean up bookmark URL
	request.URL, err = core.RemoveUTMParams(request.URL)
	if err != nil {
		panic(fmt.Errorf("failed to clean URL: %v", err))
	}

    // QRS: for local http (just for me)
    is_local_page := strings.Contains(request.URL, "://theta")
    if is_local_page {
        request.URL = strings.Replace(request.URL, "://theta", "://127.0.0.1", 1)
    }

    fmt.Println(request.URL)

	// Check if bookmark already exists.
	book, exist := h.DB.GetBookmark(0, request.URL)

	// If it already exists, we need to set ID and tags.
	if exist {
		book.HTML = request.HTML

		mapOldTags := map[string]model.Tag{}
		for _, oldTag := range book.Tags {
			mapOldTags[oldTag.Name] = oldTag
		}

		for _, newTag := range request.Tags {
			if _, tagExist := mapOldTags[newTag.Name]; !tagExist {
				book.Tags = append(book.Tags, newTag)
			}
		}
	} else {
		book = request
		book.ID, err = h.DB.CreateNewID("bookmark")
		if err != nil {
			panic(fmt.Errorf("failed to create ID: %v", err))
		}
	}

	// Since we are using extension, the extension might send the HTML content
	// so no need to download it again here. However, if it's empty, it might be not HTML file
	// so we download it here.
	var contentType string
	var contentBuffer io.Reader

	if book.HTML == "" {
		contentBuffer, contentType, _ = core.DownloadBookmark(book.URL)
	} else {
		contentType = "text/html; charset=UTF-8"
		contentBuffer = bytes.NewBufferString(book.HTML)
	}

    fmt.Println(book.ID)

	// At this point the web page already downloaded.
	// Time to process it.
	if contentBuffer != nil {
        book.CreateArchive = true
		request := core.ProcessRequest{
			DataDir:     h.DataDir,
			Bookmark:    book,
			Content:     contentBuffer,
			ContentType: contentType,
		}
        if is_local_page {
            request.LogArchival = true
        }

		var isFatalErr bool
		book, isFatalErr, err = core.ProcessBookmark(request)

		if tmp, ok := contentBuffer.(io.ReadCloser); ok {
			tmp.Close()
		}

		if err != nil && isFatalErr {
			panic(fmt.Errorf("failed to process bookmark: %v", err))
		}
        fmt.Println("add book ok")
	}

	// Save bookmark to database
	results, err := h.DB.SaveBookmarks(book)
	if err != nil || len(results) == 0 {
        fmt.Println("failed to save bookmark")
		panic(fmt.Errorf("failed to save bookmark: %v", err))
	}
	book = results[0]
    fmt.Println("save bookmark ok")

	// Return the new bookmark
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&book)
	checkError(err)
}

// apiDeleteViaExtension is handler for DELETE /api/bookmark/ext
func (h *handler) apiDeleteViaExtension(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Make sure session still valid
	err := h.validateSession(r)
	checkError(err)

	// Decode request
	request := model.Bookmark{}
	err = json.NewDecoder(r.Body).Decode(&request)
	checkError(err)

	// Check if bookmark already exists.
	book, exist := h.DB.GetBookmark(0, request.URL)
	if exist {
		// Delete bookmarks
		err = h.DB.DeleteBookmarks(book.ID)
		checkError(err)

		// Delete thumbnail image and archives from local disk
		strID := strconv.Itoa(book.ID)
		imgPath := fp.Join(h.DataDir, "thumb", strID)
		archivePath := fp.Join(h.DataDir, "archive", strID)

		os.Remove(imgPath)
		os.Remove(archivePath)
	}

	fmt.Fprint(w, 1)
}
