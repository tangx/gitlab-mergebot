package mergebot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-courier/logr"
	"github.com/pkg/errors"
	"github.com/tangx/gitlab-mergebot/pkg/tools/confgitlab"
	"gopkg.in/yaml.v3"
)

// MergeBot gitlab 合并机器人
type MergeBot struct {
	ctx    context.Context
	gitlab *confgitlab.GitlabV4
}

func New(ctx context.Context, gitlab *confgitlab.GitlabV4) *MergeBot {
	log := logr.StdLogger().WithValues("app", "mergebot")

	if ctx == nil {
		ctx = context.Background()
	}
	ctx = logr.WithLogger(ctx, log)

	bot := &MergeBot{
		gitlab: gitlab,
		ctx:    ctx,
	}
	return bot
}

func (b *MergeBot) Run() {
	for {
		b.PullMrs()
		time.Sleep(10 * time.Second)
	}
}

func (b *MergeBot) PullMrs() {
	log := logr.FromContext(b.ctx)

	mrs, err := b.gitlab.ListMergeReqeusts()
	if err != nil {
		log.Error(err)
		return
	}

	for _, mr := range mrs {

		if !confgitlab.IsVaildMRCandidate(mr) {
			log.Debug("[%d] %s is not a vaild mr cadidate", mr.IID, mr.Title)
			continue
		}

		// 获取评论
		notes, err := b.PullMRNotes(mr)
		if err != nil {
			log.Error(err)
			continue
		}

		// 获取 merge config
		mbconfig, err := b.ReadMergeBotConfig(mr)
		if err != nil {
			log.Error(err)
			continue
		}
		mbconfig.initial()

		if !b.isReadyBeMerged(notes, mbconfig) {
			log.Debug("[%d] %s is not ready", mr.IID, mr.Title)
			continue
		}

		// 如果有配置 assignee 则转交
		if len(mbconfig.Assignees) > 0 {
			assignee := mbconfig.RandomAssignee()

			if err := b.TransferAssigneeByName(mr, assignee); err != nil {
				log.Error(err)
				continue
			}

			log.Info("transfer assignee success")
			continue
		}

		// merge
		err = b.AcceptMR(mr, mbconfig)
		if err != nil {
			log.Error(err)
		}
		log.Info("Accept MR succes: %s %s", mr.References.Full, mr.Title)
	}
}

// isReadyBeMerged 是否满足评论数量
func (b *MergeBot) isReadyBeMerged(notes []*confgitlab.MergeRequestNotes, mbconfig *MergebotConfig) bool {

	log := logr.FromContext(b.ctx)
	validReviewers := map[string]struct{}{}

	for _, note := range notes {
		username := note.Author.Username
		// 用户是否在 Reviewers 列表中，且内容是否符合预期
		if mbconfig.HasReviewer(username) && note.Body == "@gitbot lgtm" {
			validReviewers[username] = struct{}{}
		}
	}

	// 是否满足最小审阅数字
	now := len(validReviewers)
	if now < mbconfig.MinReviewers {
		log.Debug("only %d valid reviewer(s), want %d MinReviewers", now, mbconfig.MinReviewers)
		return false
	}

	return true
}

func (b *MergeBot) CloseMR() error {

	return nil
}

// AcceptMR 合并
func (b *MergeBot) AcceptMR(mr *confgitlab.MergeRequest, mbconfig *MergebotConfig) error {
	squashCommit := fmt.Sprintf("%s\n\nSee merge request %s", mr.Description, mr.References.Full)
	input := confgitlab.AcceptMergeRequestInput{
		ProjectID:       mr.ProjectID,
		MergeRequestIID: mr.IID,
		Options: confgitlab.AcceptMergetRequestOptions{
			Squash:                    mbconfig.Squash,
			SquashCommitMessage:       squashCommit,
			MergeCommitMessage:        mr.Title,
			ShouldRemoveSourceBranch:  mbconfig.ShouldRemoveSourceBranch,
			MergeWhenPipelineSucceeds: mbconfig.MergeWhenPipelineSucceeds,
		},
	}
	mr, err := b.gitlab.AcceptMergeRequest(b.ctx, input)
	if err != nil {
		return err
	}

	if mr.State != "merged" {
		return errors.New("accpet mr failed")
	}

	return nil
}

// PullMRNotes 获取评论
func (b *MergeBot) PullMRNotes(mr *confgitlab.MergeRequest) ([]*confgitlab.MergeRequestNotes, error) {
	input := confgitlab.MergeReqeustNotesInput{
		ProjectID:       mr.ProjectID,
		MergeRequestIID: mr.IID,
	}
	notes, err := b.gitlab.GetMergeRequestsNotes(b.ctx, input)
	return notes, err
}

// ReadMergeBotConfig
func (b *MergeBot) ReadMergeBotConfig(mr *confgitlab.MergeRequest) (*MergebotConfig, error) {
	input := confgitlab.GetRepositoryFileRawInput{
		ProjectID: mr.ProjectID,
		FilePath:  `.mergebot.yml`,
		Ref:       mr.TargetBranch,
	}

	// log := logr.FromContext(b.ctx)

	body, err := b.gitlab.GetRepositoryFileRaw(b.ctx, input)
	if err != nil {
		// log.Error(err)
		return nil, err
	}

	config := MergebotConfig{}
	err = yaml.Unmarshal(body, &config)
	if err != nil {
		// log.Error(err)
		return nil, err
	}

	return &config, nil
}

// TransferAssignee
func (b *MergeBot) TransferAssigneeByName(mr *confgitlab.MergeRequest, assignee string) error {

	// check new assignee exists
	users, err := b.gitlab.GetUserByName(assignee)
	if err != nil {
		return err
	}

	if len(users) != 1 {
		return fmt.Errorf("transfer assignee failed, invalid assignee, get %d user(s), want only 1", len(users))
	}

	user := users[0]
	// update mr
	input := confgitlab.UpdateMergeRequestInput{
		ProjectID:       mr.ProjectID,
		MergeRequestIID: mr.IID,
		Options: confgitlab.UpdateMergeRequestInputOptions{
			AssigneeID: user.ID,
		},
	}

	mr, err = b.gitlab.UpdateMergeRequest(input)
	if err != nil {
		return err
	}

	if mr.Assignee.Username != assignee {
		return fmt.Errorf("transfer assignee failed, got new assignee %s , want %s", mr.Assignee.Username, assignee)
	}

	return nil
}
