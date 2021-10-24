package bitwarden

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"math/big"
)

type resourceSecureNoteType struct{}

func (r resourceSecureNoteType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	// TODO: HERE
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"object": {
				Type: types.StringType,
				Computed: true,
			},
			"id": {
				Type: types.StringType,
				Computed: true,
			},
			"organization_id": {
				Type: types.StringType,
				Required: true,
			},
			"folder_id": {
				Type: types.StringType,
				Optional: true,
			},
			"type": {
				Type: types.NumberType,
				Computed: true,
			},
			"reprompt": {
				Type: types.NumberType,
				Optional: true,
			},
			"name": {
				Type: types.StringType,
				Required: true,
			},
			"notes": {
				Type: types.StringType,
				Required: true,
			},
			"favorite": {
				Type: types.BoolType,
				Optional: true,
			},
			"collection_ids": {
				Type: types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"revision_date": {
				Type: types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (r resourceSecureNoteType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceSecureNote{
		p: *(p.(*provider)),
	}, nil
}

type resourceSecureNote struct {
	p provider
}

func (r resourceSecureNote) ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {

}

// Create a new resource
func (r resourceSecureNote) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan SecureNote
	diags := req.Plan.Get(ctx, &plan)
	log.Printf("%+v\n", plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secureNote, err := r.p.client.CreateSecureNote(plan)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: HERE
	var result = SecureNote{
		Object:         types.String{Value: secureNote.Object},
		ID:             types.String{Value: secureNote.ID},
		OrganizationId: types.String{Value: secureNote.OrganizationId},
		FolderID:       types.String{Value: secureNote.FolderID},
		Type:           types.Number{Value: big.NewFloat(float64(secureNote.Type))},
		Reprompt:       types.Number{Value: big.NewFloat(float64(secureNote.Reprompt))},
		Name:           types.String{Value: secureNote.Name},
		Notes:          types.String{Value: secureNote.Notes},
		Favorite:       types.Bool{Value: secureNote.Favorite},
		CollectionIDs:  secureNote.CollectionIDs,
		RevisionDate:   types.String{Value: secureNote.RevisionDate.String()},
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceSecureNote) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var plan SecureNote
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secureNoteId := plan.ID.Value

	// Get order current value
	secureNote, err := r.p.client.GetItem(secureNoteId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading order",
			"Could not read secure note ID " + secureNoteId + ": "+err.Error(),
		)
		return
	}

	// TODO: HERE
	state := SecureNote{
		Object:         types.String{Value: secureNote.Object},
		ID:             types.String{Value: secureNote.ID},
		OrganizationId: types.String{Value: secureNote.OrganizationId},
		FolderID:       types.String{Value: secureNote.FolderID},
		Type:           types.Number{Value: big.NewFloat(float64(secureNote.Type))},
		Reprompt:       types.Number{Value: big.NewFloat(float64(secureNote.Reprompt))},
		Name:           types.String{Value: secureNote.Name},
		Notes:          types.String{Value: secureNote.Notes},
		Favorite:       types.Bool{Value: secureNote.Favorite},
		CollectionIDs:  secureNote.CollectionIDs,
		RevisionDate:   types.String{Value: secureNote.RevisionDate.String()},
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceSecureNote) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

}

// Delete resource
func (r resourceSecureNote) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

}
