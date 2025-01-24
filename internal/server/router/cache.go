package router

import (
	"io"
	"os"

	"github.com/did-server/config"
)

func (h *handle) createCacheFile(filepath string, r io.Reader) error {
	path := config.CachePath + "/" + filepath

	f, err := os.Create(path)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		h.logger.Error(err)
		return err
	}

	return nil
}
