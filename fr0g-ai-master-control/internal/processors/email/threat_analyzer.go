package email

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ThreatAnalyzer performs threat analysis on email messages
type ThreatAnalyzer struct {
	config           *EmailConfig
	spamKeywords     []string
	phishingPatterns []*regexp.Regexp
	malwareSignatures []string
	suspiciousDomains []string
}

// NewThreatAnalyzer creates a new threat analyzer
func NewThreatAnalyzer(config *EmailConfig) *ThreatAnalyzer {
	analyzer := &ThreatAnalyzer{
		config: config,
	}
	
	analyzer.initializeDetectionRules()
	return analyzer
}

// initializeDetectionRules initializes threat detection rules
func (ta *ThreatAnalyzer) initializeDetectionRules() {
	// Spam keywords
	ta.spamKeywords = []string{
		"viagra", "cialis", "lottery", "winner", "congratulations",
		"free money", "click here", "act now", "limited time",
		"urgent", "immediate", "guaranteed", "risk-free",
		"no obligation", "call now", "order now", "buy now",
		"discount", "save money", "cash", "credit", "loan",
		"debt", "refinance", "mortgage", "investment",
		"make money", "work from home", "business opportunity",
	}

	// Phishing patterns
	phishingRegexes := []string{
		`(?i)verify.*account`,
		`(?i)suspend.*account`,
		`(?i)click.*link.*immediately`,
		`(?i)update.*payment.*information`,
		`(?i)confirm.*identity`,
		`(?i)security.*alert`,
		`(?i)unusual.*activity`,
		`(?i)login.*attempt`,
		`(?i)account.*compromised`,
		`(?i)immediate.*action.*required`,
	}

	for _, pattern := range phishingRegexes {
		if regex, err := regexp.Compile(pattern); err == nil {
			ta.phishingPatterns = append(ta.phishingPatterns, regex)
		}
	}

	// Malware signatures (simplified - in production, use proper antivirus engines)
	ta.malwareSignatures = []string{
		"X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*", // EICAR test
		"eval(", "javascript:", "vbscript:", "activex",
		"<script", "</script>", "document.write", "window.open",
	}

	// Suspicious domains
	ta.suspiciousDomains = []string{
		"bit.ly", "tinyurl.com", "t.co", "goo.gl", "ow.ly",
		"is.gd", "buff.ly", "adf.ly", "short.link",
	}
}

// AnalyzeEmail performs comprehensive threat analysis on an email
func (ta *ThreatAnalyzer) AnalyzeEmail(email *EmailMessage) (*EmailThreatAnalysis, error) {
	analysis := &EmailThreatAnalysis{
		Recommendations: make([]string, 0),
	}

	// Perform various threat analyses
	analysis.SpamScore = ta.calculateSpamScore(email)
	analysis.PhishingScore = ta.calculatePhishingScore(email)
	analysis.MalwareScore = ta.calculateMalwareScore(email)
	analysis.SuspiciousLinks = ta.extractSuspiciousLinks(email)
	analysis.SuspiciousWords = ta.extractSuspiciousWords(email)
	analysis.AttachmentThreats = ta.analyzeAttachments(email)
	analysis.DomainReputation = ta.checkDomainReputation(email.From)
	
	// Email authentication checks
	analysis.SPFResult = ta.checkSPF(email)
	analysis.DKIMResult = ta.checkDKIM(email)
	analysis.DMARCResult = ta.checkDMARC(email)

	// Generate recommendations
	analysis.Recommendations = ta.generateRecommendations(analysis)

	return analysis, nil
}

// calculateSpamScore calculates the spam probability score
func (ta *ThreatAnalyzer) calculateSpamScore(email *EmailMessage) float64 {
	score := 0.0
	totalChecks := 0

	// Check for spam keywords in subject and body
	content := strings.ToLower(email.Subject + " " + email.Body)
	keywordCount := 0
	for _, keyword := range ta.spamKeywords {
		if strings.Contains(content, keyword) {
			keywordCount++
		}
	}
	
	if len(ta.spamKeywords) > 0 {
		score += float64(keywordCount) / float64(len(ta.spamKeywords)) * 0.4
		totalChecks++
	}

	// Check for excessive capitalization
	if ta.hasExcessiveCapitalization(email.Subject + " " + email.Body) {
		score += 0.2
	}
	totalChecks++

	// Check for suspicious sender patterns
	if ta.hasSuspiciousSender(email.From) {
		score += 0.3
	}
	totalChecks++

	// Check for suspicious links
	if len(ta.extractSuspiciousLinks(email)) > 0 {
		score += 0.1
	}
	totalChecks++

	// Normalize score
	if totalChecks > 0 {
		score = score / float64(totalChecks) * 4 // Scale to 0-1 range
		if score > 1.0 {
			score = 1.0
		}
	}

	return score
}

// calculatePhishingScore calculates the phishing probability score
func (ta *ThreatAnalyzer) calculatePhishingScore(email *EmailMessage) float64 {
	score := 0.0
	content := email.Subject + " " + email.Body

	// Check for phishing patterns
	patternMatches := 0
	for _, pattern := range ta.phishingPatterns {
		if pattern.MatchString(content) {
			patternMatches++
		}
	}

	if len(ta.phishingPatterns) > 0 {
		score += float64(patternMatches) / float64(len(ta.phishingPatterns)) * 0.6
	}

	// Check for suspicious links
	suspiciousLinks := ta.extractSuspiciousLinks(email)
	if len(suspiciousLinks) > 0 {
		score += 0.3
	}

	// Check for urgency indicators
	urgencyWords := []string{"urgent", "immediate", "expires", "deadline", "act now"}
	urgencyCount := 0
	lowerContent := strings.ToLower(content)
	for _, word := range urgencyWords {
		if strings.Contains(lowerContent, word) {
			urgencyCount++
		}
	}
	if urgencyCount > 0 {
		score += 0.1
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateMalwareScore calculates the malware probability score
func (ta *ThreatAnalyzer) calculateMalwareScore(email *EmailMessage) float64 {
	score := 0.0
	content := email.Body

	// Check for malware signatures
	signatureMatches := 0
	for _, signature := range ta.malwareSignatures {
		if strings.Contains(content, signature) {
			signatureMatches++
		}
	}

	if len(ta.malwareSignatures) > 0 {
		score += float64(signatureMatches) / float64(len(ta.malwareSignatures)) * 0.7
	}

	// Check attachments for potential threats
	for _, attachment := range email.Attachments {
		if ta.isSuspiciousAttachment(attachment) {
			score += 0.3
			break // Don't double-count multiple suspicious attachments
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// extractSuspiciousLinks extracts potentially suspicious links from the email
func (ta *ThreatAnalyzer) extractSuspiciousLinks(email *EmailMessage) []string {
	var suspiciousLinks []string
	
	// Simple URL extraction regex
	urlRegex := regexp.MustCompile(`https?://[^\s<>"]+`)
	urls := urlRegex.FindAllString(email.Body, -1)

	for _, urlStr := range urls {
		if parsedURL, err := url.Parse(urlStr); err == nil {
			// Check if domain is in suspicious list
			for _, suspiciousDomain := range ta.suspiciousDomains {
				if strings.Contains(parsedURL.Host, suspiciousDomain) {
					suspiciousLinks = append(suspiciousLinks, urlStr)
					break
				}
			}

			// Check for IP addresses instead of domains
			if ta.isIPAddress(parsedURL.Host) {
				suspiciousLinks = append(suspiciousLinks, urlStr)
			}

			// Check for suspicious URL patterns
			if ta.hasSuspiciousURLPattern(urlStr) {
				suspiciousLinks = append(suspiciousLinks, urlStr)
			}
		}
	}

	return suspiciousLinks
}

// extractSuspiciousWords extracts suspicious words found in the email
func (ta *ThreatAnalyzer) extractSuspiciousWords(email *EmailMessage) []string {
	var suspiciousWords []string
	content := strings.ToLower(email.Subject + " " + email.Body)

	for _, keyword := range ta.spamKeywords {
		if strings.Contains(content, keyword) {
			suspiciousWords = append(suspiciousWords, keyword)
		}
	}

	return suspiciousWords
}

// analyzeAttachments analyzes email attachments for threats
func (ta *ThreatAnalyzer) analyzeAttachments(email *EmailMessage) []AttachmentThreat {
	var threats []AttachmentThreat

	for _, attachment := range email.Attachments {
		if threat := ta.analyzeAttachment(attachment); threat != nil {
			threats = append(threats, *threat)
		}
	}

	return threats
}

// analyzeAttachment analyzes a single attachment for threats
func (ta *ThreatAnalyzer) analyzeAttachment(attachment EmailAttachment) *AttachmentThreat {
	// Check file extension
	if ta.isSuspiciousAttachment(attachment) {
		return &AttachmentThreat{
			Filename:    attachment.Filename,
			ThreatType:  "suspicious_extension",
			Confidence:  0.8,
			Description: "File has potentially dangerous extension",
		}
	}

	// Check file content for malware signatures
	content := string(attachment.Content)
	for _, signature := range ta.malwareSignatures {
		if strings.Contains(content, signature) {
			return &AttachmentThreat{
				Filename:    attachment.Filename,
				ThreatType:  "malware_signature",
				Confidence:  0.9,
				Description: "File contains known malware signature",
			}
		}
	}

	return nil
}

// Helper functions

func (ta *ThreatAnalyzer) hasExcessiveCapitalization(text string) bool {
	if len(text) == 0 {
		return false
	}
	
	upperCount := 0
	for _, char := range text {
		if char >= 'A' && char <= 'Z' {
			upperCount++
		}
	}
	
	return float64(upperCount)/float64(len(text)) > 0.5
}

func (ta *ThreatAnalyzer) hasSuspiciousSender(sender string) bool {
	// Check for suspicious sender patterns
	suspiciousPatterns := []string{
		"noreply", "no-reply", "donotreply", "admin", "support",
		"security", "alert", "notification", "update",
	}
	
	lowerSender := strings.ToLower(sender)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerSender, pattern) {
			return true
		}
	}
	
	return false
}

func (ta *ThreatAnalyzer) isSuspiciousAttachment(attachment EmailAttachment) bool {
	suspiciousExtensions := []string{
		".exe", ".scr", ".bat", ".cmd", ".com", ".pif", ".vbs",
		".js", ".jar", ".zip", ".rar", ".7z", ".doc", ".docx",
		".xls", ".xlsx", ".ppt", ".pptx", ".pdf",
	}
	
	filename := strings.ToLower(attachment.Filename)
	for _, ext := range suspiciousExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	
	return false
}

func (ta *ThreatAnalyzer) isIPAddress(host string) bool {
	// Simple IP address detection
	ipRegex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	return ipRegex.MatchString(host)
}

func (ta *ThreatAnalyzer) hasSuspiciousURLPattern(url string) bool {
	suspiciousPatterns := []string{
		"login", "signin", "verify", "update", "secure",
		"account", "banking", "paypal", "amazon", "microsoft",
	}
	
	lowerURL := strings.ToLower(url)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerURL, pattern) {
			return true
		}
	}
	
	return false
}

func (ta *ThreatAnalyzer) checkDomainReputation(sender string) string {
	// Simplified domain reputation check
	// In production, this would query reputation databases
	if strings.Contains(sender, "@") {
		domain := strings.Split(sender, "@")[1]
		
		// Check against known bad domains (simplified)
		badDomains := []string{"example.com", "test.com", "spam.com"}
		for _, badDomain := range badDomains {
			if strings.Contains(domain, badDomain) {
				return "bad"
			}
		}
		
		// Check against known good domains
		goodDomains := []string{"gmail.com", "outlook.com", "yahoo.com"}
		for _, goodDomain := range goodDomains {
			if strings.Contains(domain, goodDomain) {
				return "good"
			}
		}
	}
	
	return "unknown"
}

func (ta *ThreatAnalyzer) checkSPF(email *EmailMessage) string {
	// Simplified SPF check - in production, implement proper SPF validation
	if spf, exists := email.Headers["Received-SPF"]; exists && len(spf) > 0 {
		if strings.Contains(strings.ToLower(spf[0]), "pass") {
			return "pass"
		} else if strings.Contains(strings.ToLower(spf[0]), "fail") {
			return "fail"
		}
	}
	return "none"
}

func (ta *ThreatAnalyzer) checkDKIM(email *EmailMessage) string {
	// Simplified DKIM check - in production, implement proper DKIM validation
	if dkim, exists := email.Headers["DKIM-Signature"]; exists && len(dkim) > 0 {
		return "pass"
	}
	return "none"
}

func (ta *ThreatAnalyzer) checkDMARC(email *EmailMessage) string {
	// Simplified DMARC check - in production, implement proper DMARC validation
	spfResult := ta.checkSPF(email)
	dkimResult := ta.checkDKIM(email)
	
	if spfResult == "pass" || dkimResult == "pass" {
		return "pass"
	}
	return "fail"
}

func (ta *ThreatAnalyzer) generateRecommendations(analysis *EmailThreatAnalysis) []string {
	var recommendations []string

	if analysis.SpamScore > 0.7 {
		recommendations = append(recommendations, "High spam probability - consider blocking sender")
	}

	if analysis.PhishingScore > 0.7 {
		recommendations = append(recommendations, "High phishing probability - do not click links or provide information")
	}

	if analysis.MalwareScore > 0.5 {
		recommendations = append(recommendations, "Potential malware detected - scan attachments before opening")
	}

	if len(analysis.SuspiciousLinks) > 0 {
		recommendations = append(recommendations, "Suspicious links detected - verify URLs before clicking")
	}

	if len(analysis.AttachmentThreats) > 0 {
		recommendations = append(recommendations, "Dangerous attachments detected - do not open")
	}

	if analysis.SPFResult == "fail" {
		recommendations = append(recommendations, "SPF authentication failed - sender may be spoofed")
	}

	if analysis.DMARCResult == "fail" {
		recommendations = append(recommendations, "DMARC authentication failed - email may be fraudulent")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Email appears safe but always exercise caution")
	}

	return recommendations
}
