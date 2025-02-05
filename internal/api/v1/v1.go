package v1

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/api/v1/chart"
	"github.com/KubeOperator/kubepi/internal/api/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/api/v1/proxy"
	"github.com/KubeOperator/kubepi/internal/api/v1/role"
	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	"github.com/KubeOperator/kubepi/internal/api/v1/user"
	"github.com/KubeOperator/kubepi/internal/api/v1/webkubectl"
	"github.com/KubeOperator/kubepi/internal/api/v1/ws"
	v1Role "github.com/KubeOperator/kubepi/internal/model/v1/role"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	v1RoleService "github.com/KubeOperator/kubepi/internal/service/v1/role"
	v1RoleBindingService "github.com/KubeOperator/kubepi/internal/service/v1/rolebinding"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/KubeOperator/kubepi/pkg/collectons"
	"github.com/KubeOperator/kubepi/pkg/i18n"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

func authHandler() iris.Handler {
	return func(ctx *context.Context) {
		p := sessions.Get(ctx).Get("profile")
		if p == nil {
			ctx.Values().Set("message", "please login")
			ctx.StopWithStatus(iris.StatusUnauthorized)
			return
		}
		ctx.Values().Set("profile", p)
		ctx.Next()
	}
}

func langHandler() iris.Handler {
	return func(ctx *context.Context) {
		s := sessions.Get(ctx)
		p := s.Get("profile")
		u, ok := p.(session.UserProfile)
		lang := i18n.LanguageZhCN
		if ok {
			lang = u.Language
		}
		ctx.Values().Set("language", lang)
		ctx.Next()
	}
}

func pageHandler() iris.Handler {
	return func(ctx *context.Context) {
		if ctx.URLParamExists(pkgV1.PageSize) && ctx.URLParamExists(pkgV1.PageNum) {
			pageNum, err := ctx.URLParamInt(pkgV1.PageNum)
			if err != nil {
				ctx.Values().Set("message", fmt.Sprintf("page num format err %s", err.Error()))
				ctx.StopWithStatus(iris.StatusBadRequest)
				return
			}
			pageSize, err := ctx.URLParamInt(pkgV1.PageSize)
			if err != nil {
				ctx.Values().Set("message", fmt.Sprintf("page size format err %s", err.Error()))
				ctx.StopWithStatus(iris.StatusBadRequest)
				return
			}
			ctx.Values().Set(pkgV1.PageNum, pageNum)
			ctx.Values().Set(pkgV1.PageSize, pageSize)
		}
		ctx.Next()
	}
}

func resourceExtractHandler() iris.Handler {
	return func(ctx *context.Context) {
		path := ctx.Request().URL.Path
		ss := strings.Split(path, "/")
		// "" "api" "v1" "resource"
		if len(ss) >= 4 {
			ctx.Values().Set("resource", ss[3])
		}
		ctx.Next()
	}
}

func roleHandler() iris.Handler {
	return func(ctx *context.Context) {
		// 查询当前用户的角色
		// 查询角色的 rolebinding 获取 roles
		p := sessions.Get(ctx).Get("profile")
		u := p.(session.UserProfile)

		if u.IsAdministrator {
			ctx.Next()
			return
		}
		roleBindingService := v1RoleBindingService.NewService()
		rbs, err := roleBindingService.GetRoleBindingBySubject(v1Role.Subject{
			Kind: "User",
			Name: u.Name,
		}, common.DBOptions{})
		if err != nil {
			if !errors.As(err, &storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		roleNameHash := map[string]struct{}{}
		for i := range rbs {
			roleName := rbs[i].RoleRef
			roleNameHash[roleName] = struct{}{}
		}
		var roleNames []string
		for key := range roleNameHash {
			roleNames = append(roleNames, key)
		}

		roleService := v1RoleService.NewService()
		rs, err := roleService.GetByNames(roleNames, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		ctx.Values().Set("roles", rs)
		ctx.Next()
	}
}

func getVerbByRoute(path, method string) string {
	switch strings.ToLower(method) {
	case "put":
		return "update"
	case "delete":
		return "delete"
	case "get":
		if strings.Contains(path, "/:name") {
			return "get"
		} else {
			return "list"
		}
	case "post":
		if strings.HasSuffix(path, "search") {
			return "list"
		}
		return "create"

	}
	return ""
}

func apiResourceHandler(party iris.Party) iris.Handler {
	return func(ctx *context.Context) {
		//1. 确定所有的api资源有哪些
		apiBuilder := party.(*router.APIBuilder)
		routes := apiBuilder.GetRoutes()
		resourceMap := map[string]*collectons.StringSet{}
		for i := range routes {
			if strings.HasPrefix(routes[i].Path, "/api/v1/") {
				ss := strings.Split(routes[i].Path, "/")
				if len(ss) >= 4 {
					resourceName := ss[3]
					//过滤session资源
					if resourceName == "sessions" || resourceName == "proxy" || resourceName == "ws" || resourceName == "charts" || resourceName == "webkubectl" {
						continue
					}
					if _, ok := resourceMap[resourceName]; !ok {
						resourceMap[resourceName] = collectons.NewStringSet()
					}
					resourceMap[resourceName].Add(getVerbByRoute(routes[i].Path, routes[i].Method))
				}
			}
		}
		displayMap := map[string][]string{}
		for k, _ := range resourceMap {
			verbs := resourceMap[k]
			if len(verbs.ToSlice()) > 0 {
				displayMap[k] = verbs.ToSlice()
			}
		}
		ctx.Values().Set("data", displayMap)
	}
}

func roleAccessHandler() iris.Handler {
	return func(ctx *context.Context) {
		//// 查询角色的 resources
		//// 通过api resource 过滤出来资源主体,method 过滤操作
		p := sessions.Get(ctx).Get("profile")
		u := p.(session.UserProfile)
		if !strings.Contains(ctx.Request().URL.Path, "/proxy") && !strings.Contains(ctx.Request().URL.Path, "/ws") &&
			!strings.Contains(ctx.Request().URL.Path, "/charts") && !strings.Contains(ctx.Request().URL.Path, "/apps") {
			// 放通admin权限
			if u.IsAdministrator {
				ctx.Next()
				return
			}
			rs := ctx.Values().Get("roles")
			roles := rs.([]v1Role.Role)
			requestResource := ctx.Values().GetString("resource")
			if requestResource != "" {
				currentRoute := ctx.GetCurrentRoute()
				requestVerb := getVerbByRoute(currentRoute.Path(), currentRoute.Method())
				resourceMatched, methodMatch := matchRoles(requestResource, requestVerb, roles)
				if !(resourceMatched && methodMatch) {
					ctx.StopWithStatus(iris.StatusForbidden)
					ctx.Values().Set("message", []string{"user %s can not access resource %s %s", u.Name, requestResource, requestVerb})
					return
				}
			}
		}

		ctx.Next()
	}
}

func matchRoles(requestResource, requestMethod string, rs []v1Role.Role) (bool, bool) {
	resourceMatch := false
	methodMatch := false
	for i := range rs {
		for j := range rs[i].Rules {
			for k := range rs[i].Rules[j].Resource {
				if rs[i].Rules[j].Resource[k] == requestResource || rs[i].Rules[j].Resource[k] == "*" {
					resourceMatch = true
					for x := range rs[i].Rules[j].Verbs {
						if rs[i].Rules[j].Verbs[x] == requestMethod || rs[i].Rules[j].Verbs[x] == "*" {
							methodMatch = true
						}
					}
				}
			}
		}
	}
	return resourceMatch, methodMatch
}

func resourceNameInvalidHandler() iris.Handler {
	return func(ctx *context.Context) {
		r := ctx.GetCurrentRoute()
		if strings.Contains(r.Path(), "/:name") {
			resourceName := ctx.Params().GetString("name")
			if resourceName == "" {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", resourceName))
				return
			}
		}
		ctx.Next()
	}
}

func AddV1Route(app iris.Party) {
	v1Party := app.Party("/v1")
	v1Party.Use(langHandler())
	v1Party.Use(pageHandler())
	session.Install(v1Party)
	authParty := v1Party.Party("")
	authParty.Use(resourceExtractHandler())
	authParty.Use(authHandler())
	authParty.Use(roleHandler())
	authParty.Use(roleAccessHandler())
	authParty.Use(resourceNameInvalidHandler())
	authParty.Get("/", apiResourceHandler(authParty))
	user.Install(authParty)
	cluster.Install(authParty)
	role.Install(authParty)
	proxy.Install(authParty)
	ws.Install(authParty)
	chart.Install(authParty)
	webkubectl.Install(authParty, v1Party)
}
