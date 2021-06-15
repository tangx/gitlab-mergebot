package mergebot

import (
	"math/rand"
	"time"
)

type MergebotConfig struct {
	MinReviewers              int      `yaml:"minReviewers"`
	Assignees                 []string `yaml:"assignees,omitempty"`
	Reviewers                 []string `yaml:"reviewers"`
	reviewersSet              map[string]bool
	Squash                    bool `yaml:"squash,omitempty"`
	ShouldRemoveSourceBranch  bool `yaml:"shouldRemoveSourceBranch,omitempty"`
	MergeWhenPipelineSucceeds bool `yaml:"mergeWhenPipelineSucceeds,omitempty"`
}

func (mbc *MergebotConfig) initial() {
	set := make(map[string]bool)
	for _, reviewer := range mbc.Reviewers {
		set[reviewer] = true
	}
	mbc.reviewersSet = set
}
func (mbc *MergebotConfig) HasReviewer(name string) bool {
	return mbc.reviewersSet[name]
}

func (mbc *MergebotConfig) RandomAssignee() string {
	rand.Seed(time.Now().Unix())

	i := rand.Intn(len(mbc.Assignees))
	return mbc.Assignees[i]
}
