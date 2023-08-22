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

func resourceDynamicIP() *schema.Resource {
	return &schema.Resource{
		Description:   "`device42_dynamic_ip` data resource can be used to generate a new IP",
		CreateContext: resourceDynamicIPSet,
		ReadContext:   resourceDynamicIPRead,
		UpdateContext: resourceDynamicIPSet,
		DeleteContext: resourceDynamicIPDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Description: "The last time this resource was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id": &schema.Schema{
				Description: "The `id` of this IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"address": &schema.Schema{
				Description: "The `address` of this IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mask_bits": &schema.Schema{
				Description: "The `mask_bits` for the IP.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"subnet": &schema.Schema{
				Description: "The `subnet` for the IP.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"subnet_id": &schema.Schema{
				Description: "The `subnet_id` for the IP.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"vrf_group": &schema.Schema{
				Description: "The `vrf_group` for the IP.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"vrf_group_id": &schema.Schema{
				Description: "The `vrf_group_id` for the IP.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}
func resourceDynamicIPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	ipID := d.Get("id").(string)
	ipSubnetID := d.Get("subnet_id").(int)
	ipVRFGroup := d.Get("vrf_group").(string)
	ipVRFGroupID := d.Get("vrf_group_id").(int)

	var id int
	_, err = fmt.Sscan(ipID, &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to suggest ip",
			Detail:   err.Error(),
		})
		return diags
	}

	ip, err := c.GetIPByID(id)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to suggest ip",
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] ip : %v", ip))

	_ = d.Set("subnet_id", ipSubnetID)
	_ = d.Set("vrf_group_id", ipVRFGroupID)
	_ = d.Set("vrf_group", ipVRFGroup)

	return diags
}

func resourceDynamicIPSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	log.Println(fmt.Sprintf("[DEBUG] ip : %s", d.Get("ip")))

	ipMaskBits := d.Get("mask_bits").(int)
	subnetID := d.Get("subnet_id").(int)
	ipVRFGroupID := d.Get("vrf_group_id").(int)

	ip := &device42.IP{}

	ip, err = c.SuggestIPWithVRFGroupID(ipVRFGroupID, subnetID, ipMaskBits, true)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to suggest ip",
			Detail:   err.Error(),
		})
		return diags
	}

	ip.VRFGroupID = ipVRFGroupID
	ip.SubnetID = subnetID

	log.Println(fmt.Sprintf("[DEBUG] ip : %v", ip))
	ip, err = c.UpdateIP(ip)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to suggest ip",
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] ip : %v", ip))

	_ = d.Set("label", ip.Label)
	_ = d.Set("address", ip.Address)
	_ = d.Set("subnet_id", ip.SubnetID)
	_ = d.Set("vrf_group_id", ipVRFGroupID)
	_ = d.Set("vrf_group", ip.VRFGroup)

	d.SetId(strconv.Itoa(ip.ID))

	resourceDynamicIPRead(ctx, d, m)

	return diags
}

func resourceDynamicIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	if ipID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ip id is 0",
			Detail:   "the current id of this ip 0. not sure why.",
		})
		return diags
	}
	ip, err := c.GetIPByID(ipID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get ip with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] ip : %v", ip))

	_ = d.Set("label", ip.Label)
	_ = d.Set("address", ip.Address)
	_ = d.Set("subnet", ip.Subnet)
	_ = d.Set("subnet_id", ip.SubnetID)
	_ = d.Set("vrf_group", ip.VRFGroup)

	return diags
}

// delete ip
func resourceDynamicIPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	_, err := fmt.Sscan(d.Id(), &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get ip id",
			Detail:   err.Error(),
		})
		return diags
	}

	err = c.DeleteIP(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete ip with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
