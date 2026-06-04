package handler

import (
	"errors"
	"net/http"
	"strconv"

	"blog/server/internal/middleware"
	"blog/server/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) ListPublished(c *gin.Context) {
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))
	authorID, _ := strconv.Atoi(c.Query("authorId"))

	items, pagination, err := h.projectService.ListPublished(service.ProjectFilter{
		FeaturedOnly: c.Query("featured") == "true",
		Keyword:      c.Query("keyword"),
		Stack:        c.Query("stack"),
		AuthorID:     uint(authorID),
		Sort:         c.DefaultQuery("sort", "featured"),
		Page:         page,
		PageSize:     pageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载项目列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *ProjectHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var viewerID uint
	var role string
	if authUser := middleware.GetAuthUser(c); authUser != nil {
		viewerID = authUser.ID
		role = authUser.Role
	}

	item, err := h.projectService.GetByID(uint(id), viewerID, role)
	if err != nil {
		status := http.StatusNotFound
		if errors.Is(err, service.ErrProjectNotPublished) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *ProjectHandler) Mine(c *gin.Context) {
	authUser := middleware.GetAuthUser(c)
	page, pageSize := service.ParsePagination(c.Query("page"), c.Query("pageSize"))

	items, pagination, err := h.projectService.ListMine(authUser.ID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "加载我的项目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items, "pagination": pagination})
}

func (h *ProjectHandler) MineDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)

	item, err := h.projectService.GetByID(uint(id), authUser.ID, authUser.Role)
	if err != nil {
		status := http.StatusNotFound
		if errors.Is(err, service.ErrProjectNotPublished) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var payload service.ProjectPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	item, err := h.projectService.Create(authUser.ID, authUser.Role, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": item})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var payload service.ProjectPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	authUser := middleware.GetAuthUser(c)
	item, err := h.projectService.Update(uint(id), authUser.ID, authUser.Role, payload)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrProjectNoPermission) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *ProjectHandler) Submit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	authUser := middleware.GetAuthUser(c)

	item, err := h.projectService.Submit(uint(id), authUser.ID, authUser.Role)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrProjectNoPermission) {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": item})
}
