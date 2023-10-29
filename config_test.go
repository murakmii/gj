package gojiai

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
	"sync"
	"testing"
)

func TestX(t *testing.T) {
	m := &sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	m.Lock()
	fmt.Println("yes double lock")
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		before func()
		after  func()
		expect *Config
	}{
		{
			name:   "parse simple config",
			in:     "{\"class_path\":[\"foo\",\"bar\"],\"system_properties\":{\"hoge\":\"fuga\"}}",
			expect: &Config{ClassPath: []string{"foo", "bar"}, SysProps: map[string]string{"hoge": "fuga"}},
		},
		{
			name: "parse config uses env var",
			before: func() {
				os.Setenv("SAMPLE_PATH", "baz")
			},
			in:     "{\"class_path\":[\"{{ .SAMPLE_PATH }}\",\"bar\"]}",
			expect: &Config{ClassPath: []string{"baz", "bar"}},
			after: func() {
				os.Unsetenv("SAMPLE_PATH")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.before != nil {
				test.before()
			}
			if test.after != nil {
				defer test.after()
			}

			got, gotErr := ReadConfig(bytes.NewBufferString(test.in))
			if gotErr != nil {
				t.Errorf("ReadConfig() returned unexpected error: %s", gotErr)
				return
			}

			if diff := cmp.Diff(test.expect, got); len(diff) > 0 {
				t.Errorf("ReadConfig() returned unexpected config: %s", diff)
				return
			}
		})
	}
}
