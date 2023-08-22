package device42

import (
	"context"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBuildings() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_buildings` data source can be used to retrieve all buildings.",
		ReadContext: dataSourceBuildingsRead,
		Schema: map[string]*schema.Schema{
			"buildings": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The `id` of this building.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "The `name` of this building.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"address": &schema.Schema{
							Description: "The `address` of this building.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"notes": &schema.Schema{
							Description: "`notes` on this building.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// get buildings
func dataSourceBuildingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics

	buildings, err := c.GetBuildings()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get a list of buildings",
			Detail:   err.Error(),
		})
		return diags
	}

	bs := flattenBuildingsData(buildings)
	if err := d.Set("buildings", bs); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to set buildings",
			Detail:   err.Error(),
		})
		return diags
	}

	ids := make([]int, len((*buildings)))
	for _, i := range *buildings {
		ids = append(ids, i.BuildingID)
	}

	checksum := idsChecksum(ids)

	d.SetId(checksum)

	return diags
}

// flatten buildings to a map
func flattenBuildingsData(buildings *[]device42.Building) []interface{} {
	if buildings != nil {
		bs := make([]interface{}, len(*buildings))

		for i, building := range *buildings {
			b := make(map[string]interface{})

			b["id"] = building.BuildingID
			b["name"] = building.Name
			b["address"] = building.Address
			b["notes"] = building.Notes

			bs[i] = b
		}

		return bs
	}

	return make([]interface{}, 0)
}
