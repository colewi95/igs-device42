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

func resourceBuilding() *schema.Resource {
	return &schema.Resource{
		Description:   "`device42_building` resource can be used to create, update and delete buildings.",
		CreateContext: resourceBuildingSet,
		ReadContext:   resourceBuildingRead,
		UpdateContext: resourceBuildingSet,
		DeleteContext: resourceBuildingDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Description: "When the resource was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "The `name` of the building.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"address": &schema.Schema{
				Description: "The `address` of the building.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"notes": &schema.Schema{
				Description: "`notes` for the building.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBuildingSet(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Println(fmt.Sprintf("[DEBUG] building name : %s", d.Get("name").(string)))

	building, err := c.SetBuilding(&device42.Building{
		Name:    d.Get("name").(string),
		Address: d.Get("address").(string),
		Notes:   d.Get("notes").(string),
	})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create building with name " + d.Get("name").(string),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] building : %v", building))

	d.SetId(strconv.Itoa(building.BuildingID))

	resourceBuildingRead(ctx, d, m)

	return diags
}

func resourceBuildingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	buildingID, err := strconv.Atoi(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read building id",
			Detail:   err.Error(),
		})
		return diags
	}
	building, err := c.GetBuildingByID(buildingID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get building with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	log.Println(fmt.Sprintf("[DEBUG] building : %v", building))

	_ = d.Set("name", building.Name)
	_ = d.Set("address", building.Address)
	_ = d.Set("notes", building.Notes)

	return diags
}

func resourceBuildingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*device42.API)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var id int
	_, err := fmt.Sscan(d.Id(), &id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get building id",
			Detail:   err.Error(),
		})
		return diags
	}

	err = c.DeleteBuilding(id)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to delete building with id " + d.Id(),
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
