package api

// health_checks_http_security.go contains per-header validation for the HTTP
// security-header checks attached to health-check HTTP endpoint tests.

import (
	"net/http"
	"strings"
)

// checkSecurityHeaders evaluates security headers in the response.
func checkSecurityHeaders(headers http.Header, isHTTPS bool) *SecurityHeaders {
	result := &SecurityHeaders{}
	score := 0
	maxScore := 0

	// HSTS (Strict-Transport-Security) - only for HTTPS
	if isHTTPS {
		maxScore += 20
		result.HSTS = checkHSTSHeader(headers.Get("Strict-Transport-Security"))
		if result.HSTS.Present && result.HSTS.Status == statusSuccess {
			score += 20
		}
	}

	// Content-Security-Policy
	maxScore += 20
	result.CSP = checkCSPHeader(headers.Get("Content-Security-Policy"))
	if result.CSP.Present && result.CSP.Status == statusSuccess {
		score += 20
	}

	// X-Frame-Options
	maxScore += 15
	result.XFrameOptions = checkXFrameOptionsHeader(headers.Get("X-Frame-Options"))
	if result.XFrameOptions.Present && result.XFrameOptions.Status == statusSuccess {
		score += 15
	}

	// X-Content-Type-Options
	maxScore += 15
	result.XContentType = checkXContentTypeHeader(headers.Get("X-Content-Type-Options"))
	if result.XContentType.Present && result.XContentType.Status == statusSuccess {
		score += 15
	}

	// X-XSS-Protection (deprecated but still checked)
	maxScore += 10
	result.XSSProtection = checkXSSProtectionHeader(headers.Get("X-XSS-Protection"))
	if result.XSSProtection.Present {
		score += 10
	}

	// Referrer-Policy
	maxScore += 10
	result.ReferrerPolicy = checkReferrerPolicyHeader(headers.Get("Referrer-Policy"))
	if result.ReferrerPolicy.Present && result.ReferrerPolicy.Status == statusSuccess {
		score += 10
	}

	// Permissions-Policy
	maxScore += 10
	result.PermissionsPolicy = checkPermissionsPolicyHeader(headers.Get("Permissions-Policy"))
	if result.PermissionsPolicy.Present {
		score += 10
	}

	// Calculate overall score and status
	if maxScore > 0 {
		result.Score = (score * percentMultiplier) / maxScore
	}

	switch {
	case result.Score >= scoreThresholdGood:
		result.OverallStatus = statusSuccess
	case result.Score >= scoreThresholdWarn:
		result.OverallStatus = statusWarning
	default:
		result.OverallStatus = statusError
	}

	return result
}

func checkHSTSHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusError
		check.Message = "Missing: Add Strict-Transport-Security header"
		return check
	}
	check.Present = true
	// Check for max-age directive
	if strings.Contains(strings.ToLower(value), "max-age=") {
		check.Status = statusSuccess
		if strings.Contains(strings.ToLower(value), "includesubdomains") {
			check.Message = "Good: HSTS enabled with includeSubDomains"
		} else {
			check.Message = "OK: Consider adding includeSubDomains"
		}
	} else {
		check.Status = statusWarning
		check.Message = "Warning: max-age directive not found"
	}
	return check
}

func checkCSPHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusError
		check.Message = "Missing: Add Content-Security-Policy header"
		return check
	}
	check.Present = true
	// Check for unsafe directives
	lowerVal := strings.ToLower(value)
	switch {
	case strings.Contains(lowerVal, "'unsafe-inline'") || strings.Contains(lowerVal, "'unsafe-eval'"):
		check.Status = statusWarning
		check.Message = "Warning: Contains unsafe directives"
	case strings.Contains(lowerVal, "default-src") || strings.Contains(lowerVal, "script-src"):
		check.Status = statusSuccess
		check.Message = "Good: CSP policy defined"
	default:
		check.Status = statusWarning
		check.Message = "Warning: Consider adding script-src directive"
	}
	return check
}

func checkXFrameOptionsHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusError
		check.Message = "Missing: Add X-Frame-Options (DENY or SAMEORIGIN)"
		return check
	}
	check.Present = true
	upperVal := strings.ToUpper(value)
	if upperVal == "DENY" || upperVal == "SAMEORIGIN" {
		check.Status = statusSuccess
		check.Message = "Good: Clickjacking protection enabled"
	} else {
		check.Status = statusWarning
		check.Message = "Warning: Use DENY or SAMEORIGIN"
	}
	return check
}

func checkXContentTypeHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusError
		check.Message = "Missing: Add X-Content-Type-Options: nosniff"
		return check
	}
	check.Present = true
	if strings.EqualFold(value, "nosniff") {
		check.Status = statusSuccess
		check.Message = "Good: MIME type sniffing protection enabled"
	} else {
		check.Status = statusWarning
		check.Message = "Warning: Value should be 'nosniff'"
	}
	return check
}

func checkXSSProtectionHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusWarning
		check.Message = "Not present (deprecated header, CSP preferred)"
		return check
	}
	check.Present = true
	check.Status = statusSuccess
	check.Message = "Present (deprecated, rely on CSP instead)"
	return check
}

func checkReferrerPolicyHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusWarning
		check.Message = "Missing: Consider adding Referrer-Policy"
		return check
	}
	check.Present = true
	lowerVal := strings.ToLower(value)
	goodPolicies := []string{"strict-origin", "strict-origin-when-cross-origin", "no-referrer", "same-origin"}
	for _, policy := range goodPolicies {
		if strings.Contains(lowerVal, policy) {
			check.Status = statusSuccess
			check.Message = "Good: Secure referrer policy"
			return check
		}
	}
	check.Status = statusWarning
	check.Message = "Warning: Consider stricter referrer policy"
	return check
}

func checkPermissionsPolicyHeader(value string) *HeaderCheck {
	check := &HeaderCheck{Value: value}
	if value == "" {
		check.Present = false
		check.Status = statusWarning
		check.Message = "Missing: Consider adding Permissions-Policy"
		return check
	}
	check.Present = true
	check.Status = statusSuccess
	check.Message = "Good: Feature policy defined"
	return check
}
