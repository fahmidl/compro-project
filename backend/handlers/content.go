package handlers

import (
	"errors"
	"net/http"

	"compro-backend/db"
	"compro-backend/models"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	DB *db.DB
}

func NewContentHandler(database *db.DB) *ContentHandler {
	return &ContentHandler{DB: database}
}

func (h *ContentHandler) GetContent(c *gin.Context) {
	content, err := h.DB.GetContent(c.Request.Context())
	if errors.Is(err, db.ErrNotFound) {
		// seed default content
		def := models.DefaultSiteContent()
		if err := h.DB.PutContent(c.Request.Context(), &def); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to seed default content"})
			return
		}
		c.JSON(http.StatusOK, def)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch content"})
		return
	}
	c.JSON(http.StatusOK, content)
}

func (h *ContentHandler) UpdateContent(c *gin.Context) {
	var updated models.SiteContent
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if existing content exists; if so, preserve the ID
	existing, err := h.DB.GetContent(c.Request.Context())
	if err == nil && existing != nil {
		updated.ID = existing.ID
	} else if !errors.Is(err, db.ErrNotFound) && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if err := h.DB.PutContent(c.Request.Context(), &updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update content"})
		return
	}
	c.JSON(http.StatusOK, updated)
}
