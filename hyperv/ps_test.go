package hyperv

import "testing"

func TestPs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "''",
		},
		{
			name:  "plain string",
			input: "hello",
			want:  "'hello'",
		},
		{
			name:  "string with single quote",
			input: "it's",
			want:  "'it''s'",
		},
		{
			name:  "string with two single quotes",
			input: "it's a test",
			want:  "'it''s a test'",
		},
		{
			name:  "injection attempt",
			input: "evil'; Stop-VM; echo '",
			want:  "'evil''; Stop-VM; echo '''",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ps(tc.input)
			if got != tc.want {
				t.Errorf("ps(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}
