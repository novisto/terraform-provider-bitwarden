package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"terraform-bitwarden-sync/bitwarden"
)

func main() {
	tfsdk.Serve(context.Background(), bitwarden.New, tfsdk.ServeOpts{
		Name: "bitwarden",
	})
}
