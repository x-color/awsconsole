package awsconsole

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	fzf "github.com/junegunn/fzf/src"
)

func GenerateAccountsRolesFile(profileName string) error {
	info, err := extractSSOInfo(config.DefaultSharedConfigFilename(), profileName)
	if err != nil {
		return err
	}
	fp, err := cachedTokenFilepath(info)
	if err != nil {
		return err
	}

	accessToken, err := extractAccessToken(fp)
	if err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	cli := sso.NewFromConfig(cfg)
	ar, err := listAccountsRoles(cli, accessToken)
	if err != nil {
		return err
	}

	instance := Instance{
		SsoURL:       info.url,
		AccountRoles: ar,
	}
	return updateInstancesCache(instance)
}

func updateInstancesCache(instance Instance) error {
	cache, err := loadStore()
	if err != nil {
		cache = &Store{Instances: []Instance{}}
	}
	cache.addAndUpdateInstance(instance)

	return saveStore(cache)
}

func Jump(profileName string) error {
	url, err := extractSSOStartURL(config.DefaultSharedConfigFilename(), profileName)
	if err != nil {
		return err
	}

	store, err := loadStore()
	if err != nil {
		return err
	}

	accountRoles := store.findInstance(url).AccountRoles
	list := make([]string, len(accountRoles))
	for i, ar := range accountRoles {
		list[i] = ar.String()
	}

	var jumpErr error
	err = fuzzyFinderForAccountRoles(list, func(s string) {
		if err := jump(parseAccountRole(s).URL(url)); err != nil {
			jumpErr = err
		}
	})
	if err != nil {
		return err
	}
	return jumpErr
}

func fuzzyFinderForAccountRoles(inputs []string, f func(string)) error {
	options, err := fzf.ParseOptions(true, []string{})
	if err != nil {
		return err
	}
	options.Printer = f

	options.Input = make(chan string, len(inputs))
	for _, in := range inputs {
		options.Input <- in
	}

	_, err = fzf.Run(options)
	return err
}

func jump(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
