package output

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/roboll/autoscale"
	"github.com/roboll/autoscale-etc-hosts/provider"
)

const (
	start    = "# --begin autoscale-etc-hosts output--"
	end      = "# --end autoscale-etc-hosts output--"
	filename = "hosts"
)

//DoRemove removes entries originally created by this tool, as marked by the begin and end delimiters.
func DoRemove() error {
	desc, err := os.Stat(filename)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	output := bytes.NewBuffer(make([]byte, 0, desc.Size()))

	started := false
	for scanner.Scan() {
		text := scanner.Text()
		if !started {
			if text == start {
				started = true
			} else {
				output.WriteString(text + "\n")
			}
		} else {
			if text == end {
				started = false
			}
		}
	}

	err = file.Close()
	if err != nil {
		return err
	}
	err = os.Rename(filename, filename+".bk")
	if err != nil {
		return err
	}

	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()
	output.WriteTo(outfile)

	return nil
}

//DoCreate creates entries associated with the group and provider.
func DoCreate(p autoscale.Provider, group string, domain string, toStdout bool, publicIP bool) error {
	if len(domain) == 0 {
		var err error
		domain, err = getDomain()
		if err != nil {
			return fmt.Errorf("Unable to get domain for instance. %s", err)
		}
	}

	hosts, err := provider.InstanceHostnames(p, &group, publicIP)
	if err != nil {
		return fmt.Errorf("Failed to get hosts for group. %s", err)
	}
	log.Println("got the hostnames and ips")

	output := bytes.Buffer{}
	output.WriteString(start + "\n")
	for hostname, ip := range hosts {
		output.WriteString(fmt.Sprintf("%s %s.%s\n", ip, hostname, domain))
		log.Printf("wrote this string: %s", fmt.Sprintf("%s %s.%s\n", ip, hostname, domain))
	}
	output.WriteString(end + "\n")

	if toStdout {
		output.WriteTo(os.Stdout)
	} else {
		var file *os.File
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Println("thinking i should create the file")
			file, err = os.Create(filename)
			if err != nil {
				return err
			}
		} else {
			log.Println("naw, dont create, use an existing file")
			file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
		}
		output.WriteTo(file)
		defer file.Close()
	}
	return nil
}

func getDomain() (string, error) {
	cmd := exec.Command("hostname", "-d")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
