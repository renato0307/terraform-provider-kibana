package kibana

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gk "github.com/renato0307/go-kibana/kibana"
)

// Provider - Kibana Terraform provider definition
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_PASSWORD", nil),
			},
			"space": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_SPACE", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kibana_actions_connector": resourceActionsConnector(),
			"kibana_alerting_rule":     resourceAlertingRule(),
		},
		// DataSourcesMap: map[string]*schema.Resource{
		// 	"hashicups_coffees":     dataSourceCoffees(),
		// 	"hashicups_ingredients": dataSourceIngredients(),
		// 	"hashicups_order":       dataSourceOrder(),
		// },
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	space := d.Get("space").(string)
	host := d.Get("host").(string)

	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c, err := gk.NewClient(&host, &username, &password, &space)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Kibana client",
			Detail:   "Unable to create Kibana client with basic authentication",
		})
		return nil, diags
	}

	return c, diags
}
