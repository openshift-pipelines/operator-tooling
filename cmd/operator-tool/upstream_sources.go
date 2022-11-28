package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	// "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"sigs.k8s.io/yaml"
)

type upstreamSources struct {
	Sources []source `json:"git"`
}

/*
  - automerge: 'never'
    update_policy: 'static'
    branch: main
    commit: 1ed8154d2fa5e2cfc07af10930182f854f946080
    url: https://github.com/openshift-pipelines/pipelines-as-code
    dest_formats:
      branch:
        gen_source_repos: true
        push_url: https://gitlab.cee.redhat.com/tekton/pipelines-as-code
*/
type source struct {
	Branch       string `json:"branch"`
	Commit       string `json:"commit"`
	URL          string `json:"url"`
	Automerge    string `json:"automerge"`
	UpdatePolicy string `json:"update_policy"`
	// FIXME: support DestFormats
}

func generateUpstreamSources(filename string) error {
	components, err := readCompoments(filename)
	if err != nil {
		return err
	}
	us := &upstreamSources{Sources: []source{}}
	for _, component := range components {
		url := "https://github.com/" + component.Github
		branch := component.Version
		if branch == "nightly" {
			branch = "main" // FIXME: we may want to use remote head here… but today we assume it's main
		}
		commitHash := ""

		repo, err := git.Init(memory.NewStorage(), nil)
		if err != nil {
			return err
		}
		rem, err := repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{url},
		})
		if err != nil {
			return err
		}
		refs, err := rem.List(&git.ListOptions{})
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Fetching %s : %s...\n", component.Github, branch)
		for _, r := range refs {
			if !r.Name().IsTag() && !r.Name().IsBranch() {
				continue
			}
			if r.Name().Short() == branch {
				commitHash = r.Hash().String()
				if r.Name().IsTag() {
					err = repo.Fetch(&git.FetchOptions{
						RemoteName: "origin",
					})
					if err != nil {
						return err
					}
					tags, err := repo.TagObjects()
					if err != nil {
						return err
					}
					err = tags.ForEach(func(t *object.Tag) error {
						if t.Name == branch {
							commit, err := t.Commit()
							if err != nil {
								return err
							}
							commitHash = commit.Hash.String()
						}
						return nil
					})
					if err != nil {
						return err
					}
				}
				break
			}
		}
		source := source{
			Automerge:    "never",
			UpdatePolicy: "static",
			URL:          url,
			Branch:       branch,
			Commit:       commitHash,
		}
		us.Sources = append(us.Sources, source)
	}
	data, err := yaml.Marshal(us)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
