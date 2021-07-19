package kibana

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
			"config": {
				Type:     schema.TypeString,
				Required: true,
				Computed: false,
			},
			"secrets": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  false,
				Sensitive: true,
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
	}

	log.Printf("Unmarshalling config: %s", d.Get("config").(string))
	var config map[string]interface{}
	configValue, ok := d.GetOk("config")
	if ok {
		err := json.Unmarshal([]byte(configValue.(string)), &config)
		connector.Config = config
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("Unmarshalling secrets: %s", d.Get("secrets").(string))
	var secrets map[string]interface{}
	secretsValue, ok := d.GetOk("secrets")
	if ok {
		err := json.Unmarshal([]byte(secretsValue.(string)), &secrets)
		connector.Secrets = secrets
		if err != nil {
			return diag.FromErr(err)
		}
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

	if connector.Config != nil {
		configValue, err := json.Marshal(connector.Config)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("config", string(configValue))
	}

	if connector.Secrets != nil {
		secretsValue, err := json.Marshal(connector.Secrets)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("secrets", string(secretsValue))
	}

	return diags
}

func resourceActionsConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*gk.Client)

	connector := gk.UpdateConnector{
		Name: d.Get("name").(string),
	}

	log.Printf("Unmarshalling config: %s", d.Get("config").(string))
	var config map[string]interface{}
	configValue, ok := d.GetOk("config")
	if ok {
		err := json.Unmarshal([]byte(configValue.(string)), &config)
		connector.Config = config
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("Unmarshalling secrets: %s", d.Get("secrets").(string))
	var secrets map[string]interface{}
	secretsValue, ok := d.GetOk("secrets")
	if ok {
		err := json.Unmarshal([]byte(secretsValue.(string)), &secrets)
		connector.Secrets = secrets
		if err != nil {
			return diag.FromErr(err)
		}
	}

	connectorID := d.Id()
	_, err := c.UpdateConnector(connectorID, connector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceActionsConnectorRead(ctx, d, m)
}

func resourceActionsConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c := m.(*gk.Client)

	connectorID := d.Id()

	err := c.DeleteConnector(connectorID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
