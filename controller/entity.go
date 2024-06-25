package controller

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"password-lock/models"
	"password-lock/validations"
	"strings"
)

func (c Controller) CreateEntity(ctx *gin.Context) {

	var entity models.Entity

	err := json.NewDecoder(strings.NewReader(ctx.PostForm("entity"))).Decode(&entity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	err = validations.ValidateCreateEntityRequest(entity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	if entity.Type < 6 && entity.Type > 0 {
		entity.IconPath = c.service.GetEntityIconPath(entity.Type)
	} else {
		entity.IconPath = c.service.Cfg.DefaultEntityIconPath
	}

	me := ctx.Value("me").(string)

	entity.UserUuid = me

	if entity.Type != 6 && entity.Name == "" {
		entity.Name = models.TypeMap[entity.Type]
	}

	createdEntity, err := c.service.CreateEntity(ctx, entity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if entity.Type == 6 {
		var path string
		file, err := ctx.FormFile("file")
		if err != nil || file == nil {
			data := map[string]interface{}{
				"entity": createdEntity,
			}

			response, err := json.Marshal(data)
			if err != nil {
				c.SendResponse(ctx, Response{
					Status: http.StatusInternalServerError,
					Error:  err.Error(),
				})
			}

			encryptedResponse, err := c.encryptResponse(string(response))
			if err != nil {
				c.SendResponse(ctx, Response{
					Status: http.StatusInternalServerError,
					Error:  err.Error(),
				})
			}

			ctx.JSON(http.StatusCreated, encryptedResponse)
			return
		} else {
			path = strings.Join([]string{me, createdEntity.Uuid, file.Filename}, "/")
		}

		err = c.service.UploadIconToBucket(ctx, path, file)
		if err != nil {
			c.SendResponse(ctx, Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			})
			return
		}

		createdEntity.IconPath = path

		err = c.service.UpdateEntity(ctx, createdEntity)
		if err != nil {
			c.SendResponse(ctx, Response{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			})
			return
		}
	}

	data := map[string]interface{}{
		"entity": createdEntity,
	}

	response, err := json.Marshal(data)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	encryptedResponse, err := c.encryptResponse(string(response))
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, encryptedResponse)
	return
}

func (c Controller) UpdateEntity(ctx *gin.Context) {
	var updatedEntity *models.Entity

	// decoding json message to user model
	err := json.NewDecoder(strings.NewReader(ctx.PostForm("entity"))).Decode(&updatedEntity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	file, _ := ctx.FormFile("file")
	if file != nil {
		var path string

		me := ctx.Value("me").(string)
		path = strings.Join([]string{me, updatedEntity.Uuid, file.Filename}, "/")

		if file != nil {
			err = c.service.UploadIconToBucket(ctx, path, file)
			if err != nil {
				c.SendResponse(ctx, Response{
					Status: http.StatusInternalServerError,
					Error:  err.Error(),
				})
				return
			}
		}

		updatedEntity.IconPath = path
	}

	err = c.service.UpdateEntity(ctx, updatedEntity)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	data := map[string]interface{}{
		"entity": updatedEntity,
	}

	response, err := json.Marshal(data)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	encryptedResponse, err := c.encryptResponse(string(response))
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, encryptedResponse)
	return
}

func (c Controller) DeleteEntity(ctx *gin.Context) {

	me := ctx.Value("me").(string)

	entityUuid := ctx.Param("entity_uuid")

	entity, err := c.service.GetEntityByUuid(ctx, entityUuid)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	if entity.UserUuid != me {
		c.SendResponse(ctx, Response{
			Status: http.StatusForbidden,
			Error:  err.Error(),
		})
		return
	}

	err = c.service.DeleteEntity(ctx, entityUuid)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	c.SendResponse(ctx, Response{
		Status:  http.StatusOK,
		Message: "entity successfully deleted",
	})
	return

}

func (c Controller) GetEntity(ctx *gin.Context) {

	entity, err := c.service.GetEntityByUuid(ctx, ctx.Query("entity_uuid"))
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	data := map[string]interface{}{
		"entity": entity,
	}

	response, err := json.Marshal(data)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	encryptedResponse, err := c.encryptResponse(base64.StdEncoding.EncodeToString(response))

	ctx.JSON(200, encryptedResponse)
}

func (c Controller) ListEntities(ctx *gin.Context) {
	entities, err := c.service.ListEntities(ctx)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	data := map[string]interface{}{
		"entities": entities,
	}

	response, err := json.Marshal(data)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	encryptedResponse, err := c.encryptResponse(string(response))

	ctx.JSON(200, encryptedResponse)
}

func (c Controller) GetEntityIconSignedUrl(ctx *gin.Context) {

	entityUuid := ctx.Param("entity_uuid")

	entity, err := c.service.GetEntityByUuid(ctx, entityUuid)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusNotFound,
			Error:  err.Error(),
		})
	}

	signedUrl, err := c.service.GetEntityIconSignedUrl(ctx, entity.IconPath)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	data := map[string]interface{}{
		"signed_url": signedUrl,
	}

	response, err := json.Marshal(data)
	if err != nil {
		c.SendResponse(ctx, Response{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
	}

	encryptedResponse, err := c.encryptResponse(string(response))

	ctx.JSON(200, encryptedResponse)
}
