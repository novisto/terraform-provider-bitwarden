package bitwarden

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"strconv"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *Client
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"password": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  false,
				Sensitive: true,
			},
			"bw_serve_port": {
				Type:     types.Int64Type,
				Optional: true,
				Computed: false,
			},
		},
	}, nil
}

type providerData struct {
	Password    types.String `tfsdk:"password"`
	BwServePort types.Int64  `tfsdk:"bw_serve_port"`
}

func (p *provider) Configure(
	ctx context.Context,
	request tfsdk.ConfigureProviderRequest,
	response *tfsdk.ConfigureProviderResponse,
) {
	// Retrieve provider data from configuration
	var config providerData
	diags := request.Config.Get(ctx, &config)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// User must provide a password to the provider
	var password string
	if config.Password.Unknown {
		// Cannot connect to client with an unknown value
		response.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	// If password is not provided in the config, try to get it from the environment
	if config.Password.Null {
		password = os.Getenv("BW_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if password == "" {
		// Cannot continue without a password
		response.Diagnostics.AddError(
			"Unable to find password",
			"password cannot be an empty string",
		)
		return
	}

	bwServePort := int64(0)
	if config.BwServePort.Unknown {
		// Cannot connect to client with an unknown value
		response.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as bwServePort",
		)
		return
	}

	if config.BwServePort.Null {
		port := os.Getenv("BW_SERVE_PORT")
		if port != "" {
			p, err := strconv.Atoi(port)
			if err != nil {
				response.Diagnostics.AddError(
					"Could not parse environment variable BW_SERVE_PORT to int.",
					err.Error(),
				)
			}
			bwServePort = int64(p)
		}
	} else {
		bwServePort = config.BwServePort.Value
	}

	// Create a new BitWarden client and set it to the provider client
	c, err := NewClient(password, bwServePort)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to create client",
			fmt.Sprintf("Unable to create BitWarden client: %s", err.Error()),
		)
		return
	}

	p.client = c
	p.configured = true
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"bitwarden_secure_note": resourceSecureNoteType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
