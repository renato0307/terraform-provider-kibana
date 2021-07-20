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

func TestAccKibanaActionsConnector_basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy: testAccCheckKibanaActionsConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKibanaActionsConnector(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKibanaActionsConnectorExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckKibanaActionsConnectorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["kibana_actions_connector." + resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		c := testAccProvider.Meta().(*gk.Client)
		_, err := c.GetConnector(rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccKibanaActionsConnector(resourceName string) string {
	return fmt.Sprintf(`
			resource "kibana_actions_connector" "%s" {
			  connector_type_id = ".index"
			  name              = "test_connector_%s"
			
			  config = jsonencode(
				{
				  "index" : "test-index",
				  "refresh" : true
				  "executionTimeField" : null
				}
			  )
			}`, resourceName, resourceName)
}

func testAccCheckKibanaActionsConnectorDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*gk.Client)

	// loop through the resources in state, verifying each widget
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kibana_actions_connector" {
			continue
		}

		_, err := c.GetRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("actions connector (%s) still exists", rs.Primary.ID)
		}

		// If the error is equivalent to 404 not found, the widget is destroyed.
		// Otherwise return the error
		if !strings.Contains(err.Error(), "404") {
			return err
		}
	}
	return nil
}