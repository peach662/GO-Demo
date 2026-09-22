package todo

import "testing"

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
		want bool
	}{
		{name: "pending to processing", from: StatusPending, to: StatusProcessing, want: true},
		{name: "processing to completed", from: StatusProcessing, to: StatusCompleted, want: true},
		{name: "processing to failed", from: StatusProcessing, to: StatusFailed, want: true},
		{name: "failed to processing", from: StatusFailed, to: StatusProcessing, want: true},
		{name: "same status is idempotent", from: StatusCompleted, to: StatusCompleted, want: true},
		{name: "completed to processing is invalid", from: StatusCompleted, to: StatusProcessing, want: false},
		{name: "pending to completed is invalid", from: StatusPending, to: StatusCompleted, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransition(tt.from, tt.to); got != tt.want {
				t.Fatalf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
