// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/spilliams/terraform-provider-tree-example/internal/blocks"
	"github.com/spilliams/tree-terraform-provider/pkg/storage/dynamodb"
)

const (
	providerAttrAWSProfile = "profile"
	providerAttrAWSRegion  = "region"
	providerAttrTableName  = "table_name"
	providerAttrKeyARN     = "kms_key_arn"
)

// ScaffoldingProvider defines the provider implementation.
type ScaffoldingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	commit  string
}

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}

func New(version, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
			commit:  commit,
		}
	}
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	AWSProfile types.String `tfsdk:"profile"`
	AWSRegion  types.String `tfsdk:"region"`
	TableName  types.String `tfsdk:"table_name"`
	KMSKeyARN  types.String `tfsdk:"kms_key_arn"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tree"
	resp.Version = fmt.Sprintf("%s-%s", p.version, p.commit)
}

func (p *ScaffoldingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with a tree of information.",
		Attributes: map[string]schema.Attribute{
			providerAttrAWSProfile: schema.StringAttribute{
				Description: "The AWS profile to use for DynamoDB storage.",
				Required:    true,
			},
			providerAttrAWSRegion: schema.StringAttribute{
				Description: "The AWS region to use for DynamoDB storage.",
				Required:    true,
			},
			providerAttrTableName: schema.StringAttribute{
				Description: "The table name to use for DynamoDB storage.",
				Required:    true,
			},
			providerAttrKeyARN: schema.StringAttribute{
				Description: "The ARN of the KMS key to use for encrypting the DynamoDB storage.",
				Required:    true,
			},
		},
	}
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Tree client")

	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.AWSProfile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(providerAttrAWSProfile),
			"Unknown profile",
			"Cannot configure the provider client with an unknown profile.",
		)
	}
	if data.AWSRegion.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(providerAttrAWSRegion),
			"Unknown region",
			"Cannot configure the provider client with an unknown region.",
		)
	}
	ctx = tflog.SetField(ctx, providerAttrAWSRegion, data.AWSRegion.ValueString())

	if data.TableName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(providerAttrTableName),
			"Unknown table name",
			"Cannot configure the provider client with an unknown DynamoDB storage table name.",
		)
	}
	ctx = tflog.SetField(ctx, providerAttrTableName, data.TableName.ValueString())

	if data.KMSKeyARN.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(providerAttrKeyARN),
			"Unknown KMS Key ARN",
			"Cannot configure the provider client with an unknown KMS Key ARN.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Tree DynamoDB client")

	client, err := dynamodb.NewClient(ctx,
		data.AWSProfile.ValueString(),
		data.AWSRegion.ValueString(),
		data.TableName.ValueString(),
		data.KMSKeyARN.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create provider client",
			"An unexpected error occurred when creating the provider client.\n\n"+
				err.Error(),
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Tree client", map[string]any{"success": true})
}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
