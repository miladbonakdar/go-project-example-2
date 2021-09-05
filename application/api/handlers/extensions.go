package handlers

import (
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/indraframework"
	"net/http"

	"github.com/gin-gonic/gin"
)

func tryActions(c GinContext, actions ...func() (error error, dto dto.Dto)) (Success bool) {
	for _, action := range actions {
		if err, data := action(); err != nil {
			logger.WithName(logtags.ApisPreActionError).
				ErrorException(err, "error wile trying some actions")
			jsonBadRequest(c, data, err)
			return false
		}
	}
	return true
}

func jsonError(c GinContext, data dto.Dto, err *indraframework.IndraException) {
	data.SetError(err)
	c.JSON(err.ErrorCode, data)
}

func jsonBadRequest(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.BadRequestException(err.Error(), "bad request"))
}

func jsonForbiddenRequest(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.ForbiddenException(err.Error(), "forbidden request"))
}

func jsonNotFound(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.NotFoundException(err.Error(), "not found"))
}

func jsonInternalServerError(c GinContext, data dto.Dto, err error) {
	jsonError(c, data, indraframework.InternalServerException(err.Error(), "internal server error"))
}

func jsonSuccess(c GinContext, value interface{}) {
	c.JSON(http.StatusOK, value)
}

func success(c GinContext) {
	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}
