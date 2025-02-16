package handlers

import (
	"net/http"
	"theGoodCompany/internal/services"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
    Service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
    return &DocumentHandler{
        Service: service,
    }
}

func (h *DocumentHandler) GetDocument(c *gin.Context) {
    // Get ID from path parameter
    idParam := c.Param("id")

    // Call service layer
    doc, err := h.Service.GetDocumentByID(c.Request.Context(), idParam)
    if err != nil {
        switch err {
        case services.ErrInvalidID:
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        case services.ErrNotFound:
            c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        }
        return
    }

    // Return the document with embedded data
    response := gin.H{
        "_id":  doc.ID,
        "data": doc.Data,
    }
    c.JSON(http.StatusOK, response)
}