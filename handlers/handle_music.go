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

	musicFile, musicFileHandler, err := r.FormFile("music_file")
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer musicFile.Close()

	if _, err := os.Stat("./music"); os.IsNotExist(err) {
		if err := os.Mkdir("./music", os.ModePerm); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	thumbnailFile, thumbnailFileHandler, err := r.FormFile("thumbnail_file")
	if err != nil {
		return http.StatusBadRequest, err
	}
	defer thumbnailFile.Close()

	if _, err := os.Stat("./music"); os.IsNotExist(err) {
		if err := os.Mkdir("./music", os.ModePerm); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	if _, err := os.Stat("./music/thumbnails"); os.IsNotExist(err) {
		if err := os.Mkdir("./music/thumbnails", os.ModePerm); err != nil {
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

	musicFileExt := filepath.Ext(musicFileHandler.Filename)
	thumbnailFileExt := filepath.Ext(thumbnailFileHandler.Filename)
	baseFileName := fmt.Sprintf("%s-%s-%s-(%d)", sanitizeFileName(artist), sanitizeFileName(title), sanitizeFileName(album), year)
	newMusicFileName := fmt.Sprintf("%s%s", baseFileName, musicFileExt)
	newThumbnailFileName := fmt.Sprintf("%s-thumnail%s", baseFileName, thumbnailFileExt)
	musicFilePath := filepath.Join("./music", newMusicFileName)
	thumbnailFilePath := filepath.Join("./music/thumbnails", newThumbnailFileName)

	musicDst, err := os.Create(musicFilePath)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	thumbnailDst, err := os.Create(thumbnailFilePath)
	if err != nil {
		return http.StatusInternalServerError, err

	}

	defer musicDst.Close()
	defer thumbnailDst.Close()

	if _, err := io.Copy(thumbnailDst, thumbnailFile); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := io.Copy(musicDst, musicFile); err != nil {
		return http.StatusInternalServerError, err
	}

	newMusic := &models.Music{
		Title:             title,
		Artist:            artist,
		Album:             album,
		Location:          musicFilePath,
		ThumbnailLocation: thumbnailFilePath,
		Year:              int32(year),
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

	http.ServeFile(w, r, music.Location)

	return http.StatusOK, nil
}

func (s *APIServer) handleGetThumbnail(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	music, err := s.store.GetMusicByID(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.ServeFile(w, r, music.ThumbnailLocation)

	return http.StatusOK, nil
}

func sanitizeFileName(name string) string {
	return strings.NewReplacer(" ", "_", "/", "-", "\\", "-", ":", "-", "*", "-", "?", "-", "\"", "-", "<", "-", ">", "-", "|", "-").Replace(name)
}
