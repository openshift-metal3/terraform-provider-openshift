package openshift

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vincent-petithory/dataurl"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func dataSourceInlineIgnitionV2_2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInlineIgnitionV2_2Read,
		Schema: map[string]*schema.Schema{
			"ignition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inlined": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceInlineIgnitionV2_2Read(d *schema.ResourceData, meta interface{}) error {
	var ignition types.Config

	ignitionData := d.Get("ignition").(string)

	if err := json.Unmarshal([]byte(ignitionData), &ignition); err != nil {
		return err
	}

	transport := &http.Transport{}

	if len(ignition.Ignition.Security.TLS.CertificateAuthorities) > 0 {
		caCertPool := x509.NewCertPool()
		transport.TLSClientConfig = &tls.Config{RootCAs: caCertPool}

		for _, caCertEncoded := range ignition.Ignition.Security.TLS.CertificateAuthorities {
			dataURL, err := dataurl.DecodeString(caCertEncoded.Source)
			if err != nil {
				return fmt.Errorf("could not decode CA certificate: %s", err.Error())
			}

			caCertPool.AppendCertsFromPEM(dataURL.Data)
		}
	}
	client := retryablehttp.NewClient()
	client.HTTPClient.Transport = transport

	for idx, append := range ignition.Ignition.Config.Append {
		if url, err := url.Parse(append.Source); err == nil {
			resp, err := client.Get(url.String())
			if err != nil {
				return fmt.Errorf("could not fetch append: %s", err.Error())
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Errorf("could not read ignition data: %s", err.Error())
			}

			ignition.Ignition.Config.Append[idx].Source = dataurl.EncodeBytes(data)
		}
	}

	result, err := json.Marshal(ignition)
	if err != nil {
		return fmt.Errorf("could not convert ignition back to JSON: %s", err.Error())
	}

	d.SetId(time.Now().UTC().String())
	d.Set("inlined", string(result))
	return nil
}
