package lib

import (
	"testing"
)

func TestInterface(t *testing.T) {
	var mgen *myGenerator
	var gen Generator = mgen
	gen.Start()
}
