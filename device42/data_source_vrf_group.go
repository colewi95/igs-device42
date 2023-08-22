package device42

import (
	"context"
	"fmt"
	"strconv"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVRFGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_vrf_group` data source can be used to retrieve a single VRF group by its `id` or its `name`.",
		ReadContext: dataSourceVRFGroupRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description:  "The `id` of a VRF group.",
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"id", "name"},
			},
			"name": &schema.Schema{
				Description:  "The `name` of a VRF group.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"id", "name"},
			},
			"description": &schema.Schema{
				Description: "The `description` of the VRF group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"building_ids": &schema.Schema{
				Description: "The `building_ids` of the VRF group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

// get a vrf group by id
func dataSourceVRFGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	vrfGroupID := d.Get("id").(int)
	vrfGroupName := d.Get("name").(string)
	vrfGroup := &device42.VRFGroup{}

	if vrfGroupID != 0 {
		fmt.Printf("[DEBUG] vrf group id : %d\n", vrfGroupID)
		vrfGroup, err = c.GetVRFGroupByID(vrfGroupID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get vrf group with id " + strconv.Itoa(vrfGroupID),
				Detail:   err.Error(),
			})
			return diags
		}
	} else if vrfGroupName != "" {
		fmt.Printf("[DEBUG] vrf group name %s\n", vrfGroupName)
		vrfGroup, err = c.GetVRFGroupByName(vrfGroupName)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get vrf group with name " + vrfGroupName,
				Detail:   err.Error(),
			})
			return diags
		}
	}

	fmt.Printf("[DEBUG] vrf group : %v\n", vrfGroup)

	c.WriteToDebugLog(fmt.Sprintf("%v", vrfGroup))

	buildings := make([]int, len(d.Get("building_ids").([]interface{})))

	for i, v := range d.Get("building_ids").([]interface{}) {
		b, err := c.GetBuildingByID(v.(int))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get building with id " + strconv.Itoa(v.(int)),
				Detail:   err.Error(),
			})
			return diags
		}
		buildings[i] = (*b).BuildingID
	}

	_ = d.Set("name", vrfGroup.Name)
	_ = d.Set("description", vrfGroup.Description)
	_ = d.Set("building_ids", buildings)

	d.SetId(strconv.Itoa(vrfGroup.ID))

	return diags
}
