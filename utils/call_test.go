package utils

import (
	"fmt"
	"testing"
)

ackage utils

import (
"fmt"
"testing"
)

func funcA(arg0 string, arg1 string) string {
	return fmt.Sprintf("Function one! %v %v", arg0, arg1)
}

func TestCall(t *testing.T) {
	StubStorage["funcA"] = funcA
	resA, _ := Call("funcA", "value", "keyword1")
	expected := "Function one! value keyword1"
	actual, ok := resA.(string)
	if !ok {
		t.Errorf("Test failed, expected: '%v', got:  '%v'", true, ok)
	}
	if expected != actual {
		t.Errorf("Test failed, expected: '%v', got:  '%v'", expected, actual)
	}

}

func TestCallSprintf(t *testing.T) {
	StubStorage["spf"] = fmt.Sprintf
	resA, _ := Call("spf", "Hello %d", 5)
	s, ok := resA.(string)
	if !ok {
		t.Errorf("Test failed, expected: '%v', got:  '%v'", true, ok)
	}
	if s != "Hello 5" {
		t.Errorf("Test failed, expected: '%v', got:  '%v'", "Hello 5", s)
	}
}

func returnList() []string {
	return []string{"a", "b", "hello"}
}

func TestCall2(t *testing.T) {
	StubStorage["returnList"] = returnList
	res, err := Call("returnList")
	if err != nil {
		t.Errorf("Test failed, expected: '%v', got:  '%v'", nil, err)
	}
	list, _ := res.([]string)
	fmt.Println(list)
}
