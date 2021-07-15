package istio

import (
	"context"
	"gotest.tools/assert"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"testing"
)

func Test_CreateEnvoyFilter(t *testing.T) {
	t.Parallel()
	var fakeIstioClient = getFakeIstioClient()

	_, err := fakeIstioClient.CreateEnvoyFilter(context.Background(), "default", &v1alpha3.EnvoyFilter{})
	if err != nil {
		assert.Error(t, err, "error occurred when creating envoyfilter")
	}

}

func getFakeIstioClient() Istio {

	return FakeClient()
}
