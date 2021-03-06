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

// resourceAlertingRuleCreate - creates an alerting Rule
func resourceAlertingRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// maps the resource data to an RuleCreate struct
	rule, err := expandCreateRule(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// calls API to create the rule
	c := m.(*gk.Client)
	newRule, err := c.CreateRule(rule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newRule.ID)

	// reads the created rule
	return resourceAlertingRuleRead(ctx, d, m)
}

// resourceAlertingRuleRead - reads an alerting rule
func resourceAlertingRuleRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	ruleID := d.Id()

	// reads the rule from Kibana
	c := m.(*gk.Client)
	rule, err := c.GetRule(ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	// maps Rule to the resource data
	err = flattenRule(d, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	
	return nil
}

// resourceAlertingRuleUpdate - updates an alerting rule
func resourceAlertingRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// maps the resource data to an RuleUpdate struct
	ruleID := d.Id()
	rule, err := expandRuleUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// calls API to update the rule
	c := m.(*gk.Client)
	updatedRule, err := c.UpdateRule(ruleID, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	// sets common fields
	d.SetId(updatedRule.ID)
	_ = d.Set("last_updated", time.Now().Format(time.RFC850))

	// reads the updated rule and returns
	return resourceAlertingRuleRead(ctx, d, m)
}

// resourceAlertingRuleDelete - deletes an alerting rule
func resourceAlertingRuleDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	ruleID := d.Id()

	c := m.(*gk.Client)
	err := c.DeleteRule(ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness.
	d.SetId("")

	return nil
}

// Expand and flatten functions

// flattenRule - fills the resource data from a Rule
func flattenRule(d *schema.ResourceData, rule *gk.Rule) error {

	if len(rule.Actions) > 0 {
		actions, err := flattenRuleActions(rule)
		if err != nil {
			return err
		}
		err = d.Set("action", actions)
		if err != nil {
			return err
		}
	}

	_ = d.Set("api_key_owner", rule.ApiKeyOwner)
	_ = d.Set("consumer", rule.Consumer)
	_ = d.Set("created_at", rule.CreatedAt)
	_ = d.Set("created_by", rule.CreatedBy)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("last_execution_date", rule.ExecutionStatus.LastExecutionDate)
	_ = d.Set("last_execution_status", rule.ExecutionStatus.Status)
	_ = d.Set("id", rule.ID)
	_ = d.Set("mute_all", rule.MuteAll)
	_ = d.Set("muted_alert_ids", rule.MutedAlertIDs)
	_ = d.Set("name", rule.Name)
	_ = d.Set("notify_when", rule.NotifyWhen)
	_ = d.Set("param_agg_field", rule.Params.AggField)
	_ = d.Set("param_agg_type", rule.Params.AggType)
	_ = d.Set("param_es_query", rule.Params.ESQuery)
	_ = d.Set("param_group_by", rule.Params.GroupBy)
	_ = d.Set("param_index", rule.Params.Index)
	_ = d.Set("param_term_field", rule.Params.TermField)
	_ = d.Set("param_term_size", rule.Params.TermSize)
	_ = d.Set("param_threshold", rule.Params.Threshold)
	_ = d.Set("param_time_field", rule.Params.TimeField)
	_ = d.Set("param_time_window_size", rule.Params.TimeWindowSize)
	_ = d.Set("param_time_window_unit", rule.Params.TimeWindowUnit)
	_ = d.Set("rule_type_id", rule.RuleTypeID)
	_ = d.Set("schedule_interval", rule.Schedule.Interval)
	_ = d.Set("scheduled_task_id", rule.ScheduledTaskId)
	_ = d.Set("tags", rule.Tags)
	_ = d.Set("throttle", rule.Throttle)
	_ = d.Set("updated_at", rule.UpdatedAt)
	_ = d.Set("updated_by", rule.UpdatedBy)

	return nil
}

func flattenRuleActions(rule *gk.Rule) ([]interface{}, error) {
	log.Printf("flattenRule - number of actions found: %d", len(rule.Actions))
	var actions []interface{}
	for _, a := range rule.Actions {
		params, err := json.Marshal(a.Params)
		if err != nil {
			return nil, err
		}

		paramsString := string(params)
		m := map[string]interface{}{
			"id":     a.ID,
			"group":  a.Group,
			"params": paramsString,
		}
		actions = append(actions, m)
	}
	return actions, nil
}

func expandRuleUpdate(d *schema.ResourceData) (gk.UpdateRule, error) {
	// sets rule actions
	var actions []gk.RuleAction
	if v, ok := d.GetOk("action"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			v := v.(map[string]interface{})
			var b map[string]interface{}
			err := json.Unmarshal([]byte(v["params"].(string)), &b)
			if err != nil {
				return gk.UpdateRule{}, err
			}
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
	return rule, nil
}

func expandCreateRule(d *schema.ResourceData) (gk.CreateRule, error) {
	// sets rule actions
	var actions []gk.RuleAction
	if v, ok := d.GetOk("action"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			v := v.(map[string]interface{})
			var b map[string]interface{}
			err := json.Unmarshal([]byte(v["params"].(string)), &b)
			if err != nil {
				return gk.CreateRule{}, err
			}
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
	return rule, nil
}
