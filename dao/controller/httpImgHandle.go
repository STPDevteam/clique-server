package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"stp_dao_v2/models"
	"time"
)

// @Summary update img
// @Tags img
// @version 0.0.1
// @description update img
// @Produce json
// @Param file formData file true "file" "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/img/upload [post]
func (svc *Service) httpUploadImg(c *gin.Context) {
	var (
		w = c.Writer
		r = c.Request
	)

	r.Body = http.MaxBytesReader(w, r.Body, svc.appConfig.MaxUpdateImgSize)
	if err := r.ParseMultipartForm(svc.appConfig.MaxUpdateImgSize); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "The uploaded file is too big. Please choose an file that's less than 1MB in size",
		})
		return
	}

	// The argument to FormFile must match the name attribute of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "FormFile key must 'file'",
		})
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "The provided file format is not allowed. Please upload a JPEG or PNG image",
		})
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll("./static", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	// Create a new file in the uploads directory
	path := fmt.Sprintf("/static/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	paths := fmt.Sprintf(".%s", path)
	dst, err := os.Create(paths)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var url string
	if r.ProtoMajor == 2 {
		url = fmt.Sprintf("https://%s%s", r.Host, path)
	} else {
		url = fmt.Sprintf("http://%s%s", r.Host, path)
	}
	//if strings.Contains(r.Proto, "HTTPS") {
	//
	//} else {
	//
	//}

	c.JSON(http.StatusOK, models.Response{
		Code: http.StatusOK,
		Data: models.ResUploadImgPath{
			Path: url,
		},
		Message: "ok",
	})
}
