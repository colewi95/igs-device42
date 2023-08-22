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

func resourceIP() *schema.Resource {
	return &schema.Resource{
		Description:   "`device42_ip` resource can be used to create, update or delete a IP.",
		CreateContext: resourceIPSet,
		ReadContext:   resourceIPRead,
		UpdateContext: resourceIPSet,
		DeleteContext: resourceIPDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Description: "The last time this resource was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id": &schema.Schema{
				Description: "The `id` of the IP.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"address": &schema.Schema{
				Description: "The IP `address`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"label": &schema.Schema{
				Description: "The IP `label`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vrf_group": &schema.Schema{
				Description: "The `vrf_group` for the IP IP.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"subnet": &schema.Schema{
				Description: "The `subnet` for the IP IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subnet_id": &schema.Schema{
				Description: "The `subnet_id` for the IP IP.",
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
			},
		},
	}
}

func resourceIPSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Println(fmt.Sprintf("[DEBUG] IP : %s", d.Get("ip")))

	address := d.Get("address").(string)
	label := d.Get("label").(string)
	subnetID := d.Get("subnet_id").(int)

	ip, err := c.SetIP(&device42.IP{
		IPAddress: address,
		Label:     label,
		SubnetID:  subnetID,
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create IP with address " + d.Get("address").(string),
			Detail:   err.Error(),
		})
	}

	log.Println(fmt.Sprintf("[DEBUG] IP : %v", ip))

	d.SetId(strconv.Itoa(ip.ID))

	resourceIPRead(ctx, d, m)

	return diags
}

func resourceIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	ipID, err := strconv.Atoi(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read id",
			Detail:   err.Error(),
		})
		return diags
	}
	ip, err := c.GetIPByID(ipID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get IP with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] IP : %v", ip))

	_ = d.Set("id", ip.ID)
	_ = d.Set("address", ip.Address)
	_ = d.Set("label", ip.Label)
	_ = d.Set("subnet", ip.Subnet)
	_ = d.Set("subnet_id", ip.SubnetID)
	_ = d.Set("vrf_group", ip.VRFGroup)

	return diags
}

// delete ip
func resourceIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	_, err := fmt.Sscan(d.Id(), &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get IP id",
			Detail:   err.Error(),
		})
		return diags
	}

	err = c.DeleteIP(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete IP with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
