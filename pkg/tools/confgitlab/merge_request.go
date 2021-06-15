package confgitlab

import (
	"context"

	"github.com/go-courier/httptransport/httpx"
	"github.com/pkg/errors"
	pkgerrors "github.com/pkg/errors"
)

// MergeReqeustsInput list merge requests
// https://docs.gitlab.com/ee/api/merge_requests.html#list-merge-requests
type MergeReqeustsInput struct {
	httpx.MethodGet
	State string `name:"state,omitempty" in:"query"`
	Scope string `name:"scope,omitempty" in:"query"`
}

func (MergeReqeustsInput) Path() string {
	return "/api/v4/merge_requests"
}

func (g *GitlabV4) ListMergeReqeusts() (mrs []*MergeRequest, err error) {
	req := MergeReqeustsInput{
		State: "opened",
		Scope: "assigned_to_me",
	}

	_, err = g.c.Do(context.Background(), &req, g.metas).Into(&mrs)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "get Merge Requests failed: %v", err)
	}

	return mrs, nil
}

// MergeReqeustNotesInput list all mr notes
// https://docs.gitlab.com/ee/api/notes.html#list-all-merge-request-notes
type MergeReqeustNotesInput struct {
	httpx.MethodGet
	ProjectID       int    `name:"project_id" in:"path"`
	MergeRequestIID int    `name:"merge_request_iid" in:"path"`
	Sort            string `name:"sort,omitempty" in:"query"`
	OrderBy         string `name:"order_by,omitempty" in:"query"`
}

func (MergeReqeustNotesInput) Path() string {
	return "/api/v4/projects/:project_id/merge_requests/:merge_request_iid/notes"
}

func (g *GitlabV4) GetMergeRequestsNotes(ctx context.Context, input MergeReqeustNotesInput) ([]*MergeRequestNotes, error) {

	notes := []*MergeRequestNotes{}
	_, err := g.c.Do(ctx, &input, g.metas).Into(&notes)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "get Merge Request's notes failed: %v", err)
	}

	return notes, nil
}

// AcceptMergeRequestInput apply a merge request
// https://docs.gitlab.com/ee/api/merge_requests.html#accept-mr
type AcceptMergeRequestInput struct {
	httpx.MethodPut
	ProjectID       int                        `name:"project_id" in:"path"`
	MergeRequestIID int                        `name:"merge_request_iid" in:"path"`
	Options         AcceptMergetRequestOptions `in:"body"`
}
type AcceptMergetRequestOptions struct {
	MergeCommitMessage        string `json:"merge_commit_message,omitempty"`
	SquashCommitMessage       string `json:"squash_commit_message,omitempty"`
	Squash                    bool   `json:"squash,omitempty" default:"true"`
	ShouldRemoveSourceBranch  bool   `json:"should_remove_source_branch" default:"true"`
	MergeWhenPipelineSucceeds bool   `json:"merge_when_pipeline_succeeds,omitempty"`
	Sha                       string `json:"sha,omitempty"`
}

func (AcceptMergeRequestInput) Path() string {
	return "/api/v4/projects/:project_id/merge_requests/:merge_request_iid/merge"
}

func (g *GitlabV4) AcceptMergeRequest(ctx context.Context, input AcceptMergeRequestInput) (*MergeRequest, error) {

	mr := &MergeRequest{}
	_, err := g.c.Do(ctx, &input, g.metas).Into(&mr)

	if err != nil {
		return nil, pkgerrors.Wrapf(err, "accpet Merge Request failed: %v", err)
	}

	return mr, nil
}

// UpdateMergeRequestInput 更新 Merge Request
type UpdateMergeRequestInput struct {
	httpx.MethodPut
	ProjectID       int                            `name:"project_id" in:"path"`
	MergeRequestIID int                            `name:"merge_request_iid" in:"path"`
	Options         UpdateMergeRequestInputOptions `in:"body"`
}
type UpdateMergeRequestInputOptions struct {
	Title        string `json:"title,omitempty"`
	AssigneeID   int    `json:"assignee_id,omitempty"`
	AssigneeIDs  []int  `json:"assignee_ids,omitempty"`
	TargetBranch string `json:"target_branch,omitempty"`
	ReviewersIDs []int  `json:"reviewers_ids,omitempty"`
}

func (UpdateMergeRequestInput) Path() string {
	return "/api/v4/projects/:project_id/merge_requests/:merge_request_iid"
}

func (g *GitlabV4) UpdateMergeRequest(input UpdateMergeRequestInput) (*MergeRequest, error) {
	mr := &MergeRequest{}
	ctx := context.Background()
	_, err := g.c.Do(ctx, &input, g.metas).Into(mr)
	if err != nil {
		return nil, errors.Wrapf(err, "gitlab update merge request failed: %v", err)
	}

	return mr, nil
}
