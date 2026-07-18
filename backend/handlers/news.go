package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"compro-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsHandler struct {
	Collection *mongo.Collection
}

func NewNewsHandler(db *mongo.Database) *NewsHandler {
	return &NewsHandler{Collection: db.Collection("news_posts")}
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric characters (keep hyphens)
	var result []byte
	for i := 0; i < len(slug); i++ {
		c := slug[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			result = append(result, c)
		}
	}
	return string(result)
}

type createPostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Summary string `json:"summary"`
	Image   string `json:"image"`
}

type updatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Summary string `json:"summary"`
	Image   string `json:"image"`
}

func (h *NewsHandler) GetPosts(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "publishedAt", Value: -1}})
	cursor, err := h.Collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}
	defer cursor.Close(context.Background())

	var posts []models.NewsPost
	if err := cursor.All(context.Background(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode posts"})
		return
	}

	if posts == nil {
		posts = []models.NewsPost{}
	}

	c.JSON(http.StatusOK, posts)
}

func (h *NewsHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	// Try to find by ObjectID first, then by slug
	objectID, err := primitive.ObjectIDFromHex(id)
	var post models.NewsPost
	if err == nil {
		err = h.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&post)
	} else {
		err = h.Collection.FindOne(context.Background(), bson.M{"slug": id}).Decode(&post)
	}

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *NewsHandler) CreatePost(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")
	now := primitive.NewDateTimeFromTime(time.Now())

	post := models.NewsPost{
		Title:       req.Title,
		Slug:        generateSlug(req.Title),
		Content:     req.Content,
		Summary:     req.Summary,
		Image:       req.Image,
		Author:      username.(string),
		PublishedAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result, err := h.Collection.InsertOne(context.Background(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, post)
}

func (h *NewsHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build update fields
	update := bson.M{"updatedAt": primitive.NewDateTimeFromTime(time.Now())}
	if req.Title != "" {
		update["title"] = req.Title
		update["slug"] = generateSlug(req.Title)
	}
	if req.Content != "" {
		update["content"] = req.Content
	}
	if req.Summary != "" {
		update["summary"] = req.Summary
	}
	if req.Image != "" {
		update["image"] = req.Image
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.NewsPost
	err = h.Collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": update},
		opts,
	).Decode(&updated)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *NewsHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	result, err := h.Collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}
