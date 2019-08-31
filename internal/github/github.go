// Copyright 2019 Francisco Souza. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package github

import (
	"context"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

func NewClient(token string) *Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	hc := oauth2.NewClient(context.Background(), tokenSource)
	return &Client{client: github.NewClient(hc)}
}

func (c *Client) LoadRepoInfo(ctx context.Context, owner, name string) (*github.Repository, error) {
	repo, _, err := c.client.Repositories.Get(ctx, owner, name)
	return repo, err
}
