# awsconsole

This is a simple tool to open AWS Management Console from your terminal by SSO login.

It opens the console with specified AWS Account and IAM Role.

## Prerequisites

- Golang 1.21 or later
- `~/.aws/config` is configured for SSO login

## Installation

```bash
$ go install github.com/x-color/awsconsole@latest
```

## Usage

First, you need to login AWS by SSO to get your AWS Accounts information.

```bash
$ aws sso login
```

Next, run the following comamnd to generate AWS account list file(`~/.aws/cli/awsconsole.json`).

```bash
$ awsconsole -update
```

After generating the file, you can open AWS console by the following command. This command opens interactive prompt to select AWS account and IAM role.

```bash
$ awsconsole
```

*If you want to use another profile, you can specify it by `AWS_PROFILE` environment variable.

### Optional setting

You can simplify accessing the AWS Management Console by an [alias](https://docs.aws.amazon.com/cli/latest/userguide/cli-usage-alias.html) for `awsconsole` command in AWS CLI, making it quicker.

If you add the following setting to your `~/.aws/cli/alias`,

```ini
[toplevel]
console = !f() {
    awsconsole $@
}; f
```

You can access the AWS Management Console by the following command.

```bash
$ aws console
```
