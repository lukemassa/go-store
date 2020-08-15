package store

import (
	"context"
	"time"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// GitRepoStore a store to hold the value
type GitRepoStore struct {
	store    Store
	worktree *git.Worktree
	gitURL   string
}

// NewGitRepo store the result of an http get call
func NewGitRepo(gitURL string, config Config) GitRepoStore {
	ret := GitRepoStore{
		gitURL: gitURL,
	}
	store := New(func(ctx context.Context) (interface{}, error) {
		err := ret.setupRepo(ctx)
		if err != nil {
			return nil, err
		}

		err = ret.worktree.PullContext(ctx, &git.PullOptions{})

		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}
		return ret.worktree, nil
	}, config)
	ret.store = store

	return ret
}

func (s *GitRepoStore) setupRepo(ctx context.Context) error {
	if s.worktree != nil {
		return nil
	}

	fs := memfs.New()
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: s.gitURL,
	})
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	s.worktree = worktree
	return nil
}

// Get gets the result of the http call
func (s *GitRepoStore) Get() (*git.Worktree, error) {
	res, err := s.store.Get()
	if err != nil {
		return nil, err
	}
	return res.(*git.Worktree), err
}

// Wait until value is available
func (s *GitRepoStore) Wait(maxWait time.Duration) error {
	return s.store.Wait(maxWait)
}
