package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"sigs.k8s.io/yaml"
)

type upstream_sources struct {
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
	Url          string `json:"url"`
	Automerge    string `json:"automerge"`
	UpdatePolicy string `json:"update_policy"`
	// FIXME: support DestFormats
}

func generateUpstreamSources(filename string) error {
	components, err := readCompoments(filename)
	if err != nil {
		return err
	}
	us := &upstream_sources{Sources: []source{}}
	for name, component := range components {
		fmt.Println("component", name, component)
		url := "https://github.com/" + component.Github
		branch := component.Version
		commit := ""

		rem := git.NewRemote(memory.NewStorage(), &gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{url},
		})
		refs, err := rem.List(&git.ListOptions{})
		if err != nil {
			return err
		}
		for _, r := range refs {
			if !r.Name().IsTag() {
				continue
			}
			if r.Name().Short() == component.Version {
				commit = r.Hash().String()
				break
			}
		}
		source := source{
			Automerge:    "never",
			UpdatePolicy: "static",
			Url:          url,
			Branch:       branch,
			Commit:       commit,
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
