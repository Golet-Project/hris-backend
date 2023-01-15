package primitive

import (
	"log"

	"github.com/pkg/errors"
)

// CommonResult provides data struct
// that identifies how response should be
// returned, either success or fail
type CommonResult struct {
	code    int    `json:"-"`
	message string `json:"-"`
}

// SetResponse set the response code, and error if exists
func (c *CommonResult) SetResponse(code int, err error, message ...string) {
	c.code = code

	if code >= 400 && code < 500 { // client error
		if err != nil {
			c.message = err.Error()
		} else {
			c.message = "client error"
		}
		return
	}

	if len(message) > 0 {
		c.message = message[0]
	}

	if code >= 500 { // server error
		// send to logger
		log.Println("ERROR", errors.WithStack(err))
		c.message = "internal server error"
	}
}

// GetCode return response status code
func (c CommonResult) GetCode() int {
	return c.code
}

// GetMessage return message of response
func (c CommonResult) GetMessage() string {
	return c.message
}

// BaseResponse is a template for
// how the response data structure
// should be returned
type BaseResponse[T interface{}] struct {
	Message string `json:"message"`
	Data    *T
}

// BaseResponseArray is a template for
// how the response array data structure
// should be returned
type BaseResponseArray[T interface{}] struct {
	Message string
	Data    []T
}
