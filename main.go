package main

import (
	"flag"
	"log"

	"github.com/roboll/autoscale"
	"github.com/roboll/autoscale-etc-hosts/output"
)

var providerName string
var groupName string
var domain string

var region string
var usePublicIP bool

var remove bool
var toStdout bool

func init() {
	flag.StringVar(&region, "region", "", "cloud region; defaults to this instance region, if possible")
	flag.BoolVar(&usePublicIP, "use-public-ip", false, "use public ip: default false")

	flag.StringVar(&providerName, "provider", "", "cloud provider")
	flag.StringVar(&groupName, "group", "", "autoscale group name; defaults to this instance autoscale group, if possible")
	flag.StringVar(&domain, "domain", "", "domain override: if empty, defaults to $(hostname -d)")

	flag.BoolVar(&remove, "remove", false, "remove entries created by autoscale-etc-hosts")
	flag.BoolVar(&toStdout, "stdout", false, "output to stdout")
}

func main() {
	flag.Parse()

	if remove {
		err := output.DoRemove()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var p autoscale.Provider
		switch providerName {
		case "aws":
			if region != "" {
				p = &autoscale.AWS{
					Region: &region,
				}
			} else {
				p = &autoscale.AWS{}
			}
		default:
			log.Fatal("Must specify a valid provider(aws).")
		}
		err := output.DoCreate(p, groupName, domain, toStdout, usePublicIP)
		if err != nil {
			log.Fatal(err)
		}
	}
}
