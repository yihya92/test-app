package Employees

import (
	"net/http"
	"github.com/gorilla/mux"

)

type Route struct {
	Name               string
	Method             string
	Pattern            string
	HandlerFunc        http.HandlerFunc
	AddToAccessEntry   bool
	DisplayName        string //create new
	DisplayOrder       int64
	Module             string //inventory
	ModuleDisplayOrder int64
	Level1             string //transactions
	Level1DisplayOrder int64
	Level2             string //purchase request
	Level2DisplayOrder int64
	Level3             string
	Level3DisplayOrder int64
}

type Routes []Route

func CreateRoutes(UC *UserControl) Routes {
	return Routes{
		// Route{
		// 	"HTTP_ReceiveSSR",
		// 	"POST",
		// 	"/HTTP_ReceiveSSR/",
		// 	UC.HTTP_ReceiveSSR,
		// 	false,
		// 	"Receive SSR", // DisplayName
		// 	1,             // DisplayOrder
		// 	"",            // Module
		// 	0,             //ModuleDisplayOrder
		// 	"",            // Level1
		// 	0,             // Level1DisplayOrder
		// 	"",            // Level2
		// 	0,             // Level2DisplayOrder
		// 	"",            // Level3
		// 	0,             // Level3DisplayOrder
		// },
	}
}

func (Uc *UserControl) AddToRouter(router *mux.Router, UC *UserControl) {
	// When StrictSlash is set to true, if the route path is "/path/", accessing "/path" will redirect
	// to the former and vice versa
	routes := CreateRoutes(UC)

	Uc.AddSubscriptionRoutes(&routes)


	for _, route := range routes {

		//var handler http.Handler
		handler := route.HandlerFunc
		//handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		//add to OKAPI AUC

	}


	// create default access group

	for _, route := range routes {
		//var handler http.Handler
		handler := route.HandlerFunc
		//handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
}

func Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}


