package config

import (
	"net/http"
	_ "net/http/pprof"
)

func StartPprof(address string) {

	http.ListenAndServe(address, nil)

}
