package lib

import (
	"testing"

	"github.com/tinywell/loadgen/model"
)

func TestInterface(t *testing.T) {
	var mgen *myGenerator
	var gen model.Generator = mgen
	gen.Start()
}

func TestStop(t *testing.T) {

}

func TestStatus(t *testing.T) {
	var mgen *myGenerator
	var gen model.Generator = mgen
	sta := gen.Status()
	t.Log(sta)
}
