package device42

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVRFGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "`device42_vrf_group` resource can be used to create, update, or delete a VRF group.",
		CreateContext: resourceVRFGroupSet,
		ReadContext:   resourceVRFGroupRead,
		UpdateContext: resourceVRFGroupSet,
		DeleteContext: resourceVRFGroupDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Description: "The last time this resource was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "The `name` of the VRF group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": &schema.Schema{
				Description: "The `description` of the VRF group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"building_ids": &schema.Schema{
				Description: "The `building_ids` of the VRF group.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceVRFGroupSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Println(fmt.Sprintf("[DEBUG] buildings : %s", d.Get("buildings")))

	buildings := make([]string, len(d.Get("building_ids").([]interface{})))

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
		buildings[i] = (*b).Name
	}

	vrfGroup, err := c.SetVRFGroup(&device42.VRFGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Buildings:   buildings,
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create vrf group with name " + d.Get("name").(string),
			Detail:   err.Error(),
		})
	}

	log.Println(fmt.Sprintf("[DEBUG] vrf group : %v", vrfGroup))

	d.SetId(strconv.Itoa(vrfGroup.ID))

	resourceVRFGroupRead(ctx, d, m)

	return diags
}

func resourceVRFGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vrfGroupID, err := strconv.Atoi(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read id",
			Detail:   err.Error(),
		})
		return diags
	}
	vrfGroup, err := c.GetVRFGroupByID(vrfGroupID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get vrf group with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] vrf group : %v", vrfGroup))

	_ = d.Set("name", vrfGroup.Name)
	_ = d.Set("description", vrfGroup.Description)
	_ = d.Set("buildings", vrfGroup.Buildings)

	return diags
}

// delete vrf group
func resourceVRFGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	_, err := fmt.Sscan(d.Id(), &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get vrf group id",
			Detail:   err.Error(),
		})
		return diags
	}

	err = c.DeleteVRFGroup(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete vrf group with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
