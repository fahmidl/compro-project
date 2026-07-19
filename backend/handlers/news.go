package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"compro-backend/db"
	"compro-backend/models"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	DB *db.DB
}

func NewNewsHandler(database *db.DB) *NewsHandler {
	return &NewsHandler{DB: database}
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
	posts, err := h.DB.ListNews(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}
	if posts == nil {
		posts = []models.NewsPost{}
	}
	c.JSON(http.StatusOK, posts)
}

func (h *NewsHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	// Try to find by ID first, then fall back to slug lookup
	post, err := h.DB.GetNewsByID(c.Request.Context(), id)
	if errors.Is(err, db.ErrNotFound) {
		// Not found by ID — try by slug
		post, err = h.DB.GetNewsBySlug(c.Request.Context(), id)
	}
	if errors.Is(err, db.ErrNotFound) {
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
	now := time.Now().UnixMilli()

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

	if err := h.DB.CreateNews(c.Request.Context(), &post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *NewsHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.DB.GetNewsByID(c.Request.Context(), id)
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
		return
	}

	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Apply updates
	if req.Title != "" {
		existing.Title = req.Title
		existing.Slug = generateSlug(req.Title)
	}
	if req.Content != "" {
		existing.Content = req.Content
	}
	if req.Summary != "" {
		existing.Summary = req.Summary
	}
	if req.Image != "" {
		existing.Image = req.Image
	}
	existing.UpdatedAt = time.Now().UnixMilli()

	if err := h.DB.UpdateNews(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *NewsHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	// Verify the post exists before deleting
	_, err := h.DB.GetNewsByID(c.Request.Context(), id)
	if errors.Is(err, db.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
		return
	}

	if err := h.DB.DeleteNews(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}
