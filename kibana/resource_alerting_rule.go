package kibana

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gk "github.com/renato0307/go-kibana/kibana"
)

func resourceAlertingRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertingRuleCreate,
		ReadContext:   resourceAlertingRuleRead,
		UpdateContext: resourceAlertingRuleUpdate,
		DeleteContext: resourceAlertingRuleDelete,
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeSet,
				Required: true,
				Computed: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							Computed: false,
						},
						"group": {
							Type:     schema.TypeString,
							Required: true,
							Computed: false,
						},
						"params": {
							Type:     schema.TypeString,
							Required: true,
							Computed: false,
						},
					},
				},
			},
			"api_key_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consumer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"last_execution_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_execution_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mute_all": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"muted_alert_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notify_when": {
				Type:     schema.TypeString,
				Required: true,
			},
			"param_agg_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_agg_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_es_query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_group_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_index": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"param_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"param_term_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_term_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"param_threshold": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"param_threshold_comparator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_time_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param_time_window_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"param_time_window_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_interval": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scheduled_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"throttle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlertingRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c := m.(*gk.Client)

	// sets rule actions
	actions := []gk.RuleAction{}
	if v, ok := d.GetOk("action"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			v := v.(map[string]interface{})
			var b map[string]interface{}
			json.Unmarshal([]byte(v["params"].(string)), &b)
			action := gk.RuleAction{
				ID:     v["id"].(string),
				Group:  v["group"].(string),
				Params: b,
			}
			actions = append(actions, action)
		}
	}

	// sets the Params.Index
	var paramsIndexSlice []string
	for _, param := range d.Get("param_index").([]interface{}) {
		paramsIndexSlice = append(paramsIndexSlice, param.(string))
	}

	// sets the Params.Threshold
	var paramsThresholdSlice []int
	for _, param := range d.Get("param_threshold").([]interface{}) {
		paramsThresholdSlice = append(paramsThresholdSlice, param.(int))
	}

	// sets the Tags
	var tagsSlice []string
	for _, param := range d.Get("tags").([]interface{}) {
		tagsSlice = append(tagsSlice, param.(string))
	}

	// sets the rest of the fields for the rule
	rule := gk.CreateRule{
		Actions:    actions,
		Consumer:   d.Get("consumer").(string),
		Name:       d.Get("name").(string),
		NotifyWhen: d.Get("notify_when").(string),
		Params: gk.RuleParams{
			AggField:            d.Get("param_agg_field").(string),
			AggType:             d.Get("param_agg_type").(string),
			ESQuery:             d.Get("param_es_query").(string),
			GroupBy:             d.Get("param_group_by").(string),
			Index:               paramsIndexSlice,
			Size:                d.Get("param_size").(int),
			TermField:           d.Get("param_term_field").(string),
			TermSize:            d.Get("param_term_size").(int),
			Threshold:           paramsThresholdSlice,
			ThresholdComparator: d.Get("param_threshold_comparator").(string),
			TimeField:           d.Get("param_time_field").(string),
			TimeWindowSize:      d.Get("param_time_window_size").(int),
			TimeWindowUnit:      d.Get("param_time_window_unit").(string),
		},
		RuleTypeID: d.Get("rule_type_id").(string),
		Schedule:   gk.RuleSchedule{Interval: d.Get("schedule_interval").(string)},
		Tags:       tagsSlice,
	}

	// calls API to create the rule
	newRule, err := c.CreateRule(rule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newRule.ID)

	// reads the created rule
	resourceAlertingRuleRead(ctx, d, m)

	return diags
}

func resourceAlertingRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c := m.(*gk.Client)

	ruleID := d.Id()

	rule, err := c.GetRule(ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("actions", rule.Actions)
	d.Set("api_key_owner", rule.ApiKeyOwner)
	d.Set("consumer", rule.Consumer)
	d.Set("created_at", rule.CreatedAt)
	d.Set("created_by", rule.CreatedBy)
	d.Set("enabled", rule.Enabled)
	d.Set("last_execution_date", rule.ExecutionStatus.LastExecutionDate)
	d.Set("last_execution_status", rule.ExecutionStatus.Status)
	d.Set("id", rule.ID)
	d.Set("mute_all", rule.MuteAll)
	d.Set("muted_alert_ids", rule.MutedAlertIDs)
	d.Set("name", rule.Name)
	d.Set("notify_when", rule.NotifyWhen)
	d.Set("param_agg_field", rule.Params.AggField)
	d.Set("param_agg_type", rule.Params.AggType)
	d.Set("param_es_query", rule.Params.ESQuery)
	d.Set("param_group_by", rule.Params.GroupBy)
	d.Set("param_index", rule.Params.Index)
	d.Set("param_term_field", rule.Params.TermField)
	d.Set("param_term_size", rule.Params.TermSize)
	d.Set("param_threshold", rule.Params.Threshold)
	d.Set("param_time_field", rule.Params.TimeField)
	d.Set("param_time_window_size", rule.Params.TimeWindowSize)
	d.Set("param_time_window_unit", rule.Params.TimeWindowUnit)
	d.Set("rule_type_id", rule.RuleTypeID)
	d.Set("schedule_interval", rule.Schedule.Interval)
	d.Set("scheduled_task_id", rule.ScheduledTaskId)
	d.Set("tags", rule.Tags)
	d.Set("throttle", rule.Throttle)
	d.Set("updated_at", rule.UpdatedAt)
	d.Set("updated_by", rule.UpdatedBy)

	return diags
}

func resourceAlertingRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*gk.Client)

	ruleID := d.Id()

	// sets rule actions
	actions := []gk.RuleAction{}
	if v, ok := d.GetOk("action"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			v := v.(map[string]interface{})
			var b map[string]interface{}
			json.Unmarshal([]byte(v["params"].(string)), &b)
			action := gk.RuleAction{
				ID:     v["id"].(string),
				Group:  v["group"].(string),
				Params: b,
			}
			actions = append(actions, action)
		}
	}

	// sets the Params.Index
	var paramsIndexSlice []string
	for _, param := range d.Get("param_index").([]interface{}) {
		paramsIndexSlice = append(paramsIndexSlice, param.(string))
	}

	// sets the Params.Threshold
	var paramsThresholdSlice []int
	for _, param := range d.Get("param_threshold").([]interface{}) {
		paramsThresholdSlice = append(paramsThresholdSlice, param.(int))
	}

	// sets the Tags
	var tagsSlice []string
	for _, param := range d.Get("tags").([]interface{}) {
		tagsSlice = append(tagsSlice, param.(string))
	}

	// sets the rest of the fields for the rule
	rule := gk.UpdateRule{
		Actions:    actions,
		Name:       d.Get("name").(string),
		NotifyWhen: d.Get("notify_when").(string),
		Params: gk.RuleParams{
			AggField:            d.Get("param_agg_field").(string),
			AggType:             d.Get("param_agg_type").(string),
			ESQuery:             d.Get("param_es_query").(string),
			GroupBy:             d.Get("param_group_by").(string),
			Index:               paramsIndexSlice,
			Size:                d.Get("param_size").(int),
			TermField:           d.Get("param_term_field").(string),
			TermSize:            d.Get("param_term_size").(int),
			Threshold:           paramsThresholdSlice,
			ThresholdComparator: d.Get("param_threshold_comparator").(string),
			TimeField:           d.Get("param_time_field").(string),
			TimeWindowSize:      d.Get("param_time_window_size").(int),
			TimeWindowUnit:      d.Get("param_time_window_unit").(string),
		},
		Throttle: d.Get("throttle").(string),
		Schedule: gk.RuleSchedule{Interval: d.Get("schedule_interval").(string)},
		Tags:     tagsSlice,
	}

	// calls API to update the rule
	updatedRule, err := c.UpdateRule(ruleID, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(updatedRule.ID)
	d.Set("last_updated", time.Now().Format(time.RFC850))

	// reads the updated rule and returns
	return resourceAlertingRuleRead(ctx, d, m)
}

func resourceAlertingRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics // Warning or errors can be collected in a slice type
	c := m.(*gk.Client)

	ruleID := d.Id()

	err := c.DeleteRule(ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
