/*
 * Copyright 2018 Foolin.  All rights reserved.
 *
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 *
 */

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
		})
	})

	e.GET("/block", func(c echo.Context) error {
		//render only file, must full name with extension
		return c.Render(http.StatusOK, "block", echo.Map{"title": "Block file title!!"})
	})

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}
