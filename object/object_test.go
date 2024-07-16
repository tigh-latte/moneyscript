package object_test

import (
	"testing"

	"git.tigh.dev/tigh-latte/monkeyscript/object"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello world"}
	hello2 := &object.String{Value: "Hello world"}

	diff1 := &object.String{Value: "My name is johnny"}
	diff2 := &object.String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
