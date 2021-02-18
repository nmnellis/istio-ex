package common

import (
	"github.com/nmnellis/istio-ex/pkg/test/packr"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/scopes"
)

type DeploymentContext struct {
	echoContext *EchoDeploymentContext
}

type EchoDeploymentContext struct {
	Deployments echo.Instances
	AppNamespace namespace.Instance
	SubsetNamespace namespace.Instance
	NoMeshNamespace namespace.Instance
}

func IstioSetupFunc(operatorFile string) func(ctx resource.Context, cfg *istio.Config) {
	return func(ctx resource.Context, cfg *istio.Config) {
		// load custom istio control plane value
		istioGatewayConfig, err := packr.RenderOperator(operatorFile, nil)
		if err != nil {
			scopes.Framework.Errorf("error rendering istio operator configuration file %s %w", operatorFile, err)
		}
		cfg.ControlPlaneValues = istioGatewayConfig
	}
}
