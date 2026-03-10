package root

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceModel describes the resource data model.
type ResourceModel struct {
	// required
	Name types.String `tfsdk:"name"`
	// computed
	ID types.String `tfsdk:"id"`
	// optional
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Defaulted             types.String `tfsdk:"defaulted"`
	Colors                types.List   `tfsdk:"colors"`
}

func (rm *ResourceModel) GetAttributes(ctx context.Context) (map[string]interface{}, diag.Diagnostics) {
	attrs := make(map[string]interface{})
	var diags diag.Diagnostics

	if !rm.Colors.IsNull() {
		colors := make([]string, len(rm.Colors.Elements()))
		diags = rm.Colors.ElementsAs(ctx, &colors, false)
		if diags.HasError() {
			return nil, diags
		}
		attrs["colors"] = colors
	}

	return attrs, diags
}

func (rm *ResourceModel) SetAttributes(attrs map[string]interface{}) {
	colorsIface, present := attrs["colors"]
	if present {
		colorsList, typed := colorsIface.([]string)
		if typed {
			colorValues := make([]attr.Value, len(colorsList))
			for i, az := range colorsList {
				colorValues[i] = types.StringValue(az)
			}
			rm.Colors, _ = types.ListValue(types.StringType, colorValues)
		}
	}
}
