package testpackage

import (
	"testing"
	even "github.com/KanybekMomukeyev/goDockerCompose/grpc/multithread/even"
)

func TestEven(t *testing.T) {

	if !even.Even(10) {
		t.Log("10 must beeven!")
		t.Fail()
	}

	if even.Even(7) {
		t.Log("7 is not even!")
		t.Fail()
	}
}

func TestOdd(t *testing.T) {

	if !even.Odd(11) {
		t.Log("11 must be odd!")
		t.Fail()
	}

	if even.Odd(10) {
		t.Log("10 is not odd!")
		t.Fail()
	}
}

func TestEven222(t *testing.T) {
	if even.Even(10) {
		t.Log("Everything OK: 10 is even, just a test to see failed output!")
		t.Fail()
	}
}

