package main

import (
	"flag"
	"log"

	"github.com/roboll/autoscale-etc-hosts/output"
	"github.com/roboll/autoscale-etc-hosts/provider"
)

var config provider.Config

var providerName string
var groupName string
var domain string

var remove bool
var toStdout bool

func init() {
	config = provider.Config{}
	flag.StringVar(&config.Region, "region", "", "cloud region; defaults to this instance region, if possible")
	flag.BoolVar(&config.UsePublicIP, "use-public-ip", false, "use public ip: default false")

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
		var p provider.Provider
		switch providerName {
		case "aws":
			p = &provider.AWS{
				Config: &config,
			}
		default:
			log.Fatal("Must specify a valid provider(aws).")
		}
		err := output.DoCreate(p, groupName, domain, toStdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
