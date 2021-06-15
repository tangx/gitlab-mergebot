package confgitlab

import "time"

type MergeRequest struct {
	IID             int        `json:"iid,omitempty"`
	ProjectID       int        `json:"project_id,omitempty"`
	Author          User       `json:"author,omitempty"`
	Assignee        User       `json:"assignee,omitempty"`
	Title           string     `json:"title,omitempty"`
	Description     string     `json:"description,omitempty"`
	Message         string     `json:"message,omitempty"`
	TargetProjectID int        `json:"target_project_id,omitempty"`
	SourceBranch    string     `json:"source_branch,omitempty"`
	TargetBranch    string     `json:"target_branch,omitempty"`
	SourceProjectID int        `json:"source_project_id,omitempty"`
	References      References `json:"references,omitempty"`

	// 是否被 merged
	State string `json:"state,omitempty"`
	// 能否被 merge
	MergeStatus string `json:"merge_status,omitempty"`

	MergeBy User          `json:"merge_by,omitempty"`
	MergeAt time.Duration `json:"merge_at,omitempty"`

	ClosedBy User          `json:"closed_by,omitempty"`
	ClouedAt time.Duration `json:"cloued_at,omitempty"`
}

type References struct {
	Full     string `json:"full,omitempty"`
	Short    string `json:"short,omitempty"`
	Relative string `json:"relative,omitempty"`
}
type MergeRequestNotes struct {
	Body   string `json:"body"`
	Author User   `json:"author"`
}

type User struct {
	ID       int    `name:"id,omitempty" in:"query" json:"id,omitempty"`
	Username string `name:"username,omitempty" in:"query" json:"username,omitempty"`
}

// merge
const (
	MergeStatus_CannotBeMerged = "cannot_be_merged"
	MergeStatus_CanBeMerged    = "can_be_merged"
)
