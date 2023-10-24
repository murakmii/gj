package class_file

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestMethodDescriptor_Params(t *testing.T) {
	tests := []struct {
		sut    MethodDescriptor
		expect []FieldType
	}{
		{sut: "()V", expect: nil},
		{sut: "(I)V", expect: []FieldType{"I"}},
		{sut: "(IBCD)V", expect: []FieldType{"I", "B", "C", "D"}},
		{sut: "(Ljava/lang/Object;)V", expect: []FieldType{"Ljava/lang/Object;"}},
		{sut: "(Ljava/lang/Object;Ljava/lang/Object;)V", expect: []FieldType{"Ljava/lang/Object;", "Ljava/lang/Object;"}},
		{sut: "(IBLjava/lang/Object;CZ)V", expect: []FieldType{"I", "B", "Ljava/lang/Object;", "C", "Z"}},
		{sut: "([I)V", expect: []FieldType{"[I"}},
		{sut: "([IBC)V", expect: []FieldType{"[I", "B", "C"}},
		{sut: "(J[I)V", expect: []FieldType{"J", "[I"}},
		{sut: "([[[[Ljava/lang/Object;)V", expect: []FieldType{"[[[[Ljava/lang/Object;"}},
		{sut: "(I[[[[Ljava/lang/Object;C)V", expect: []FieldType{"I", "[[[[Ljava/lang/Object;", "C"}},
	}

	for _, test := range tests {
		t.Run(string(test.sut), func(t *testing.T) {
			got := test.sut.Params()

			if diff := cmp.Diff(got, test.expect); len(diff) > 0 {
				t.Errorf("Params() return unexpected field types = %s", diff)
			}
		})
	}
}
