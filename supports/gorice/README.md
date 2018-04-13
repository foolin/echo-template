# Basic
Support template for [go.rice](github.com/GeertJohan/go.rice/rice)


# Install
```go
go get github.com/foolin/echo-template/supports/gorice
```

# Useage

```go
 echo.Renderer = gorice.New(rice.MustFindBox("views"))
```

# Example
```go

func main() {

	// Echo instance
	e := echo.New()

	// servers other static files
	staticBox := rice.MustFindBox("static")
	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
	e.GET("/static/*", echo.WrapHandler(staticFileServer))

	//Set Renderer
	e.Renderer = gorice.New(rice.MustFindBox("views"))

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}

```

[gorice example](https://github.com/foolin/echo-template/tree/master/examples/gorice)

# Links

[echo template](https://github.com/foolin/gin-template)

[echo framework](https://github.com/labstack/echo)

[go.rice](https://github.com/GeertJohan/go.rice)
