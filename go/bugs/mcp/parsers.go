package mcp

import (
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/triage"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"strings"
	"time"
)

// --- enum parsers ---

func parsePriority(s string) l8bugs.Priority {
	switch strings.ToLower(s) {
	case "critical":
		return l8bugs.Priority_PRIORITY_CRITICAL
	case "high":
		return l8bugs.Priority_PRIORITY_HIGH
	case "medium":
		return l8bugs.Priority_PRIORITY_MEDIUM
	case "low":
		return l8bugs.Priority_PRIORITY_LOW
	}
	return l8bugs.Priority_PRIORITY_UNSPECIFIED
}

func parseSeverity(s string) l8bugs.Severity {
	switch strings.ToLower(s) {
	case "blocker":
		return l8bugs.Severity_SEVERITY_BLOCKER
	case "major":
		return l8bugs.Severity_SEVERITY_MAJOR
	case "minor":
		return l8bugs.Severity_SEVERITY_MINOR
	case "trivial":
		return l8bugs.Severity_SEVERITY_TRIVIAL
	}
	return l8bugs.Severity_SEVERITY_UNSPECIFIED
}

func parseBugStatus(s string) l8bugs.BugStatus {
	switch strings.ToLower(s) {
	case "open":
		return l8bugs.BugStatus_BUG_STATUS_OPEN
	case "triaged":
		return l8bugs.BugStatus_BUG_STATUS_TRIAGED
	case "in_progress":
		return l8bugs.BugStatus_BUG_STATUS_IN_PROGRESS
	case "in_review":
		return l8bugs.BugStatus_BUG_STATUS_IN_REVIEW
	case "resolved":
		return l8bugs.BugStatus_BUG_STATUS_RESOLVED
	case "closed":
		return l8bugs.BugStatus_BUG_STATUS_CLOSED
	case "reopened":
		return l8bugs.BugStatus_BUG_STATUS_REOPENED
	case "wont_fix":
		return l8bugs.BugStatus_BUG_STATUS_WONT_FIX
	case "duplicate":
		return l8bugs.BugStatus_BUG_STATUS_DUPLICATE
	case "cannot_reproduce":
		return l8bugs.BugStatus_BUG_STATUS_CANNOT_REPRODUCE
	}
	return l8bugs.BugStatus_BUG_STATUS_UNSPECIFIED
}

func parseFeatureStatus(s string) l8bugs.FeatureStatus {
	switch strings.ToLower(s) {
	case "proposed":
		return l8bugs.FeatureStatus_FEATURE_STATUS_PROPOSED
	case "triaged":
		return l8bugs.FeatureStatus_FEATURE_STATUS_TRIAGED
	case "approved":
		return l8bugs.FeatureStatus_FEATURE_STATUS_APPROVED
	case "in_progress":
		return l8bugs.FeatureStatus_FEATURE_STATUS_IN_PROGRESS
	case "in_review":
		return l8bugs.FeatureStatus_FEATURE_STATUS_IN_REVIEW
	case "done":
		return l8bugs.FeatureStatus_FEATURE_STATUS_DONE
	case "closed":
		return l8bugs.FeatureStatus_FEATURE_STATUS_CLOSED
	case "rejected":
		return l8bugs.FeatureStatus_FEATURE_STATUS_REJECTED
	case "deferred":
		return l8bugs.FeatureStatus_FEATURE_STATUS_DEFERRED
	}
	return l8bugs.FeatureStatus_FEATURE_STATUS_UNSPECIFIED
}

func parseResolution(s string) l8bugs.Resolution {
	switch strings.ToLower(s) {
	case "fixed":
		return l8bugs.Resolution_RESOLUTION_FIXED
	case "wont_fix":
		return l8bugs.Resolution_RESOLUTION_WONT_FIX
	case "duplicate":
		return l8bugs.Resolution_RESOLUTION_DUPLICATE
	case "cannot_reproduce":
		return l8bugs.Resolution_RESOLUTION_CANNOT_REPRODUCE
	case "by_design":
		return l8bugs.Resolution_RESOLUTION_BY_DESIGN
	case "obsolete":
		return l8bugs.Resolution_RESOLUTION_OBSOLETE
	}
	return l8bugs.Resolution_RESOLUTION_UNSPECIFIED
}

func parseAuthorType(s string) l8bugs.AuthorType {
	switch strings.ToLower(s) {
	case "human":
		return l8bugs.AuthorType_AUTHOR_TYPE_HUMAN
	case "ai_agent":
		return l8bugs.AuthorType_AUTHOR_TYPE_AI_AGENT
	case "system":
		return l8bugs.AuthorType_AUTHOR_TYPE_SYSTEM
	}
	return l8bugs.AuthorType_AUTHOR_TYPE_AI_AGENT
}

func parseDigestPeriod(s string) l8bugs.DigestPeriod {
	switch strings.ToLower(s) {
	case "daily":
		return l8bugs.DigestPeriod_DIGEST_PERIOD_DAILY
	case "weekly":
		return l8bugs.DigestPeriod_DIGEST_PERIOD_WEEKLY
	case "custom":
		return l8bugs.DigestPeriod_DIGEST_PERIOD_CUSTOM
	}
	return l8bugs.DigestPeriod_DIGEST_PERIOD_UNSPECIFIED
}

// --- generate_digest ---

func (s *Server) handleGenerateDigest(args map[string]interface{}) (*CallToolResult, error) {
	projectID := getStr(args, "project_id")
	if projectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}

	periodStr := getStr(args, "period")
	period := parseDigestPeriod(periodStr)

	var startDate, endDate int64
	now := time.Now()
	switch periodStr {
	case "daily":
		startDate = now.AddDate(0, 0, -1).Unix()
		endDate = now.Unix()
	case "weekly":
		startDate = now.AddDate(0, 0, -7).Unix()
		endDate = now.Unix()
	default:
		startDate = int64(getInt(args, "start_date", 0))
		endDate = int64(getInt(args, "end_date", 0))
	}

	triager := triage.Get()
	if triager == nil || !triager.Available() {
		return nil, fmt.Errorf("AI digest generation unavailable")
	}

	digest, err := triager.GenerateDigest(projectID, period, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return textResult(protoToText(digest)), nil
}
