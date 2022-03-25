package lungo

import "net/http"

const (
	// Version of Lungo
	Version = "0.0.1"

	// Website of Lungo
	Website = "https://github.com/felix-kaestner/lungo"

	// Header - RFC 4229
	HeaderAIM                     = "A-IM"
	HeaderAccept                  = "Accept"
	HeaderAcceptAdditions         = "Accept-Additions"
	HeaderAcceptCharset           = "Accept-Charset"
	HeaderAcceptEncoding          = "Accept-Encoding"
	HeaderAcceptFeatures          = "Accept-Features"
	HeaderAcceptLanguage          = "Accept-Language"
	HeaderAcceptRanges            = "Accept-Ranges"
	HeaderAge                     = "Age"
	HeaderAllow                   = "Allow"
	HeaderAlternates              = "Alternates"
	HeaderAuthenticationInfo      = "Authentication-Info"
	HeaderAuthorization           = "Authorization"
	HeaderCExt                    = "C-Ext"
	HeaderCMan                    = "C-Man"
	HeaderCOpt                    = "C-Opt"
	HeaderCPEP                    = "C-PEP"
	HeaderCPEPInfo                = "C-PEP-Info"
	HeaderCacheControl            = "Cache-Control"
	HeaderConnection              = "Connection"
	HeaderContentBase             = "Content-Base"
	HeaderContentDisposition      = "Content-Disposition"
	HeaderContentEncoding         = "Content-Encoding"
	HeaderContentID               = "Content-ID"
	HeaderContentLanguage         = "Content-Language"
	HeaderContentLength           = "Content-Length"
	HeaderContentLocation         = "Content-Location"
	HeaderContentMD5              = "Content-MD5"
	HeaderContentRange            = "Content-Range"
	HeaderContentScriptType       = "Content-Script-Type"
	HeaderContentStyleType        = "Content-Style-Type"
	HeaderContentType             = "Content-Type"
	HeaderContentVersion          = "Content-Version"
	HeaderCookie                  = "Cookie"
	HeaderCookie2                 = "Cookie2"
	HeaderDAV                     = "DAV"
	HeaderDate                    = "Date"
	HeaderDefaultStyle            = "Default-Style"
	HeaderDeltaBase               = "Delta-Base"
	HeaderDepth                   = "Depth"
	HeaderDerivedFrom             = "Derived-From"
	HeaderDestination             = "Destination"
	HeaderDifferentialID          = "Differential-ID"
	HeaderDigest                  = "Digest"
	HeaderETag                    = "ETag"
	HeaderExpect                  = "Expect"
	HeaderExpires                 = "Expires"
	HeaderExt                     = "Ext"
	HeaderFrom                    = "From"
	HeaderGetProfile              = "GetProfile"
	HeaderHost                    = "Host"
	HeaderIM                      = "IM"
	HeaderIf                      = "If"
	HeaderIfMatch                 = "If-Match"
	HeaderIfModifiedSince         = "If-Modified-Since"
	HeaderIfNoneMatch             = "If-None-Match"
	HeaderIfRange                 = "If-Range"
	HeaderIfUnmodifiedSince       = "If-Unmodified-Since"
	HeaderKeepAlive               = "Keep-Alive"
	HeaderLabel                   = "Label"
	HeaderLastModified            = "Last-Modified"
	HeaderLink                    = "Link"
	HeaderLocation                = "Location"
	HeaderLockToken               = "Lock-Token"
	HeaderMIMEVersion             = "MIME-Version"
	HeaderMan                     = "Man"
	HeaderMaxForwards             = "Max-Forwards"
	HeaderMeter                   = "Meter"
	HeaderNegotiate               = "Negotiate"
	HeaderOpt                     = "Opt"
	HeaderOrderingType            = "Ordering-Type"
	HeaderOrigin                  = "Origin"
	HeaderOverwrite               = "Overwrite"
	HeaderP3P                     = "P3P"
	HeaderPEP                     = "PEP"
	HeaderPICSLabel               = "PICS-Label"
	HeaderPepInfo                 = "Pep-Info"
	HeaderPosition                = "Position"
	HeaderPragma                  = "Pragma"
	HeaderProfileObject           = "ProfileObject"
	HeaderProtocol                = "Protocol"
	HeaderProtocolInfo            = "Protocol-Info"
	HeaderProtocolQuery           = "Protocol-Query"
	HeaderProtocolRequest         = "Protocol-Request"
	HeaderProxyAuthenticate       = "Proxy-Authenticate"
	HeaderProxyAuthenticationInfo = "Proxy-Authentication-Info"
	HeaderProxyAuthorization      = "Proxy-Authorization"
	HeaderProxyFeatures           = "Proxy-Features"
	HeaderProxyInstruction        = "Proxy-Instruction"
	HeaderPublic                  = "Public"
	HeaderRange                   = "Range"
	HeaderReferer                 = "Referer"
	HeaderRetryAfter              = "Retry-After"
	HeaderSafe                    = "Safe"
	HeaderSecurityScheme          = "Security-Scheme"
	HeaderServer                  = "Server"
	HeaderSetCookie               = "Set-Cookie"
	HeaderSetCookie2              = "Set-Cookie2"
	HeaderSetProfile              = "SetProfile"
	HeaderSoapAction              = "SoapAction"
	HeaderStatusURI               = "Status-URI"
	HeaderSurrogateCapability     = "Surrogate-Capability"
	HeaderSurrogateControl        = "Surrogate-Control"
	HeaderTCN                     = "TCN"
	HeaderTE                      = "TE"
	HeaderTimeout                 = "Timeout"
	HeaderTrailer                 = "Trailer"
	HeaderTransferEncoding        = "Transfer-Encoding"
	HeaderURI                     = "URI"
	HeaderUpgrade                 = "Upgrade"
	HeaderUserAgent               = "User-Agent"
	HeaderVariantVary             = "Variant-Vary"
	HeaderVary                    = "Vary"
	HeaderVia                     = "Via"
	HeaderWWWAuthenticate         = "WWW-Authenticate"
	HeaderWantDigest              = "Want-Digest"
	HeaderWarning                 = "Warning"
	HeaderXForwardedFor           = "X-Forwarded-For"
	HeaderXForwardedProto         = "X-Forwarded-Proto"
	HeaderXForwardedProtocol      = "X-Forwarded-Protocol"
	HeaderXForwardedSsl           = "X-Forwarded-Ssl"
	HeaderXUrlScheme              = "X-Url-Scheme"
	HeaderXHTTPMethodOverride     = "X-HTTP-Method-Override"
	HeaderXRealIP                 = "X-Real-IP"
	HeaderXRequestID              = "X-Request-ID"
	HeaderXRequestedWith          = "X-Requested-With"

	// Access control
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"

	// Security
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderReferrerPolicy                  = "Referrer-Policy"
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXCSRFToken                      = "X-CSRF-Token"

	// Charset
	CharsetUTF8 = "charset=utf-8"

	// MIME types
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + CharsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + CharsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + CharsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + CharsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + CharsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + CharsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

type (
	// Map is a shorthand of type `map[string]any`.
	Map map[string]any
)

// Common HTTP methods.
//
// Provided as a slice to loop over.
var methods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}
