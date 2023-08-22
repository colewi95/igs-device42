package device42

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_subnet` data source can be used to retrieve a single subnet using its `name` and `network` or just by its `id`.",
		ReadContext: dataSourceSubnetRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description: "The `id` of a subnet.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"name": &schema.Schema{
				Description: "The `name` of the subnet.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"network": &schema.Schema{
				Description: "The `network` of the subnet. (e.g., 192.168.0.0/24)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mask_bits": &schema.Schema{
				Description: "The `mask_bits` of the subnet.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"mask": &schema.Schema{
				Description: "The `mask` of the subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"gateway": &schema.Schema{
				Description: "The `gateway` of the subnet.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vrf_group_id": &schema.Schema{
				Description: "The `vrf_group_id` of the subnet.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"is_supernet": &schema.Schema{
				Description: "Is this subnet a supernet?",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"tags": &schema.Schema{
				Description: "All`tags` for a subnet.",
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
func dataSourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	subnetID := d.Get("id").(int)
	vrfGroupID := d.Get("vrf_group_id").(int)
	subnetName := d.Get("name").(string)
	network := d.Get("network").(string)
	subnet := &device42.Subnet{}

	if subnetID != 0 {
		log.Printf("[DEBUG] subnet id: %d\n", subnetID)

		subnet, err = c.GetSubnetByID(subnetID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get subnet with id " + strconv.Itoa(subnetID),
				Detail:   err.Error(),
			})
			return diags
		}
	} else if subnetName != "" && network != "" {
		log.Printf("[DEBUG] subnet name: %s\n", subnetName)

		subnet, err = c.GetSubnetByNameWithNetwork(subnetName, network)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get subnet with name " + subnetName,
				Detail:   err.Error(),
			})
			return diags
		}
	} else if subnetName != "" {
		log.Printf("[DEBUG] subnet name: %s\n", subnetName)

		subnet, err = c.GetSubnetByNameWithVRFGroupID(subnetName, vrfGroupID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get subnet with name " + subnetName,
				Detail:   err.Error(),
			})
			return diags
		}
	}

	c.WriteToDebugLog(fmt.Sprintf("%v", subnet))

	_, ipv4Net, err := net.ParseCIDR(subnet.Network + "/" + strconv.Itoa(subnet.MaskBits))
	_ = d.Set("mask", ipv4MaskString(ipv4Net.Mask))

	_ = d.Set("is_supernet", d.Get("is_supernet").(bool))
	_ = d.Set("gateway", subnet.Gateway)
	_ = d.Set("name", subnet.Name)
	_ = d.Set("network", subnet.Network)
	_ = d.Set("mask_bits", subnet.MaskBits)
	_ = d.Set("vrf_group_id", subnet.VrfGroupID)
	_ = d.Set("tags", subnet.Tags)

	d.SetId(strconv.Itoa(subnet.SubnetID))

	return diags
}
