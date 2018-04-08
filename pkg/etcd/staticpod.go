package etcd

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/appscode/kutil"
	"github.com/appscode/kutil/meta"
	"github.com/etcd-manager/etcd-discovery/pkg/config"
	"github.com/etcd-manager/etcd-discovery/pkg/constants"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	etcdVolumeName  = "etcd-data"
	certsVolumeName = "etcd-certs"
)

type etcdStaticPod struct {
	manifestDir string
	cfg         *config.EtcdFlags
}

var _ Process = &etcdStaticPod{}

func (p *etcdStaticPod) Type() ProcessType {
	return ProcessTypeStaticPod
}

func (p *etcdStaticPod) Start() error {
	// gets etcd StaticPodSpec, actualized for the current MasterConfiguration
	pathType := v1.HostPathDirectoryOrCreate
	etcdMounts := map[string]v1.Volume{
		etcdVolumeName:  NewVolume(etcdVolumeName, p.cfg.DataDir, &pathType),
		certsVolumeName: NewVolume(certsVolumeName, p.cfg.CertificatesDir+"/etcd", &pathType),
	}
	cmds, err := p.getEtcdCommand()
	if err != nil {
		return err
	}
	spec := ComponentPod(v1.Container{
		Name:            constants.Etcd,
		Command:         cmds,
		Image:           p.cfg.Version.GetDockerImage(),
		ImagePullPolicy: v1.PullIfNotPresent,
		// Mount the etcd datadir path read-write so etcd can store data in a more persistent manner
		VolumeMounts: []v1.VolumeMount{
			NewVolumeMount(etcdVolumeName, p.cfg.DataDir, false),
			NewVolumeMount(certsVolumeName, p.cfg.CertificatesDir+"/etcd", false),
		},
		LivenessProbe: p.EtcdProbe(
			constants.EtcdCACertName, constants.EtcdHealthcheckClientCertName, constants.EtcdHealthcheckClientKeyName,
		),
	}, etcdMounts)

	// writes etcd StaticPod to disk
	if err := p.WriteStaticPodToDisk(constants.Etcd, spec); err != nil {
		return err
	}

	fmt.Printf("[etcd] Wrote Static Pod manifest for a local etcd instance to %q\n", p.GetStaticPodFilepath())
	return nil
}

func (p *etcdStaticPod) Stop() error {
	filename := p.GetStaticPodFilepath()
	return os.Remove(filename)
}

func (p *etcdStaticPod) ExitState() (error, *os.ProcessState) {
	return kutil.ErrUnknown, nil
}

// getEtcdCommand builds the right etcd command from the given config object
func (p *etcdStaticPod) getEtcdCommand() ([]string, error) {
	args, err := p.cfg.ToArgs()
	if err != nil {
		return nil, err
	}
	command := []string{"etcd"}
	command = append(command, args...)
	return command, nil
}

// GetStaticPodFilepath returns the location on the disk where the Static Pod should be present
func (p *etcdStaticPod) GetStaticPodFilepath() string {
	return filepath.Join(p.manifestDir, constants.Etcd+".yaml")
}

// ComponentPod returns a Pod object from the container and volume specifications
func ComponentPod(container v1.Container, volumes map[string]v1.Volume) v1.Pod {
	return v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        container.Name,
			Namespace:   metav1.NamespaceSystem,
			Annotations: map[string]string{"scheduler.alpha.kubernetes.io/critical-pod": ""},
			// The component and tier labels are useful for quickly identifying the control plane Pods when doing a .List()
			// against Pods in the kube-system namespace. Can for example be used together with the WaitForPodsWithLabel function
			Labels: map[string]string{"component": container.Name, "tier": "control-plane"},
		},
		Spec: v1.PodSpec{
			Containers:  []v1.Container{container},
			HostNetwork: true,
			Volumes:     VolumeMapToSlice(volumes),
		},
	}
}

// EtcdProbe is a helper function for building a shell-based, etcdctl v1.Probe object to healthcheck etcd
func (p *etcdStaticPod) EtcdProbe(CACertName string, CertName string, KeyName string) *v1.Probe {
	tlsFlags := fmt.Sprintf("--cacert=%[1]s/%[2]s --cert=%[1]s/%[3]s --key=%[1]s/%[4]s", p.cfg.CertificatesDir, CACertName, CertName, KeyName)
	// etcd pod is alive if a linearizable get succeeds.
	cmd := fmt.Sprintf("ETCDCTL_API=3 etcdctl --endpoints=%s:%d %s get foo", p.GetProbeAddress(), config.ClientPort, tlsFlags)

	return &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{"/bin/sh", "-ec", cmd},
			},
		},
		InitialDelaySeconds: 15,
		TimeoutSeconds:      15,
		FailureThreshold:    8,
	}
}

// NewVolume creates a v1.Volume with a hostPath mount to the specified location
func NewVolume(name, path string, pathType *v1.HostPathType) v1.Volume {
	return v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: path,
				Type: pathType,
			},
		},
	}
}

// NewVolumeMount creates a v1.VolumeMount to the specified location
func NewVolumeMount(name, path string, readOnly bool) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      name,
		MountPath: path,
		ReadOnly:  readOnly,
	}
}

// VolumeMapToSlice returns a slice of volumes from a map's values
func VolumeMapToSlice(volumes map[string]v1.Volume) []v1.Volume {
	v := make([]v1.Volume, 0, len(volumes))

	for _, vol := range volumes {
		v = append(v, vol)
	}

	return v
}

// WriteStaticPodToDisk writes a static pod file to disk
func (p *etcdStaticPod) WriteStaticPodToDisk(componentName string, pod v1.Pod) error {
	// creates target folder if not already exists
	if err := os.MkdirAll(p.manifestDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %q: %v", p.manifestDir, err)
	}

	// writes the pod to disk
	serialized, err := meta.MarshalToYAML(&pod, v1.SchemeGroupVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest for %q to YAML: %v", componentName, err)
	}

	filename := p.GetStaticPodFilepath()

	if err := ioutil.WriteFile(filename, serialized, 0600); err != nil {
		return fmt.Errorf("failed to write static pod manifest file for %q (%q): %v", componentName, filename, err)
	}

	return nil
}

// GetProbeAddress returns an IP address or 127.0.0.1 to use for liveness probes
// in static pod manifests.
func (p *etcdStaticPod) GetProbeAddress() string {
	hosts := p.cfg.ListenClientURLs.Hosts.List()
	if len(hosts) > 0 {
		host := hosts[0]

		// Return the IP if the URL contains an address instead of a name.
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
		// Use the local resolver to try resolving the name within the URL.
		// If the name can not be resolved, return an IPv4 loopback address.
		// Otherwise, select the first valid IPv4 address.
		// If the name does not resolve to an IPv4 address, select the first valid IPv6 address.
		if addrs, err := net.LookupIP(host); err == nil {
			var ip net.IP
			for _, addr := range addrs {
				if addr.To4() != nil {
					ip = addr
					break
				}
				if addr.To16() != nil && ip == nil {
					ip = addr
				}
			}
			return ip.String()
		}
	}
	return "127.0.0.1"
}
