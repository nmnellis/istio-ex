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
		if deploymentCtx == nil {
			deploymentCtx = &common.DeploymentContext{}
		}
		echoCtx := &common.EchoDeploymentContext{}

		if err := createNamespaces(ctx, echoCtx); err != nil {
			return err
		}
		if err := generateTLSCertificates(ctx, "echo-certs", echoCtx.AppNamespace); err != nil {
			return err
		}
		if err := generateTLSCertificates(ctx, "echo-certs", echoCtx.SubsetNamespace); err != nil {
			return err
		}

		apps, err := deployApplications(ctx, echoCtx)
		if err != nil {
			return err
		}

		echoCtx.Deployments = apps
		return nil
	}
}

func deployApplications(ctx resource.Context, echoCtx *common.EchoDeploymentContext) (echo.Instances, error) {
	builder := echoboot.NewBuilder(ctx)
	frontendApp := newEchoConfig("frontend", echoCtx.AppNamespace, ctx.Clusters()[0], true, false)
	if _, err := builder.With(nil, frontendApp).Build(); err != nil {
		scopes.Framework.Errorf("error setting up frontend echo %v", err.Error())
		return nil, err
	}

	backendApp := newEchoConfig("backend", echoCtx.AppNamespace, ctx.Clusters()[0], true, false)
	if _, err := builder.With(nil, backendApp).Build(); err != nil {
		scopes.Framework.Errorf("error setting up backend echo %v", err.Error())
		return nil, err
	}

	subsetApp := newEchoConfig("subset", echoCtx.SubsetNamespace, ctx.Clusters()[0], true, true)
	if _, err := builder.With(nil, subsetApp).Build(); err != nil {
		scopes.Framework.Errorf("error setting up subset echos %v", err.Error())
		return nil, err
	}

	nonMeshApp := newEchoConfig("no-mesh", echoCtx.NoMeshNamespace, ctx.Clusters()[0], false, false)
	if _, err := builder.With(nil, nonMeshApp).Build(); err != nil {
		scopes.Framework.Errorf("error setting up no mesh echo %v", err.Error())
		return nil, err
	}

	apps, err := builder.Build()
	if err != nil {
		scopes.Framework.Errorf("error setting up echo apps %v", err.Error())
		return nil, err
	}
	return apps, nil
}

func createNamespaces(ctx resource.Context,  echoCtx *common.EchoDeploymentContext) error {
	var err error

	if echoCtx.AppNamespace, err = namespace.New(ctx, namespace.Config{Prefix: "app", Inject: true}); err != nil {
		return err
	}

	if echoCtx.SubsetNamespace, err = namespace.New(ctx, namespace.Config{Prefix: "subset", Inject: true}); err != nil {
		return err
	}

	if echoCtx.NoMeshNamespace, err = namespace.New(ctx, namespace.Config{Prefix: "no-mesh", Inject: false}); err != nil {
		return err
	}

	return nil
}

func generateTLSCertificates(ctx resource.Context, secretName string, ns namespace.Instance) error {
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
		Name:      secretName,
		CACrt:     string(echoCA),
		TLSKey:    string(echoKey),
		TLSCert:   string(echoCrt),
		Cluster:   ctx.Clusters()[0],
	})
	if err != nil {
		return err
	}
	return nil
}

func newEchoConfig(service string, ns namespace.Instance, cluster cluster.Cluster, hasSidecar bool, useSubsets bool) echo.Config {
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
	} else if useSubsets {
		subset = []echo.SubsetConfig{
			{
				Version: "v1",
			},
			{
				Version: "v2",
			},
		}
	}

	return echo.Config{
		Namespace: ns,
		Service:   service,
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
			{
				// TCP port
				Name:         "tcp",
				Protocol:     protocol.TCP,
				ServicePort:  9000,
				InstancePort: 9000,
				TLS:          false,
			},
		},
		Subsets:     subset,
		TLSSettings: tlsSettings,
		Cluster:     cluster,
	}
}
