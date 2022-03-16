package secure

// Config defines the configuration options for the secure middleware
type Config struct {
	// XSSProtection provides protection against cross-site scripting attack (XSS)
	// by setting the `X-XSS-Protection` header.
	//
	// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
	//
	// Optional. Default value "1; mode=block".
	//
	// Possible values:
	// - "0" - Disables XSS filtering.
	// - "1" - Enables XSS filtering. If a cross-site scripting attack is detected, the browser will sanitize the page.
	// - "1; mode=block" - Enables XSS filtering. If a cross-site scripting attack is detected, the browser will prevent rendering.
	// - "1; report=<reporting-URI>" - Enables XSS filtering (Chromium only). If a cross-site scripting attack is detected, the browser will report the violation using the functionality of the CSP `report-uri` directive to send a report.
	XSSProtection string

	// XFrameOptions can be used to indicate whether or not a browser should
	// be allowed to render a page in a <frame>, <iframe> or <object> .
	// Sites can use this to avoid clickjacking attacks, by ensuring that their
	// content is not embedded into other sites by setting the `X-Frame-Options` header.
	//
	// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	//
	// Optional. Default value "SAMEORIGIN".
	//
	// Possible values:
	// - "SAMEORIGIN" - The page can only be displayed in a frame on the same origin as the page itself.
	// - "DENY" - The page cannot be displayed in a frame, regardless of the site attempting to do so.
	// - "ALLOW-FROM uri" - The page can only be displayed in a frame on the specified origin.
	XFrameOptions string

	// ContentSecurityPolicy provides protection against cross-site scripting (XSS),
	// clickjacking and other code injection attacks resulting from execution of malicious
	// content in the trusted web page context by setting the `Content-Security-Policy` header.
	//
	// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy
	//
	// Optional. Default value "".
	ContentSecurityPolicy string

	// ContentTypeNosniff provides protection against overriding Content-Type
	// header by setting the `X-Content-Type-Options` header.
	//
	// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options
	//
	// Optional. Default value "nosniff".
	ContentTypeNosniff string

	// ReferrerPolicy provides protection against leaking potentially sensitive request paths
	// to third parties by setting the `Referrer-Policy` header.
	//
	// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
	//
	// Optional. Default value "".
	ReferrerPolicy string
}

// DefaultConfig contains the default value for the
// secure middleware configuration
var DefaultConfig = &Config{
	XSSProtection:         "1; mode=block",
	XFrameOptions:         "SAMEORIGIN",
	ContentSecurityPolicy: "",
	ContentTypeNosniff:    "nosniff",
	ReferrerPolicy:        "",
}
