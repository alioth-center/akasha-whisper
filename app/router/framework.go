package router

import "github.com/alioth-center/infrastructure/network/http"

var frameworkRouter = http.NewRouter("framework")

var FrameworkRouterGroup = []http.EndPointInterface{
	http.NewEndPointBuilder[http.NoBody, http.NoResponse]().
		SetAllowMethods(http.GET).
		SetRouter(frameworkRouter.Group("/health")).
		SetHandlerChain(http.NewChain(http.EmptyHandler[http.NoBody, http.NoResponse]())).
		Build(),
}
