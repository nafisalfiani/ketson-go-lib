package codes

import (
	"math"
)

type Code uint32

type Message struct {
	StatusCode int    `json:"status_code"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

type AppMessage map[Code]Message

// * Important: For new codes, please add them to the bottom of corresponding list to avoid changing existing codes and potentially breaking existing flow

const NoCode Code = math.MaxUint32

const (
	// Success code
	CodeSuccess = Code(iota + 10)
	CodeAccepted
)

const (
	// common errors
	CodeInvalidValue = Code(iota + 1000)
	CodeContextDeadlineExceeded
	CodeContextCanceled
	CodeInternalServerError
	CodeServerUnavailable
	CodeNotImplemented
	CodeBadRequest
	CodeNotFound
	CodeConflict
	CodeUnauthorized
	CodeTooManyRequest
	CodeMarshal
	CodeUnmarshal
)

const (
	// SQL errors
	CodeSQL = Code(iota + 1300)
	CodeSQLInit
	CodeSQLTx
	CodeSQLRead
	CodeSQLRowScan
	CodeSQLRecordDoesNotExist
	CodeSQLUniqueConstraint
	CodeSQLConflict
	CodeSQLNoRowsAffected
	CodeSQLInsert
	CodeSQLUpdate
	CodeSQLDelete
)

const (
	// NoSQL errors
	CodeNoSQL = Code(iota + 1400)
	CodeNoSQLInit
	CodeNoSQLRead
	CodeNoSQLRecordDoesNotExist
	CodeNoSQLUniqueConstraint
	CodeNoSQLConflict
	CodeNoSQLNoRowsAffected
	CodeNoSQLInsert
	CodeNoSQLUpdate
	CodeNoSQLDelete
)

const (
	// third party/client errors
	CodeClient = Code(iota + 1500)
	CodeClientMarshal
	CodeClientUnmarshal
	CodeClientErrorOnRequest
	CodeClientErrorOnReadBody
)

const (
	// general file I/O errors
	CodeFile = Code(iota + 1600)
	CodeFilePathOpenFailed
)

const (
	// auth errors
	CodeAuth = Code(iota + 1700)
	CodeAuthWrongPassword
	CodeAuthRefreshTokenExpired
	CodeAuthAccessTokenExpired
	CodeAuthFailure
	CodeAuthInvalidToken
	CodeForbidden
	CodeAuthRevokeRefreshTokenFailed
)

const (
	// JSON encoding errors
	CodeJSONSchema = Code(iota + 1900)
	CodeJSONSchemaInvalid
	CodeJSONSchemaNotFound
	CodeJSONStructInvalid
	CodeJSONRawInvalid
	CodeJSONValidationError
	CodeJSONMarshalError
	CodeJSONUnmarshalError
)

const (
	// data conversion error
	CodeConvert = Code(iota + 2200)
	CodeConvertTime
)

const (
	// Reset Password Error
	CodePasswordDoesNotMatch = Code(iota + 3800)
	CodeFailedResetPassword
	CodeResetPasswordTokenExpired
	CodeEmptyEmail
	CodeInvalidEmail
	CodeSameCurrentPassword
	CodePasswordIsNotFilled
	CodeResetPasswordTokenInvalid
	CodePasswordIsWeak
)

const (
	// Redis Cache Error
	CodeRedisGet = Code(iota + 3900)
	CodeRedisSetex
	CodeFailedLock
	CodeFailedReleaseLock
	CodeLockExist
	CodeCacheMarshal
	CodeCacheUnmarshal
	CodeCacheGetSimpleKey
	CodeCacheSetSimpleKey
	CodeCacheDeleteSimpleKey
	CodeCacheGetHashKey
	CodeCacheSetHashKey
	CodeCacheDeleteHashKey
	CodeCacheSetExpiration
	CodeCacheDecode
	CodeCacheLockNotAcquired
	CodeCacheInvalidCastType
	CodeCacheNotFound
)

const (
	CodeErrorHttpNewRequest = Code(iota + 4000)
	CodeErrorHttpDo
	CodeErrorIoutilReadAll
	CodeHttpUnmarshal
	CodeHttpMarshal
)

const (
	// Code Feature Flag Retriever Errors
	CodeFeatureFlagRetrieverFailed = Code(iota + 4100)
)

const (
	CodeEmail = Code(iota + 4200)
	CodeSendEmailFailed
	CodeParseHTMlTemplateFailed
	CodeConvertMJMLToHTMLFailed
)

// Error messages only
var ErrorMessages = AppMessage{
	CodeInvalidValue:            ErrMsgBadRequest,
	CodeContextDeadlineExceeded: ErrMsgContextTimeout,
	CodeContextCanceled:         ErrMsgContextTimeout,
	CodeInternalServerError:     ErrMsgInternalServerError,
	CodeServerUnavailable:       ErrMsgServiceUnavailable,
	CodeNotImplemented:          ErrMsgNotImplemented,
	CodeBadRequest:              ErrMsgBadRequest,
	CodeNotFound:                ErrMsgNotFound,
	CodeConflict:                ErrMsgConflict,
	CodeUnauthorized:            ErrMsgUnauthorized,
	CodeTooManyRequest:          ErrMsgTooManyRequest,
	CodeMarshal:                 ErrMsgBadRequest,
	CodeUnmarshal:               ErrMsgBadRequest,
	CodeJSONMarshalError:        ErrMsgBadRequest,
	CodeJSONUnmarshalError:      ErrMsgBadRequest,
	CodeJSONValidationError:     ErrMsgBadRequest,

	CodeSQL:                   ErrMsgInternalServerError,
	CodeSQLInit:               ErrMsgInternalServerError,
	CodeSQLTx:                 ErrMsgInternalServerError,
	CodeSQLRead:               ErrMsgInternalServerError,
	CodeSQLRowScan:            ErrMsgInternalServerError,
	CodeSQLRecordDoesNotExist: ErrMsgNotFound,
	CodeSQLUniqueConstraint:   ErrMsgConflict,
	CodeSQLConflict:           ErrMsgConflict,
	CodeSQLNoRowsAffected:     ErrMsgInternalServerError,
	CodeSQLInsert:             ErrMsgInternalServerError,
	CodeSQLUpdate:             ErrMsgInternalServerError,
	CodeSQLDelete:             ErrMsgInternalServerError,

	CodeNoSQL:                   ErrMsgInternalServerError,
	CodeNoSQLInit:               ErrMsgInternalServerError,
	CodeNoSQLRead:               ErrMsgInternalServerError,
	CodeNoSQLRecordDoesNotExist: ErrMsgNotFound,
	CodeNoSQLUniqueConstraint:   ErrMsgConflict,
	CodeNoSQLConflict:           ErrMsgConflict,
	CodeNoSQLNoRowsAffected:     ErrMsgInternalServerError,
	CodeNoSQLInsert:             ErrMsgInternalServerError,
	CodeNoSQLUpdate:             ErrMsgInternalServerError,
	CodeNoSQLDelete:             ErrMsgInternalServerError,

	CodeClient:                ErrMsgInternalServerError,
	CodeClientMarshal:         ErrMsgInternalServerError,
	CodeClientUnmarshal:       ErrMsgInternalServerError,
	CodeClientErrorOnRequest:  ErrMsgBadRequest,
	CodeClientErrorOnReadBody: ErrMsgBadRequest,

	CodeAuth:                         ErrMsgUnauthorized,
	CodeAuthWrongPassword:            ErrMsgPasswordDoesNotMatch,
	CodeAuthRefreshTokenExpired:      ErrMsgRefreshTokenExpired,
	CodeAuthAccessTokenExpired:       ErrMsgAccessTokenExpired,
	CodeAuthFailure:                  ErrMsgUnauthorized,
	CodeAuthInvalidToken:             ErrMsgInvalidToken,
	CodeForbidden:                    ErrMsgForbidden,
	CodeAuthRevokeRefreshTokenFailed: ErrMsgRevokeRefreshTokenFailed,

	CodeConvert:     ErrMsgInternalServerError,
	CodeConvertTime: ErrMsgInternalServerError,

	CodeRedisGet:             ErrMsgInternalServerError,
	CodeRedisSetex:           ErrMsgInternalServerError,
	CodeFailedLock:           ErrMsgInternalServerError,
	CodeFailedReleaseLock:    ErrMsgInternalServerError,
	CodeCacheMarshal:         ErrMsgInternalServerError,
	CodeCacheUnmarshal:       ErrMsgInternalServerError,
	CodeCacheGetSimpleKey:    ErrMsgInternalServerError,
	CodeCacheSetSimpleKey:    ErrMsgInternalServerError,
	CodeCacheDeleteSimpleKey: ErrMsgInternalServerError,
	CodeCacheGetHashKey:      ErrMsgInternalServerError,
	CodeCacheSetHashKey:      ErrMsgInternalServerError,
	CodeCacheDeleteHashKey:   ErrMsgInternalServerError,
	CodeCacheSetExpiration:   ErrMsgInternalServerError,
	CodeCacheDecode:          ErrMsgInternalServerError,
	CodeCacheLockNotAcquired: ErrMsgInternalServerError,
	CodeCacheInvalidCastType: ErrMsgInternalServerError,
	CodeCacheNotFound:        ErrMsgInternalServerError,

	CodeErrorHttpNewRequest: ErrMsgInternalServerError,
	CodeErrorHttpDo:         ErrMsgInternalServerError,
	CodeErrorIoutilReadAll:  ErrMsgInternalServerError,
	CodeHttpMarshal:         ErrMsgInternalServerError,
	CodeHttpUnmarshal:       ErrMsgInternalServerError,

	CodeFeatureFlagRetrieverFailed: ErrMsgInternalServerError,

	// File I/O error
	CodeFile:               ErrMsgInternalServerError,
	CodeFilePathOpenFailed: ErrMsgInternalServerError,
}

// Successful messages only
var ApplicationMessages = AppMessage{
	CodeSuccess:  SuccessDefault,
	CodeAccepted: SuccessAccepted,
}

func Compile(code Code) Message {
	if appMsg, ok := ApplicationMessages[code]; ok {
		return appMsg
	}

	return Message{
		StatusCode: SuccessDefault.StatusCode,
		Title:      SuccessDefault.Title,
		Body:       SuccessDefault.Body,
	}
}
