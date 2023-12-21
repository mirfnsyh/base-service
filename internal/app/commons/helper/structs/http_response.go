package structs

// Meta defines meta format format for api format
type Meta struct {
	APIEnv  string `json:"api_env" mapstructure:"api_env"`
	Version string `json:"version" mapstructure:"version"`
}

type Response struct {
	ResponseCode string       `json:"response_code" mapstructure:"response_code"`
	ResponseDesc ResponseDesc `json:"response_desc" mapstructure:"response_desc"`
	Meta         Meta         `json:"meta" mapstructure:"meta"`
}

type SuccessResponse struct {
	Response
	Data interface{} `json:"data,omitempty" mapstructure:"data,omitempty"`
}

// error Response
type ErrorResponse struct {
	Response
	HttpStatus int `json:"-"`
}

func (e *ErrorResponse) Error() string {
	return e.ResponseDesc.EN
}

// ResponseDesc defines details data response
type ResponseDesc struct {
	ID string `json:"id" mapstructure:"id"`
	EN string `json:"en" mapstructure:"en"`
}
