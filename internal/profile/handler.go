package profile

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
)

// service define business logic for profile.
type service struct {
	profile.Upserter
	profile.Deleter
	profile.Getter
}

// Handler rest api for profile.
type Handler struct {
	service *service
}

// NewHandler creates a new handler for profile.
func NewHandler(upserter profile.Upserter, deleter profile.Deleter, getter profile.Getter) *Handler {
	s := &service{
		Upserter: upserter,
		Deleter:  deleter,
		Getter:   getter,
	}
	return &Handler{service: s}
}

// Create manages the creation of a new profile.
func (h *Handler) Create(c *gin.Context) {
	var p profile.Profile
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	createdProfile, err := h.service.Create(ctx, &p)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdProfile)
}

// Update manages the update of an existing profile.
func (h *Handler) Update(c *gin.Context) {
	var profile profile.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
		return
	}

	ctx := c.Request.Context()
	updatedProfile, err := h.service.Update(ctx, id, &profile)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

// Delete manages the deletion of a profile.
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ID"})
		return
	}

	ctx := c.Request.Context()
	err := h.service.Delete(ctx, id)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error more specifically if needed.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}

// GetByID manages fetching a profile by its ID.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ID"})
		return
	}

	ctx := c.Request.Context()
	profile, err := h.service.GetByID(ctx, id)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error more specifically if needed.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetAll manages fetching all profiles with pagination.
func (h *Handler) GetAll(c *gin.Context) {
	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "50"))

	if currentPage < 1 {
		currentPage = 1
	}
	if perPage <= 0 {
		perPage = 50
	}

	ctx := c.Request.Context()
	profiles, totalItems, err := h.service.GetAll(ctx, currentPage, perPage)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination := profile.NewPagination(currentPage, perPage, totalItems)

	response := gin.H{
		"profiles":   profiles,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Query(c *gin.Context) {
	var query map[string]interface{}
	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "50"))

	if currentPage < 1 {
		currentPage = 1
	}
	if perPage <= 0 {
		perPage = 50
	}

	ctx := c.Request.Context()
	results, totalItems, err := h.service.Query(ctx, query, currentPage, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination := profile.NewPagination(currentPage, perPage, totalItems)

	response := gin.H{
		"results":    results,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetKeys(c *gin.Context) {
	ctx := c.Request.Context()
	keys, err := h.service.GetKeys(ctx)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keys)
}
