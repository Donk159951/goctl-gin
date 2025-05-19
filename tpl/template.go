package tpl

import _ "embed"

var (
	//go:embed routes.tpl
	RoutesTemplate string

	//go:embed routes_setup.tpl
	RoutesSetupTemplate string
)
