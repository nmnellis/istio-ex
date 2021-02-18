package main

import (
	"github.com/nmnellis/istio-ex/pkg/test/apps"
	"github.com/nmnellis/istio-ex/pkg/test/common"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/resource"
	"net/http"
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
)

var (
	i             istio.Instance
	deploymentCtx common.DeploymentContext
)

func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		Setup(istio.Setup(&i, common.IstioSetupFunc("ingressgateway-ports.yaml"))).
		Setup(apps.DeployEchos(&deploymentCtx)).
		Run()
}

// Run the API tests
func TestRouting(t *testing.T) {
	framework.
		NewTest(t).
		Run(func(ctx framework.TestContext) {

			tgs := []common.TestGroup{
				{
					Name: "routing",
					Cases: []common.TestCase{
						{
							Name:        "prefix-1",
							Description: "HTTP/HTTPS prefix based routing",
							Test:        testPrefixMatch,
							Namespace:   deploymentCtx.EchoContext.AppNamespace.Name(),
							FileName:    "prefix-1.yaml",
						},
					},
				},
			}
			for _, tg := range tgs {
				tg.Run(ctx, t, &deploymentCtx)
			}
		})
}

// testPrefixMatch makes a call from frontend to backend application
func testPrefixMatch(ctx resource.Context, t *testing.T, deploymentCtx *common.DeploymentContext) {
	src := deploymentCtx.EchoContext.Deployments.GetOrFail(t, echo.Service("frontend"))

	backendHost := "backend." + deploymentCtx.EchoContext.AppNamespace.Name() + ".svc.cluster.local"

	src.CallOrFail(t, echo.CallOptions{
		Port: &echo.Port{
			Protocol:    "http",
			ServicePort: 8090,
		},
		Scheme:    "http",
		Address:   backendHost,
		Method:    http.MethodGet,
		Path:      "/route1",
		Count:     1,
		Validator: echo.ExpectOK(),
	})

	src.CallOrFail(t, echo.CallOptions{
		Port: &echo.Port{
			Protocol:    "http",
			ServicePort: 8090,
		},
		Scheme:    "http",
		Address:   backendHost,
		Method:    http.MethodGet,
		Path:      "/route2",
		Count:     1,
		Validator: echo.ExpectOK(),
	})

	src.CallOrFail(t, echo.CallOptions{
		Port: &echo.Port{
			Protocol:    "http",
			ServicePort: 8090,
		},
		Scheme:    "http",
		Address:   backendHost,
		Method:    http.MethodGet,
		Path:      "/bad-route",
		Count:     1,
		Validator: echo.ExpectCode("404"),
	})
}
