package laundry

import (
	"fmt"
	"net/http"
)

type laundry struct {
	R *router
}

type Config struct {
	BasketDir string
}

func (l *laundry) Run(addr string) error {
	return http.ListenAndServe(addr, l.R)
}

func New(c ...Config) *laundry {

	l := &laundry{
		R: &router{
			handlers: make(map[string]map[string]handler),
		},
	}

	var baskets []basket
	if len(c) > 0 {
		baskets = loadBasket(c[0].BasketDir)
	} else {
		baskets = loadBasket("./basket")
	}

	for _, basket := range baskets {
		for _, endpoint := range basket.Endpoints {
			fmt.Println(endpoint)
			l.R.HandleFunc(endpoint.Method, endpoint.Pattern, basket)
		}
	}

	return l
}
