package handlers

import (
	"context"
	"net/http"

	"compro-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContentHandler struct {
	Collection *mongo.Collection
}

func NewContentHandler(db *mongo.Database) *ContentHandler {
	return &ContentHandler{Collection: db.Collection("site_content")}
}

func (h *ContentHandler) GetContent(c *gin.Context) {
	var content models.SiteContent
	err := h.Collection.FindOne(context.Background(), bson.M{}).Decode(&content)
	if err == mongo.ErrNoDocuments {
		// seed default content
		def := models.DefaultSiteContent()
		res, _ := h.Collection.InsertOne(context.Background(), def)
		def.ID = res.InsertedID.(primitive.ObjectID)
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

	// Find existing document
	var existing models.SiteContent
	err := h.Collection.FindOne(context.Background(), bson.M{}).Decode(&existing)
	if err == mongo.ErrNoDocuments {
		// Insert new
		res, err := h.Collection.InsertOne(context.Background(), updated)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create content"})
			return
		}
		updated.ID = res.InsertedID.(primitive.ObjectID)
		c.JSON(http.StatusOK, updated)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// Update existing document (preserve ID)
	updated.ID = existing.ID
	_, err = h.Collection.ReplaceOne(context.Background(), bson.M{"_id": existing.ID}, updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update content"})
		return
	}
	c.JSON(http.StatusOK, updated)
}
