package dto

const (
	// Success
	CODE_SUCCESS = "200.0"
	// Failed
	CODE_FAILED_REQUEST_BODY_NOT_VALID = "400.10"
	CODE_FAILED_REQUEST_PATH_NOT_VALID = "400.11"
	CODE_FAILED_INTERNAL_PROCESS       = "500.10"
)

const (
	// Success
	MESSAGE_SUCCESS_USER_GET_BALANCE    = "success user get balance"
	MESSAGE_SUCCESS_USER_DO_TRANSACTION = "success user do transaction"
	// Failed
	MESSAGE_FAILED_GET_REQUEST_BODY    = "failed get data from body"
	MESSAGE_FAILED_GET_PATH_PARAM      = "failed get data from path"
	MESSAGE_FAILED_USER_GET_BALANCE    = "failed user get balance"
	MESSAGE_FAILED_USER_DO_TRANSACTION = "failed user do transaction"
)
