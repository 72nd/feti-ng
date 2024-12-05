package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

type DeploymentWatcher struct {
	fsnotify.Watcher
	deploy Deploy
}

func NewDeploymentWatcher(deploy Deploy) (*DeploymentWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer watcher.Close()

	err = watcher.Add(deploy.Config.Path(deploy.Config.AssetsDir))
	if err != nil {
		return nil, err
	}

	return &DeploymentWatcher{
		Watcher: *watcher,
		deploy:  deploy,
	}, nil
}

func (w DeploymentWatcher) Run() {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				// TODO Handle
				return
			}
			w.handleChanges(event)
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			fmt.Println(err)
		}
	}
}

func (w DeploymentWatcher) handleChanges(event fsnotify.Event) error {
	return nil
}
