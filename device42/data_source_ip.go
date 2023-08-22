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

func dataSourceIP() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_ip` data source can be used to retrieve a single IP using its `address` and a `subnet_id`.",
		ReadContext: dataSourceIPRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description: "The `id` of an IP.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"label": &schema.Schema{
				Description: "The `lablel` of the IP.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"address": &schema.Schema{
				Description:  "The `address` of the IP.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"subnet_id"},
			},
			"mac_address": &schema.Schema{
				Description: "The `mac_address` of the IP.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subnet_id": &schema.Schema{
				Description:  "The `subnet_id` of the IP.",
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"address"},
			},
			"subnet": &schema.Schema{
				Description: "The `subnet` of the IP. (e.g., 192.168.0.0/24)",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// get a building by id
func dataSourceIPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	ipID := d.Get("id").(string)
	ipAddress := d.Get("address").(string)
	ipSubnetID := d.Get("subnet_id").(int)
	ip := &device42.IP{}

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

	if ipID != "" {
		log.Printf("[DEBUG] ip id: %s\n", ipID)

		ip, err = c.GetIPByID(id)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get ip with id " + ipID,
				Detail:   err.Error(),
			})
			return diags
		}
	} else if ipAddress != "" && ipSubnetID != 0 {
		log.Printf("[DEBUG] ip address: %s\n", ipAddress)

		ip, err = c.GetIPByAddressWithSubnetID(ipAddress, ipSubnetID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get subnet with address " + ipAddress + "and subnet id " + strconv.Itoa(ipSubnetID),
				Detail:   err.Error(),
			})
			return diags
		}
	}

	c.WriteToDebugLog(fmt.Sprintf("ip : %v", ip))

	_ = d.Set("address", ip.Address)
	_ = d.Set("subnet", ip.Subnet)
	_ = d.Set("subnet_id", ip.SubnetID)
	_ = d.Set("label", ip.Label)
	_ = d.Set("mac_address", ip.MacAddress)

	d.SetId(strconv.Itoa(ip.ID))

	return diags
}
