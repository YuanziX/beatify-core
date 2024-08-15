package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yuanzix/beatify-core/models"
	"github.com/yuanzix/beatify-core/utils"
)

func (s *APIServer) UploadMusicHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return http.StatusBadRequest, err
	}

	file, handler, err := r.FormFile("music_file")
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer file.Close()

	if _, err := os.Stat("./music"); os.IsNotExist(err) {
		if err := os.Mkdir("./music", os.ModePerm); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	title := r.FormValue("title")
	artist := r.FormValue("artist")
	album := r.FormValue("album")
	yearStr := r.FormValue("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return http.StatusBadRequest, err
	}

	ext := filepath.Ext(handler.Filename)
	newFileName := fmt.Sprintf("%s-%s-%s-(%d)%s", sanitizeFileName(artist), sanitizeFileName(title), sanitizeFileName(album), year, ext)
	filePath := filepath.Join("./music", newFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return http.StatusInternalServerError, err
	}

	newMusic := &models.Music{
		Title:    title,
		Artist:   artist,
		Album:    album,
		Location: filePath,
		Year:     int32(year),
	}

	createdMusic, err := s.store.CreateMusic(newMusic)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdMusic)
}

func (s *APIServer) handleGetMusicList(w http.ResponseWriter, r *http.Request) (int, error) {
	pageNo, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || pageNo <= 0 {
		return http.StatusBadRequest, fmt.Errorf("invalid page number: %v", err)
	}

	musicList, err := s.store.GetMusicList(pageNo)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(*musicList) == 0 {
		return http.StatusRequestedRangeNotSatisfiable, errors.New("reached end of content")
	}

	return utils.WriteJSON(w, http.StatusOK, musicList)
}

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

func sanitizeFileName(name string) string {
	return strings.NewReplacer(" ", "_", "/", "-", "\\", "-", ":", "-", "*", "-", "?", "-", "\"", "-", "<", "-", ">", "-", "|", "-").Replace(name)
}
