package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saichler/l8bugs/go/bugs/common"
	l8bugs "github.com/saichler/l8bugs/go/types/l8bugs"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8web/go/web/server"
	"github.com/saichler/l8web/go/web/webhook"
	"github.com/saichler/l8web/go/web/webhook/github"
	"github.com/saichler/l8web/go/web/webhook/gitlab"
)

const (
	bugService     = "Bug"
	featureService = "Feature"
	projectService = "Project"
	serviceArea    = byte(20)
)

// Register registers the GitHub and GitLab webhook handlers on the REST server.
func Register(svr *server.RestServer, vnic ifs.IVNic) {
	h := &webhookHandler{vnic: vnic}

	// GitHub webhook
	ghProvider := &github.Provider{}
	ghHandler := webhook.NewHandler(ghProvider, h.handleGitHubEvent, h.secretFuncGitHub)
	svr.RegisterHandler("webhook/github", ghHandler)
	fmt.Println("[webhook] GitHub webhook registered")

	// GitLab webhook
	glProvider := &gitlab.Provider{}
	glHandler := webhook.NewHandler(glProvider, h.handleGitLabEvent, h.secretFuncGitLab)
	svr.RegisterHandler("webhook/gitlab", glHandler)
	fmt.Println("[webhook] GitLab webhook registered")
}

type webhookHandler struct {
	vnic ifs.IVNic
}

// --- GitHub handlers ---

func (h *webhookHandler) handleGitHubEvent(eventType string, payload []byte) int {
	switch eventType {
	case "pull_request":
		h.handleGitHubPR(payload)
	case "push":
		h.handleGitHubPush(payload)
	default:
		fmt.Printf("[webhook] ignoring GitHub event: %s\n", eventType)
	}
	return http.StatusOK
}

func (h *webhookHandler) secretFuncGitHub(payload []byte) string {
	repoURL := github.RepoURL(payload)
	project := h.findProject(repoURL)
	if project == nil {
		fmt.Printf("[webhook] no project found for repo: %s\n", repoURL)
		return ""
	}
	return project.WebhookSecret
}

func (h *webhookHandler) handleGitHubPR(body []byte) {
	var ev github.PullRequestEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		fmt.Printf("[webhook] failed to parse GitHub PR event: %s\n", err)
		return
	}

	if ev.Action != "closed" || !ev.PR.Merged {
		return
	}

	fmt.Printf("[webhook] GitHub PR merged: %s\n", ev.PR.Title)

	refs := webhook.ExtractIssueRefs(ev.PR.Title + " " + ev.PR.Body)
	for _, ref := range refs {
		h.linkCommitToIssue(ref, ev.PR.MergeCommit, ev.PR.HTMLURL, true)
	}
}

func (h *webhookHandler) handleGitHubPush(body []byte) {
	var ev github.PushEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		fmt.Printf("[webhook] failed to parse GitHub push event: %s\n", err)
		return
	}

	for _, commit := range ev.Commits {
		refs := webhook.ExtractIssueRefs(commit.Message)
		for _, ref := range refs {
			h.linkCommitToIssue(ref, commit.ID, "", false)
		}
	}
}

// --- GitLab handlers ---

func (h *webhookHandler) handleGitLabEvent(eventType string, payload []byte) int {
	switch eventType {
	case "Push Hook":
		h.handleGitLabPush(payload)
	case "Merge Request Hook":
		h.handleGitLabMR(payload)
	default:
		fmt.Printf("[webhook] ignoring GitLab event: %s\n", eventType)
	}
	return http.StatusOK
}

func (h *webhookHandler) secretFuncGitLab(payload []byte) string {
	repoURL := gitlab.RepoURL(payload)
	project := h.findProject(repoURL)
	if project == nil {
		fmt.Printf("[webhook] no project found for GitLab repo: %s\n", repoURL)
		return ""
	}
	return project.WebhookSecret
}

func (h *webhookHandler) handleGitLabPush(body []byte) {
	var ev gitlab.PushEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		fmt.Printf("[webhook] failed to parse GitLab push event: %s\n", err)
		return
	}

	for _, commit := range ev.Commits {
		refs := webhook.ExtractIssueRefs(commit.Message)
		for _, ref := range refs {
			h.linkCommitToIssue(ref, commit.ID, "", false)
		}
	}
}

func (h *webhookHandler) handleGitLabMR(body []byte) {
	var ev gitlab.MergeRequestEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		fmt.Printf("[webhook] failed to parse GitLab MR event: %s\n", err)
		return
	}

	if ev.ObjectAttrs.Action != "merge" {
		return
	}

	fmt.Printf("[webhook] GitLab MR merged: %s\n", ev.ObjectAttrs.Title)

	refs := webhook.ExtractIssueRefs(ev.ObjectAttrs.Title + " " + ev.ObjectAttrs.Description)
	for _, ref := range refs {
		h.linkCommitToIssue(ref, ev.ObjectAttrs.MergeCommit, ev.ObjectAttrs.URL, true)
	}
}

// --- Shared logic ---

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
	query := fmt.Sprintf("select * from BugsProject where repositoryUrl='%s'", repoURL)
	projects, err := common.QueryEntities[l8bugs.BugsProject](projectService, serviceArea, query, h.vnic)
	if err != nil || len(projects) == 0 {
		return nil
	}
	return projects[0]
}
