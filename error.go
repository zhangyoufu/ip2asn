package ip2asn

import "fmt"

func reportError(errc chan<- error, pattern string, varargs ...interface{}) {
	if errc == nil {
		return
	}
	errc <- fmt.Errorf(pattern, varargs...)
}
