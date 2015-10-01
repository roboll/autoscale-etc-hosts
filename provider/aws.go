package provider

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWS struct {
	Config *Config
}

//GetInstaneMap maps hostnames to ips.
func (a *AWS) GetInstanceMap(group *string) (map[string]string, error) {
	meta := ec2metadata.New(&ec2metadata.Config{})
	region := a.Config.Region
	if len(region) == 0 {
		var err error
		region, err = meta.Region()
		if err != nil {
			return nil, fmt.Errorf("Unable to get region from metadata service. %s", err)
		}
	}

	awsConfig := &aws.Config{Region: &region}
	ec := ec2.New(awsConfig)
	as := autoscaling.New(awsConfig)

	if group == nil {
		g, err := getGroupFromMeta(meta, ec)
		if err != nil {
			return nil, err
		}
		group = g
	}

	m := map[string]string{}
	output, err := as.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{group},
	})
	if err != nil {
		return nil, err
	}
	if len(output.AutoScalingGroups) != 1 {
		return nil, errors.New("Expected 1 autoscaling group for name.")
	}
	instances := []*string{}
	for _, group := range output.AutoScalingGroups {
		for _, instance := range group.Instances {
			instances = append(instances, instance.InstanceId)
		}
	}

	out, err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: instances,
	})
	if err != nil {
		return nil, err
	}

	for _, res := range out.Reservations {
		for _, instance := range res.Instances {
			var ip string
			if a.Config.UsePublicIP {
				ip = *instance.PublicIpAddress
			} else {
				ip = *instance.PrivateIpAddress
			}
			m[ipToHostname(ip)] = ip
		}
	}

	return m, nil
}

func ipToHostname(ip string) string {
	return strings.Replace(ip, ".", "-", -1)
}

func getGroupFromMeta(meta *ec2metadata.Client, ec *ec2.EC2) (*string, error) {
	id, err := meta.GetMetadata("instance-id")
	if err != nil {
		return nil, fmt.Errorf("Unable to get instance id from metadata service. %s", err)
	}
	output, err := ec.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{&id},
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to get instance data from ec2. %s", err)
	}

	if len(output.Reservations) != 1 {
		return nil, errors.New("Expected one reservation when searching by instance-id.")
	}
	for _, res := range output.Reservations {
		if len(res.Instances) != 1 {
			return nil, errors.New("Expected one instance when searching by instance-id.")
		}
		for _, instance := range res.Instances {
			for _, tag := range instance.Tags {
				if "aws:autoscaling:groupName" == *tag.Key {
					return tag.Value, nil
				}
			}
		}
	}
	return nil, errors.New("Unable to find autoscaling group tag.")
}
