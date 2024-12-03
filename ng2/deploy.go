package main

type Deploy struct {
	Config      Config
	OutputDir   string
	LiveServe   bool
	RebuildSass bool
}

func (d Deploy) Run() {
}
