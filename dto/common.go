package dto

const (
	// Success
	CODE_SUCCESS = "200.0"
	// Failed
	CODE_FAILED_REQUEST_BODY_NOT_VALID     = "400.10"
	CODE_FAILED_REQUEST_PATH_NOT_VALID     = "400.11"
	CODE_FAILED_HEADER_SOURCE_TYPE_WRONG   = "400.12"
	CODE_FAILED_TRANSACTION_ALREADY_EXISTS = "400.20"
	CODE_FAILED_BALANCE_NOT_ENOUGH         = "400.21"
	CODE_FAILED_HEADER_CONTENT_TYPE_WRONG  = "415.10"
	CODE_FAILED_INTERNAL_PROCESS           = "500.10"
)

const (
	// Success
	MESSAGE_SUCCESS_USER_GET_BALANCE    = "success get user balance"
	MESSAGE_SUCCESS_USER_DO_TRANSACTION = "success update user balance"
	// Failed
	MESSAGE_FAILED_HEADER_SOURCE_TYPE_WRONG  = "header source-type wrong"
	MESSAGE_FAILED_HEADER_CONTENT_TYPE_WRONG = "header content-type wrong"
	MESSAGE_FAILED_GET_REQUEST_BODY          = "failed get data from body"
	MESSAGE_FAILED_GET_PATH_PARAM            = "failed get data from path"
	MESSAGE_FAILED_USER_GET_BALANCE          = "failed get user balance"
	MESSAGE_FAILED_USER_DO_TRANSACTION       = "failed update user balance"
)
