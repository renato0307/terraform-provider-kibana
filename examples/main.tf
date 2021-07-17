
resource "kibana_actions_connector" "sample_connector" {
  config_execution_time_field = null
  config_index                = "test_connector_created_with_custom_provider_index"
  config_refresh              = true
  connector_type_id           = ".index"
  name                        = "test_connector_created_with_custom_provider"
}
