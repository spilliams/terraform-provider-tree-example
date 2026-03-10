// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package root

// func TestAccExampleDataSource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Read testing
// 			{
// 				Config: testAccExampleDataSourceConfig,
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					statecheck.ExpectKnownValue(
// 						"data.scaffolding_example.test",
// 						tfjsonpath.New("id"),
// 						knownvalue.StringExact("example-id"),
// 					),
// 				},
// 			},
// 		},
// 	})
// }

// const testAccExampleDataSourceConfig = `
// data "scaffolding_example" "test" {
//   configurable_attribute = "example"
// }
// `
