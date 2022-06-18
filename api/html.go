package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
)

// Headers godoc
// @Summary Get html properties
// @Description Get html properties
// @ID html-properties
// @Accept   json
// @Tags	HTML
// @Produce  json
// @Success 200 {object} types.HTMLType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /v1/html-properties [get]
func (ctl *controller) getHTMLProperties(c *gin.Context) {

	htmlProps := config.Config.HTML
	c.JSON(http.StatusOK, htmlProps)

}
