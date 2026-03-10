package blocks

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/spilliams/terraform-provider-tree-example/internal/blocks/root"
)

func AllDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// root.NewDataSource,
	}
}

func AllResources() []func() resource.Resource {
	return []func() resource.Resource{
		root.NewResource,
	}
}
