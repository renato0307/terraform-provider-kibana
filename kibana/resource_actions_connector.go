package kibana

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gk "github.com/renato0307/go-kibana/kibana"
)

func resourceActionsConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActionsConnectorCreate,
		ReadContext:   resourceActionsConnectorRead,
		UpdateContext: resourceActionsConnectorUpdate,
		DeleteContext: resourceActionsConnectorDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"connector_type_id": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"is_preconfigured": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
			},
			"is_missing_secrets": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
			},
			"config_execution_time_field": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"config_index": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"config_refresh": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
			},
		},
	}
}

func resourceActionsConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c := m.(*gk.Client)

	connector := gk.CreateConnector{
		Name:            d.Get("name").(string),
		ConnectorTypeId: d.Get("connector_type_id").(string),
		Config: gk.ConnectorConfig{
			ExecutionTimeField: d.Get("config_execution_time_field").(string),
			Index:              d.Get("config_index").(string),
			Refresh:            d.Get("config_refresh").(bool),
		},
	}

	newConnector, err := c.CreateConnector(connector)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newConnector.ID)

	resourceActionsConnectorRead(ctx, d, m)

	return diags
}

func resourceActionsConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c := m.(*gk.Client)

	connectorId := d.Id()

	connector, err := c.GetConnector(connectorId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", connector.Name)
	d.Set("connector_type_id", connector.ConnectorTypeId)
	d.Set("is_preconfigured", connector.IsPreconfigured)
	d.Set("is_missing_secrets", connector.IsMissingSecrets)
	d.Set("config_execution_time_field", connector.Config.ExecutionTimeField)
	d.Set("config_index", connector.Config.Index)
	d.Set("config_refresh", connector.Config.Refresh)

	return diags
}

func resourceActionsConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// c := m.(*hc.Client)

	// orderID := d.Id()

	// if d.HasChange("items") {
	// 	items := d.Get("items").([]interface{})
	// 	ois := []hc.OrderItem{}

	// 	for _, item := range items {
	// 		i := item.(map[string]interface{})

	// 		co := i["coffee"].([]interface{})[0]
	// 		coffee := co.(map[string]interface{})

	// 		oi := hc.OrderItem{
	// 			Coffee: hc.Coffee{
	// 				ID: coffee["id"].(int),
	// 			},
	// 			Quantity: i["quantity"].(int),
	// 		}
	// 		ois = append(ois, oi)
	// 	}

	// 	_, err := c.UpdateOrder(orderID, ois)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}

	// 	d.Set("last_updated", time.Now().Format(time.RFC850))
	// }

	return resourceActionsConnectorRead(ctx, d, m)
}

func resourceActionsConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	// c := m.(*hc.Client)

	// orderID := d.Id()

	// err := c.DeleteOrder(orderID)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
