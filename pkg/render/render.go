package render

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	selfservicev1 "github.com/AlonaKaplan/selfserviceoverlay/api/v1"
)

func NetAttachDef(overlayNet *selfservicev1.OverlayNetwork) (*netv1.NetworkAttachmentDefinition, error) {
	cniNetConf, err := renderCNINetworkConfig(overlayNet)
	if err != nil {
		return nil, err
	}
	cniNetConfRaw, err := json.Marshal(cniNetConf)
	if err != nil {
		return nil, err
	}

	const netAttachDefKind = "NetworkAttachmentDefinition"
	const netAttachDefAPIVer = "v1"
	blockOwnerDeletion := true
	return &netv1.NetworkAttachmentDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: netAttachDefAPIVer,
			Kind:       netAttachDefKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      overlayNet.Name,
			Namespace: overlayNet.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         overlayNet.APIVersion,
					Kind:               overlayNet.Kind,
					Name:               overlayNet.Name,
					UID:                overlayNet.UID,
					BlockOwnerDeletion: &blockOwnerDeletion,
				},
			},
		},
		Spec: netv1.NetworkAttachmentDefinitionSpec{
			Config: string(cniNetConfRaw),
		},
	}, nil
}

func renderCNINetworkConfig(overlayNet *selfservicev1.OverlayNetwork) (map[string]interface{}, error) {
	const (
		cniVersionKey       = "cniVersion"
		cniVersion          = "0.3.1"
		topologyKey         = "topology"
		topologyLayer2      = "layer2"
		typeKey             = "type"
		ovnK8sCniOverlay    = "ovn-k8s-cni-overlay"
		nameKey             = "name"
		netAttachDefNameKey = "netAttachDefName"
	)
	cniNetConf := map[string]interface{}{
		cniVersionKey:       cniVersion,
		typeKey:             ovnK8sCniOverlay,
		nameKey:             overlayNet.Namespace + "-" + overlayNet.Spec.Name,
		netAttachDefNameKey: overlayNet.Namespace + "/" + overlayNet.Name,
		topologyKey:         topologyLayer2,
	}

	if overlayNet.Spec.Mtu != "" {
		mtu, err := strconv.Atoi(overlayNet.Spec.Mtu)
		if err != nil {
			return nil, err
		}
		const mtuKey = "mtu"
		cniNetConf[mtuKey] = mtu
	}

	if overlayNet.Spec.Subnets != "" {
		if err := validateSubnets(overlayNet.Spec.Subnets); err != nil {
			return nil, err
		}
		const subnetsKey = "subnets"
		cniNetConf[subnetsKey] = overlayNet.Spec.Subnets
	}

	if overlayNet.Spec.ExcludeSubnets != "" {
		if err := validateSubnets(overlayNet.Spec.ExcludeSubnets); err != nil {
			return nil, err
		}
		const excludeSubnetsKey = "excludeSubnets"
		cniNetConf[excludeSubnetsKey] = overlayNet.Spec.ExcludeSubnets
	}

	return cniNetConf, nil
}

func validateSubnets(subnets string) error {
	subentsSlice := strings.Split(subnets, ",")
	for _, subent := range subentsSlice {
		if _, _, err := net.ParseCIDR(subent); err != nil {
			return err
		}
	}
	return nil
}
