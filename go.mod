module github.com/nmnellis/istio-ex

go 1.15

replace github.com/spf13/viper => github.com/istio/viper v1.3.3-0.20190515210538-2789fed3109c

require (
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/rogpeppe/go-internal v1.7.0 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/tools v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	istio.io/istio v0.0.0-20210211212916-be8d7d33ae67
	istio.io/pkg v0.0.0-20201230223204-2d0a1c8bd9e5
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
)

// Pending https://github.com/kubernetes/kube-openapi/pull/220
replace k8s.io/kube-openapi => github.com/howardjohn/kube-openapi v0.0.0-20210104181841-c0b40d2cb1c8
