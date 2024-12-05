package main

import (
	"fmt"
	"path/filepath"

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

	moduleDir, err := ModuleDir()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(filepath.Join(moduleDir, "static"))
	if err != nil {
		return nil, err
	}
	err = watcher.Add(filepath.Join(moduleDir, "tpl"))
	if err != nil {
		return nil, err
	}
	err = watcher.Add(deploy.Config.configDir)
	if err != nil {
		return nil, err
	}

	return &DeploymentWatcher{
		Watcher: *watcher,
		deploy:  deploy,
	}, nil
}

func (w *DeploymentWatcher) Run() {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					// TODO: Handle properly
					fmt.Println("will return")
					return
				}
				w.handleEvent(event)
			case err, ok := <-w.Errors:
				// TODO: Handle properly
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	<-done
}

func (w DeploymentWatcher) handleEvent(event fsnotify.Event) error {
	fmt.Printf("file '%s' has changed", event.Name)
	w.deploy.Build()
	return nil
}
