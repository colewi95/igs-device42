package device42

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description:   "`device42_subnet` resource can be used to create, update or delete a subnet.",
		CreateContext: resourceSubnetSet,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetSet,
		DeleteContext: resourceSubnetDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Description: "The last time this resource was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "The `name` of the subnet.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"network": &schema.Schema{
				Description: "The `network` of the subnet.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"gateway": &schema.Schema{
				Description: "The `gateway` of the subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mask_bits": &schema.Schema{
				Description: "The `mask_bits` of the subnet.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"mask": &schema.Schema{
				Description: "The `mask` of the subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vrf_group_id": &schema.Schema{
				Description: "The `vrf_group_id` of the subnet.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"is_supernet": &schema.Schema{
				Description: "Is this subnet a supernet?",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tags": &schema.Schema{
				Description: "The `tags` for this subnet.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSubnetSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Println(fmt.Sprintf("[DEBUG] subnet : %s", d.Get("subnets")))

	tags := interfaceSliceToStringSlice(d.Get("tags").([]interface{}))

	subnet, err := c.SetSubnet(&device42.Subnet{
		Name:       d.Get("name").(string),
		Network:    d.Get("network").(string),
		MaskBits:   d.Get("mask_bits").(int),
		VrfGroupID: d.Get("vrf_group_id").(int),
		Tags:       tags,
	})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create subnet with name " + d.Get("name").(string),
			Detail:   err.Error(),
		})
	}

	if !d.Get("is_supernet").(bool) {
		subnet.Gateway = ipv4GatewayFromNetwork(subnet.Network)
	}
	subnet, err = c.SetSubnet(subnet)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create subnet with name " + d.Get("name").(string),
			Detail:   err.Error(),
		})
	}

	log.Println(fmt.Sprintf("[DEBUG] subnet : %v", subnet))

	d.SetId(strconv.Itoa(subnet.SubnetID))

	resourceSubnetRead(ctx, d, m)

	return diags
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	subnetID, err := strconv.Atoi(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read id",
			Detail:   err.Error(),
		})
		return diags
	}
	subnet, err := c.GetSubnetByID(subnetID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get subnet with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] subnet : %v", subnet))

	_, ipv4Net, err := net.ParseCIDR(subnet.Network + "/" + strconv.Itoa(subnet.MaskBits))
	_ = d.Set("mask", ipv4MaskString(ipv4Net.Mask))

	_ = d.Set("is_supernet", d.Get("is_supernet").(bool))
	_ = d.Set("gateway", subnet.Gateway)
	_ = d.Set("name", subnet.Name)
	_ = d.Set("network", subnet.Network)
	_ = d.Set("mask_bits", subnet.MaskBits)
	_ = d.Set("vrf_group_id", subnet.VrfGroupID)
	_ = d.Set("tags", subnet.Tags)

	return diags
}

// delete subnet
func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	_, err := fmt.Sscan(d.Id(), &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get subnet id",
			Detail:   err.Error(),
		})
		return diags
	}

	err = c.DeleteSubnet(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete subnet with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
