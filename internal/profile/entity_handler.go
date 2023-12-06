package profile

import (
	"github.com/dportaluppi/customer-profiles-api/pkg"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
	"github.com/gin-gonic/gin"
	gojsonlogicmongodb "github.com/kubeesio/go-jsonlogic-mongodb"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

// service define business logic for entity.
type service struct {
	profile.Saver
	profile.Deleter
	profile.Getter
}

// Handler rest api for entity.
type Handler struct {
	service *service
}

// NewHandler creates a new handler for entity.
func NewHandler(upserter profile.Saver, deleter profile.Deleter, getter profile.Getter) *Handler {
	s := &service{
		Saver:   upserter,
		Deleter: deleter,
		Getter:  getter,
	}
	return &Handler{service: s}
}

// Create manages the creation of a new entity.
func (h *Handler) Create(c *gin.Context) {
	var e profile.Entity
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountId := c.Param("accountId")

	ctx := c.Request.Context()
	createdUser, err := h.service.Create(ctx, accountId, &e)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdUser)
}

// Update manages the update of an existing entity.
func (h *Handler) Update(c *gin.Context) {
	var entity profile.Entity
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountId := c.Param("accountId")
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
		return
	}

	ctx := c.Request.Context()
	updatedEntity, err := h.service.Update(ctx, accountId, id, &entity)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEntity)
}

// Delete manages the deletion of a entity.
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ID"})
		return
	}

	ctx := c.Request.Context()
	err := h.service.Delete(ctx, c.Param("accountId"), id)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error more specifically if needed.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted"})
}

// GetByID manages fetching a entity by its ID.
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing ID"})
		return
	}

	ctx := c.Request.Context()
	entity, err := h.service.GetByID(ctx, c.Param("accountId"), id)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		// TODO: Handle error more specifically if needed.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity)
}

// GetAll manages fetching all entities with pagination.
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
	entities, totalItems, err := h.service.GetAll(ctx, c.Param("accountId"), currentPage, perPage)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination := pkg.NewPagination(currentPage, perPage, totalItems)

	response := gin.H{
		"entities":   entities,
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
	results, totalItems, err := h.service.Query(ctx, c.Param("accountId"), query, currentPage, perPage)
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
	results, totalItems, err := h.service.Pipeline(ctx, c.Param("accountId"), mongoQuery.Map(), currentPage, perPage)
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

func (h *Handler) CreateRelationship(context *gin.Context) {
	var relationship profile.Relationship
	if err := context.ShouldBindJSON(&relationship); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountId := context.Param("accountId")
	entityId := context.Param("id")

	ctx := context.Request.Context()
	entityWithRelationships, err := h.service.AddRelationship(ctx, accountId, entityId, relationship)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, entityWithRelationships)
}

func (h *Handler) ReplaceRelationships(context *gin.Context) {
	var newRelationships []profile.Relationship
	if err := context.BindJSON(&newRelationships); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	accountId := context.Param("accountId")
	entityId := context.Param("id")

	ctx := context.Request.Context()
	entityWithRelationships, err := h.service.ReplaceRelationships(ctx, accountId, entityId, newRelationships)
	if err != nil {
		err = errors.WithStack(err)
		log.Printf("%+v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, entityWithRelationships)
}
