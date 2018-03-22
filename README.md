# echo-template

[![GoDoc](https://godoc.org/github.com/foolin/echo-template?status.png)](https://godoc.org/github.com/foolin/echo-template)

Golang template for [echo framework](https://github.com/labstack/echo)!


# Feature
* Easy and simple to use for echo framework.
* Use golang html/template syntax.
* Support configure master layout file.
* Support configure template file extension.
* Support configure templates directory.
* Support configure cache template.
* Support include file.
* Support dynamic reload template(disable cache mode).
* Support multiple templates for fontend and backend.

# Docs
See https://www.godoc.org/github.com/foolin/echo-template

# Install
```bash
go get github.com/foolin/echo-template
```

# Usage
```go
package main

import (
	"github.com/labstack/echo"
	"github.com/foolin/echo-template"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Renderer = echotemplate.Default()

	e.GET("/page", func(c echo.Context) error {
		//render only file, must full name with extension
		return c.Render(http.StatusOK, "page.html", echo.Map{"title": "Page file title!!"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}
```

# Configure

```go
    TemplateConfig{
		Root:      "views", //template root path
		Extension: ".tpl", //file extension
		Master:    "layouts/master", //master layout file
		Partials:  []string{"partials/head"}, //partial files
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			// more funcs
		},
		DisableCache: false, //if disable cache, auto reload template file for debug.
	}
```


# Render

### Render with master
The `ctx` is instance of  `echo.Context`
```go
//use name without extension `.html`
ctx.Render(http.StatusOK, "index", echo.Map{})
```

### Render only file(not use master layout)
```go
//use full name with extension `.html`
ctx.Render(http.StatusOK, "page.html", echo.Map{})
```


# Include syntax
```go
//template file
{{include "layouts/footer"}}
```

# Examples

### Basic example
```go

package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/foolin/echo-template"
)

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = echotemplate.Default()

	// Routes
	e.GET("/", func(c echo.Context) error {
		//render with master
		return c.Render(http.StatusOK, "index", echo.Map{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	e.GET("/page", func(c echo.Context) error {
		//render only file, must full name with extension
		return c.Render(http.StatusOK, "page.html", echo.Map{"title": "Page file title!!"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}


```

Project structure:
```go
|-- app/views/
    |--- index.html          
    |--- page.html
    |-- layouts/
        |--- footer.html
        |--- master.html
    

See in "examples/basic" folder
```

[Basic example](https://github.com/foolin/echo-template/tree/master/examples/basic)

### Advance example
```go

package main

import (
	"net/http"
	"html/template"
	"time"
	"github.com/labstack/echo"
	"github.com/foolin/echo-template"
	"github.com/labstack/echo/middleware"
)

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = echotemplate.New(echotemplate.TemplateConfig{
		Root:      "views",
		Extension: ".tpl",
		Master:    "layouts/master",
		Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			"copy": func() string{
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	e.GET("/", func(c echo.Context) error {
		//render with master
		return c.Render(http.StatusOK, "index", echo.Map{
			"title": "Index title!",
			"add": func(a int, b int) int {
				return a + b
			},
		})
	})

	e.GET("/page", func(c echo.Context) error {
		//render only file, must full name with extension
		return c.Render(http.StatusOK, "page.tpl", echo.Map{"title": "Page file title!!"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}

```

Project structure:
```go
|-- app/views/
    |--- index.tpl          
    |--- page.tpl
    |-- layouts/
        |--- footer.tpl
        |--- head.tpl
        |--- master.tpl
    |-- partials/
        |--- ad.tpl
    

See in "examples/advance" folder
```

[Advance example](https://github.com/foolin/echo-template/tree/master/examples/advance)

### Multiple example
```go

package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/foolin/echo-template"
)

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//new template engine
	e.Renderer = echotemplate.New(echotemplate.TemplateConfig{
		Root:      "views/fontend",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	e.GET("/", func(ctx echo.Context) error {
		// `HTML()` is a helper func to deal with multiple TemplateEngine's.
		// It detects the suitable TemplateEngine for each path automatically.
		return echotemplate.Render(ctx, http.StatusOK, "index", echo.Map{
			"title": "Fontend title!",
		})
	})

	//=========== Backend ===========//

	//new middleware
	mw := echotemplate.NewMiddleware(echotemplate.TemplateConfig{
		Root:      "views/backend",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	// You should use helper func `Middleware()` to set the supplied
	// TemplateEngine and make `HTML()` work validly.
	backendGroup := e.Group("/admin", mw)

	backendGroup.GET("/", func(ctx echo.Context) error {
		// With the middleware, `HTML()` can detect the valid TemplateEngine.
		return echotemplate.Render(ctx, http.StatusOK, "index", echo.Map{
			"title": "Backend title!",
		})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}


```

Project structure:
```go
|-- app/views/
    |-- fontend/
        |--- index.html
        |-- layouts/
            |--- footer.html
            |--- head.html
            |--- master.html
        |-- partials/
        |--- ad.html
    |-- backend/
        |--- index.html
        |-- layouts/
            |--- footer.html
            |--- head.html
            |--- master.html
        
See in "examples/multiple" folder
```

[Multiple example](https://github.com/foolin/echo-template/tree/master/examples/multiple)
