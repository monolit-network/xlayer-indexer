package evmclient

import (
	"errors"
	"testing"
)

func TestIsRetryableErrorConnectionIssues(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "connection reset by peer",
			err:  errors.New(`Post "http://rpc": read tcp 127.0.0.1:1->127.0.0.1:2: read: connection reset by peer`),
			want: true,
		},
		{
			name: "connection refused",
			err:  errors.New(`Post "http://rpc": dial tcp 127.0.0.1:8545: connect: connection refused`),
			want: true,
		},
		{
			name: "non retryable",
			err:  errors.New("some other error"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isRetryableError(tc.err); got != tc.want {
				t.Fatalf("isRetryableError(%q) = %v, want %v", tc.err.Error(), got, tc.want)
			}
		})
	}
}
