package _const

import (
	"fmt"
	"net/http"
)

var (
	RedisKeyNotFound       = fmt.Errorf("Can not find key from Redis")
	ErrorUserHaveBeenLock  = fmt.Errorf("User have been lock")
	ErrorUserHaveBeenBlock = fmt.Errorf("User have been block")
)

var (
	CodeSuccess             = customCode{code: 0, message: "Success", detail: nil, httpStatus: http.StatusOK}
	CodeUsernameExisted     = customCode{code: 100, message: "Username already exists", detail: nil, httpStatus: http.StatusOK}
	CodeWrongPassword       = customCode{code: 101, message: "Wrong password", detail: nil, httpStatus: http.StatusOK}
	CodeExceedLoginAttempts = customCode{code: 102, message: "Exceed login attempts", detail: nil, httpStatus: http.StatusOK}
	CodeUserNotActivated    = customCode{code: 103, message: "User not verify email", detail: nil, httpStatus: http.StatusOK}
	CodeLockingAccount      = customCode{code: 104, message: "This account have been lock", detail: nil, httpStatus: http.StatusOK}
	CodePassMissUpperCase   = customCode{code: 105, message: "Password should have 1 letter upper case", detail: nil, httpStatus: http.StatusOK}
	CodePassMissLowerCase   = customCode{code: 106, message: "Password should have 1 letter lower case", detail: nil, httpStatus: http.StatusOK}
	CodePassMissNumber      = customCode{code: 107, message: "Password should have 1 number", detail: nil, httpStatus: http.StatusOK}
	CodePassMissSpecial     = customCode{code: 108, message: "Password should have 1 special char", detail: nil, httpStatus: http.StatusOK}
	CodeWhiteSpaceInAccount = customCode{code: 109, message: "White space in account", detail: nil, httpStatus: http.StatusOK}
	CodeAccountHaveSpecial  = customCode{code: 110, message: "Account shouldn't have special char", detail: nil, httpStatus: http.StatusOK}
	CodeAccountBeginNumber  = customCode{code: 111, message: "Account shouldn't begin with number", detail: nil, httpStatus: http.StatusOK}
	CodeWrongOldPassword    = customCode{code: 113, message: "Wrong old password", detail: nil, httpStatus: http.StatusOK}
	CodeTenantExisted       = customCode{code: 114, message: "Tenant already exists", detail: nil, httpStatus: http.StatusOK}
	CodeCompanyExisted      = customCode{code: 115, message: "Company already exists", detail: nil, httpStatus: http.StatusOK}
	CodeInvalidPhoneNumber  = customCode{code: 116, message: "Invalid phone number", detail: nil, httpStatus: http.StatusOK}
	CodeOnlyOneCompany      = customCode{code: 117, message: "Only one company is allowed", detail: nil, httpStatus: http.StatusOK}

	CodeInvalidToken              = customCode{code: 201, message: "Invalid token", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeTokenExpired              = customCode{code: 202, message: "Token expired", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeUserNotFound              = customCode{code: 203, message: "User not found", detail: nil, httpStatus: http.StatusBadRequest}
	CodeBadgateway                = customCode{code: 204, message: "Bad gateway", detail: nil, httpStatus: http.StatusBadGateway}
	CodeSyncTokenInvalid          = customCode{code: 301, message: "Invalid Sync Token", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeSyncTokenExpired          = customCode{code: 302, message: "Sync token is expired", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeSyncUserLimit             = customCode{code: 303, message: "Your data reach limit synced users", detail: nil, httpStatus: http.StatusBadRequest}
	CodeTenantNotFound            = customCode{code: 304, message: "Tenant is not found", detail: nil, httpStatus: http.StatusBadRequest}
	CodeTokenNotFound             = customCode{code: 305, message: "Token not found", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeCompanyIdNotFound         = customCode{code: 306, message: "Company id not found", detail: nil, httpStatus: http.StatusBadRequest}
	CodeCompanyNotFound           = customCode{code: 307, message: "Company not found", detail: nil, httpStatus: http.StatusBadRequest}
	CodeCompanyIdInHeaderMisMatch = customCode{code: 331, message: "Company ID in header and path is mismatch", detail: nil, httpStatus: http.StatusOK}
	CodeCompanyCtxNotFound        = customCode{code: 308, message: "Company Ctx not found", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeUserCtxNotFound           = customCode{code: 309, message: "User Ctx not found", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeNot1stTimeLogin           = customCode{code: 310, message: "This user is not 1st time login", detail: nil, httpStatus: http.StatusOK}
	CodeEmailExists               = customCode{code: 311, message: "This email is exists", detail: nil, httpStatus: http.StatusOK}
	CodePhoneExists               = customCode{code: 312, message: "This phone is exists", detail: nil, httpStatus: http.StatusOK}
	CodeLoginLimit                = customCode{code: 313, message: "Your login reach limit call", detail: nil, httpStatus: http.StatusOK}
	CodeLoginRate                 = customCode{code: 314, message: "Your login reach rate call", detail: nil, httpStatus: http.StatusOK}
	CodeUserNotInCompany          = customCode{code: 315, message: "This user param is not belong to company", detail: nil, httpStatus: http.StatusBadRequest}
	CodeBlockingAccount           = customCode{code: 316, message: "Your account is blocking", detail: nil, httpStatus: http.StatusOK}
	CodeStatusNotAvailable        = customCode{code: 317, message: "Status input is not available", detail: nil, httpStatus: http.StatusOK}
	CodeUserNotInTenant           = customCode{code: 318, message: "This user is not belong to tenant", detail: nil, httpStatus: http.StatusBadRequest}
	CodeUserHaveBeenLock          = customCode{code: 319, message: "This user is already lock", detail: nil, httpStatus: http.StatusOK}
	CodeUserHaveBeenUnlock        = customCode{code: 320, message: "This user is already unlock", detail: nil, httpStatus: http.StatusOK}
	CodeListCompaniesNotFound     = customCode{code: 321, message: "List companies not found", detail: nil, httpStatus: http.StatusBadRequest}
	CodeRoleNotFound              = customCode{code: 322, message: "Role not found", detail: nil, httpStatus: http.StatusBadRequest}
	CodeAgencyConnectionNotFound  = customCode{code: 323, message: "Agency connection not found", detail: nil, httpStatus: http.StatusOK}

	CodeUserAndCompanyNotSameTenant = customCode{code: 315, message: "Company and user not same tenant", detail: nil, httpStatus: http.StatusBadRequest}
	CodeApiKeyCtxNotFound           = customCode{code: 323, message: "Api key ctx not found", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeCompanyLogoEmpty            = customCode{code: 324, message: "The company logo is empty", detail: nil, httpStatus: http.StatusOK}
	CodeDuplicateCompany            = customCode{code: 325, message: "The company domain already existed", detail: nil, httpStatus: http.StatusOK}
	CodeDuplicateCompanyMerchant    = customCode{code: 326, message: "The company merchant already existed", detail: nil, httpStatus: http.StatusOK}
	UploadFileSizeMaxExceedLimit    = customCode{code: 330, message: "Maximum file size is 5Mb", detail: nil, httpStatus: http.StatusOK}
	UploadFileExtensionNotSupported = customCode{code: 329, message: "Only accept png/jpge file", detail: nil, httpStatus: http.StatusOK}
	CodeNickAndNickOrTenantName     = customCode{code: 331, message: "Nick and NickOrTenantName can not be used at the same time", detail: nil, httpStatus: http.StatusOK}
	CodeAgencyConnectionExisted     = customCode{code: 332, message: "Agency already existed", detail: nil, httpStatus: http.StatusOK}
	CodeInvalidCompanyDomain        = customCode{code: 333, message: "Invalid company domain", detail: nil, httpStatus: http.StatusOK}

	CodePermissionNotAllowed = customCode{code: 401, message: "Permission not allowed", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeActionNotAllowed     = customCode{code: 402, message: "Action not allowed", detail: nil, httpStatus: http.StatusBadRequest}
	CodeDBError              = customCode{code: 500, message: "Database can not execute", detail: nil, httpStatus: http.StatusInternalServerError}
	CodeBadRequest           = customCode{code: 400, message: "Bad request", detail: nil, httpStatus: http.StatusOK}
	CodeInternalError        = customCode{code: 501, message: "Internal Error", detail: nil, httpStatus: http.StatusExpectationFailed}
	CodeUpdateUserFailed     = customCode{code: 502, message: "Update user fail", detail: nil, httpStatus: http.StatusOK}

	CodeIPRestricted = customCode{code: 1080, message: "IP restricted", detail: nil, httpStatus: http.StatusNotFound}

	CodeCompanyDomainAlreadyExists = customCode{code: 1000, message: "Company domain already exists", detail: nil, httpStatus: http.StatusOK}
)

type customCode struct {
	code       int
	message    string
	detail     interface{}
	httpStatus int
}

// Code returns the integer number of current error code.
func (c customCode) Code() int {
	return c.code
}

// Message returns the brief message for current error code.
func (c customCode) Message() string {
	return c.message
}

// Detail returns the detailed information of current error code,
// which is mainly designed as an extension field for error code.
func (c customCode) Detail() interface{} {
	return c.detail
}

// String returns current error code as a string.
func (c customCode) String() string {
	if c.detail != nil {
		return fmt.Sprintf(`%d:%s %v`, c.code, c.message, c.detail)
	}
	if c.message != "" {
		return fmt.Sprintf(`%d:%s`, c.code, c.message)
	}
	return fmt.Sprintf(`%d`, c.code)
}

func (c customCode) HttpStatus() int {
	return c.httpStatus
}
