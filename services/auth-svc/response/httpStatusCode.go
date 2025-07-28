package response

// httpStatusCode defines common HTTP status codes for responses.
type (
	httpStatusCode int
	httpMesssage   string
)

const (
	StatusOK                  httpStatusCode = 200
	StatusCreated             httpStatusCode = 201
	StatusBadRequest          httpStatusCode = 400
	StatusUnauthorized        httpStatusCode = 401
	StatusForbidden           httpStatusCode = 403
	StatusNotFound            httpStatusCode = 404
	StatusInternalServerError httpStatusCode = 500
)
