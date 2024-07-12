package awsconsole

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseAccountRole(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want AccountRole
	}{
		{
			name: "Parse account role string",
			args: args{
				s: "account (123456789012) - roleA",
			},
			want: AccountRole{
				AccountID:   "123456789012",
				AccountName: "account",
				RoleName:    "roleA",
			},
		},
		{
			name: "Parse generated string from AccountRole",
			args: args{
				s: AccountRole{
					AccountID:   "123456789012",
					AccountName: "account",
					RoleName:    "roleA",
				}.String(),
			},
			want: AccountRole{
				AccountID:   "123456789012",
				AccountName: "account",
				RoleName:    "roleA",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseAccountRole(tt.args.s)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("parseAccountRole() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
