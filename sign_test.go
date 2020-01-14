package linepay

import (
	"testing"
)

func Test_calculate(t *testing.T) {

	in := "d7N2zcBCDH7EXw28ym/ppeNqa/Gp/9Xv/hO40MNjtI8="
	ss := calculate("A", "BODY")

	if ss != in {
		t.Errorf("failed to test singer calculate, want '%s', but got '%s'", in, ss)
	}

}
