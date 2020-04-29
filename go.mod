module github.com/openshift-metal3/terraform-provider-openshift

go 1.14

require (
	github.com/ajeddeloh/go-json v0.0.0-20200220154158-5ae607161559 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/ignition v0.35.1-0.20200406191649-73ff2d799cb3
	github.com/gophercloud/gophercloud v0.8.0
	github.com/hashicorp/go-retryablehttp v0.6.6
	github.com/hashicorp/terraform-plugin-sdk v1.10.0
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	go4.org v0.0.0-20200411211856-f5505b9728dd // indirect
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
