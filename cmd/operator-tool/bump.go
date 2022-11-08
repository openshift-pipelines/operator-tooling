// Package main is the main package :D
package main

import (
	"sort"
)

func bump(filename string, bugfix bool) error {
	newComponents := map[string]component{}
	components, err := readCompoments(filename)
	if err != nil {
		return err
	}
	for name, component := range components {
		newComponent, err := bumpComponent(name, component, bugfix)
		if err != nil {
			return err
		}
		newComponents[name] = newComponent
	}
	return writeComponents(filename, newComponents)
}

func bumpComponent(name string, c component, bugfix bool) (component, error) {
	newVersion := c.Version
	newerVersions, err := checkComponentNewerVersions(c, bugfix)
	if err != nil {
		return component{}, err
	}
	if len(newerVersions) > 0 {
		// Get the latest one
		sort.Sort(newerVersions) // sort just in case
		newVersion = "v" + newerVersions[len(newerVersions)-1].String()
	}
	return component{
		Github:  c.Github,
		Version: newVersion,
	}, nil
}
