package analyzer

// Typed views over the kubectl JSON that the K8S-* rules consume.
//
// These mirror only the fields the rules actually read. Decoding into narrow
// structs rather than a generic map keeps the rule code readable, and turns a
// schema change into a compile error instead of a silently empty string.

type podList struct {
	Items []pod `json:"items"`
}

type pod struct {
	Metadata objectMeta `json:"metadata"`
	Spec     podSpec    `json:"spec"`
	Status   podStatus  `json:"status"`
}

type objectMeta struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type podSpec struct {
	NodeName        string      `json:"nodeName"`
	HostNetwork     bool        `json:"hostNetwork"`
	HostPID         bool        `json:"hostPID"`
	Containers      []container `json:"containers"`
	SecurityContext *podSecCtx  `json:"securityContext"`
}

type podSecCtx struct {
	RunAsUser    *int64 `json:"runAsUser"`
	RunAsNonRoot *bool  `json:"runAsNonRoot"`
}

type container struct {
	Name            string           `json:"name"`
	SecurityContext *containerSecCtx `json:"securityContext"`
	Resources       resources        `json:"resources"`
}

type containerSecCtx struct {
	Privileged               *bool  `json:"privileged"`
	AllowPrivilegeEscalation *bool  `json:"allowPrivilegeEscalation"`
	RunAsUser                *int64 `json:"runAsUser"`
	RunAsNonRoot             *bool  `json:"runAsNonRoot"`
	ReadOnlyRootFilesystem   *bool  `json:"readOnlyRootFilesystem"`
}

type resources struct {
	Limits map[string]string `json:"limits"`
}

type podStatus struct {
	Phase                 string            `json:"phase"`
	Reason                string            `json:"reason"`
	Message               string            `json:"message"`
	ContainerStatuses     []containerStatus `json:"containerStatuses"`
	InitContainerStatuses []containerStatus `json:"initContainerStatuses"`
}

type containerStatus struct {
	Name         string         `json:"name"`
	Ready        bool           `json:"ready"`
	RestartCount int            `json:"restartCount"`
	State        containerState `json:"state"`
	LastState    containerState `json:"lastState"`
}

type containerState struct {
	Waiting    *stateWaiting    `json:"waiting"`
	Terminated *stateTerminated `json:"terminated"`
}

type stateWaiting struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type stateTerminated struct {
	Reason   string `json:"reason"`
	ExitCode int    `json:"exitCode"`
	Message  string `json:"message"`
}

type endpointsList struct {
	Items []endpoints `json:"items"`
}

type endpoints struct {
	Metadata objectMeta       `json:"metadata"`
	Subsets  []endpointSubset `json:"subsets"`
}

type endpointSubset struct {
	Addresses         []endpointAddress `json:"addresses"`
	NotReadyAddresses []endpointAddress `json:"notReadyAddresses"`
}

type endpointAddress struct {
	IP string `json:"ip"`
}

// ReadyAddresses counts endpoints actually eligible to receive traffic.
func (e endpoints) ReadyAddresses() int {
	total := 0
	for _, s := range e.Subsets {
		total += len(s.Addresses)
	}
	return total
}

// NotReadyAddresses counts backing pods that exist but are failing readiness.
func (e endpoints) NotReadyAddresses() int {
	total := 0
	for _, s := range e.Subsets {
		total += len(s.NotReadyAddresses)
	}
	return total
}

type nodeList struct {
	Items []node `json:"items"`
}

type node struct {
	Metadata objectMeta `json:"metadata"`
	Status   nodeStatus `json:"status"`
}

type nodeStatus struct {
	Conditions []nodeCondition `json:"conditions"`
}

type nodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type pvcList struct {
	Items []pvc `json:"items"`
}

type pvc struct {
	Metadata objectMeta `json:"metadata"`
	Spec     pvcSpec    `json:"spec"`
	Status   pvcStatus  `json:"status"`
}

type pvcSpec struct {
	StorageClassName *string `json:"storageClassName"`
	VolumeName       string  `json:"volumeName"`
}

type pvcStatus struct {
	Phase string `json:"phase"`
}

type ingressList struct {
	Items []ingress `json:"items"`
}

type ingress struct {
	Metadata objectMeta  `json:"metadata"`
	Spec     ingressSpec `json:"spec"`
}

type ingressSpec struct {
	TLS   []ingressTLS  `json:"tls"`
	Rules []ingressRule `json:"rules"`
}

type ingressTLS struct {
	Hosts      []string `json:"hosts"`
	SecretName string   `json:"secretName"`
}

type ingressRule struct {
	Host string       `json:"host"`
	HTTP *ingressHTTP `json:"http"`
}

type ingressHTTP struct {
	Paths []ingressPath `json:"paths"`
}

type ingressPath struct {
	Path    string         `json:"path"`
	Backend ingressBackend `json:"backend"`
}

type ingressBackend struct {
	Service *ingressBackendService `json:"service"`
}

type ingressBackendService struct {
	Name string `json:"name"`
}

// BackendServices lists every service name this ingress routes to.
func (i ingress) BackendServices() []string {
	names := make([]string, 0, 4)
	for _, rule := range i.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, p := range rule.HTTP.Paths {
			if p.Backend.Service != nil && p.Backend.Service.Name != "" {
				names = append(names, p.Backend.Service.Name)
			}
		}
	}
	return names
}

type serviceList struct {
	Items []service `json:"items"`
}

type service struct {
	Metadata objectMeta `json:"metadata"`
}

type secretList struct {
	Items []secretMeta `json:"items"`
}

type secretMeta struct {
	Metadata objectMeta `json:"metadata"`
}
