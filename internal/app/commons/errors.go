package commons

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mirfnsyh/base-service/internal/app/commons/helper/structs"
	"github.com/sirupsen/logrus"
	"net/http"

	httpHelper "github.com/mirfnsyh/base-service/internal/app/commons/helper/http"
)

type Response struct {
	Code       string
	HttpStatus int
	Error      error
	EN         string
}

type TraceError struct {
	TrueError string
	Traces    []string
	Err       error
	ErrString string
}

type ErrorWithDetails struct {
	TrueError string
	Traces    []string
	Err       error
	ErrString string
}

type Error struct {
	Code       string
	HttpStatus int
	EN         string
	ErrorWithDetails
}

func (e TraceError) Error() string {
	return fmt.Sprintf("%v", e.Err.Error())
}

func (e Error) Error() string {
	return e.Code
}

const (
	ErrDBConn              = "ErrDBConnection"
	ErrDataNotFound        = "ErrDataNotFound"
	ErrBadRequest          = "ErrBadRequest"
	ErrValidateAttribute   = "ErrValidationAttributeError"
	ErrDataNotValid        = "ErrDataNotValid"
	ErrUnprocessableEntity = "ErrUnprocessableEntity"
	ErrInternalServerError = "ErrInternalServerError"
	ErrUnauthorized        = "ErrUnauthorized"
)

func ResponseError(err string, errParam ...string) Error {
	errorMapResponse := make(map[string]Error)

	maxErrParam := 3
	startDummyErrParam := len(errParam)
	for i := startDummyErrParam; i < maxErrParam; i++ {
		errParam = append(errParam, "")
	}

	errorMapResponse[ErrDBConn] = Error{
		Code:       "1001",
		EN:         "Database connection error",
		HttpStatus: http.StatusInternalServerError,
	}
	return errorMapResponse[err]
}

func PrintError(err interface{}, c *gin.Context, http *httpHelper.HttpHandlerContext) {
	if err, ok := err.(TraceError); ok {
		ErrorLogPrint(err.Traces, err.TrueError)

		WriteError(c, err.Err, http)
		return
	}

	WriteError(c, err.(error), http)
}

func ErrorLogPrint(traces []string, trueError string) {
	var tracesString string
	for i := len(traces) - 1; i >= 0; i-- {
		tracesString = fmt.Sprintf("%s|%s", tracesString, traces[i])
	}

	logrus.Errorf("%s : %s", tracesString, trueError)
}

// WriteError sending error response based on err type
func WriteError(c *gin.Context, err error, http *httpHelper.HttpHandlerContext) {
	errd, ok := err.(Error)
	if !ok {
		var errorResponse *structs.ErrorResponse = &structs.ErrorResponse{}
		errorResponse = structs.ErrUnknown
		errorResponse.Meta = http.M
		writeErrorResponse(c, errorResponse)
	} else {
		var errorResponse *structs.ErrorResponse = &structs.ErrorResponse{}
		errorResponse.HttpStatus = errd.HttpStatus
		errorResponse.ResponseCode = errd.Code
		errorResponse.ResponseDesc.EN = errd.EN
		errorResponse.Meta = http.M
		writeErrorResponse(c, errorResponse)
	}
}

func writeErrorResponse(c *gin.Context, errorResponse *structs.ErrorResponse) {
	c.JSON(errorResponse.HttpStatus, errorResponse)
}
