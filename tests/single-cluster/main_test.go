package main

import (
	"github.com/nmnellis/istio-ex/pkg/test/apps"
	"github.com/nmnellis/istio-ex/pkg/test/common"
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
)

var (
	i istio.Instance
)

func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		Setup(istio.Setup(&i, common.IstioSetupFunc("ingressgateway-ports.yaml"))).
		Setup(apps.DeployEchos(&common.DeploymentContext{})).
		Run()
}

// Run the API tests
func TestMultiClusterAPI(t *testing.T) {
	return
}
