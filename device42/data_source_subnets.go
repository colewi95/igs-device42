package device42

import (
	"context"
	"fmt"
	"net"
	"strconv"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSubnets() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_vrf_groups` can be used to retrieve all VRF groups.",
		ReadContext: dataSourceSubnetsRead,
		Schema: map[string]*schema.Schema{
			"vrf_group_id": &schema.Schema{
				Description: "Filter by `vrf_group_id`",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"parent_subnet_id": &schema.Schema{
				Description: "Filter by `parent_subnet_id`",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"subnets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All `subnets`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The `id` of the subnet.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "The `name` of the subnet.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"network": &schema.Schema{
							Description: "The `network` of the subnet.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"gateway": &schema.Schema{
							Description: "The `gateway` of the subnet.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"vlan": &schema.Schema{
							Description: "The `vlan` of the subnet.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"mask_bits": &schema.Schema{
							Description: "The `mask bits` of the subnet.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"mask": &schema.Schema{
							Description: "The `mask` of the subnet.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"is_supernet": &schema.Schema{
							Description: "Is this subnet a supernet?",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"parent_subnet_id": &schema.Schema{
							Description: "The `mask bits` of the subnet.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// get subnets
func dataSourceSubnetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	vrfGroupID := d.Get("vrf_group_id").(int)
	parentSubnetID := d.Get("parent_subnet_id").(int)

	var diags diag.Diagnostics
	var err error
	var subnets *[]device42.Subnet

	if vrfGroupID != 0 && parentSubnetID != 0 {
		subnets, err = c.GetSubnetsByParentSubnetIDWithVRFGroupID(parentSubnetID, vrfGroupID)
	} else if vrfGroupID != 0 {
		subnets, err = c.GetSubnetsByVRFGroupID(vrfGroupID)
	} else if parentSubnetID != 0 {
		subnets, err = c.GetSubnetsByParentSubnetID(parentSubnetID)
	} else {
		subnets, err = c.GetSubnets()
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get a list of subnets",
			Detail:   err.Error(),
		})
		return diags
	}

	s := flattenSubnetsData(subnets, d)
	if err := d.Set("subnets", s); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to set subnets",
			Detail:   err.Error(),
		})
		return diags
	}

	ids := make([]int, len((*subnets)))
	for _, i := range *subnets {
		ids = append(ids, i.SubnetID)
	}

	checksum := idsChecksum(ids)

	d.SetId(checksum)

	return diags
}

// flatten vrf groups to a map
func flattenSubnetsData(subnets *[]device42.Subnet, d *schema.ResourceData) []interface{} {
	if subnets != nil {
		ss := make([]interface{}, len(*subnets))

		for i, subnet := range *subnets {
			s := make(map[string]interface{})

			_, ipv4Net, _ := net.ParseCIDR(subnet.Network + "/" + strconv.Itoa(subnet.MaskBits))
			s["mask"] = ipv4MaskString(ipv4Net.Mask)

			s["gateway"] = subnet.Gateway
			s["id"] = subnet.SubnetID
			s["name"] = subnet.Name
			s["network"] = subnet.Network
			s["mask_bits"] = subnet.MaskBits
			s["vlan"] = fmt.Sprintf("%v", subnet.ParentVlanNumber)
			s["parent_subnet_id"] = subnet.ParentSubnetID

			ss[i] = s
		}

		return ss
	}

	return make([]interface{}, 0)
}
