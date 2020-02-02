// MIT License

// Copyright (c) (2020) Alok Parlikar <alok@parlikar.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type tailscaleCfg struct {
	LoginName string
}

func main() {
	login := flag.String("login", "", "Email address to connect with tailscale")
	flag.Parse()

	if *login == "" {
		flag.Usage()
	}

	if os.Geteuid() != 0 {
		log.Fatalf("please run this command as root (or with sudo)")
	}

	cfgPath := "/var/lib/tailscale/relay.conf"

	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("unable to read tailscale config: %v", err)
	}

	var cfg tailscaleCfg
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("unable to read tailscale configuration: %v", err)
	}

	// store this config in it's own copy. Overwrite if already exists.
	err = ioutil.WriteFile(fmt.Sprintf("%s.%s", cfgPath, cfg.LoginName), data, 0644)
	if err != nil {
		log.Fatalf("unable to save config file for %s: %v", *login, err)
	}

	if *login == cfg.LoginName {
		log.Printf("already running tailscale as %q", *login)
		return
	}

	// check if a config file already exists for the given login.
	loginCfgPath := fmt.Sprintf("%s.%s", cfgPath, *login)
	data, err = ioutil.ReadFile(loginCfgPath)
	if err != nil {
		log.Fatalf("could not read config for %q: %v. please run tailscale-login first.", *login, err)
	}

	// now overwrite the main config file
	err = ioutil.WriteFile(cfgPath, data, 0644)
	if err != nil {
		log.Fatalf("unable to write tailscale config file: %v", err)
	}

	// now restart the tailscale service
	cmd := exec.Command("systemctl", "restart", "tailscale-relay")
	err = cmd.Run()
	if err != nil {
		log.Fatal("unable to restart tailscale service: %v", err)
	}
}
