// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package root

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spilliams/tree-terraform-provider/pkg/storage/dynamodb"
)

// RootDataSource defines the data source implementation.
type RootDataSource struct {
	client *dynamodb.Client
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RootDataSource{}

func NewDataSource() datasource.DataSource {
	return &RootDataSource{}
}

func (d *RootDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*dynamodb.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dynamodb.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *RootDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root"
}

func (d *RootDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Root data source",

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
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The root's identifier",
				Computed:            true,
			},
			"colors": schema.ListAttribute{
				Optional:    true,
				Description: "Root's colors.",
				ElementType: types.StringType,
			},
		},
	}
}

func (d *RootDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	real, err := d.client.GetEntity(ctx, "root", state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting root",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(real.ID())
	resp.Diagnostics.Append(state.SetAttributes(real.Attributes())...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
