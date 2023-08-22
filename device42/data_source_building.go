package device42

import (
	"context"
	"log"
	"strconv"

	device42 "github.com/chopnico/device42-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBuilding() *schema.Resource {
	return &schema.Resource{
		Description: "`device42_building` data source can be used to retrieve a single building using its `id` or `name`.",
		ReadContext: dataSourceBuildingRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				AtLeastOneOf: []string{"id", "name"},
				Description:  "The `id` of a building.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"id", "name"},
				Description:  "The `name` of a building.",
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The `address` of a building.",
			},
			"notes": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "`notes` for a building.",
			},
		},
	}
}

// get a building by id
func dataSourceBuildingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	var diags diag.Diagnostics
	var err error

	buildingID := d.Get("id").(int)
	buildingName := d.Get("name").(string)
	building := &device42.Building{}

	if buildingID != 0 {
		log.Printf("[DEBUG] building id : %d", buildingID)
		building, err = c.GetBuildingByID(buildingID)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get building with id " + strconv.Itoa(buildingID),
				Detail:   err.Error(),
			})
			return diags
		}
	} else if buildingName != "" {
		log.Printf("[DEBUG] building name : %s", buildingName)
		building, err = c.GetBuildingByName(buildingName)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to get building with name " + buildingName,
				Detail:   err.Error(),
			})
			return diags
		}
	}

	log.Printf("[DEBUG] building : %v", building)

	_ = d.Set("name", building.Name)
	_ = d.Set("address", building.Address)
	_ = d.Set("notes", building.Notes)

	d.SetId(strconv.Itoa(building.BuildingID))

	return diags
}
