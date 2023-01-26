package processing_test

import (
	"fmt"

	apirulev1beta1 "github.com/kyma-incubator/api-gateway/api/v1beta1"
	"github.com/kyma-incubator/api-gateway/internal/processing"
	"github.com/onsi/gomega"
	rulev1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	"istio.io/api/networking/v1beta1"
	networkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	ApiName                     = "test-apirule"
	ApiUID            types.UID = "eab0f1c8-c417-11e9-bf11-4ac644044351"
	ApiNamespace                = "some-namespace"
	ApiAPIVersion               = "gateway.kyma-project.io/v1alpha1"
	ApiKind                     = "ApiRule"
	ApiPath                     = "/.*"
	HeadersApiPath              = "/headers"
	ImgApiPath                  = "/img"
	JwtIssuer                   = "https://oauth2.example.com/"
	JwksUri                     = "https://oauth2.example.com/.well-known/jwks.json"
	JwtIssuer2                  = "https://oauth2.another.example.com/"
	JwksUri2                    = "https://oauth2.another.example.com/.well-known/jwks.json"
	OathkeeperSvc               = "fake.oathkeeper"
	OathkeeperSvcPort uint32    = 1234
	TestLabelKey                = "key"
	TestLabelValue              = "value"
	DefaultDomain               = "myDomain.com"
	TestSelectorKey             = "app"
)

var (
	ApiMethods                     = []string{"GET"}
	ApiScopes                      = []string{"write", "read"}
	ServicePort             uint32 = 8080
	ApiGateway                     = "some-gateway"
	ServiceName                    = "example-service"
	ServiceHostWithNoDomain        = "myService"
	ServiceHost                    = ServiceHostWithNoDomain + "." + DefaultDomain

	TestAllowOrigin  = []*v1beta1.StringMatch{{MatchType: &v1beta1.StringMatch_Regex{Regex: ".*"}}}
	TestAllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	TestAllowHeaders = []string{"header1", "header2"}
	TestCors         = &processing.CorsConfig{
		AllowOrigins: TestAllowOrigin,
		AllowMethods: TestAllowMethods,
		AllowHeaders: TestAllowHeaders,
	}

	TestAdditionalLabels = map[string]string{TestLabelKey: TestLabelValue}
)

func GetTestConfig() processing.ReconciliationConfig {
	return processing.ReconciliationConfig{
		OathkeeperSvc:     OathkeeperSvc,
		OathkeeperSvcPort: OathkeeperSvcPort,
		CorsConfig:        TestCors,
		AdditionalLabels:  TestAdditionalLabels,
		DefaultDomainName: DefaultDomain,
	}
}

func GetFakeClient(objs ...client.Object) client.Client {
	scheme := runtime.NewScheme()
	err := networkingv1beta1.AddToScheme(scheme)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = rulev1alpha1.AddToScheme(scheme)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = securityv1beta1.AddToScheme(scheme)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func GetRuleFor(path string, methods []string, mutators []*apirulev1beta1.Mutator, accessStrategies []*apirulev1beta1.Authenticator) apirulev1beta1.Rule {
	return apirulev1beta1.Rule{
		Path:             path,
		Methods:          methods,
		Mutators:         mutators,
		AccessStrategies: accessStrategies,
	}
}

func GetRuleWithServiceFor(path string, methods []string, mutators []*apirulev1beta1.Mutator, accessStrategies []*apirulev1beta1.Authenticator, service *apirulev1beta1.Service) apirulev1beta1.Rule {
	return apirulev1beta1.Rule{
		Path:             path,
		Methods:          methods,
		Mutators:         mutators,
		AccessStrategies: accessStrategies,
		Service:          service,
	}
}

func GetJwtRuleWithService(jwtIssuer, jwksUri, serviceName string, namespace ...string) apirulev1beta1.Rule {
	jwtConfigJSON := fmt.Sprintf(`{"authentications": [{"issuer": "%s", "jwksUri": "%s"}]}`, jwtIssuer, jwksUri)
	jwt := []*apirulev1beta1.Authenticator{
		{
			Handler: &apirulev1beta1.Handler{
				Name: "jwt",
				Config: &runtime.RawExtension{
					Raw: []byte(jwtConfigJSON),
				},
			},
		},
	}

	port := uint32(8080)
	jwtRuleService := &apirulev1beta1.Service{
		Name: &serviceName,
		Port: &port,
	}
	if len(namespace) > 0 {
		jwtRuleService.Namespace = &namespace[0]
	}

	return GetRuleWithServiceFor("path", ApiMethods, []*apirulev1beta1.Mutator{}, jwt, jwtRuleService)
}

func GetAPIRuleFor(rules []apirulev1beta1.Rule) *apirulev1beta1.APIRule {
	return &apirulev1beta1.APIRule{
		ObjectMeta: v1.ObjectMeta{
			Name:      ApiName,
			UID:       ApiUID,
			Namespace: ApiNamespace,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: ApiAPIVersion,
			Kind:       ApiKind,
		},
		Spec: apirulev1beta1.APIRuleSpec{
			Gateway: &ApiGateway,
			Service: &apirulev1beta1.Service{
				Name: &ServiceName,
				Port: &ServicePort,
			},
			Host:  &ServiceHost,
			Rules: rules,
		},
	}
}

func ToCSVList(input []string) string {
	if len(input) == 0 {
		return ""
	}

	res := `"` + input[0] + `"`

	for i := 1; i < len(input); i++ {
		res = res + "," + `"` + input[i] + `"`
	}

	return res
}

var ActionToString = func(a processing.Action) string { return a.String() }
