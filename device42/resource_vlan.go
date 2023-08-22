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

func resourceVLAN() *schema.Resource {
	return &schema.Resource{
		Description:   "`device42_vlan` resource can be used to create, update or delete a VLAN.",
		CreateContext: resourceVLANSet,
		ReadContext:   resourceVLANRead,
		UpdateContext: resourceVLANSet,
		DeleteContext: resourceVLANDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Description: "The last time this resource was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "The `name` of the VLAN.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"number": &schema.Schema{
				Description: "The VLAN `number.`",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"tags": &schema.Schema{
				Description: "The `tags` for this VLAN.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceVLANSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Println(fmt.Sprintf("[DEBUG] VLAN : %s", d.Get("vlan")))

	tags := interfaceSliceToStringSlice(d.Get("tags").([]interface{}))

	vlan, err := c.SetVLAN(&device42.VLAN{
		Name:   d.Get("name").(string),
		Number: d.Get("number").(int),
		Tags:   tags,
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create VLAN with name " + d.Get("name").(string),
			Detail:   err.Error(),
		})
	}

	log.Println(fmt.Sprintf("[DEBUG] VLAN : %v", vlan))

	d.SetId(strconv.Itoa(vlan.VlanID))

	resourceVLANRead(ctx, d, m)

	return diags
}

func resourceVLANRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	vlanID, err := strconv.Atoi(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read id",
			Detail:   err.Error(),
		})
		return diags
	}
	vlan, err := c.GetVLANByID(vlanID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get VLAN with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] VLAN : %v", vlan))

	_ = d.Set("name", vlan.Name)
	_ = d.Set("number", vlan.Number)
	_ = d.Set("tags", vlan.Tags)

	return diags
}

// delete vlan
func resourceVLANDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	_, err := fmt.Sscan(d.Id(), &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get VLAN id",
			Detail:   err.Error(),
		})
		return diags
	}

	err = c.DeleteVLAN(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete VLAN with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
