package v1alpha1

import (
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	snowv1 "github.com/aws/eks-anywhere/pkg/providers/snow/api/v1beta1"
)

func TestSnowMachineConfigSetDefaults(t *testing.T) {
	tests := []struct {
		name   string
		before *SnowMachineConfig
		after  *SnowMachineConfig
	}{
		{
			name: "optional fields all empty",
			before: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{},
			},
			after: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					InstanceType:             DefaultSnowInstanceType,
					PhysicalNetworkConnector: DefaultSnowPhysicalNetworkConnectorType,
				},
			},
		},
		{
			name: "instance type exists",
			before: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					InstanceType: "instance-type-1",
				},
			},
			after: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					InstanceType:             "instance-type-1",
					PhysicalNetworkConnector: DefaultSnowPhysicalNetworkConnectorType,
				},
			},
		},
		{
			name: "ssh key name exists",
			before: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					SshKeyName: "ssh-name",
				},
			},
			after: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					SshKeyName:               "ssh-name",
					InstanceType:             DefaultSnowInstanceType,
					PhysicalNetworkConnector: DefaultSnowPhysicalNetworkConnectorType,
				},
			},
		},
		{
			name: "physical network exists",
			before: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					PhysicalNetworkConnector: "network-1",
				},
			},
			after: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					PhysicalNetworkConnector: "network-1",
					InstanceType:             DefaultSnowInstanceType,
				},
			},
		},
		{
			name: "os family exists",
			before: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					OSFamily: "ubuntu",
				},
			},
			after: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					InstanceType:             DefaultSnowInstanceType,
					PhysicalNetworkConnector: DefaultSnowPhysicalNetworkConnectorType,
					OSFamily:                 Ubuntu,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			tt.before.SetDefaults()
			g.Expect(tt.before).To(Equal(tt.after))
		})
	}
}

func TestSnowMachineConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		obj     *SnowMachineConfig
		wantErr string
	}{
		{
			name: "valid config with amiID, instance type, devices, osFamily",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					AMIID:        "ami-1",
					InstanceType: DefaultSnowInstanceType,
					Devices:      []string{"1.2.3.4"},
					OSFamily:     Bottlerocket,
				},
			},
			wantErr: "",
		},
		{
			name: "valid without ami",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					InstanceType: DefaultSnowInstanceType,
					Devices:      []string{"1.2.3.4"},
					OSFamily:     Ubuntu,
				},
			},
			wantErr: "",
		},
		{
			name: "invalid instance type",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					AMIID:        "ami-1",
					InstanceType: "invalid-instance-type",
					Devices:      []string{"1.2.3.4"},
					OSFamily:     Bottlerocket,
				},
			},
			wantErr: "InstanceType invalid-instance-type is not supported",
		},
		{
			name: "empty devices",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					AMIID:        "ami-1",
					InstanceType: DefaultSnowInstanceType,
					OSFamily:     Bottlerocket,
				},
			},
			wantErr: "Devices must contain at least one device IP",
		},
		{
			name: "invalid container volume",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					AMIID:        "ami-1",
					InstanceType: DefaultSnowInstanceType,
					Devices:      []string{"1.2.3.4"},
					ContainersVolume: &snowv1.Volume{
						Size: 7,
					},
					OSFamily: Bottlerocket,
				},
			},
			wantErr: "ContainersVolume.Size must be no smaller than 8 Gi",
		},
		{
			name: "invalid os family",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					AMIID:        "ami-1",
					InstanceType: DefaultSnowInstanceType,
					Devices:      []string{"1.2.3.4"},
					OSFamily:     "invalidOS",
				},
			},
			wantErr: "SnowMachineConfig OSFamily invalidOS is not supported",
		},
		{
			name: "empty os family",
			obj: &SnowMachineConfig{
				Spec: SnowMachineConfigSpec{
					AMIID:        "ami-1",
					InstanceType: DefaultSnowInstanceType,
					Devices:      []string{"1.2.3.4"},
					OSFamily:     "",
				},
			},
			wantErr: "SnowMachineConfig OSFamily must be specified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			err := tt.obj.Validate()
			if tt.wantErr == "" {
				g.Expect(err).To(BeNil())
			} else {
				g.Expect(err).To(MatchError(ContainSubstring(tt.wantErr)))
			}
		})
	}
}

func TestSnowMachineConfigSetControlPlaneAnnotation(t *testing.T) {
	g := NewWithT(t)
	m := &SnowMachineConfig{}
	m.SetControlPlaneAnnotation()
	g.Expect(m.Annotations).To(Equal(map[string]string{"anywhere.eks.amazonaws.com/control-plane": "true"}))
}

func TestSnowMachineConfigSetEtcdAnnotation(t *testing.T) {
	g := NewWithT(t)
	m := &SnowMachineConfig{}
	m.SetEtcdAnnotation()
	g.Expect(m.Annotations).To(Equal(map[string]string{"anywhere.eks.amazonaws.com/etcd": "true"}))
}

func TestNewSnowMachineConfigGenerate(t *testing.T) {
	g := NewWithT(t)
	want := &SnowMachineConfigGenerate{
		TypeMeta: metav1.TypeMeta{
			Kind:       SnowMachineConfigKind,
			APIVersion: SchemeBuilder.GroupVersion.String(),
		},
		ObjectMeta: ObjectMeta{
			Name: "snow-cluster",
		},
		Spec: SnowMachineConfigSpec{
			AMIID:                    "",
			InstanceType:             DefaultSnowInstanceType,
			SshKeyName:               DefaultSnowSshKeyName,
			PhysicalNetworkConnector: DefaultSnowPhysicalNetworkConnectorType,
		},
	}
	g.Expect(NewSnowMachineConfigGenerate("snow-cluster")).To(Equal(want))
}
