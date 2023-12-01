package user

import (
	"github.com/dportaluppi/customer-profiles-api/pkg"
	"github.com/dportaluppi/customer-profiles-api/pkg/user"
	"github.com/gin-gonic/gin"
	gojsonlogicmongodb "github.com/kubeesio/go-jsonlogic-mongodb"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

// service define business logic for user.
type service struct {
	user.Upserter
	user.Deleter
	user.Getter
}

// Handler rest api for user.
type Handler struct {
	service *service
}

// NewHandler creates a new handler for user.
func NewHandler(upserter user.Upserter, deleter user.Deleter, getter user.Getter) *Handler {
	s := &service{
		Upserter: upserter,
		Deleter:  deleter,
		Getter:   getter,
	}
	return &Handler{service: s}
}

// Create manages the creation of a new user.
func (h *Handler) Create(c *gin.Context) {
	var p user.User
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	createdUser, err := h.service.Create(ctx, &p)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdUser)
}

// Update manages the update of an existing user.
func (h *Handler) Update(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
		return
	}

	ctx := c.Request.Context()
	updatedProfile, err := h.service.Update(ctx, id, &user)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
}

// Delete manages the deletion of a user.
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

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// GetByID manages fetching a user by its ID.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ID"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.service.GetByID(ctx, id)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error more specifically if needed.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
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

	pagination := pkg.NewPagination(currentPage, perPage, totalItems)

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

	pagination := pkg.NewPagination(currentPage, perPage, totalItems)

	response := gin.H{
		"results":    results,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) QueryJsonLogic(c *gin.Context) {
	mongoQuery, err := gojsonlogicmongodb.Convert(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting JSON logic: " + err.Error()})
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
	results, totalItems, err := h.service.Pipeline(ctx, mongoQuery.Map(), currentPage, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination := pkg.NewPagination(currentPage, perPage, totalItems)

	response := gin.H{
		"results":    results,
		"pagination": pagination,
	}

	c.JSON(http.StatusOK, response)
}
