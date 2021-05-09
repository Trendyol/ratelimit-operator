package istio

import (
	"istio.io/client-go/pkg/clientset/versioned/fake"
	"testing"
)

//TODO: Istio tests
func Test_CreateEnvoyFilter(t *testing.T) {
	t.Parallel()
	//var fakeIstioClient = fakeClient()

	//data, error := fakeIstioClient.CreateEnvoyFilter()
}

func fakeClient() IstioClient {
	return &istioClient{client: fake.NewSimpleClientset()}
}
