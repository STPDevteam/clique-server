package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"stp_dao_v2/models"
	"strings"
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
func HttpUploadImg(c *gin.Context) {
	var (
		w = c.Writer
		r = c.Request
	)

	r.Body = http.MaxBytesReader(w, r.Body, viper.GetInt64("app.max_upload_img_size"))
	if err := r.ParseMultipartForm(viper.GetInt64("app.max_upload_img_size")); err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    400,
			Message: "The uploaded file is too big. Please choose an file that's less than 1MB in size",
		})
		return
	}

	// The argument to FormFile must match the name attribute of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    400,
			Message: "FormFile key must 'file'",
		})
		return
	}
	defer file.Close()

	suffixArr := strings.Split(fileHeader.Filename, ".")
	suffix := suffixArr[len(suffixArr)-1]
	if strings.ToLower(suffix) != "jpeg" && strings.ToLower(suffix) != "jpg" && strings.ToLower(suffix) != "png" { // && strings.ToLower(suffix) != "svg"
		c.JSON(http.StatusOK, models.Response{
			Code:    400,
			Message: "The provided file format is not allowed. Please upload a JPEG or PNG image",
		})
		return
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" { //&& fileHeader.Header.Get("Content-Type") != "image/svg+xml"
		c.JSON(http.StatusOK, models.Response{
			Code:    400,
			Message: "The provided file format is not allowed. Please upload a JPEG or PNG image",
		})
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	// Create the uploads folder if it doesn't already exist
	err = os.MkdirAll("./static", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	// Create a new file in the uploads directory
	path := fmt.Sprintf("/static/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	paths := fmt.Sprintf(".%s", path)
	dst, err := os.Create(paths)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var url string
	if !strings.Contains(c.Request.Host, "localhost") {
		url = fmt.Sprintf("https://%s%s", r.Host, path)
	} else {
		url = fmt.Sprintf("http://%s%s", r.Host, path)
	}

	c.JSON(http.StatusOK, models.Response{
		Code: http.StatusOK,
		Data: models.ResUploadImgPath{
			Path: url,
		},
		Message: "ok",
	})
}
