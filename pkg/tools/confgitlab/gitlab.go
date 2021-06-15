package confgitlab

import (
	"context"
	"time"

	"github.com/go-courier/courier"
	"github.com/go-courier/httptransport/client"
	"github.com/go-courier/logr"
	"github.com/tangx/gitlab-mergebot/pkg/tools/datatypes"
)

type GitlabV4 struct {
	Endpoint     string        `env:""`
	PrivateToken string        `env:""`
	Timeout      time.Duration `env:""`
	c            *client.Client
	metas        courier.Metadata
	ctx          context.Context
}

func New(ctx context.Context, endpoint string, privateToken string) *GitlabV4 {
	log := logr.StdLogger().WithValues("conftool", "gitlabv4")
	ctx = logr.WithLogger(ctx, log)

	g := &GitlabV4{
		Endpoint:     endpoint,
		PrivateToken: privateToken,
		ctx:          ctx,
	}
	g.Init()

	return g
}

func (g *GitlabV4) Init() {
	ep, err := datatypes.ParseEndpoint(g.Endpoint)

	if err != nil {
		panic(err)
	}

	if g.Timeout == 0 {
		g.Timeout = 5 * time.Second
	}

	c := &client.Client{
		Protocol: ep.Scheme,
		Host:     ep.Hostname,
		Port:     ep.Port,
		Timeout:  g.Timeout,
	}

	c.SetDefaults()
	g.c = c
	g.metas = make(courier.Metadata)
	g.metas.Set("PRIVATE-TOKEN", g.PrivateToken)
}
