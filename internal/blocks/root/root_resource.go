// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package root

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spilliams/tree-terraform-provider/pkg/storage/dynamodb"
)

// RootResource defines the resource implementation.
type RootResource struct {
	client *dynamodb.Client
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RootResource{}
var _ resource.ResourceWithImportState = &RootResource{}

func NewResource() resource.Resource {
	return &RootResource{}
}

func (r *RootResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dynamodb.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *dynamodb.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *RootResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root"
}

func (r *RootResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Root resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the root.",
			},
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "Root's configurable attribute",
				Optional:            true,
			},
			"defaulted": schema.StringAttribute{
				MarkdownDescription: "Root's configurable attribute with default value",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("example value when not configured"),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The root's identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"colors": schema.ListAttribute{
				Optional:    true,
				Description: "Root's colors.",
				ElementType: types.StringType,
			},
		},
	}
}

func (r *RootResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attrs, diags := plan.GetAttributes(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	saved, err := r.client.CreateEntity(ctx,
		"root",
		plan.Name.ValueString(),
		attrs,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating root",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(saved.ID())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RootResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	real, err := r.client.GetEntity(ctx, "root", state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting root",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(real.ID())
	state.SetAttributes(real.Attributes())

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RootResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// identity changes
	if !plan.Name.Equal(state.Name) {
		updated, err := r.client.UpdateEntity(ctx,
			"root",
			plan.ID.ValueString(),
			plan.Name.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating root",
				err.Error(),
			)
			return
		}
		plan.Name = types.StringValue(updated.Label())
	}

	if !plan.Colors.Equal(state.Colors) {
		colors := make([]string, len(plan.Colors.Elements()))
		diags = plan.Colors.ElementsAs(ctx, &colors, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		err := r.client.UpdateAttribute(ctx,
			"root",
			plan.ID.ValueString(),
			"colors",
			colors,
		)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("colors"),
				"Error updating root's colors",
				err.Error(),
			)
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RootResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete root, got error: %s", err))
	//     return
	// }
}

func (r *RootResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
