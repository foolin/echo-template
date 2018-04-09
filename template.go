/*
 * Copyright 2018 Foolin.  All rights reserved.
 *
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 *
 */

/*
Golang template for echo framework, Use golang html/template syntax,
Easy and simple to use for echo framework, See https://github.com/foolin/echo-template
for more information.
*/
package echotemplate

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"github.com/labstack/echo"
)

const templateEngineKey = "echo-template_templateEngine"

type TemplateEngine struct {
	config   TemplateConfig
	tplMap   map[string]*template.Template
	tplMutex sync.RWMutex
}

type TemplateConfig struct {
	Root         string           //view root
	Extension    string           //template extension
	Master       string           //template master
	Partials     []string         //template partial, such as head, foot
	Funcs        template.FuncMap //template functions
	DisableCache bool             //disable cache, debug mode
	Delims       Delims           //delimeters
}

type Delims struct {
	Left  string
	Right string
}

func New(config TemplateConfig) *TemplateEngine {
	return &TemplateEngine{
		config:   config,
		tplMap:   make(map[string]*template.Template),
		tplMutex: sync.RWMutex{},
	}
}

func Default() *TemplateEngine {
	return New(TemplateConfig{
		Root:         "views",
		Extension:    ".html",
		Master:       "layouts/master",
		Partials:     []string{},
		Funcs:        make(template.FuncMap),
		DisableCache: false,
		Delims:       Delims{Left: "{{", Right: "}}"},
	})
}

func (e *TemplateEngine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return e.executeRender(w, name, data)
}

func (e *TemplateEngine) executeRender(out io.Writer, name string, data interface{}) error {
	useMaster := true
	if filepath.Ext(name) == e.config.Extension {
		useMaster = false
		name = strings.TrimRight(name, e.config.Extension)

	}
	return e.executeTemplate(out, name, data, useMaster)
}

func (e *TemplateEngine) executeTemplate(out io.Writer, name string, data interface{}, useMaster bool) error {
	var tpl *template.Template
	var err error
	var ok bool

	allFuncs := make(template.FuncMap, 0)
	allFuncs["include"] = func(layout string) (template.HTML, error) {
		buf := new(bytes.Buffer)
		err := e.executeTemplate(buf, layout, data, false)
		return template.HTML(buf.String()), err
	}

	// Get the plugin collection
	for k, v := range e.config.Funcs {
		allFuncs[k] = v
	}

	e.tplMutex.RLock()
	tpl, ok = e.tplMap[name]
	e.tplMutex.RUnlock()

	exeName := name
	if useMaster && e.config.Master != "" {
		exeName = e.config.Master
	}

	if !ok || e.config.DisableCache {
		tplList := make([]string, 0)
		if useMaster {
			//render()
			if e.config.Master != "" {
				tplList = append(tplList, e.config.Master)
			}
		}
		tplList = append(tplList, name)
		tplList = append(tplList, e.config.Partials...)

		// Loop through each template and test the full path
		tpl = template.New(name).Funcs(allFuncs).Delims(e.config.Delims.Left, e.config.Delims.Right)
		for _, v := range tplList {
			// Get the absolute path of the root template
			path, err := filepath.Abs(e.config.Root + string(os.PathSeparator) + v + e.config.Extension)
			if err != nil {
				return fmt.Errorf("TemplateEngine path:%v error: %v", path, err)
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("TemplateEngine render read name:%v, path:%v, error: %v", v, path, err)
			}
			var tmpl *template.Template
			if v == name {
				tmpl = tpl
			} else {
				tmpl = tpl.New(v)
			}
			_, err = tmpl.Parse(string(data))
			if err != nil {
				return fmt.Errorf("TemplateEngine render parser name:%v, path:%v, error: %v", v, path, err)
			}
		}
		e.tplMutex.Lock()
		e.tplMap[name] = tpl
		e.tplMutex.Unlock()
	}

	// Display the content to the screen
	err = tpl.Funcs(allFuncs).ExecuteTemplate(out, exeName, data)
	if err != nil {
		return fmt.Errorf("TemplateEngine execute template error: %v", err)
	}

	return nil
}

// You should use helper func `Middleware()` to set the supplied
// TemplateEngine and make `Render()` work validly.
func Render(ctx echo.Context, code int, name string, data interface{}) error {
	if val := ctx.Get(templateEngineKey); val != nil {
		if e, ok := val.(*TemplateEngine); ok {
			return e.Render(ctx.Response().Writer, name, data, ctx)
		}
	}
	return ctx.Render(code, name, data)
}

//New gin middleware for func `gintemplate.Render()`
func NewMiddleware(config TemplateConfig) echo.MiddlewareFunc {
	return Middleware(New(config))
}

func Middleware(e *TemplateEngine) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(templateEngineKey, e)
			return next(c)
		}
	}
}
