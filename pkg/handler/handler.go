package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type Handler struct {
	Files    []File
	Messages binding.ExternalStringList
}

func NewHandler() *Handler {
	h := new(Handler)
	h.Messages = binding.BindStringList(&[]string{})
	return h
}

func (h *Handler) Browse(uri fyne.ListableURI, _ error) {
	defer func() {
		if err := recover(); err != nil {
			if strings.Contains(err.(error).Error(), "runtime error") {
				return
			} else {
				panic(err)
			}
		}
	}()
	fileList, err := uri.List()
	if err != nil {
		panic(err)
	}
	for _, file := range fileList {
		fileInfo, err := os.Stat(file.Path())
		if err != nil {
			panic(err)
		}
		if !strings.Contains(file.Name(), "watermark") && !fileInfo.IsDir() && strings.HasPrefix(file.MimeType(), "image") {
			f := File{
				FullPath: file.Path(),
				Dir:      filepath.Dir(file.Path()),
				Base:     file.Name(),
				Ext:      strings.TrimLeft(file.Extension(), "."),
			}
			h.Files = append(h.Files, f)
			h.Log(fmt.Sprintf("Found: %s in: %s", f.Base, f.Dir))
		}
	}
}

func (h *Handler) Log(msg string) {
	now := time.Now()
	h.Messages.Append(now.Format("2006-01-02 15:04:05") + ": " + msg)
}

func (h *Handler) Clear() {
	h.Files = nil
}
