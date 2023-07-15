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
	err     error  `json:"-"`
}

// SetResponse set the response code, and error if exists
func (c *CommonResult) SetResponse(code int, responseMessage string, err ...error) {
	c.code = code
	c.message = responseMessage

	if len(err) > 0 {
		c.err = err[0]

		if code >= 500 { // server error
			// send to logger
			log.Println("ERROR", errors.WithStack(err[0]))
		}
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

// GetError return error of response
func (c CommonResult) GetError() error {
	return c.err
}

// BaseResponse is a template for
// how the response data structure
// should be returned
type BaseResponse struct {
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BaseResponseArray is a template for
// how the response array data structure
// should be returned
type BaseResponseArray struct {
	Message string        `json:"message"`
	Data    []interface{} `json:"data,omitempty"`
}
