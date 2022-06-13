# Terraform Bitwarden Provider

This provider allows managing BitWarden resources in Terraform.

The supported features are based on our internal needs, so this plugin does
not intend to support all features of BitWarden.

## Prerequisites

- Install the [`bw` CLI](https://bitwarden.com/help/article/cli/)
- Login to BitWarden with `bw login`

## Usage

Provide the password for your BitWarden account using either the `BW_PASSWORD` environment
variable or through the provider configuration. You are now ready to use the provider.

You can also run `bw serve` yourself and provide the port on which it is running either the
`BW_SERVE_PORT` environment variable or through the provier configuration.

## Running locally

Local setup for development, you will need Go 1.18 and Terraform 1.0.3+

1. Copy the `.terraformrc.example` file into your HOME and change 
   the name to `.terraformrc` and the path inside to your own username
2. From the root of the repo run `make install`
3. In the `examples` directory run `TF_LOG=TRACE t apply -auto-approve`
