/*
Copyright 2023 AmidaWare Inc.

Licensed under the Tactical RMM License Version 1.0 (the “License”).
You may only use the Licensed Software in accordance with the License.
A copy of the License is available at:

https://license.tacticalrmm.com

*/

package agent

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func createAgentConfig(baseurl, agentid, apiurl, token, agentpk, cert, proxy, meshdir, natsport string, insecure bool, tmpdir string) {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SOFTWARE\TacticalRMM`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatalln("Error creating registry key:", err)
	}
	defer k.Close()

	err = k.SetStringValue("BaseURL", baseurl)
	if err != nil {
		log.Fatalln("Error creating BaseURL registry key:", err)
	}

	err = k.SetStringValue("AgentID", agentid)
	if err != nil {
		log.Fatalln("Error creating AgentID registry key:", err)
	}

	err = k.SetStringValue("ApiURL", apiurl)
	if err != nil {
		log.Fatalln("Error creating ApiURL registry key:", err)
	}

	err = k.SetStringValue("Token", token)
	if err != nil {
		log.Fatalln("Error creating Token registry key:", err)
	}

	err = k.SetStringValue("AgentPK", agentpk)
	if err != nil {
		log.Fatalln("Error creating AgentPK registry key:", err)
	}

	if len(cert) > 0 {
		err = k.SetStringValue("Cert", cert)
		if err != nil {
			log.Fatalln("Error creating Cert registry key:", err)
		}
	}

	if len(proxy) > 0 {
		err = k.SetStringValue("Proxy", proxy)
		if err != nil {
			log.Fatalln("Error creating Proxy registry key:", err)
		}
	}

	if len(meshdir) > 0 {
		err = k.SetStringValue("MeshDir", meshdir)
		if err != nil {
			log.Fatalln("Error creating MeshDir registry key:", err)
		}
	}

	if len(natsport) > 0 {
		err = k.SetStringValue("NatsStandardPort", natsport)
		if err != nil {
			log.Fatalln("Error creating NatsStandardPort registry key:", err)
		}
	}

	if insecure {
		err = k.SetStringValue("Insecure", "true")
		if err != nil {
			log.Fatalln("Error creating Insecure registry key:", err)
		}
	}

	if len(tmpdir) > 0 {
		err = k.SetStringValue("TmpDir", tmpdir)
		if err != nil {
			log.Fatalln("Error creating TmpDir registry key:", err)
		}
	}
}

func (a *Agent) checkExistingAndRemove(silent bool) {
	hasReg := false
	_, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\TacticalRMM`, registry.ALL_ACCESS)
	if err == nil {
		hasReg = true
	}
	if hasReg {
		tacUninst := filepath.Join(a.ProgramDir, a.GetUninstallExe())
		tacUninstArgs := [2]string{tacUninst, "/VERYSILENT"}

		// fork: never display a GUI window during install.
		// The interactive "remove existing installation" prompt has been
		// intentionally removed so installs are fully silent.
		fmt.Println("Existing installation found and must be removed before attempting to reinstall.")
		fmt.Println("Run the following command to uninstall, and then re-run this installer.")
		fmt.Printf(`"%s" %s `, tacUninstArgs[0], tacUninstArgs[1])
		os.Exit(0)
	}
}

func (a *Agent) installerMsg(msg, alert string, silent bool) {
	// fork: never display a GUI window during install.
	// All install status/error messages are written to stdout/log instead.
	fmt.Println(msg)

	if alert == "error" {
		a.Logger.Fatalln(msg)
	}
}
