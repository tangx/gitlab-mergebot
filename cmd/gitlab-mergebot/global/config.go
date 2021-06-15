package global

import (
	"context"
	"os"

	"github.com/tangx/gitlab-mergebot/pkg/mergebot"
	"github.com/tangx/gitlab-mergebot/pkg/tools/confgitlab"
)

var (
	gitlab_Endpoint     = os.Getenv("GITLAB_Endpoint")
	gitlab_PrivateToken = os.Getenv("GITLAB_PrivateToken")
)

var (
	ctx = context.Background()

	gitlab = confgitlab.New(ctx, gitlab_Endpoint, gitlab_PrivateToken)

	mbot = mergebot.New(ctx, gitlab)
)

func MergebotRun() {
	mbot.Run()
}
