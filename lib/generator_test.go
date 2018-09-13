package lib

import (
	"testing"
	"tinywell/loadgen/model"
)

func TestInterface(t *testing.T) {
	var mgen *myGenerator
	var gen model.Generator = mgen
	gen.Start()
}

func TestStop(t *testing.T) {

}
