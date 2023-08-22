package device42

import (
	"context"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVRFGroups() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_vrf_groups` can be used to retrieve all VRF groups.",
		ReadContext: dataSourceVRFGroupsRead,
		Schema: map[string]*schema.Schema{
			"vrf_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All `vrf_groups`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The `id` of the VRF group.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "The `name` of the VRF group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": &schema.Schema{
							Description: "The `description` of the VRF group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"buildings": &schema.Schema{
							Description: "The `buildings` of the VRF group.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

// get vrf groups
func dataSourceVRFGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics

	vrfGroups, err := c.GetVRFGroups()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get a list of vrf groups",
			Detail:   err.Error(),
		})
		return diags
	}

	vgs := flattenVRFGroupsData(vrfGroups)
	if err := d.Set("vrf_groups", vgs); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to set vrf groups",
			Detail:   err.Error(),
		})
		return diags
	}

	ids := make([]int, len((*vrfGroups)))
	for _, i := range *vrfGroups {
		ids = append(ids, i.ID)
	}

	checksum := idsChecksum(ids)

	d.SetId(checksum)

	return diags
}

// flatten vrf groups to a map
func flattenVRFGroupsData(vrfGroups *[]device42.VRFGroup) []interface{} {
	if vrfGroups != nil {
		vgs := make([]interface{}, len(*vrfGroups))

		for i, vrfGroup := range *vrfGroups {
			vg := make(map[string]interface{})

			vg["id"] = vrfGroup.ID
			vg["name"] = vrfGroup.Name
			vg["description"] = vrfGroup.Description
			vg["buildings"] = vrfGroup.Buildings

			vgs[i] = vg
		}

		return vgs
	}

	return make([]interface{}, 0)
}
