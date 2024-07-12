package awsconsole

import (
	"os"
	"testing"
)

func Test_extractAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantErr bool
	}{
		{
			name: "Extract token from credentials cache file",
			content: `{
				"accessToken":"token",
				"expiresAt":"2099-01-01T00:00:00Z"
			}`,
			want:    "token",
			wantErr: false,
		},
		{
			name: "Expired token",
			content: `{
				"accessToken":"token",
				"expiresAt":"2000-01-01T00:00:00Z"
			}`,
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid file format",
			content: `{
				"key":"value",
			}`,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatalf("temp file creation: %v", err)
			}
			defer os.Remove(f.Name())

			_, err = f.Write([]byte(tt.content))
			if err != nil {
				t.Fatalf("temp file write: %v", err)
			}

			got, err := extractAccessToken(f.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("extractAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractSSOStartURL(t *testing.T) {
	type args struct {
		profileName string
	}
	tests := []struct {
		name    string
		args    args
		content string
		want    string
		wantErr bool
	}{
		{
			name: "Extract SSO start URL from AWS config file",
			args: args{
				profileName: "user1",
			},
			content: `
[profile user1]
sso_session = my-sso
sso_account_id = 444455556666
sso_role_name = readOnly
region = us-east-1
output = json

[sso-session my-sso]
sso_region = us-east-1
sso_start_url = https://my-sso-portal.awsapps.com/start
sso_registration_scopes = sso:account:access`,
			want:    "https://my-sso-portal.awsapps.com/start",
			wantErr: false,
		},
		{
			name: "Extract SSO start URL from legacy AWS config file",
			args: args{
				profileName: "default",
			},
			content: `
[default]
sso_start_url = https://my-sso-portal.awsapps.com/start
sso_region = us-east-1
sso_account_id = 111122223333
sso_role_name = readOnly
region = us-west-2
output = text
`,
			want:    "https://my-sso-portal.awsapps.com/start",
			wantErr: false,
		},
		{
			name: "Profile not found in AWS config file",
			args: args{
				profileName: "not-found",
			},
			content: `
[profile user1]
sso_session = my-sso
sso_account_id = 444455556666
sso_role_name = readOnly
region = us-east-1
output = json

[sso-session my-sso]
sso_region = us-east-1
sso_start_url = https://my-sso-portal.awsapps.com/start
sso_registration_scopes = sso:account:access`,
			want:    "",
			wantErr: true,
		},
		{
			name: "SSO start URL not found in AWS config file",
			args: args{
				profileName: "default",
			},
			content: `
[profile user1]
sso_session = my-sso
sso_account_id = 444455556666
sso_role_name = readOnly
region = us-east-1
output = json

[sso-session my-sso]
sso_region = us-east-1
sso_registration_scopes = sso:account:access`,
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid file format",
			args: args{
				profileName: "default",
			},
			content: `
sso_start_url = https://my-sso-portal.awsapps.com/start
`,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatalf("temp file creation: %v", err)
			}
			defer os.Remove(f.Name())

			_, err = f.Write([]byte(tt.content))
			if err != nil {
				t.Fatalf("temp file write: %v", err)
			}

			got, err := extractSSOStartURL(f.Name(), tt.args.profileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractSSOStartURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractSSOStartURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
