package store

import (
	"github.com/dportaluppi/customer-profiles-api/pkg"
	"github.com/dportaluppi/customer-profiles-api/pkg/store"
	"github.com/gin-gonic/gin"
	gojsonlogicmongodb "github.com/kubeesio/go-jsonlogic-mongodb"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

// service define business logic for store.
type service struct {
	store.Upserter
	store.Deleter
	store.Getter
}

// Handler rest api for store.
type Handler struct {
	service *service
}

// NewHandler creates a new handler for store.
func NewHandler(upserter store.Upserter, deleter store.Deleter, getter store.Getter) *Handler {
	s := &service{
		Upserter: upserter,
		Deleter:  deleter,
		Getter:   getter,
	}
	return &Handler{service: s}
}

// Create manages the creation of a new store.
func (h *Handler) Create(c *gin.Context) {
	var s store.Store
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	createdUser, err := h.service.Create(ctx, &s)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdUser)
}

// Update manages the update of an existing store.
func (h *Handler) Update(c *gin.Context) {
	var store store.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
		return
	}

	ctx := c.Request.Context()
	updatedStore, err := h.service.Update(ctx, id, &store)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedStore)
}

// Delete manages the deletion of a store.
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

	c.JSON(http.StatusOK, gin.H{"message": "Store deleted"})
}

// GetByID manages fetching a store by its ID.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ID"})
		return
	}

	ctx := c.Request.Context()
	store, err := h.service.GetByID(ctx, id)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error more specifically if needed.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, store)
}

// GetAll manages fetching all stores with pagination.
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
	stores, totalItems, err := h.service.GetAll(ctx, currentPage, perPage)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination := pkg.NewPagination(currentPage, perPage, totalItems)

	response := gin.H{
		"stores":     stores,
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
