// Code generated by gotestmd DO NOT EDIT.
package sriov

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"

	"github.com/networkservicemesh/integration-k8s-packet/suites/spire"
)

type Suite struct {
	shell.Suite
	spireSuite spire.Suite
}

func (s *Suite) SetupSuite() {
	suite.Run(s.T(), &s.spireSuite)
	r := s.Runner("../deployments-k8s/examples/sriov")
	s.T().Cleanup(func() {
		r.Run(`kubectl -n nsm-system get pods`)
		r.Run(`kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`kubectl apply -k .`)
}
func (s *Suite) TestVFIOConnection() {
	r := s.Runner("../deployments-k8s/examples/VFIOConnection")
	s.T().Cleanup(func() {
		r.Run(`kubectl -n ${NAMESPACE} get pods`)
		r.Run(`NSE_POD=$(kubectl -n ${NAMESPACE} get pods -l app=nse |` + "\n" + `  grep -v "NAME" |` + "\n" + `  sed -E "s/([.]*) .*/\1/g")`)
		r.Run(`kubectl -n ${NAMESPACE} exec ${NSE_POD} --container ponger -- /bin/bash -c '                  \` + "\n" + `  sleep 10 && kill $(ps -A | grep "pingpong" | sed -E "s/ *([0-9]*).*/\1/g") 1>/dev/null 2>&1 & \` + "\n" + `'`)
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../apps/vfio-nsc` + "\n" + `- ../../apps/vfio-nse` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=5m pod -l app=nsc`)
	r.Run(`kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=5m pod -l app=nse`)
	r.Run(`NSC_POD=$(kubectl -n ${NAMESPACE} get pods -l app=nsc |` + "\n" + `  grep -v "NAME" |` + "\n" + `  sed -E "s/([.]*) .*/\1/g")`)
	r.Run(`kubectl -n ${NAMESPACE} logs ${NSC_POD} sidecar |` + "\n" + `  grep "All client init operations are done."`)
	r.Run(`PING_RESULTS=$(kubectl -n ${NAMESPACE} exec ${NSC_POD} --container pinger -- /bin/bash -c ' \` + "\n" + `  /root/dpdk-pingpong/build/app/pingpong                                                    \` + "\n" + `    --no-huge                                                                               \` + "\n" + `    --                                                                                      \` + "\n" + `    -n 500                                                                                  \` + "\n" + `    -c                                                                                      \` + "\n" + `    -C 0a:11:22:33:44:55                                                                    \` + "\n" + `    -S 0a:55:44:33:22:11                                                                    \` + "\n" + `' 2>&1) || (echo "${PING_RESULTS}" 1>&2 && false)`)
	r.Run(`PONG_PACKETS="$(echo "${PING_RESULTS}"                      |` + "\n" + `                grep "rx .* pong packets"                   |` + "\n" + `                sed -E 's/rx ([0-9]*) pong packets/\1/g')"  \` + "\n" + `  || (echo "${PING_RESULTS}" 1>&2 && false)`)
	r.Run(`test "${PONG_PACKETS}" -ne 0 \` + "\n" + `  || (echo "${PING_RESULTS}" 1>&2 && false)`)
}
