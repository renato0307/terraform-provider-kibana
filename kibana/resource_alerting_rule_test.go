package kibana

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gk "github.com/renato0307/go-kibana/kibana"
	"strings"
	"testing"
)

func TestAccKibanaAlertingRule_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckKibanaAlertingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKibanaAlertingRule(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaAlertingRuleExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckKibanaAlertingRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["kibana_alerting_rule." + resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		c := testAccProvider.Meta().(*gk.Client)
		_, err := c.GetRule(rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccKibanaAlertingRule(resourceName string) string {
	return fmt.Sprintf(`
			resource "kibana_actions_connector" "%s" {
			  connector_type_id = ".index"
			  name              = "test_connector_created_with_custom_provider"
			
			  config = jsonencode(
				{
				  "index" : "test-index",
				  "refresh" : true
				  "executionTimeField" : null
				}
			  )
			}

			resource "kibana_alerting_rule" "%s" {
			  action {
				id    = kibana_actions_connector.%s.id
				group = "query matched"
				params = jsonencode(
				  {
					"documents" : [
					  {
						"@timestamp" : "{{context.date}}",
						"tags" : "{{rule.tags}}",
						"rule" : {
						  "id" : "{{rule.id}}",
						  "name" : "{{rule.name}}",
						  "params" : { "{{rule.type}}" : "{{params}}" },
						  "space" : "{{rule.spaceId}}",
						  "type" : "{{rule.type}}"
						},
						"kibana" : {
						  "alert" : {
							"id" : "{{alert.id}}",
							"context" : { "{{rule.type}}" : "{{context}}" },
							"actionGroup" : "{{alert.actionGroup}}",
							"actionGroupName" : "{{alert.actionGroupName}}"
						  }
						},
						"event" : { "kind" : "alert" }
					  }
					]
				  }
				)
			  }
			
			  consumer    = "alerts"
			  enabled     = true
			  name        = "my-terraform-rule"
			  notify_when = "onActiveAlert"
			  param_es_query = jsonencode(
				{
				  "query" : {
					"bool" : {
					  "filter" : [
						{
						  "bool" : {
							"should" : [{ "range" : { "value.count" : { "gt" : "0" } } }],
							"minimum_should_match" : 1
						  }
						},
						{ "match_phrase" : { "namespace.keyword" : "AWS/SQS" } },
						{
						  "match_phrase" : { "metric_name.keyword" : "NumberOfMessagesReceived" }
						}
					  ]
					}
				  }
				}
			  )
			  param_index                = ["my-index*"]
			  param_size                 = 1
			  param_threshold            = [1]
			  param_threshold_comparator = ">"
			  param_time_field           = "timestamp"
			  param_time_window_size     = 5
			  param_time_window_unit     = "m"
			  rule_type_id               = ".es-query"
			  schedule_interval          = "5m"
			  tags                       = ["tag1", "tag2", "tag3"]
			}`, resourceName, resourceName, resourceName)
}

func testAccCheckKibanaAlertingRuleDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*gk.Client)

	// loop through the resources in state, verifying each widget
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_alerting_rule" {
			continue
		}

		_, err := c.GetRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("alerting rule (%s) still exists", rs.Primary.ID)
		}

		// If the error is equivalent to 404 not found, the widget is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}
	return nil
}