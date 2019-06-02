package cmd

import "testing"

// TODO: Move this to long running e2e integration test; can be flaky :(
func Test_isApplicationPageActive(t *testing.T) {
	type args struct {
		pageURL string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ok #1", args{"Proj1_Info.cfm?Name=776091&S=S"}, true},
		{"fail #1", args{"Proj1_Info.cfm?Name=776130&S=S"}, false},
		{"fail #2", args{"Proj1_Info.cfm?Name=776131&S=S"}, false},
		{"fail #3", args{"Proj1_Info.cfm?Name=776137&S=S"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isApplicationPageActive(tt.args.pageURL); got != tt.want {
				t.Errorf("isApplicationPageActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
