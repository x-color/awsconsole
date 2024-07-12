package awsconsole

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStore_addAndUpdateInstance(t *testing.T) {
	type fields struct {
		Instances []Instance
	}
	type args struct {
		instance Instance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Instance
	}{
		{
			name: "Add new instance for the same SSO URL",
			fields: fields{
				Instances: []Instance{
					{
						SsoURL: "https://sso.aws.com",
						AccountRoles: []AccountRole{
							{
								AccountID:   "123456789012",
								AccountName: "account",
								RoleName:    "roleA",
							},
						},
					},
				},
			},
			args: args{
				instance: Instance{
					SsoURL: "https://sso.aws.com",
					AccountRoles: []AccountRole{
						{
							AccountID:   "123456789012",
							AccountName: "account",
							RoleName:    "roleB",
						},
					},
				},
			},
			want: []Instance{
				{
					SsoURL: "https://sso.aws.com",
					AccountRoles: []AccountRole{
						{
							AccountID:   "123456789012",
							AccountName: "account",
							RoleName:    "roleB",
						},
					},
				},
			},
		},
		{
			name: "Add new instance for the different SSO URL",
			fields: fields{
				Instances: []Instance{
					{
						SsoURL: "https://sso.aws.com",
						AccountRoles: []AccountRole{
							{
								AccountID:   "123456789012",
								AccountName: "account",
								RoleName:    "roleA",
							},
						},
					},
				},
			},
			args: args{
				instance: Instance{
					SsoURL: "https://sso2.aws.com",
					AccountRoles: []AccountRole{
						{
							AccountID:   "123456789012",
							AccountName: "account",
							RoleName:    "roleB",
						},
					},
				},
			},
			want: []Instance{
				{
					SsoURL: "https://sso.aws.com",
					AccountRoles: []AccountRole{
						{
							AccountID:   "123456789012",
							AccountName: "account",
							RoleName:    "roleA",
						},
					},
				},
				{
					SsoURL: "https://sso2.aws.com",
					AccountRoles: []AccountRole{
						{
							AccountID:   "123456789012",
							AccountName: "account",
							RoleName:    "roleB",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Store{
				Instances: tt.fields.Instances,
			}
			c.addAndUpdateInstance(tt.args.instance)
			if diff := cmp.Diff(tt.want, c.Instances); diff != "" {
				t.Errorf("After addAndUpdateInstance(), Instances mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStore_findInstance(t *testing.T) {
	type fields struct {
		Instances []Instance
	}
	type args struct {
		ssoURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Instance
	}{
		{
			name: "Find instance",
			fields: fields{
				Instances: []Instance{
					{
						SsoURL: "https://sso.aws.com",
						AccountRoles: []AccountRole{
							{
								AccountID:   "123456789012",
								AccountName: "account",
								RoleName:    "roleA",
							},
						},
					},
					{
						SsoURL: "https://sso2.aws.com",
						AccountRoles: []AccountRole{
							{
								AccountID:   "123456789012",
								AccountName: "account",
								RoleName:    "roleB",
							},
						},
					},
				},
			},
			args: args{
				ssoURL: "https://sso.aws.com",
			},
			want: &Instance{
				SsoURL: "https://sso.aws.com",
				AccountRoles: []AccountRole{
					{
						AccountID:   "123456789012",
						AccountName: "account",
						RoleName:    "roleA",
					},
				},
			},
		},
		{
			name: "Not found",
			fields: fields{
				Instances: []Instance{},
			},
			args: args{
				ssoURL: "https://sso.aws.com",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Store{
				Instances: tt.fields.Instances,
			}
			got := c.findInstance(tt.args.ssoURL)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findInstance() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
