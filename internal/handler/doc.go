package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/paudarco/doc-storage/internal/errors"
	"github.com/paudarco/doc-storage/internal/handler/response"
	"github.com/paudarco/doc-storage/internal/service"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type DocHandler struct {
	doc service.Doc
	log *logrus.Logger
}

func NewDocHandler(doc service.Doc, log *logrus.Logger) *DocHandler {
	return &DocHandler{
		doc: doc,
		log: log,
	}
}

func (h *DocHandler) UploadDoc(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.NewErrorResponse(c, h.log, errors.ErrUnauthorized)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	var meta map[string]interface{}
	if metaValues := form.Value["meta"]; len(metaValues) > 0 {
		if err := json.Unmarshal([]byte(metaValues[0]), &meta); err != nil {
			response.NewErrorResponse(c, h.log, err)
			return
		}
	} else {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	var jsonData json.RawMessage
	if jsonValues := form.Value["json"]; len(jsonValues) > 0 {
		jsonData = json.RawMessage(jsonValues[0])
	}

	var fileData []byte
	var fileHeader *multipart.FileHeader
	if fileHeaders := form.File["file"]; len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
		file, err := fileHeader.Open()
		if err != nil {
			response.NewErrorResponse(c, h.log, err)
			return
		}
		defer file.Close()

		fileData, err = io.ReadAll(file)
		if err != nil {
			response.NewErrorResponse(c, h.log, err)
			return
		}
	}

	doc, err := h.doc.Create(c.Request.Context(), userID, meta, jsonData, fileData)
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	respData := gin.H{}
	if doc.IsFile {
		respData["file"] = doc.Name
	}
	if len(jsonData) > 0 {
		var jsonResp interface{}
		_ = json.Unmarshal(jsonData, &jsonResp)
		respData["json"] = jsonResp
	}

	c.JSON(http.StatusOK, gin.H{
		"data": respData,
	})
}

func (h *DocHandler) ListDocs(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	loginFilter, keyFilter, valueFilter, limit := getQueryParams(c)

	// Получаем список документов
	docs, err := h.doc.List(c.Request.Context(), userID, loginFilter, keyFilter, valueFilter, limit)
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	// Преобразуем в формат ответа
	docList := make([]gin.H, len(docs))
	for i, doc := range docs {
		docList[i] = gin.H{
			"id":      doc.ID,
			"name":    doc.Name,
			"file":    doc.IsFile,
			"public":  doc.Public,
			"created": doc.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if doc.Mime != "" {
			docList[i]["mime"] = doc.Mime
		}
		if len(doc.Grant) > 0 {
			docList[i]["grant"] = doc.Grant
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"docs": docList,
		},
	})
}

func (h *DocHandler) GetDoc(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	docID := c.Param("id")
	if docID == "" {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	doc, err := h.doc.GetByID(c.Request.Context(), userID, docID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			response.NewErrorResponse(c, h.log, err)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			response.NewErrorResponse(c, h.log, err)
			return
		}
		response.NewErrorResponse(c, h.log, err)
		return
	}

	if c.Request.Method == "HEAD" {
		if doc.IsFile {
			c.Header("Content-Type", doc.Mime)
			c.Header("Content-Length", fmt.Sprintf("%d", len(doc.FileData)))
		}
		c.Status(http.StatusOK)
		return
	}

	// Обработка GET запроса
	if doc.IsFile {
		// Отдаем файл
		c.Data(http.StatusOK, doc.Mime, doc.FileData)
	} else {
		var jsonData interface{}
		if doc.JSONData != nil {
			_ = json.Unmarshal(doc.JSONData.([]byte), &jsonData)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": jsonData,
		})
	}
}

func (h *DocHandler) DeleteDoc(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.NewErrorResponse(c, h.log, err)
		return
	}

	docID := c.Param("id")
	if docID == "" {
		response.NewErrorResponse(c, h.log, errors.ErrInvalidRequestBody)
		return
	}

	// Удаляем документ
	err = h.doc.Delete(c.Request.Context(), userID, docID)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			response.NewErrorResponse(c, h.log, err)
			return
		}
		if strings.Contains(err.Error(), "not found") {
			response.NewErrorResponse(c, h.log, err)
			return
		}
		response.NewErrorResponse(c, h.log, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": gin.H{
			docID: true,
		},
	})
}
