package main

import (
	"log"
	"os/exec"
)

func genCertificateSimple() {
	args := []string{
		"req",
		"-subj",
		"/C=/ST=/O=/L=/CN=localhost/OU=/",
		"-x509",
		"-nodes",
		"-days",
		"3650",
		"-newkey",
		"rsa:4096",
		"-keyout",
		"generated-key.pem",
		"-out",
		"generated-cert.pem",
	}
	cmd := exec.Command("openssl", args...)
	_, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to generated certificate: ", err.Error())
	}
	log.Println("Key was written to generated.key.")
}
