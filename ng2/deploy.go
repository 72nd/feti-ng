package main

type Deploy struct {
	Config    Config
	OutputDir string
	LiveServe bool
}

func (d Deploy) Run() {
}
