package triage

import (
	"encoding/json"
	"fmt"
	l8common "github.com/saichler/l8common/go/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
	"regexp"
	"strconv"
	"strings"
)

// StackFrame represents a single frame in a stack trace.
type StackFrame struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Package  string `json:"package"`
	Line     int    `json:"line"`
}

// StackTraceInfo contains parsed stack trace information.
type StackTraceInfo struct {
	Language  string       `json:"language"`
	Frames    []StackFrame `json:"frames"`
	ErrorType string       `json:"error_type"`
	ErrorMsg  string       `json:"error_message"`
}

// RootCauseResult holds the enhanced root cause analysis output.
type RootCauseResult struct {
	RootCause          string   `json:"root_cause"`
	LikelyFiles        []string `json:"likely_files"`
	ErrorCategory      string   `json:"error_category"`
	IsRegressionLikely bool     `json:"is_regression_likely"`
	Confidence         int32    `json:"confidence"`
	SuggestedFix       string   `json:"suggested_fix"`
}

// Stack trace patterns for different languages.
var (
	goFrameRe     = regexp.MustCompile(`^\s*(.+?)\(.*\)\s*$`)
	goFileLineRe  = regexp.MustCompile(`^\s+(.+?):(\d+)`)
	goPanicRe     = regexp.MustCompile(`^panic:\s*(.+)`)
	javaFrameRe   = regexp.MustCompile(`^\s+at\s+([^(]+)\(([^:]+):(\d+)\)`)
	javaExcRe     = regexp.MustCompile(`^([a-zA-Z0-9_.]+(?:Exception|Error)):\s*(.*)`)
	pyFrameRe     = regexp.MustCompile(`^\s+File "([^"]+)", line (\d+), in (.+)`)
	pyExcRe       = regexp.MustCompile(`^([A-Z][a-zA-Z]*(?:Error|Exception)):\s*(.*)`)
	jsFrameRe     = regexp.MustCompile(`^\s+at\s+(?:(.+?)\s+)?\(?(.+?):(\d+):\d+\)?`)
	jsErrRe       = regexp.MustCompile(`^([A-Z][a-zA-Z]*Error):\s*(.*)`)
)

// ParseStackTrace detects the language and extracts structured frame info.
func ParseStackTrace(trace string) *StackTraceInfo {
	if trace == "" {
		return nil
	}

	lines := strings.Split(trace, "\n")

	// Try each language parser in order.
	if info := parseGoTrace(lines); info != nil {
		return info
	}
	if info := parseJavaTrace(lines); info != nil {
		return info
	}
	if info := parsePythonTrace(lines); info != nil {
		return info
	}
	if info := parseJSTrace(lines); info != nil {
		return info
	}

	return &StackTraceInfo{Language: "unknown"}
}

func parseGoTrace(lines []string) *StackTraceInfo {
	info := &StackTraceInfo{Language: "go"}

	for _, line := range lines {
		if m := goPanicRe.FindStringSubmatch(line); m != nil {
			info.ErrorType = "panic"
			info.ErrorMsg = m[1]
		}
	}

	for i := 0; i < len(lines); i++ {
		if m := goFrameRe.FindStringSubmatch(lines[i]); m != nil {
			frame := StackFrame{Function: m[1]}
			if idx := strings.LastIndex(frame.Function, "/"); idx >= 0 {
				frame.Package = frame.Function[:idx]
			}
			if i+1 < len(lines) {
				if fm := goFileLineRe.FindStringSubmatch(lines[i+1]); fm != nil {
					frame.File = fm[1]
					frame.Line, _ = strconv.Atoi(fm[2])
					i++
				}
			}
			info.Frames = append(info.Frames, frame)
		}
	}

	if len(info.Frames) == 0 && info.ErrorType == "" {
		return nil
	}
	return info
}

func parseJavaTrace(lines []string) *StackTraceInfo {
	info := &StackTraceInfo{Language: "java"}

	for _, line := range lines {
		if m := javaExcRe.FindStringSubmatch(line); m != nil {
			info.ErrorType = m[1]
			info.ErrorMsg = m[2]
			break
		}
	}

	for _, line := range lines {
		if m := javaFrameRe.FindStringSubmatch(line); m != nil {
			fn := m[1]
			pkg := ""
			if idx := strings.LastIndex(fn, "."); idx >= 0 {
				pkg = fn[:idx]
			}
			lineNum, _ := strconv.Atoi(m[3])
			info.Frames = append(info.Frames, StackFrame{
				Function: fn, File: m[2], Package: pkg, Line: lineNum,
			})
		}
	}

	if len(info.Frames) == 0 {
		return nil
	}
	return info
}

func parsePythonTrace(lines []string) *StackTraceInfo {
	info := &StackTraceInfo{Language: "python"}

	for _, line := range lines {
		if m := pyExcRe.FindStringSubmatch(line); m != nil {
			info.ErrorType = m[1]
			info.ErrorMsg = m[2]
		}
	}

	for _, line := range lines {
		if m := pyFrameRe.FindStringSubmatch(line); m != nil {
			lineNum, _ := strconv.Atoi(m[2])
			info.Frames = append(info.Frames, StackFrame{
				File: m[1], Function: m[3], Line: lineNum,
			})
		}
	}

	if len(info.Frames) == 0 {
		return nil
	}
	return info
}

func parseJSTrace(lines []string) *StackTraceInfo {
	info := &StackTraceInfo{Language: "javascript"}

	for _, line := range lines {
		if m := jsErrRe.FindStringSubmatch(line); m != nil {
			info.ErrorType = m[1]
			info.ErrorMsg = m[2]
			break
		}
	}

	for _, line := range lines {
		if m := jsFrameRe.FindStringSubmatch(line); m != nil {
			lineNum, _ := strconv.Atoi(m[3])
			info.Frames = append(info.Frames, StackFrame{
				Function: m[1], File: m[2], Line: lineNum,
			})
		}
	}

	if len(info.Frames) == 0 {
		return nil
	}
	return info
}

const rootCauseSystemPrompt = `You are an expert software engineer analyzing a bug with a stack trace.
Analyze the bug report and stack trace to determine the root cause.
Respond ONLY with a JSON object, no markdown fencing, no explanation.

Confidence: 0-100 (how confident you are in your root cause analysis)`

// BuildRootCausePrompt builds an enhanced prompt when a stack trace is present.
func BuildRootCausePrompt(bug *l8bugs.Bug, stackInfo *StackTraceInfo, repoURL string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Bug Title: %s\n", bug.Title))
	if bug.Description != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", bug.Description))
	}

	b.WriteString(fmt.Sprintf("\nLanguage: %s\n", stackInfo.Language))
	if stackInfo.ErrorType != "" {
		b.WriteString(fmt.Sprintf("Error Type: %s\n", stackInfo.ErrorType))
	}
	if stackInfo.ErrorMsg != "" {
		b.WriteString(fmt.Sprintf("Error Message: %s\n", stackInfo.ErrorMsg))
	}

	if len(stackInfo.Frames) > 0 {
		b.WriteString("\nStack Frames (most recent first):\n")
		for i, f := range stackInfo.Frames {
			if i >= 10 {
				b.WriteString(fmt.Sprintf("  ... and %d more frames\n", len(stackInfo.Frames)-10))
				break
			}
			b.WriteString(fmt.Sprintf("  %d. %s", i+1, f.Function))
			if f.File != "" {
				b.WriteString(fmt.Sprintf(" (%s:%d)", f.File, f.Line))
			}
			if f.Package != "" {
				b.WriteString(fmt.Sprintf(" [%s]", f.Package))
			}
			b.WriteString("\n")
		}
	}

	if repoURL != "" {
		b.WriteString(fmt.Sprintf("\nRepository: %s\n", repoURL))
	}

	b.WriteString(`
Respond with this exact JSON structure:
{
  "root_cause": "<detailed root cause explanation>",
  "likely_files": ["<file path that likely contains the bug>", ...],
  "error_category": "<category: null_reference, concurrency, resource_leak, logic_error, type_error, io_error, config_error, other>",
  "is_regression_likely": <true/false>,
  "confidence": <0-100>,
  "suggested_fix": "<brief suggested fix approach>"
}`)

	return b.String()
}

// AnalyzeRootCause performs enhanced root cause analysis on a bug with a stack trace.
func (t *Triager) AnalyzeRootCause(bug *l8bugs.Bug) (*RootCauseResult, error) {
	if bug.StackTrace == "" {
		return nil, fmt.Errorf("no stack trace present")
	}

	stackInfo := ParseStackTrace(bug.StackTrace)
	if stackInfo == nil {
		return nil, fmt.Errorf("could not parse stack trace")
	}

	// Find project repo URL if available.
	repoURL := ""
	if bug.ProjectId != "" {
		project, _ := fetchProject(t.vnic, bug.ProjectId)
		if project != nil {
			repoURL = project.RepositoryUrl
		}
	}

	prompt := BuildRootCausePrompt(bug, stackInfo, repoURL)
	response, err := t.client.Complete(rootCauseSystemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	cleaned := extractJSON(response)
	var result RootCauseResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse root cause response: %w", err)
	}

	clampRange(&result.Confidence, 0, 100)
	return &result, nil
}

func fetchProject(vnic ifs.IVNic, projectID string) (*l8bugs.BugsProject, error) {
	result, err := l8common.GetEntity(projectServiceName, serviceArea,
		&l8bugs.BugsProject{ProjectId: projectID}, vnic)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*l8bugs.BugsProject), nil
}

// formatRootCauseForStorage combines the analysis result into a human-readable string
// that gets stored in bug.AiRootCause.
func formatRootCauseForStorage(basic string, rca *RootCauseResult) string {
	var b strings.Builder
	if basic != "" {
		b.WriteString(basic)
		b.WriteString("\n\n--- Enhanced Analysis ---\n")
	}
	b.WriteString(rca.RootCause)
	if rca.ErrorCategory != "" {
		b.WriteString(fmt.Sprintf("\nCategory: %s", rca.ErrorCategory))
	}
	if len(rca.LikelyFiles) > 0 {
		b.WriteString(fmt.Sprintf("\nLikely files: %s", strings.Join(rca.LikelyFiles, ", ")))
	}
	if rca.SuggestedFix != "" {
		b.WriteString(fmt.Sprintf("\nSuggested fix: %s", rca.SuggestedFix))
	}
	if rca.IsRegressionLikely {
		b.WriteString("\nNote: This may be a regression")
	}
	return b.String()
}
