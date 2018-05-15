package gorice

import (
	"github.com/foolin/echo-template"
	"github.com/GeertJohan/go.rice"
)

/**
New echo template engine, default views root.
 */
func New(viewsRootBox *rice.Box) *echotemplate.TemplateEngine {
	return NewWithConfig(viewsRootBox, echotemplate.DefaultConfig)
}

/**
New echo template engine
Important!!! The viewsRootBox's name and config.Root must be consistent.
 */
func NewWithConfig(viewsRootBox *rice.Box, config echotemplate.TemplateConfig) *echotemplate.TemplateEngine {
	config.Root = viewsRootBox.Name()
	engine := echotemplate.New(config)
	engine.SetFileHandler(FileHandler(viewsRootBox))
	return engine
}

func FileHandler(viewsRootBox *rice.Box) echotemplate.FileHandler {
	return func(config echotemplate.TemplateConfig, tplFile string) (content string, err error) {
		// get file contents as string
		return viewsRootBox.String(tplFile + config.Extension)
	}
}
