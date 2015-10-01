package provider

import (
	"log"
	"strings"

	"github.com/roboll/autoscale"
)

//AWSInstanceHostnames maps hostnames to ips.
func InstanceHostnames(p autoscale.Provider, group *string, publicIP bool) (map[string]string, error) {
	if group == nil {
		grp, err := p.GetLocalInstanceAutoscaleGroup()
		if err != nil {
			return nil, err
		}
		group = grp
	}

	log.Println("group was passed in, getting members")
	instances, err := p.GetInstancesInGroup(group)
	if err != nil {
		return nil, err
	}
	log.Println("now getting ips")
	ips, err := p.GetInstanceIPs(instances, publicIP)
	if err != nil {
		return nil, err
	}

	log.Println("finally... mapping hostnames")
	output := map[string]string{}
	for _, ip := range *ips {
		output[ipToHostname(ip)] = ip
	}

	return output, nil
}

func ipToHostname(ip string) string {
	return strings.Replace(ip, ".", "-", -1)
}
