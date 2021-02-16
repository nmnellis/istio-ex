package apps

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/nmnellis/istio-ex/pkg/test/common"
	"github.com/nmnellis/istio-ex/pkg/test/tlssecret"
	"istio.io/istio/pkg/config/protocol"
	common2 "istio.io/istio/pkg/test/echo/common"
	"istio.io/istio/pkg/test/framework/components/cluster"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/echoboot"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/scopes"
	"strconv"
)

var (
	echoCertsBox = packr.New("certs", "./certs")
)

func DeployEchos(deploymentCtx *common.DeploymentContext) resource.SetupFn {
	return func(ctx resource.Context) error {
		var err error
		if deploymentCtx == nil {
			deploymentCtx = &common.DeploymentContext{}
		}
		echoCtx := &common.EchoDeploymentContext{}

		ns, err := namespace.New(ctx, namespace.Config{Prefix: "echo", Inject: true})
		if err != nil {
			return err
		}

		echoCtx.Namespace = ns

		echoCrt, err := echoCertsBox.Find("echo.crt")
		if err != nil {
			scopes.Framework.Error(err)
		}
		echoKey, err := echoCertsBox.Find("echo.key")
		if err != nil {
			scopes.Framework.Error(err)
		}
		echoCA, err := echoCertsBox.Find("echo-ca.crt")
		if err != nil {
			scopes.Framework.Error(err)
		}
		_, err = tlssecret.New(ctx, &tlssecret.Config{
			Namespace: ns.Name(),
			Name:      "echo-certs",
			CACrt:     string(echoCA),
			TLSKey:    string(echoKey),
			TLSCert:   string(echoCrt),
			Cluster:   ctx.Clusters()[0],
		})
		if err != nil {
			return err
		}

		scopes.Framework.Infof("Deploying echo app to cluster %v", ctx.Clusters()[0].Name())
		echoApplications := newEchoConfig("echo", ns, ctx.Clusters()[0], "us-east1", true)
		echoInstance, err := echoboot.NewBuilder(ctx).With(nil, echoApplications).Build()
		if err != nil {
			scopes.Framework.Errorf("error setting up echos %v", err.Error())
			return err
		}
		echoCtx.Deployments = echoInstance
		return nil
	}
}

func newEchoConfig(service string, ns namespace.Instance, cluster cluster.Cluster, locality string, hasSidecar bool) echo.Config {
	echoCrt, err := echoCertsBox.Find("echo.crt")
	if err != nil {
		scopes.Framework.Error(err)
	}
	echoKey, err := echoCertsBox.Find("echo.key")
	if err != nil {
		scopes.Framework.Error(err)
	}
	echoCA, err := echoCertsBox.Find("echo-ca.crt")
	if err != nil {
		scopes.Framework.Error(err)
	}

	tlsSettings := &common2.TLSSettings{
		RootCert:   string(echoCA),
		ClientCert: string(echoCrt),
		Key:        string(echoKey),
	}
	var subset []echo.SubsetConfig
	if !hasSidecar {
		subset = []echo.SubsetConfig{
			{
				Annotations: map[echo.Annotation]*echo.AnnotationValue{
					echo.SidecarInject: {
						Value: strconv.FormatBool(false)},
				},
			},
		}
	}

	return echo.Config{
		Namespace: ns,
		Service:   service,
		Locality:  locality,
		Ports: []echo.Port{
			{
				Name:     "http",
				Protocol: protocol.HTTP,
				// We use a port > 1024 to not require root
				InstancePort: 8090,
			},
			{
				// HTTPS port
				Name:         "https",
				Protocol:     protocol.HTTPS,
				ServicePort:  9443,
				InstancePort: 9443,
				TLS:          true,
			},
		},
		Subsets:     subset,
		TLSSettings: tlsSettings,
		Cluster:     cluster,
	}
}
