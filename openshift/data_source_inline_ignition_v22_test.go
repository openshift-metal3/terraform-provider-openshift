package openshift

import (
	"fmt"
	gth "github.com/gophercloud/gophercloud/testhelper"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"net/http"
	"testing"
)

func TestInlineIgnitionV22(t *testing.T) {
	gth.SetupHTTP()
	defer gth.TeardownHTTP()
	handleIgnitionRequest(t)

	testAccProvider := Provider()

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"openshift": testAccProvider,
		},
		Steps: []resource.TestStep{
			{
				Config: testIgnitionResource(gth.Server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openshift_inline_ignition_v2_2.ignition-data", "inlined", "{\"ignition\":{\"config\":{\"append\":[{\"source\":\"data:text/plain;charset=utf-8;base64,b3RoZXIgZGF0YQ==\",\"verification\":{}}]},\"security\":{\"tls\":{}},\"timeouts\":{},\"version\":\"2.2.0\"},\"networkd\":{},\"passwd\":{},\"storage\":{},\"systemd\":{}}"),
				),
			},
		},
	})
}

// Returns a resource declaration for a particular node name, and it's related introspection data source.
func testIgnitionResource(url string) string {
	return fmt.Sprintf(`
		data "openshift_inline_ignition_v2_2" "ignition-data" {
			ignition = "{\"ignition\":{\"config\":{\"append\":[{\"source\":\"%s\",\"verification\":{}}]},\"security\":{\"tls\":{\"certificateAuthorities\":[]}},\"timeouts\":{},\"version\":\"2.2.0\"},\"networkd\":{},\"passwd\":{},\"storage\":{},\"systemd\":{}}"
		}
`, url)
}

// When using the fake inspect interface, inspector isn't actually used, so we need
// to mock the inspector's API responses for returning data to us.
func handleIgnitionRequest(t *testing.T) {
	gth.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		log.Printf("[DEBUG] URL is %s", r.URL.Path)

		fmt.Fprintf(w, "other data")
	})
}
