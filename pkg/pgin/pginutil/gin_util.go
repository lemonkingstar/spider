package pginutil

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Query get query param
// ?name=xx&age=18
func Query(c *gin.Context, key string) string { return c.Query(key) }

func QueryInt(c *gin.Context, key string) int {
	if v, err := strconv.Atoi(c.Query(key)); err == nil {
		return v
	}
	return 0
}

func QueryIntDefault(c *gin.Context, key string, d int) int {
	q := c.Query(key)
	if q == "" {
		return d
	}
	if v, err := strconv.Atoi(q); err == nil {
		return v
	}
	return d
}

func QueryBool(c *gin.Context, key string) bool {
	s := strings.ToLower(strings.TrimSpace(c.Query(key)))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
}

func QueryBoolDefault(c *gin.Context, key string, d bool) bool {
	q := c.Query(key)
	if q == "" {
		return d
	}
	s := strings.ToLower(strings.TrimSpace(q))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
}

// Param get url param
// /user/:id
func Param(c *gin.Context, key string) string { return c.Param(key) }

func ParamInt(c *gin.Context, key string) int {
	if v, err := strconv.Atoi(c.Param(key)); err == nil {
		return v
	}
	return 0
}

func ParamIntDefault(c *gin.Context, key string, d int) int {
	q := c.Param(key)
	if q == "" {
		return d
	}
	if v, err := strconv.Atoi(q); err == nil {
		return v
	}
	return d
}

func ParamBool(c *gin.Context, key string) bool {
	s := strings.ToLower(strings.TrimSpace(c.Param(key)))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
}

func ParamBoolDefault(c *gin.Context, key string, d bool) bool {
	q := c.Param(key)
	if q == "" {
		return d
	}
	s := strings.ToLower(strings.TrimSpace(q))
	return !(s == "" || s == "0" || s == "no" || s == "false" || s == "none")
}

// Form get form param
// urlencoded form or multipart form
func Form(c *gin.Context, key string) string { return c.PostForm(key) }

func FormInt(c *gin.Context, key string) int {
	if v, err := strconv.Atoi(c.PostForm(key)); err == nil {
		return v
	}
	return 0
}

// FormFile get specified file
// Save to local path: c.SaveUploadedFile(file, path)
// curl -X POST http://host/api -F "file=@/path/test.txt" -H "Content-Type: multipart/form-data"
func FormFile(c *gin.Context, field string) (*multipart.FileHeader, error) {
	return c.FormFile(field)
}

// FormFiles get specified files
// curl -X POST http://host/api \
//  -F "file=@/path/test1.txt" \
//  -F "file=@/path/test2.txt" \
//  -H "Content-Type: multipart/form-data"
func FormFiles(c *gin.Context, field string) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil { return nil, err }
	return form.File[field], nil
}
