
resource "kibana_actions_connector" "sample_connector" {
  config_execution_time_field = null
  config_index                = "test_connector_created_with_custom_provider_index"
  config_refresh              = true
  connector_type_id           = ".index"
  name                        = "test_connector_created_with_custom_provider"
}

resource "kibana_alerting_rule" "sample_rule" {
  action {
    id    = kibana_actions_connector.sample_connector.id
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
  param_index                = ["phoenix-orchestration-metrics-$workspace*"]
  param_size                 = 1
  param_threshold            = [1]
  param_threshold_comparator = ">"
  param_time_field           = "timestamp"
  param_time_window_size     = 5
  param_time_window_unit     = "m"
  rule_type_id               = ".es-query"
  schedule_interval          = "5m"
  tags                       = ["tag1", "tag2", "tag3"]
}
