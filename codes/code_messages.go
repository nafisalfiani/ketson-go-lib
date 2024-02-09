package codes

import (
	"net/http"
)

// HTTP message
var (
	// 4xx
	ErrMsgBadRequest = Message{
		StatusCode: http.StatusBadRequest,
		Title:      http.StatusText(http.StatusBadRequest),
		Body:       "Invalid input. Please validate your input.",
	}
	ErrMsgUnauthorized = Message{
		StatusCode: http.StatusUnauthorized,
		Title:      http.StatusText(http.StatusUnauthorized),
		Body:       "Unauthorized access. You are not authorized to access this resource.",
	}
	ErrMsgInvalidToken = Message{
		StatusCode: http.StatusUnauthorized,
		Title:      http.StatusText(http.StatusUnauthorized),
		Body:       "Invalid token. Please renew your session by reloading.",
	}
	ErrMsgRefreshTokenExpired = Message{
		StatusCode: http.StatusUnauthorized,
		Title:      http.StatusText(http.StatusUnauthorized),
		Body:       "Session refresh token has expired. Please renew your session by reloading.",
	}
	ErrMsgAccessTokenExpired = Message{
		StatusCode: http.StatusUnauthorized,
		Title:      http.StatusText(http.StatusUnauthorized),
		Body:       "Session access token has expired. Please renew your session by reloading.",
	}
	ErrMsgForbidden = Message{
		StatusCode: http.StatusForbidden,
		Title:      http.StatusText(http.StatusForbidden),
		Body:       "Forbidden. You don't have permission to access this resource.",
	}
	ErrMsgRevokeRefreshTokenFailed = Message{
		StatusCode: http.StatusInternalServerError,
		Title:      http.StatusText(http.StatusInternalServerError),
		Body:       "Failed revoking refresh token.",
	}
	ErrMsgNotFound = Message{
		StatusCode: http.StatusNotFound,
		Title:      http.StatusText(http.StatusNotFound),
		Body:       "Record does not exist. Please validate your input or contact the administrator.",
	}
	ErrMsgContextTimeout = Message{
		StatusCode: http.StatusRequestTimeout,
		Title:      http.StatusText(http.StatusRequestTimeout),
		Body:       "Request time has been exceeded.",
	}
	ErrMsgConflict = Message{
		StatusCode: http.StatusConflict,
		Title:      http.StatusText(http.StatusConflict),
		Body:       "Record has existed. Please validate your input or contact the administrator.",
	}
	ErrMsgTooManyRequest = Message{
		StatusCode: http.StatusTooManyRequests,
		Title:      http.StatusText(http.StatusTooManyRequests),
		Body:       "Too many requests. Please wait and try again after a few moments.",
	}

	// 5xx
	ErrMsgInternalServerError = Message{
		StatusCode: http.StatusInternalServerError,
		Title:      http.StatusText(http.StatusInternalServerError),
		Body:       "Internal server error. Please contact the administrator.",
	}
	ErrMsgNotImplemented = Message{
		StatusCode: http.StatusNotImplemented,
		Title:      http.StatusText(http.StatusNotImplemented),
		Body:       "Not Implemented. Please contact the administrator.",
	}
	ErrMsgServiceUnavailable = Message{
		StatusCode: http.StatusServiceUnavailable,
		Title:      http.StatusText(http.StatusServiceUnavailable),
		Body:       "Service is unavailable. Please contact the administrator.",
	}

	// Application specific messages
	ErrMsgPasswordDoesNotMatch = Message{
		StatusCode: http.StatusBadRequest,
		Title:      "Entered Password Does Not Match",
		Body:       "",
	}

	ErrMsgPasswordIsWeak = Message{
		StatusCode: http.StatusBadRequest,
		Title:      "Entered Password Combination Is Weak. Please Choose a Stronger Password",
		Body:       "Password must be at least 8 characters and include at least one uppercase letter, one lowercase letter, and one special character",
	}

	ErrMsgInvalidEmail = Message{
		StatusCode: http.StatusBadRequest,
		Title:      "Your Email is Invalid",
		Body:       "Please input a valid email",
	}
)

// Application message
var (
	SuccessDefault = Message{
		StatusCode: http.StatusOK,
		Title:      http.StatusText(http.StatusOK),
	}

	SuccessAccepted = Message{
		StatusCode: http.StatusAccepted,
		Title:      http.StatusText(http.StatusAccepted),
	}
)
