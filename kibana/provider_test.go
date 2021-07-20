package kibana

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"os"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"kibana": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}


func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("KIBANA_PASSWORD"); v == "" {
		t.Fatal("KIBANA_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("KIBANA_SPACE"); v == "" {
		t.Fatal("KIBANA_SPACE must be set for acceptance tests")
	}
	if v := os.Getenv("KIBANA_URL"); v == "" {
		t.Fatal("KIBANA_URL must be set for acceptance tests")
	}
	if v := os.Getenv("KIBANA_USERNAME"); v == "" {
		t.Fatal("KIBANA_USERNAME must be set for acceptance tests")
	}
}
