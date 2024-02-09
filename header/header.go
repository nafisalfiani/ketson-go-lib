package header

const (
	// Header keys
	KeyRequestID      string = "x-request-id"
	KeyAuthorization  string = "authorization"
	KeyUserAgent      string = "user-agent"
	KeyContentType    string = "content-type"
	KeyContentAccept  string = "accept"
	KeyAcceptLanguage string = "accept-language"
	KeyCacheControl   string = "cache-control"
	KeyDeviceType     string = "x-device-type"
	KeyServiceName    string = "x-service-name"

	// Content type. Specifying the payload in the request
	ContentTypeJSON  string = "application/json"
	ContentTypeXML   string = "application/xml"
	ContentTypeForm  string = "application/x-www-form-urlencoded"
	ContentTypePlain string = "text/plain"

	// Cache control
	CacheControlNoCache string = "no-cache"
	CacheControlNoStore string = "no-store"

	// Authorization type
	AuthorizationBasic  string = "Basic"
	AuthorizationBearer string = "Bearer"
)
