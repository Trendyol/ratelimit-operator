module gitlab.trendyol.com/platform/base/apps/ratelimit-operator

go 1.15

require (
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/stretchr/testify v1.5.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	istio.io/api v0.0.0-20210423194545-fa0286046824
	istio.io/client-go v1.8.5-0.20210423200204-66c157dce915
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0 // indirect
	sigs.k8s.io/controller-runtime v0.7.2
)
