package awsconsole

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
)

func listAccounts(cli *sso.Client, accessToken string) (map[string]string, error) {
	accounts := make(map[string]string)
	var nextToken *string
	for {
		out, err := cli.ListAccounts(context.Background(), &sso.ListAccountsInput{
			AccessToken: aws.String(accessToken),
			NextToken:   nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, account := range out.AccountList {
			accounts[aws.ToString(account.AccountId)] = aws.ToString(account.AccountName)
		}

		if out.NextToken == nil {
			return accounts, nil
		}
		nextToken = out.NextToken
	}
}

func listAccountRoles(cli *sso.Client, accessToken, accountID string) ([]string, error) {
	roles := make([]string, 0)
	var nextToken *string
	for {
		out, err := cli.ListAccountRoles(context.Background(), &sso.ListAccountRolesInput{
			AccessToken: aws.String(accessToken),
			AccountId:   aws.String(accountID),
			NextToken:   nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, role := range out.RoleList {
			roles = append(roles, aws.ToString(role.RoleName))
		}

		if out.NextToken == nil {
			return roles, nil
		}
		nextToken = out.NextToken
	}
}

type AccountRole struct {
	AccountID   string `json:"accountID"`
	AccountName string `json:"accountName"`
	RoleName    string `json:"roleName"`
}

func (a AccountRole) String() string {
	return fmt.Sprintf("%s (%s) - %s", a.AccountName, a.AccountID, a.RoleName)
}

func (a AccountRole) URL(baseURL string) string {
	return fmt.Sprintf("%s/#/console?account_id=%s&role_name=%s", baseURL, a.AccountID, a.RoleName)
}

func parseAccountRole(s string) AccountRole {
	l := strings.Split(s, " ")
	return AccountRole{
		AccountID:   strings.TrimSuffix(strings.TrimPrefix(l[1], "("), ")"),
		AccountName: l[0],
		RoleName:    l[3],
	}
}

func listAccountsRoles(cli *sso.Client, accessToken string) ([]AccountRole, error) {
	accounts, err := listAccounts(cli, accessToken)
	if err != nil {
		return nil, err
	}
	accountsRoles := make([]AccountRole, 0)
	for accountID, accountName := range accounts {
		roles, err := listAccountRoles(cli, accessToken, accountID)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			accountsRoles = append(accountsRoles, AccountRole{
				AccountID:   accountID,
				AccountName: accountName,
				RoleName:    role,
			})
		}
	}
	return accountsRoles, nil
}
