package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/yuanzix/userAuth/utils"
)

func (s *APIServer) handleStreamAudio(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	music, err := s.store.GetMusicByID(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	file, err := os.Open(music.Location)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	fileSize := stat.Size()

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		http.ServeContent(w, r, music.Location, stat.ModTime(), file)
		return http.StatusOK, nil
	}

	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	rangeParts := strings.Split(rangeHeader, "-")

	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		return http.StatusBadRequest, err
	}

	var end int64
	if len(rangeParts) > 1 && rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			return http.StatusBadRequest, err
		}
	} else {
		end = fileSize - 1
	}

	if start < 0 || end >= fileSize || start > end {
		return http.StatusRequestedRangeNotSatisfiable, nil
	}

	w.Header().Set("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(fileSize, 10))
	w.WriteHeader(http.StatusPartialContent)

	file.Seek(start, 0)
	buf := make([]byte, end-start+1)
	file.Read(buf)
	w.Write(buf)

	return utils.WriteJSON(w, http.StatusPartialContent, map[string]string{"stream": "successful"})
}
