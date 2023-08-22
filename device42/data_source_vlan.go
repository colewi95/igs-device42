package device42

import (
	"context"
	"fmt"
	"log"
	"strconv"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVLAN() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_vlan' data source can be used to retrieve a single VLAN using its `id`",
		ReadContext: dataSourceVLANRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description: "The `id` of a VLAN.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": &schema.Schema{
				Description: "The `name` of the VLAN.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"number": &schema.Schema{
				Description: "The VLAN `number.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"tags": &schema.Schema{
				Description: "All`tags` for a VLAN.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// get a building by id
func dataSourceVLANRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	vlanID := d.Get("id").(int)
	vlan := &device42.VLAN{}

	if vlanID != 0 {
		log.Printf("[DEBUG] VLAN id: %d\n", vlanID)

		vlan, err = c.GetVLANByID(vlanID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get VLAN with id " + strconv.Itoa(vlanID),
				Detail:   err.Error(),
			})
			return diags
		}
	}

	c.WriteToDebugLog(fmt.Sprintf("%v", vlan))

	_ = d.Set("name", vlan.Name)
	_ = d.Set("number", vlan.Number)
	_ = d.Set("tags", vlan.Tags)

	d.SetId(strconv.Itoa(vlan.VlanID))

	return diags
}
