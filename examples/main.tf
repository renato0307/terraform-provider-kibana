
resource "kibana_actions_connector" "sample_connector" {
  name                        = "test_connector_created_with_custom_provider"
  connector_type_id           = ".index"
  config_index                = "test_connector_created_with_custom_provider_index"
  config_refresh              = true
  config_execution_time_field = null
}
