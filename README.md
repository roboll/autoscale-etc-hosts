# autoscale-etc-hosts

Populate `/etc/hosts` with members of an autoscaling group.

# About

When using autoscale infrastructure, ips are unpredictable. With custom hostnames, in order to perform initial bootstrap, it is sometimes necessary to discover other members of the group in order to communicate with them. When IP addresses aren't sufficient (hostnames are necessary for TLS), and dns infrastructure is not available, use `autoscale-etc-hosts` at bootstrap to discover group members addressable by hostname.

*Warning*: In order to prevent hostname collisions, be sure to run `remove` after initial bootstrap is complete.

## Support

Currently supports AWS. Extending to other providers should be a trivial task; see `autoscale.Provider`. Pull requests welcome.

### AWS

Requires `autoscaling:DescribeAutoScalingGroups` and `ec2:DescribeInstances`.

## Get

Available on github releases for some platforms, or docker `roboll/autoscale-etc-hosts`.

# Usage

```
Usage of ./autoscale-etc-hosts:
  -domain string
    	domain override: if empty, defaults to $(hostname -d)
  -group string
    	autoscale group name; defaults to this instance autoscale group, if possible
  -provider string
    	cloud provider
  -region string
    	cloud region; defaults to this instance region, if possible
  -remove
    	remove entries created by autoscale-etc-hosts
  -stdout
    	output to stdout
  -use-public-ip
    	use public ip: default false
```
