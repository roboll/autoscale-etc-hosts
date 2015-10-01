package provider

import (
	"strings"

	"github.com/roboll/autoscale"
)

//AWSInstanceHostnames maps hostnames to ips.
func InstanceHostnames(p autoscale.Provider, group string, publicIP bool) (map[string]string, error) {
	if len(group) == 0 {
		grp, err := p.GetLocalInstanceAutoscaleGroup()
		if err != nil {
			return nil, err
		}
		group = *grp
	}

	instances, err := p.GetInstancesInGroup(&group)
	if err != nil {
		return nil, err
	}
	ips, err := p.GetInstanceIPs(instances, publicIP)
	if err != nil {
		return nil, err
	}

	output := map[string]string{}
	for _, ip := range *ips {
		output[ipToHostname(ip)] = ip
	}

	return output, nil
}

func ipToHostname(ip string) string {
	return strings.Replace(ip, ".", "-", -1)
}
