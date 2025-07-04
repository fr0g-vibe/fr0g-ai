package input

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SDCardProcessor processes SD card data for threat analysis
type SDCardProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *SDCardConfig
}

// SDCardConfig holds SD card processor configuration
type SDCardConfig struct {
	MountPath         string        `yaml:"mount_path"`         // Where SD cards are mounted
	WatchDirectories  []string      `yaml:"watch_directories"`  // Directories to monitor
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
	AutoMount         bool          `yaml:"auto_mount"`
	ScanOnInsert      bool          `yaml:"scan_on_insert"`
	QuarantinePath    string        `yaml:"quarantine_path"`
	MaxFileSize       int64         `yaml:"max_file_size"`
	AllowedExtensions []string      `yaml:"allowed_extensions"`
	BlockedExtensions []string      `yaml:"blocked_extensions"`
	DeepScan          bool          `yaml:"deep_scan"`
	HashFiles         bool          `yaml:"hash_files"`
	ExtractMetadata   bool          `yaml:"extract_metadata"`
}

// SDCardData represents data found on an SD card
type SDCardData struct {
	ID           string            `json:"id"`
	DevicePath   string            `json:"device_path"`
	MountPoint   string            `json:"mount_point"`
	FileSystem   string            `json:"file_system"`
	TotalSize    int64             `json:"total_size"`
	UsedSize     int64             `json:"used_size"`
	Files        []SDCardFile      `json:"files"`
	Directories  []string          `json:"directories"`
	ScanTime     time.Time         `json:"scan_time"`
	DeviceInfo   *DeviceInfo       `json:"device_info,omitempty"`
	Metadata     map[string]string `json:"metadata"`
}

// SDCardFile represents a file found on the SD card
type SDCardFile struct {
	Path         string            `json:"path"`
	Name         string            `json:"name"`
	Extension    string            `json:"extension"`
	Size         int64             `json:"size"`
	ModTime      time.Time         `json:"mod_time"`
	Hash         string            `json:"hash,omitempty"`
	MimeType     string            `json:"mime_type,omitempty"`
	IsExecutable bool              `json:"is_executable"`
	IsHidden     bool              `json:"is_hidden"`
	Permissions  string            `json:"permissions"`
	Content      string            `json:"content,omitempty"`      // For text files
	Metadata     map[string]string `json:"metadata,omitempty"`     // EXIF, etc.
	ThreatFlags  []string          `json:"threat_flags,omitempty"` // Suspicious indicators
}

// DeviceInfo represents information about the SD card device
type DeviceInfo struct {
	Vendor       string `json:"vendor"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Capacity     int64  `json:"capacity"`
	Interface    string `json:"interface"`
	Speed        string `json:"speed"`
	Manufacturer string `json:"manufacturer"`
}

// NewSDCardProcessor creates a new SD card processor
func NewSDCardProcessor(config *SDCardConfig, aiClient AIPersonaCommunityClient) (*SDCardProcessor, error) {
	return &SDCardProcessor{
		aiClient: aiClient,
		config:   config,
	}, nil
}

// GetTag returns the processor tag
func (s *SDCardProcessor) GetTag() string {
	return "sdcard"
}

// GetDescription returns the processor description
func (s *SDCardProcessor) GetDescription() string {
	return fmt.Sprintf("SD Card Data Threat Vector Interceptor - Physical media intelligence gathering for AI community review on topic: %s", 
		s.config.CommunityTopic)
}

// ProcessWebhook processes an SD card data webhook
func (s *SDCardProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("SD Card Processor: Processing SD card threat vector webhook")
	
	// Parse SD card data from request body
	sdData, err := s.parseSDCardData(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SD card data: %w", err)
	}
	
	// Analyze SD card data for threats using AI community
	threatLevel, consensus, err := s.analyzeSDCardThreats(ctx, sdData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze SD card threats: %w", err)
	}
	
	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "SD card data threat vector submitted for community review",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"threat_level":    threatLevel,
			"consensus":       consensus,
			"review_id":       fmt.Sprintf("sdcard_review_%d", time.Now().UnixNano()),
			"device_path":     sdData.DevicePath,
			"file_count":      len(sdData.Files),
			"total_size":      sdData.TotalSize,
			"suspicious_files": s.countSuspiciousFiles(sdData.Files),
		},
		Timestamp: time.Now(),
	}
	
	log.Printf("SD Card Processor: Data analyzed - Device: %s, Files: %d, Threat Level: %s, Consensus: %.2f", 
		sdData.DevicePath, len(sdData.Files), threatLevel, consensus)
	
	return response, nil
}

// ScanSDCard scans an SD card and returns structured data
func (s *SDCardProcessor) ScanSDCard(mountPoint string) (*SDCardData, error) {
	log.Printf("SD Card Processor: Scanning SD card at %s", mountPoint)
	
	sdData := &SDCardData{
		ID:          fmt.Sprintf("sdcard_%d", time.Now().UnixNano()),
		MountPoint:  mountPoint,
		ScanTime:    time.Now(),
		Files:       []SDCardFile{},
		Directories: []string{},
		Metadata:    make(map[string]string),
	}
	
	// Get device info
	deviceInfo, err := s.getDeviceInfo(mountPoint)
	if err != nil {
		log.Printf("Warning: Could not get device info: %v", err)
	} else {
		sdData.DeviceInfo = deviceInfo
	}
	
	// Walk through all files and directories
	err = filepath.Walk(mountPoint, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Warning: Error accessing %s: %v", path, err)
			return nil // Continue scanning
		}
		
		relativePath := strings.TrimPrefix(path, mountPoint)
		if relativePath == "" {
			return nil // Skip root
		}
		
		if info.IsDir() {
			sdData.Directories = append(sdData.Directories, relativePath)
			return nil
		}
		
		// Process file
		file, err := s.processFile(path, info)
		if err != nil {
			log.Printf("Warning: Error processing file %s: %v", path, err)
			return nil // Continue scanning
		}
		
		sdData.Files = append(sdData.Files, *file)
		sdData.UsedSize += info.Size()
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to scan SD card: %w", err)
	}
	
	return sdData, nil
}

// processFile analyzes a single file
func (s *SDCardProcessor) processFile(path string, info os.FileInfo) (*SDCardFile, error) {
	file := &SDCardFile{
		Path:        path,
		Name:        info.Name(),
		Extension:   strings.ToLower(filepath.Ext(info.Name())),
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		Permissions: info.Mode().String(),
		IsHidden:    strings.HasPrefix(info.Name(), "."),
		Metadata:    make(map[string]string),
		ThreatFlags: []string{},
	}
	
	// Check if file is executable
	file.IsExecutable = (info.Mode() & 0111) != 0
	
	// Check for blocked extensions
	if s.isBlockedExtension(file.Extension) {
		file.ThreatFlags = append(file.ThreatFlags, "blocked_extension")
	}
	
	// Check for suspicious patterns
	if s.isSuspiciousFile(file) {
		file.ThreatFlags = append(file.ThreatFlags, "suspicious_pattern")
	}
	
	// Hash file if enabled and size is reasonable
	if s.config.HashFiles && file.Size < s.config.MaxFileSize {
		hash, err := s.hashFile(path)
		if err == nil {
			file.Hash = hash
		}
	}
	
	// Extract content for small text files
	if s.shouldExtractContent(file) {
		content, err := s.extractFileContent(path)
		if err == nil {
			file.Content = content
		}
	}
	
	// Extract metadata if enabled
	if s.config.ExtractMetadata {
		metadata := s.extractFileMetadata(path, file.Extension)
		for k, v := range metadata {
			file.Metadata[k] = v
		}
	}
	
	return file, nil
}

// analyzeSDCardThreats analyzes SD card data for threats using AI community
func (s *SDCardProcessor) analyzeSDCardThreats(ctx context.Context, sdData *SDCardData) (string, float64, error) {
	// Create threat analysis content
	analysisContent := fmt.Sprintf(`
SD Card Data Threat Analysis Request:
Device Path: %s
Mount Point: %s
File System: %s
Total Size: %d bytes
Used Size: %d bytes
File Count: %d
Directory Count: %d
Scan Time: %s

`, sdData.DevicePath, sdData.MountPoint, sdData.FileSystem, 
	sdData.TotalSize, sdData.UsedSize, len(sdData.Files), 
	len(sdData.Directories), sdData.ScanTime.Format(time.RFC3339))
	
	// Add device info if available
	if sdData.DeviceInfo != nil {
		analysisContent += fmt.Sprintf(`
Device Information:
- Vendor: %s
- Model: %s
- Serial Number: %s
- Capacity: %d bytes
- Interface: %s
- Speed: %s

`, sdData.DeviceInfo.Vendor, sdData.DeviceInfo.Model, 
			sdData.DeviceInfo.SerialNumber, sdData.DeviceInfo.Capacity,
			sdData.DeviceInfo.Interface, sdData.DeviceInfo.Speed)
	}
	
	// Add file analysis
	suspiciousFiles := 0
	executableFiles := 0
	hiddenFiles := 0
	
	analysisContent += "File Analysis:\n"
	for _, file := range sdData.Files {
		if len(file.ThreatFlags) > 0 {
			suspiciousFiles++
			analysisContent += fmt.Sprintf("- SUSPICIOUS: %s (%s, %d bytes) - Flags: %s\n", 
				file.Name, file.Extension, file.Size, strings.Join(file.ThreatFlags, ", "))
		}
		if file.IsExecutable {
			executableFiles++
		}
		if file.IsHidden {
			hiddenFiles++
		}
	}
	
	analysisContent += fmt.Sprintf(`
Summary Statistics:
- Suspicious Files: %d
- Executable Files: %d
- Hidden Files: %d

`, suspiciousFiles, executableFiles, hiddenFiles)
	
	// Add content samples for small text files
	textSamples := 0
	for _, file := range sdData.Files {
		if file.Content != "" && textSamples < 5 {
			analysisContent += fmt.Sprintf("Text File Content Sample (%s):\n%s\n\n", 
				file.Name, s.truncateString(file.Content, 500))
			textSamples++
		}
	}
	
	analysisContent += `
Please analyze this SD card data for potential threats including:
- Malware and malicious executables
- Data exfiltration attempts
- Suspicious file patterns and naming
- Hidden or encrypted files
- Unauthorized data storage
- Privacy violations and sensitive data
- Autorun and persistence mechanisms
- Steganography and hidden data
- Forensic artifacts and evidence tampering
- Industrial espionage indicators
- Insider threat indicators
- Compliance violations
`
	
	// Create AI community for threat analysis
	community, err := s.aiClient.CreateCommunity(ctx, s.config.CommunityTopic, s.config.PersonaCount)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to create AI community: %w", err)
	}
	
	// Submit for AI community review
	review, err := s.aiClient.SubmitForReview(ctx, community.ID, analysisContent)
	if err != nil {
		return "unknown", 0.0, fmt.Errorf("failed to submit for review: %w", err)
	}
	
	// Determine threat level based on consensus
	threatLevel := "unknown"
	consensus := 0.0
	
	if review.Consensus != nil {
		consensus = review.Consensus.OverallScore
		
		// SD card threats have different thresholds
		if consensus >= 0.9 {
			threatLevel = "critical"
		} else if consensus >= 0.8 {
			threatLevel = "high"
		} else if consensus >= 0.6 {
			threatLevel = "medium"
		} else if consensus >= 0.4 {
			threatLevel = "low"
		} else {
			threatLevel = "minimal"
		}
	}
	
	return threatLevel, consensus, nil
}

// Helper functions

func (s *SDCardProcessor) parseSDCardData(body interface{}) (*SDCardData, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}
	
	// If mount_point is provided, scan the SD card
	if mountPoint := getStringFromMap(bodyMap, "mount_point"); mountPoint != "" {
		return s.ScanSDCard(mountPoint)
	}
	
	// Otherwise parse provided data
	sdData := &SDCardData{
		ID:          getStringFromMap(bodyMap, "id"),
		DevicePath:  getStringFromMap(bodyMap, "device_path"),
		MountPoint:  getStringFromMap(bodyMap, "mount_point"),
		FileSystem:  getStringFromMap(bodyMap, "file_system"),
		ScanTime:    time.Now(),
		Files:       []SDCardFile{},
		Directories: []string{},
		Metadata:    make(map[string]string),
	}
	
	if totalSize, ok := bodyMap["total_size"].(float64); ok {
		sdData.TotalSize = int64(totalSize)
	}
	
	if usedSize, ok := bodyMap["used_size"].(float64); ok {
		sdData.UsedSize = int64(usedSize)
	}
	
	// Parse files
	if filesData, ok := bodyMap["files"].([]interface{}); ok {
		for _, fileData := range filesData {
			if fileMap, ok := fileData.(map[string]interface{}); ok {
				file := SDCardFile{
					Path:         getStringFromMap(fileMap, "path"),
					Name:         getStringFromMap(fileMap, "name"),
					Extension:    getStringFromMap(fileMap, "extension"),
					MimeType:     getStringFromMap(fileMap, "mime_type"),
					Hash:         getStringFromMap(fileMap, "hash"),
					Content:      getStringFromMap(fileMap, "content"),
					Permissions:  getStringFromMap(fileMap, "permissions"),
					IsExecutable: getBoolFromMap(fileMap, "is_executable"),
					IsHidden:     getBoolFromMap(fileMap, "is_hidden"),
					Metadata:     make(map[string]string),
					ThreatFlags:  []string{},
				}
				
				if size, ok := fileMap["size"].(float64); ok {
					file.Size = int64(size)
				}
				
				sdData.Files = append(sdData.Files, file)
			}
		}
	}
	
	return sdData, nil
}

func (s *SDCardProcessor) isBlockedExtension(ext string) bool {
	for _, blocked := range s.config.BlockedExtensions {
		if strings.EqualFold(ext, blocked) {
			return true
		}
	}
	return false
}

func (s *SDCardProcessor) isSuspiciousFile(file *SDCardFile) bool {
	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"autorun", "setup", "install", "update", "patch",
		"crack", "keygen", "serial", "license",
		"password", "pass", "pwd", "secret",
		"backup", "dump", "export", "extract",
	}
	
	nameLower := strings.ToLower(file.Name)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}
	
	// Check for double extensions
	if strings.Count(file.Name, ".") > 1 {
		return true
	}
	
	// Check for executable files with document extensions
	if file.IsExecutable && (strings.HasSuffix(file.Extension, ".pdf") || 
		strings.HasSuffix(file.Extension, ".doc") || 
		strings.HasSuffix(file.Extension, ".jpg")) {
		return true
	}
	
	return false
}

func (s *SDCardProcessor) shouldExtractContent(file *SDCardFile) bool {
	if file.Size > 10240 { // 10KB limit for content extraction
		return false
	}
	
	textExtensions := []string{".txt", ".log", ".cfg", ".conf", ".ini", ".json", ".xml", ".csv"}
	for _, ext := range textExtensions {
		if file.Extension == ext {
			return true
		}
	}
	
	return false
}

func (s *SDCardProcessor) hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (s *SDCardProcessor) extractFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	return string(content), nil
}

func (s *SDCardProcessor) extractFileMetadata(path, extension string) map[string]string {
	metadata := make(map[string]string)
	
	// Basic file metadata
	if stat, err := os.Stat(path); err == nil {
		metadata["size"] = fmt.Sprintf("%d", stat.Size())
		metadata["mod_time"] = stat.ModTime().Format(time.RFC3339)
		metadata["mode"] = stat.Mode().String()
	}
	
	// TODO: Add EXIF extraction for images, metadata for documents, etc.
	
	return metadata
}

func (s *SDCardProcessor) getDeviceInfo(mountPoint string) (*DeviceInfo, error) {
	// TODO: Implement device info extraction using system calls
	// This would involve reading from /proc, /sys, or using system commands
	
	return &DeviceInfo{
		Vendor:       "Unknown",
		Model:        "Unknown",
		SerialNumber: "Unknown",
		Interface:    "USB",
	}, nil
}

func (s *SDCardProcessor) countSuspiciousFiles(files []SDCardFile) int {
	count := 0
	for _, file := range files {
		if len(file.ThreatFlags) > 0 {
			count++
		}
	}
	return count
}

func (s *SDCardProcessor) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
