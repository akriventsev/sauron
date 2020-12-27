package main

import (
	"fmt"
	"testing"
)

func TestService_Check(t *testing.T) {
	s := Service{}
	s.Check()
	fmt.Print(s)
	t.FailNow()
}
