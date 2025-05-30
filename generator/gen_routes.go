package generator

import (
	"goctl-gin/prepare"
	"goctl-gin/tpl"
	"os"
	"path"
	"strings"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func GenRoutes() error {
	for _, group := range prepare.ApiSpec.Service.Groups {
		subDir := group.GetAnnotation(groupProperty)
		subDir, err := format.FileNamingFormat(dirStyle, subDir)
		if err != nil {
			return err
		}

		routesPkg := path.Join("routes", subDir)
		routesBase := path.Base(routesPkg)

		os.Remove(path.Join(prepare.OutputDir, routesPkg, "routes.go"))

		// handle
		handlePkg := path.Join("handler", subDir)
		handleBase := path.Base(handlePkg)

		// prefix
		prefix := group.GetAnnotation(spec.RoutePrefixKey)

		// middlewares
		var middlewares []string
		if len(group.GetAnnotation("jwt")) > 0 {
			middlewares = append(middlewares, group.GetAnnotation("jwt"))
		}

		if len(group.GetAnnotation("middleware")) > 0 {
			middlewares = append(middlewares, strings.Split(group.GetAnnotation("middleware"), ",")...)
		}

		middlewares = lo.Map(middlewares, func(item string, index int) string {
			res := strings.TrimSuffix(item, "Middleware") + "Middleware"
			return cases.Title(language.English, cases.NoLower).String(res)
		})

		// route
		var routes []map[string]string
		for _, r := range group.Routes {
			routes = append(routes, map[string]string{
				"method": strings.ToUpper(r.Method),
				"path":   r.Path,
				"handle": cases.Title(language.English, cases.NoLower).String(r.Handler),
			})
		}

		err = GenFile(
			"routes.go",
			tpl.RoutesTemplate,
			WithSubDir(routesPkg),
			WithData(map[string]any{
				"rootPkg":       prepare.RootPkg,
				"pkgName":       routesBase,
				"handlePkg":     handlePkg,
				"handleBase":    handleBase,
				"prefix":        prefix,
				"hasPrefix":     len(prefix) > 0,
				"hasMiddleware": len(middlewares) > 0,
				"middleware":    middlewares,
				"funcName":      cases.Title(language.English, cases.NoLower).String(group.GetAnnotation(groupProperty)),
				"routes":        routes,
			}),
		)

		if err != nil {
			return err
		}
	}
	return genSetup()
}

func genSetup() error {
	os.Remove(path.Join(prepare.OutputDir, "routes/setup.go"))

	var routes []map[string]string
	for _, group := range prepare.ApiSpec.Service.Groups {
		subDir := group.GetAnnotation(groupProperty)
		subDir, err := format.FileNamingFormat(dirStyle, subDir)
		if err != nil {
			return err
		}

		routesPkg := path.Join("routes", subDir)
		routesBase := path.Base(routesPkg)

		name := cases.Title(language.English, cases.NoLower).String(group.GetAnnotation(groupProperty))
		routes = append(routes, map[string]string{
			"pkg":  routesPkg,
			"base": routesBase,
			"name": name,
			"alias": func() string {
				if len(name) == 0 {
					return ""
				}
				return strings.ToLower(string(name[0])) + name[1:]
			}(),
		})

	}

	return GenFile(
		"setup.go",
		tpl.RoutesSetupTemplate,
		WithSubDir("routes"),
		WithData(map[string]any{
			"rootPkg": prepare.RootPkg,
			"routes":  routes,
		}),
	)
}
