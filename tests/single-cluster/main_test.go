package main

import (
	"github.com/nmnellis/istio-ex/pkg/test/packr"
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/scopes"
)

var (
	i istio.Instance
)

func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		Setup(istio.Setup(&i, func(ctx resource.Context, cfg *istio.Config) {
			// load custom istio control plane value
			istioGatewayConfig, err := packr.RenderOperator("ingressgateway-ports.yaml", nil)
			if err != nil {
				scopes.Framework.Errorf("error rendering istio gateway configuration %w", err)
			}
			cfg.ControlPlaneValues = istioGatewayConfig
		})).
		Run()
}

// Run the API tests
func TestMultiClusterAPI(t *testing.T) {
	return
}
