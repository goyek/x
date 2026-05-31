package cmd

import "testing"

func TestMask(t *testing.T) {
	tests := []struct {
		name    string
		cmdLine string
		want    string
	}{
		{
			name:    "no env vars",
			cmdLine: "echo hello",
			want:    "echo hello",
		},
		{
			name:    "one env var",
			cmdLine: "FOO=bar echo hello",
			want:    "FOO=[MASKED] echo hello",
		},
		{
			name:    "multiple env vars",
			cmdLine: "FOO=bar BAZ=qux echo hello",
			want:    "FOO=[MASKED] BAZ=[MASKED] echo hello",
		},
		{
			name:    "env var with space",
			cmdLine: "FOO='bar baz' echo hello",
			want:    "FOO=[MASKED] echo hello",
		},
		{
			name:    "arg with space",
			cmdLine: "echo 'hello world'",
			want:    "echo 'hello world'",
		},
		{
			name:    "complex",
			cmdLine: "FOO=bar ./script --arg='some value'",
			want:    "FOO=[MASKED] ./script '--arg=some value'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.cmdLine); got != tt.want {
				t.Errorf("Mask() = %q, want %q", got, tt.want)
			}
		})
	}
}
