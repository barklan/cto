package gitlab

import (
	"encoding/json"
	"fmt"
)

type MergeRequestOneOf struct {
	ApprovalsBeforeMerge interface{} `json:"approvals_before_merge"`
	Assignee             struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"assignee"`
	Assignees []struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"assignees"`
	Author struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved"`
	ClosedAt                    interface{} `json:"closed_at"`
	ClosedBy                    interface{} `json:"closed_by"`
	CreatedAt                   string      `json:"created_at"`
	Description                 string      `json:"description"`
	DiscussionLocked            interface{} `json:"discussion_locked"`
	Downvotes                   int64       `json:"downvotes"`
	Draft                       bool        `json:"draft"`
	ForceRemoveSourceBranch     bool        `json:"force_remove_source_branch"`
	HasConflicts                bool        `json:"has_conflicts"`
	ID                          int64       `json:"id"`
	Iid                         int64       `json:"iid"`
	Labels                      []string    `json:"labels"`
	MergeCommitSha              interface{} `json:"merge_commit_sha"`
	MergeStatus                 string      `json:"merge_status"`
	MergeWhenPipelineSucceeds   bool        `json:"merge_when_pipeline_succeeds"`
	MergedAt                    interface{} `json:"merged_at"`
	MergedBy                    interface{} `json:"merged_by"`
	Milestone                   interface{} `json:"milestone"`
	ProjectID                   int64       `json:"project_id"`
	Reference                   string      `json:"reference"`
	References                  struct {
		Full     string `json:"full"`
		Relative string `json:"relative"`
		Short    string `json:"short"`
	} `json:"references"`
	Reviewers                []interface{} `json:"reviewers"`
	Sha                      string        `json:"sha"`
	ShouldRemoveSourceBranch interface{}   `json:"should_remove_source_branch"`
	SourceBranch             string        `json:"source_branch"`
	SourceProjectID          int64         `json:"source_project_id"`
	Squash                   bool          `json:"squash"`
	SquashCommitSha          interface{}   `json:"squash_commit_sha"`
	State                    string        `json:"state"`
	TargetBranch             string        `json:"target_branch"`
	TargetProjectID          int64         `json:"target_project_id"`
	TaskCompletionStatus     struct {
		CompletedCount int64 `json:"completed_count"`
		Count          int64 `json:"count"`
	} `json:"task_completion_status"`
	TimeStats struct {
		HumanTimeEstimate   interface{} `json:"human_time_estimate"`
		HumanTotalTimeSpent interface{} `json:"human_total_time_spent"`
		TimeEstimate        int64       `json:"time_estimate"`
		TotalTimeSpent      int64       `json:"total_time_spent"`
	} `json:"time_stats"`
	Title          string `json:"title"`
	UpdatedAt      string `json:"updated_at"`
	Upvotes        int64  `json:"upvotes"`
	UserNotesCount int64  `json:"user_notes_count"`
	WebURL         string `json:"web_url"`
	WorkInProgress bool   `json:"work_in_progress"`
}

type MergeRequestSingle struct {
	ApprovalsBeforeMerge interface{}   `json:"approvals_before_merge"`
	Assignee             interface{}   `json:"assignee"`
	Assignees            []interface{} `json:"assignees"`
	Author               struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	BlockingDiscussionsResolved bool        `json:"blocking_discussions_resolved"`
	ChangesCount                string      `json:"changes_count"`
	ClosedAt                    interface{} `json:"closed_at"`
	ClosedBy                    interface{} `json:"closed_by"`
	CreatedAt                   string      `json:"created_at"`
	Description                 string      `json:"description"`
	DiffRefs                    struct {
		BaseSha  string `json:"base_sha"`
		HeadSha  string `json:"head_sha"`
		StartSha string `json:"start_sha"`
	} `json:"diff_refs"`
	DiscussionLocked            interface{} `json:"discussion_locked"`
	Downvotes                   int64       `json:"downvotes"`
	Draft                       bool        `json:"draft"`
	FirstContribution           bool        `json:"first_contribution"`
	FirstDeployedToProductionAt string      `json:"first_deployed_to_production_at"`
	ForceRemoveSourceBranch     bool        `json:"force_remove_source_branch"`
	HasConflicts                bool        `json:"has_conflicts"`
	HeadPipeline                struct {
		BeforeSha      string      `json:"before_sha"`
		CommittedAt    interface{} `json:"committed_at"`
		Coverage       string      `json:"coverage"`
		CreatedAt      string      `json:"created_at"`
		DetailedStatus struct {
			DetailsPath  string      `json:"details_path"`
			Favicon      string      `json:"favicon"`
			Group        string      `json:"group"`
			HasDetails   bool        `json:"has_details"`
			Icon         string      `json:"icon"`
			Illustration interface{} `json:"illustration"`
			Label        string      `json:"label"`
			Text         string      `json:"text"`
			Tooltip      string      `json:"tooltip"`
		} `json:"detailed_status"`
		Duration       int64  `json:"duration"`
		FinishedAt     string `json:"finished_at"`
		ID             int64  `json:"id"`
		ProjectID      int64  `json:"project_id"`
		QueuedDuration int64  `json:"queued_duration"`
		Ref            string `json:"ref"`
		Sha            string `json:"sha"`
		Source         string `json:"source"`
		StartedAt      string `json:"started_at"`
		Status         string `json:"status"`
		Tag            bool   `json:"tag"`
		UpdatedAt      string `json:"updated_at"`
		User           struct {
			AvatarURL string `json:"avatar_url"`
			ID        int64  `json:"id"`
			Name      string `json:"name"`
			State     string `json:"state"`
			Username  string `json:"username"`
			WebURL    string `json:"web_url"`
		} `json:"user"`
		WebURL     string      `json:"web_url"`
		YamlErrors interface{} `json:"yaml_errors"`
	} `json:"head_pipeline"`
	ID                        int64         `json:"id"`
	Iid                       int64         `json:"iid"`
	Labels                    []interface{} `json:"labels"`
	LatestBuildFinishedAt     string        `json:"latest_build_finished_at"`
	LatestBuildStartedAt      string        `json:"latest_build_started_at"`
	MergeCommitSha            interface{}   `json:"merge_commit_sha"`
	MergeError                interface{}   `json:"merge_error"`
	MergeStatus               string        `json:"merge_status"`
	MergeWhenPipelineSucceeds bool          `json:"merge_when_pipeline_succeeds"`
	MergedAt                  string        `json:"merged_at"`
	MergedBy                  struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"merged_by"`
	Milestone interface{} `json:"milestone"`
	Pipeline  struct {
		CreatedAt string `json:"created_at"`
		ID        int64  `json:"id"`
		ProjectID int64  `json:"project_id"`
		Ref       string `json:"ref"`
		Sha       string `json:"sha"`
		Source    string `json:"source"`
		Status    string `json:"status"`
		UpdatedAt string `json:"updated_at"`
		WebURL    string `json:"web_url"`
	} `json:"pipeline"`
	ProjectID  int64  `json:"project_id"`
	Reference  string `json:"reference"`
	References struct {
		Full     string `json:"full"`
		Relative string `json:"relative"`
		Short    string `json:"short"`
	} `json:"references"`
	Reviewers                []interface{} `json:"reviewers"`
	Sha                      string        `json:"sha"`
	ShouldRemoveSourceBranch interface{}   `json:"should_remove_source_branch"`
	SourceBranch             string        `json:"source_branch"`
	SourceProjectID          int64         `json:"source_project_id"`
	Squash                   bool          `json:"squash"`
	SquashCommitSha          interface{}   `json:"squash_commit_sha"`
	State                    string        `json:"state"`
	Subscribed               bool          `json:"subscribed"`
	TargetBranch             string        `json:"target_branch"`
	TargetProjectID          int64         `json:"target_project_id"`
	TaskCompletionStatus     struct {
		CompletedCount int64 `json:"completed_count"`
		Count          int64 `json:"count"`
	} `json:"task_completion_status"`
	TimeStats struct {
		HumanTimeEstimate   interface{} `json:"human_time_estimate"`
		HumanTotalTimeSpent interface{} `json:"human_total_time_spent"`
		TimeEstimate        int64       `json:"time_estimate"`
		TotalTimeSpent      int64       `json:"total_time_spent"`
	} `json:"time_stats"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
	Upvotes   int64  `json:"upvotes"`
	User      struct {
		CanMerge bool `json:"can_merge"`
	} `json:"user"`
	UserNotesCount int64  `json:"user_notes_count"`
	WebURL         string `json:"web_url"`
	WorkInProgress bool   `json:"work_in_progress"`
}

type ApprovedMREvent struct {
	ActionName string `json:"action_name"`
	Author     struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	AuthorID       int64  `json:"author_id"`
	AuthorUsername string `json:"author_username"`
	CreatedAt      string `json:"created_at"`
	ID             int64  `json:"id"`
	ProjectID      int64  `json:"project_id"`
	TargetID       int64  `json:"target_id"`
	TargetIid      int64  `json:"target_iid"`
	TargetTitle    string `json:"target_title"`
	TargetType     string `json:"target_type"`
}

func GetOpenMergeRequests(extraQueryParams string) ([]MergeRequestOneOf, error) {
	merge_requests := make([]MergeRequestOneOf, 1)
	dump, _ := request(fmt.Sprintf("merge_requests?state=opened%s", extraQueryParams))
	err := json.Unmarshal(dump, &merge_requests)
	if err != nil {
		return nil, err
	}
	return merge_requests, nil
}

func GetApprovedMREvents(extraQueryParams string) ([]ApprovedMREvent, error) {
	events := make([]ApprovedMREvent, 1)
	dump, _ := request(fmt.Sprintf("events?target_type=merge_request&action=approved%s", extraQueryParams))
	err := json.Unmarshal(dump, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
