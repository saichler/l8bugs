package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
	"io"
	"net/http"
	"strings"
)

const (
	bugService     = "Bug"
	featureService = "Feature"
	projectService = "Project"
	serviceArea    = byte(20)
)

// Register registers the GitHub webhook handler on the default HTTP mux.
func Register(vnic ifs.IVNic) {
	handler := &webhookHandler{vnic: vnic}
	http.HandleFunc("/bugs/webhook/github", handler.handle)
	fmt.Println("[webhook] GitHub webhook registered at /bugs/webhook/github")
}

type webhookHandler struct {
	vnic ifs.IVNic
}

// GitHub webhook payload types (minimal).

type ghPushEvent struct {
	Ref     string     `json:"ref"`
	Commits []ghCommit `json:"commits"`
	Repo    ghRepo     `json:"repository"`
}

type ghCommit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ghPREvent struct {
	Action string `json:"action"`
	PR     ghPR   `json:"pull_request"`
	Repo   ghRepo `json:"repository"`
}

type ghPR struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	Merged     bool   `json:"merged"`
	MergeCommit string `json:"merge_commit_sha"`
	HTMLURL    string `json:"html_url"`
}

type ghRepo struct {
	CloneURL string `json:"clone_url"`
	HTMLURL  string `json:"html_url"`
	FullName string `json:"full_name"`
}

func (h *webhookHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	if event == "" {
		http.Error(w, "missing X-GitHub-Event header", http.StatusBadRequest)
		return
	}

	// Extract repo URL from payload to find matching project.
	repoURL := extractRepoURL(body)
	project := h.findProject(repoURL)
	if project == nil {
		fmt.Printf("[webhook] no project found for repo: %s\n", repoURL)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Verify HMAC signature if webhook secret is configured.
	if project.WebhookSecret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if !verifySignature(body, sig, project.WebhookSecret) {
			http.Error(w, "invalid signature", http.StatusForbidden)
			return
		}
	}

	switch event {
	case "pull_request":
		h.handlePR(body)
	case "push":
		h.handlePush(body)
	default:
		fmt.Printf("[webhook] ignoring event: %s\n", event)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *webhookHandler) handlePR(body []byte) {
	var ev ghPREvent
	if err := json.Unmarshal(body, &ev); err != nil {
		fmt.Printf("[webhook] failed to parse PR event: %s\n", err)
		return
	}

	// Only process merged PRs.
	if ev.Action != "closed" || !ev.PR.Merged {
		return
	}

	fmt.Printf("[webhook] PR merged: %s\n", ev.PR.Title)

	// Extract issue refs from PR title and body.
	refs := ExtractIssueRefs(ev.PR.Title + " " + ev.PR.Body)
	for _, ref := range refs {
		h.linkCommitToIssue(ref, ev.PR.MergeCommit, ev.PR.HTMLURL, true)
	}
}

func (h *webhookHandler) handlePush(body []byte) {
	var ev ghPushEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		fmt.Printf("[webhook] failed to parse push event: %s\n", err)
		return
	}

	for _, commit := range ev.Commits {
		refs := ExtractIssueRefs(commit.Message)
		for _, ref := range refs {
			h.linkCommitToIssue(ref, commit.ID, "", false)
		}
	}
}

func (h *webhookHandler) linkCommitToIssue(ref, commitSHA, prURL string, autoTransition bool) {
	// Try as bug ID or bug number.
	bug, _ := common.GetEntity(bugService, serviceArea, &l8bugs.Bug{BugId: ref}, h.vnic)
	if bug == nil {
		bug, _ = common.GetEntity(bugService, serviceArea, &l8bugs.Bug{BugNumber: ref}, h.vnic)
	}
	if bug != nil {
		bug.LinkedCommitSha = commitSHA
		if prURL != "" {
			bug.LinkedPrUrl = prURL
		}
		if autoTransition && bug.Status == l8bugs.BugStatus_BUG_STATUS_IN_REVIEW {
			bug.Status = l8bugs.BugStatus_BUG_STATUS_RESOLVED
		}
		if err := common.PutEntity(bugService, serviceArea, bug, h.vnic); err != nil {
			fmt.Printf("[webhook] failed to update bug %s: %s\n", ref, err)
		} else {
			fmt.Printf("[webhook] linked commit %s to bug %s\n", commitSHA[:8], ref)
		}
		return
	}

	// Try as feature ID or feature number.
	feature, _ := common.GetEntity(featureService, serviceArea, &l8bugs.Feature{FeatureId: ref}, h.vnic)
	if feature == nil {
		feature, _ = common.GetEntity(featureService, serviceArea, &l8bugs.Feature{FeatureNumber: ref}, h.vnic)
	}
	if feature != nil {
		feature.LinkedCommitSha = commitSHA
		if prURL != "" {
			feature.LinkedPrUrl = prURL
		}
		if autoTransition && feature.Status == l8bugs.FeatureStatus_FEATURE_STATUS_IN_REVIEW {
			feature.Status = l8bugs.FeatureStatus_FEATURE_STATUS_DONE
		}
		if err := common.PutEntity(featureService, serviceArea, feature, h.vnic); err != nil {
			fmt.Printf("[webhook] failed to update feature %s: %s\n", ref, err)
		} else {
			fmt.Printf("[webhook] linked commit %s to feature %s\n", commitSHA[:8], ref)
		}
		return
	}

	fmt.Printf("[webhook] issue ref not found: %s\n", ref)
}

func (h *webhookHandler) findProject(repoURL string) *l8bugs.BugsProject {
	if repoURL == "" {
		return nil
	}
	projects, err := common.GetEntities(projectService, serviceArea,
		&l8bugs.BugsProject{RepositoryUrl: repoURL}, h.vnic)
	if err != nil || len(projects) == 0 {
		return nil
	}
	return projects[0]
}

func extractRepoURL(body []byte) string {
	var payload struct {
		Repo struct {
			CloneURL string `json:"clone_url"`
			HTMLURL  string `json:"html_url"`
		} `json:"repository"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	if payload.Repo.HTMLURL != "" {
		return payload.Repo.HTMLURL
	}
	return payload.Repo.CloneURL
}

func verifySignature(payload []byte, signature, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	sig, err := hex.DecodeString(signature[7:])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)
	return hmac.Equal(sig, expected)
}
