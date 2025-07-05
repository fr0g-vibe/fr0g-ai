package email

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ThreatAnalyzer analyzes emails for various threats
type ThreatAnalyzer struct {
	config           *EmailConfig
	spamKeywords     []string
	phishingPatterns []*regexp.Regexp
	malwareSignatures []string
	suspiciousDomains []string
}

// NewThreatAnalyzer creates a new threat analyzer instance
func NewThreatAnalyzer(config *EmailConfig) *ThreatAnalyzer {
	analyzer := &ThreatAnalyzer{
		config: config,
	}

	analyzer.initializeDetectionRules()
	return analyzer
}

// initializeDetectionRules initializes threat detection rules
func (ta *ThreatAnalyzer) initializeDetectionRules() {
	// Spam keywords (common spam indicators)
	ta.spamKeywords = []string{
		"viagra", "cialis", "lottery", "winner", "congratulations",
		"free money", "click here", "act now", "limited time",
		"urgent", "immediate", "guaranteed", "risk-free",
		"no obligation", "call now", "order now", "buy now",
		"discount", "save money", "cheap", "lowest price",
		"make money", "work from home", "earn extra cash",
		"weight loss", "lose weight", "diet pills",
		"casino", "gambling", "poker", "slots",
		"refinance", "mortgage", "loan", "credit",
		"rolex", "replica", "watches", "designer",
		"pharmacy", "prescription", "medication",
		"enlargement", "enhancement", "performance",
		"dating", "singles", "meet women", "meet men",
		"investment", "stock", "trading", "profit",
		"inheritance", "beneficiary", "transfer funds",
	}

	// Phishing patterns (regex patterns for common phishing attempts)
	phishingPatterns := []string{
		`(?i)verify.*account`,
		`(?i)suspend.*account`,
		`(?i)update.*payment`,
		`(?i)confirm.*identity`,
		`(?i)click.*link.*immediately`,
		`(?i)urgent.*action.*required`,
		`(?i)account.*will.*be.*closed`,
		`(?i)security.*alert`,
		`(?i)unauthorized.*access`,
		`(?i)login.*credentials`,
	}

	ta.phishingPatterns = make([]*regexp.Regexp, len(phishingPatterns))
	for i, pattern := range phishingPatterns {
		ta.phishingPatterns[i] = regexp.MustCompile(pattern)
	}

	// Malware signatures (simplified - would use more sophisticated detection in production)
	ta.malwareSignatures = []string{
		"X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*", // EICAR test string
		"malware_signature_1",
		"malware_signature_2",
	}

	// Suspicious domains (known bad domains)
	ta.suspiciousDomains = []string{
		"suspicious-domain.com",
		"phishing-site.net",
		"malware-host.org",
		"spam-sender.info",
	}
}

// AnalyzeEmail performs comprehensive threat analysis on an email
func (ta *ThreatAnalyzer) AnalyzeEmail(email *EmailMessage) (*EmailThreatAnalysis, error) {
	analysis := &EmailThreatAnalysis{
		SuspiciousLinks:   []string{},
		SuspiciousWords:   []string{},
		AttachmentThreats: []AttachmentThreat{},
		Recommendations:   []string{},
	}

	// Analyze spam indicators
	analysis.SpamScore = ta.analyzeSpam(email, analysis)

	// Analyze phishing indicators
	analysis.PhishingScore = ta.analyzePhishing(email, analysis)

	// Analyze malware indicators
	analysis.MalwareScore = ta.analyzeMalware(email, analysis)

	// Analyze domain reputation
	analysis.DomainReputation = ta.analyzeDomainReputation(email)

	// Analyze email authentication (simplified)
	ta.analyzeAuthentication(email, analysis)

	// Analyze attachments
	ta.analyzeAttachments(email, analysis)

	// Generate recommendations
	ta.generateRecommendations(analysis)

	return analysis, nil
}

// analyzeSpam analyzes the email for spam indicators
func (ta *ThreatAnalyzer) analyzeSpam(email *EmailMessage, analysis *EmailThreatAnalysis) float64 {
	score := 0.0
	content := strings.ToLower(email.Subject + " " + email.Body)

	// Check for spam keywords
	keywordCount := 0
	for _, keyword := range ta.spamKeywords {
		if strings.Contains(content, strings.ToLower(keyword)) {
			keywordCount++
			analysis.SuspiciousWords = append(analysis.SuspiciousWords, keyword)
		}
	}

	// Calculate spam score based on keyword density
	if keywordCount > 0 {
		score = float64(keywordCount) / 10.0 // Normalize to 0-1 scale
		if score > 1.0 {
			score = 1.0
		}
	}

	// Additional spam indicators
	if strings.Contains(content, "!!!") {
		score += 0.1
	}
	if strings.Contains(content, "$$") {
		score += 0.1
	}
	if len(email.Subject) > 100 {
		score += 0.1
	}
	if strings.Contains(strings.ToUpper(email.Subject), "RE:") && email.Subject == strings.ToUpper(email.Subject) {
		score += 0.2 // All caps subject with RE:
	}

	// Ensure score doesn't exceed 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// analyzePhishing analyzes the email for phishing indicators
func (ta *ThreatAnalyzer) analyzePhishing(email *EmailMessage, analysis *EmailThreatAnalysis) float64 {
	score := 0.0
	content := email.Subject + " " + email.Body

	// Check for phishing patterns
	patternMatches := 0
	for _, pattern := range ta.phishingPatterns {
		if pattern.MatchString(content) {
			patternMatches++
		}
	}

	if patternMatches > 0 {
		score = float64(patternMatches) / 5.0 // Normalize to 0-1 scale
		if score > 1.0 {
			score = 1.0
		}
	}

	// Check for suspicious links
	links := ta.extractLinks(email.Body)
	suspiciousLinkCount := 0
	for _, link := range links {
		if ta.isSuspiciousLink(link) {
			analysis.SuspiciousLinks = append(analysis.SuspiciousLinks, link)
			suspiciousLinkCount++
		}
	}

	if suspiciousLinkCount > 0 {
		score += float64(suspiciousLinkCount) / 10.0
	}

	// Check sender reputation (simplified)
	if ta.isSuspiciousSender(email.From) {
		score += 0.3
	}

	// Ensure score doesn't exceed 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// analyzeMalware analyzes the email for malware indicators
func (ta *ThreatAnalyzer) analyzeMalware(email *EmailMessage, analysis *EmailThreatAnalysis) float64 {
	score := 0.0
	content := email.Body

	// Check for malware signatures in email body
	for _, signature := range ta.malwareSignatures {
		if strings.Contains(content, signature) {
			score = 1.0 // High confidence malware detection
			break
		}
	}

	// Check attachments for malware indicators
	for _, attachment := range email.Attachments {
		if ta.isSuspiciousAttachment(attachment) {
			score += 0.5
			analysis.AttachmentThreats = append(analysis.AttachmentThreats, AttachmentThreat{
				Filename:    attachment.Filename,
				ThreatType:  "suspicious_extension",
				Confidence:  0.7,
				Description: "File extension commonly used for malware",
			})
		}
	}

	// Ensure score doesn't exceed 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// analyzeDomainReputation analyzes the sender domain reputation
func (ta *ThreatAnalyzer) analyzeDomainReputation(email *EmailMessage) string {
	domain := ta.extractDomain(email.From)
	
	for _, suspiciousDomain := range ta.suspiciousDomains {
		if strings.Contains(domain, suspiciousDomain) {
			return "bad"
		}
	}

	// Simple heuristics for domain reputation
	if strings.Contains(domain, "temp") || strings.Contains(domain, "disposable") {
		return "suspicious"
	}

	return "unknown"
}

// analyzeAuthentication analyzes email authentication headers
func (ta *ThreatAnalyzer) analyzeAuthentication(email *EmailMessage, analysis *EmailThreatAnalysis) {
	// Simplified authentication analysis
	// In production, this would parse actual SPF, DKIM, and DMARC headers
	
	analysis.SPFResult = "not_checked"
	analysis.DKIMResult = "not_checked"
	analysis.DMARCResult = "not_checked"

	// Check for authentication headers
	if headers, exists := email.Headers["Authentication-Results"]; exists && len(headers) > 0 {
		authHeader := strings.ToLower(headers[0])
		
		if strings.Contains(authHeader, "spf=pass") {
			analysis.SPFResult = "pass"
		} else if strings.Contains(authHeader, "spf=fail") {
			analysis.SPFResult = "fail"
		}
		
		if strings.Contains(authHeader, "dkim=pass") {
			analysis.DKIMResult = "pass"
		} else if strings.Contains(authHeader, "dkim=fail") {
			analysis.DKIMResult = "fail"
		}
		
		if strings.Contains(authHeader, "dmarc=pass") {
			analysis.DMARCResult = "pass"
		} else if strings.Contains(authHeader, "dmarc=fail") {
			analysis.DMARCResult = "fail"
		}
	}
}

// analyzeAttachments analyzes email attachments for threats
func (ta *ThreatAnalyzer) analyzeAttachments(email *EmailMessage, analysis *EmailThreatAnalysis) {
	for _, attachment := range email.Attachments {
		if ta.isSuspiciousAttachment(attachment) {
			threat := AttachmentThreat{
				Filename:    attachment.Filename,
				ThreatType:  "suspicious_extension",
				Confidence:  0.6,
				Description: "File has potentially dangerous extension",
			}
			analysis.AttachmentThreats = append(analysis.AttachmentThreats, threat)
		}

		// Check for executable files
		if ta.isExecutableFile(attachment.Filename) {
			threat := AttachmentThreat{
				Filename:    attachment.Filename,
				ThreatType:  "executable",
				Confidence:  0.8,
				Description: "Executable file attachment",
			}
			analysis.AttachmentThreats = append(analysis.AttachmentThreats, threat)
		}
	}
}

// generateRecommendations generates security recommendations based on analysis
func (ta *ThreatAnalyzer) generateRecommendations(analysis *EmailThreatAnalysis) {
	if analysis.SpamScore > 0.7 {
		analysis.Recommendations = append(analysis.Recommendations, "High spam score detected - consider blocking sender")
	}

	if analysis.PhishingScore > 0.6 {
		analysis.Recommendations = append(analysis.Recommendations, "Potential phishing attempt - verify sender authenticity")
	}

	if analysis.MalwareScore > 0.5 {
		analysis.Recommendations = append(analysis.Recommendations, "Malware indicators detected - quarantine immediately")
	}

	if len(analysis.SuspiciousLinks) > 0 {
		analysis.Recommendations = append(analysis.Recommendations, "Suspicious links detected - do not click")
	}

	if len(analysis.AttachmentThreats) > 0 {
		analysis.Recommendations = append(analysis.Recommendations, "Suspicious attachments detected - scan before opening")
	}

	if analysis.DomainReputation == "bad" {
		analysis.Recommendations = append(analysis.Recommendations, "Sender from known malicious domain - block immediately")
	}

	if analysis.SPFResult == "fail" || analysis.DKIMResult == "fail" || analysis.DMARCResult == "fail" {
		analysis.Recommendations = append(analysis.Recommendations, "Email authentication failed - sender may be spoofed")
	}
}

// Helper functions

// extractLinks extracts URLs from email content
func (ta *ThreatAnalyzer) extractLinks(content string) []string {
	urlRegex := regexp.MustCompile(`https?://[^\s<>"{}|\\^` + "`" + `\[\]]+`)
	return urlRegex.FindAllString(content, -1)
}

// isSuspiciousLink checks if a link is suspicious
func (ta *ThreatAnalyzer) isSuspiciousLink(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return true // Malformed URLs are suspicious
	}

	domain := strings.ToLower(u.Hostname())
	
	// Check against known suspicious domains
	for _, suspiciousDomain := range ta.suspiciousDomains {
		if strings.Contains(domain, suspiciousDomain) {
			return true
		}
	}

	// Check for URL shorteners (could be used to hide malicious links)
	shorteners := []string{"bit.ly", "tinyurl.com", "t.co", "goo.gl", "ow.ly"}
	for _, shortener := range shorteners {
		if strings.Contains(domain, shortener) {
			return true
		}
	}

	// Check for suspicious patterns
	if strings.Contains(domain, "secure") && strings.Contains(domain, "update") {
		return true
	}

	return false
}

// isSuspiciousSender checks if the sender is suspicious
func (ta *ThreatAnalyzer) isSuspiciousSender(sender string) bool {
	domain := ta.extractDomain(sender)
	
	// Check against suspicious domains
	for _, suspiciousDomain := range ta.suspiciousDomains {
		if strings.Contains(domain, suspiciousDomain) {
			return true
		}
	}

	return false
}

// extractDomain extracts domain from email address
func (ta *ThreatAnalyzer) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

// isSuspiciousAttachment checks if an attachment is suspicious
func (ta *ThreatAnalyzer) isSuspiciousAttachment(attachment EmailAttachment) bool {
	filename := strings.ToLower(attachment.Filename)
	
	// Suspicious extensions
	suspiciousExts := []string{
		".exe", ".scr", ".bat", ".cmd", ".com", ".pif", ".vbs", ".js",
		".jar", ".zip", ".rar", ".7z", ".ace", ".arj", ".cab",
	}

	for _, ext := range suspiciousExts {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}

	return false
}

// isExecutableFile checks if a file is executable
func (ta *ThreatAnalyzer) isExecutableFile(filename string) bool {
	filename = strings.ToLower(filename)
	executableExts := []string{".exe", ".scr", ".bat", ".cmd", ".com", ".pif"}
	
	for _, ext := range executableExts {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}

	return false
}
