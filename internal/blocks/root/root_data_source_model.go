package root

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceModel describes the data source data model.
type DataSourceModel struct {
	// required
	Name types.String `tfsdk:"name"`
	// computed
	ID types.String `tfsdk:"id"`
	// optional
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Defaulted             types.String `tfsdk:"defaulted"`
	Colors                types.List   `tfsdk:"colors"`
}

func (dm *DataSourceModel) SetAttributes(attrs map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	colorsIface, present := attrs["colors"]
	if present {
		colorsList, typed := colorsIface.([]string)
		if typed {
			colorValues := make([]attr.Value, len(colorsList))
			for i, az := range colorsList {
				colorValues[i] = types.StringValue(az)
			}
			dm.Colors, diags = types.ListValue(types.StringType, colorValues)
		}
	}
	return diags
}
