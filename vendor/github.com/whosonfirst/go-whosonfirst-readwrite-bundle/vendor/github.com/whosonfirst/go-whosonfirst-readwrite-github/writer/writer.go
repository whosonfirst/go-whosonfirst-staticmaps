package writer

import (
	"errors"
	wof_writer "github.com/whosonfirst/go-whosonfirst-readwrite/writer"
	"io"
)

type GitHubWriter struct {
	wof_writer.Writer
}

func NewGitHubWriter(root string) (wof_writer.Writer, error) {

	wr := GitHubWriter{}

	return &wr, nil
}

func (wr *GitHubWriter) Write(path string, fh io.ReadCloser) error {
	return errors.New("Please write me")
}

func (wr *GitHubWriter) URI(path string) string{
     return ""
}
