package confgitlab

import (
	"bytes"
	"context"

	"github.com/go-courier/httptransport/httpx"
	pkgerrors "github.com/pkg/errors"
)

type GetRepositoryFileRawInput struct {
	httpx.MethodGet
	ProjectID int    `name:"project_id" in:"path"`
	FilePath  string `name:"file_path" in:"path"`
	Ref       string `name:"ref" in:"query"`
}

func (GetRepositoryFileRawInput) Path() string {
	return "/api/v4/projects/:project_id/repository/files/:file_path/raw"
}

func (g *GitlabV4) GetRepositoryFileRaw(ctx context.Context, input GetRepositoryFileRawInput) ([]byte, error) {

	buf := bytes.NewBuffer(nil)
	_, err := g.c.Do(ctx, &input, g.metas).Into(buf)
	if err != nil {
		err = pkgerrors.Wrapf(err, "get gitlab file raw failed: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil

}
