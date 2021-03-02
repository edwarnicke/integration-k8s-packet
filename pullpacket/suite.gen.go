// Code generated by gotestmd DO NOT EDIT.
package pullpacket

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
)

type Suite struct {
	base.Suite
}

func (s *Suite) SetupSuite() {
	parents := []interface{}{&s.Suite}
	for _, p := range parents {
		if v, ok := p.(suite.TestingSuite); ok {
			v.SetT(s.T())
		}
		if v, ok := p.(suite.SetupAllSuite); ok {
			v.SetupSuite()
		}
	}
	r := s.Runner("../deployments-k8s/examples/use-cases/Kernel2Kernel")
	r.Run(`NAMESPACE=($(kubectl create -f ../namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../../apps/registry-memory` + "\n" + `- ../../../apps/nsmgr` + "\n" + `- ../../../apps/forwarder-sriov` + "\n" + `- ../../../apps/nse-vfio` + "\n" + `- ../../../apps/nsc-vfio` + "\n" + `` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=10m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=10m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=10m pod -l app=forwarder-sriov -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=10m pod -l app=nsmgr -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=10m pod -l app=nsm-registry -n ${NAMESPACE}`)
}
func (s *Suite) Test() {}