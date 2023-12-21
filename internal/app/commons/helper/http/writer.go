package http

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mirfnsyh/base-service/internal/app/commons/helper/structs"
)

type HttpHandlerContext struct {
	M structs.Meta
	E map[error]*structs.ErrorResponse
}

func NewContextHandler(meta structs.Meta) HttpHandlerContext {
	var errMap map[error]*structs.ErrorResponse = map[error]*structs.ErrorResponse{
		// register general error here, so if there are new general error you must add it here
		structs.ErrInvalidHeader: structs.ErrInvalidHeader,
		structs.ErrUnauthorized:  structs.ErrUnauthorized,
	}

	return HttpHandlerContext{
		M: meta,
		E: errMap,
	}
}

func (hctx HttpHandlerContext) AddError(key error, value *structs.ErrorResponse) {
	hctx.E[key] = value
}

func (hctx HttpHandlerContext) AddErrorMap(errMap map[error]*structs.ErrorResponse) {
	for k, v := range errMap {
		hctx.E[k] = v
	}
}

type CustomWriter struct {
	C HttpHandlerContext
}

// func (cw *CustomWriter) Write(c *gin.Context, data interface{}) {
func (cw *HttpHandlerContext) Write(c *gin.Context, data interface{}) {
	var successResp structs.SuccessResponse
	voData := reflect.ValueOf(data)
	arrayData := []interface{}{}

	if voData.Kind() != reflect.Slice {
		if voData.IsValid() {
			arrayData = []interface{}{data}
		}
		successResp.Data = arrayData
	} else {
		if voData.Len() != 0 {
			successResp.Data = data
		} else {
			successResp.Data = arrayData
		}
	}

	successResp.ResponseCode = "000000"
	successResp.Meta = cw.M

	writeSuccessResponse(c, successResp)
}

func (cw *HttpHandlerContext) WriteData(c *gin.Context, data interface{}) {
	writeResponse(c, data, "application/json", http.StatusOK)
}

// WriteError sending error response based on err type
func (cw *HttpHandlerContext) WriteError(c *gin.Context, err error) {
	if len(cw.E) > 0 {
		errorResponse := LookupError(cw.E, err)
		if errorResponse == nil {
			errorResponse = structs.ErrUnknown
		}

		errorResponse.Meta = cw.M
		writeErrorResponse(c, errorResponse)
	} else {
		var errorResponse *structs.ErrorResponse = &structs.ErrorResponse{}
		if errors.As(err, &errorResponse) {
			errorResponse.Meta = cw.M
			writeErrorResponse(c, errorResponse)
		} else {
			errorResponse = structs.ErrUnknown
			errorResponse.Meta = cw.M
			writeErrorResponse(c, errorResponse)
		}
	}
}

func writeResponse(c *gin.Context, response interface{}, contentType string, httpStatus int) {
	c.JSON(httpStatus, response)
}

func writeSuccessResponse(c *gin.Context, response structs.SuccessResponse) {
	writeResponse(c, response, "application/json", http.StatusOK)
}

func writeErrorResponse(c *gin.Context, errorResponse *structs.ErrorResponse) {
	writeResponse(c, errorResponse, "application/json", errorResponse.HttpStatus)
}

// LookupError will get error message based on error type, with variables if you want give dynamic message error
func LookupError(lookup map[error]*structs.ErrorResponse, err error) (res *structs.ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
	}

	return
}
