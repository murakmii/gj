package class_file

import "testing"

func TestMethodInfo_NumArgs(t *testing.T) {
	tests := []struct {
		desc   string
		expect int
	}{
		{desc: "()V", expect: 0},
		{desc: "(I)V", expect: 1},
		{desc: "(IBCD)V", expect: 4},
		{desc: "(Ljava/lang/Object;)V", expect: 1},
		{desc: "(IBLjava/lang/Object;CZ)V", expect: 5},
		{desc: "([I)V", expect: 1},
		{desc: "([IBC)V", expect: 3},
		{desc: "(J[I)V", expect: 2},
		{desc: "([[[[Ljava/lang/Object;)V", expect: 1},
		{desc: "(I[[[[Ljava/lang/Object;C)V", expect: 3},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			sut := &MethodInfo{desc: &test.desc}
			got := sut.NumArgs()

			if got != test.expect {
				t.Errorf("MethodInfo.NumArgs retuend = %d, expected = %d", got, test.expect)
			}
		})
	}
}
