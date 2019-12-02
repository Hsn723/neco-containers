package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/test/e2e/framework/pod"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	rookCephOSDAppName = "rook-ceph-osd"
)

// cnpvlog is for logging in this package.
var cnpvlog = logf.Log.WithName("rook-topologyspreadconstraints-mutator")

// +kubebuilder:webhook:verbs=create;update,path=/validate-projectcalico-org-networkpolicy,mutating=false,failurePolicy=fail,groups=crd.projectcalico.org,resources=networkpolicies,versions=v1,name=vnetworkpolicy.kb.io
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

// rookTopologySpreadConstraintsMutator is a mutating webhook to inject TopologySpreadConstraints to Rook OSD pods.
type rookTopologySpreadConstraintsMutator struct {
	client              client.Client
	decoder             *admission.Decoder
	namespace			string
}

// NewRookTopologySpreadConstraintsMutator creates a webhook handler to inject
// TopologySpreadConstraints to Rook OSD pods.
func NewRookTopologySpreadConstraintsMutator(c client.Client, dec *admission.Decoder, namespace string) http.Handler {
	return &webhook.Admission{Handler: rookTopologySpreadConstraintsMutator&{c, dec, namespace}}
}

// Handle implements admission.Handler interface.
func (m *rookTopologySpreadConstraintsMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	if err := json.Unmarshal(req.Object.Raw, pod); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Namespace != m.namespace || pod.Labels["app"] != rookCephOSDAppName {
		return admission.Allowed("not corresponding to Rook's OSD pod")
	}

	if pod.Spec.TopologySpreadConstraints != nil {
		return admission.Allowed("topologySpreadConstraints resource already exists")
	}

	constraint := corev1.TopologySpreadConstraint{}
	constraint.MaxSkew = 1
	constraint.TopologyKey = "topology.rook.io/rack"
	constraint.LabelSelector.MatchLabels["app"] = rookCephOSDAppName

	pod.Spec.TopologySpreadConstraints

	return admission.PatchResponseFromRaw(req.Object.Raw, mutatedPod)
}
