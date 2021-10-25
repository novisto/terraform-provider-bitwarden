package bitwarden

import "github.com/hashicorp/terraform-plugin-framework/types"

type SecureNote struct {
	Object         types.String `tfsdk:"object"`
	ID             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	FolderID       types.String `tfsdk:"folder_id"`
	Type           types.Number `tfsdk:"type"`
	Reprompt       types.Bool   `tfsdk:"reprompt"`
	Name           types.String `tfsdk:"name"`
	Notes          types.String `tfsdk:"notes"`
	Favorite       types.Bool   `tfsdk:"favorite"`
	CollectionIDs  []string     `tfsdk:"collection_ids"`
	RevisionDate   types.String `tfsdk:"revision_date"`
}
