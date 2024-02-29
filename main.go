/*
Copyright 2019 The Jetstack cert-manager contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"sync/atomic"

	"github.com/djkormo/adcs-simulator/certserv"
	"github.com/djkormo/adcs-simulator/version"

	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//"github.com/Azure/go-ntlmssp"
)

var (

	//caWorkDir = flag.Lookup("workdir").Value.(flag.Getter).Get().(string)
	caWorkDir = "/usr/local/adcs-sim"
	serverPem = caWorkDir + "/ca/server.pem"
	serverKey = caWorkDir + "/ca/server.key"
	serverCsr = caWorkDir + "/ca/server.csr"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	setupLog = ctrl.Log.WithName("adcs-sim")
)

func main() {
	port := flag.Int("port", 8443, "Port to listen on")
	dns := flag.String("dns", "", "Comma separated list of domains for the simulator server certificate")
	ips := flag.String("ips", "", "Comma separated list of IPs for the simulator server certificate")

	opts := zap.Options{
		Development: false, //was true
	}
	opts.BindFlags(flag.CommandLine)

	flag.Parse()

	// based on https://sdk.operatorframework.io/docs/building-operators/golang/references/logging/

	configLog := uzap.NewProductionEncoderConfig()
	// changing  time format to RFC3339Nano -> 2006-01-02T15:04:05.999999999Z07:00"
	configLog.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339Nano))
	}
	logfmtEncoder := zaplogfmt.NewEncoder(configLog)

	// Construct a new logr.logger.
	logger := zap.New(zap.UseDevMode(false), zap.WriteTo(os.Stdout), zap.Encoder(logfmtEncoder))
	ctrl.SetLogger(logger)

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	setupLog.Info("Starting ADCS Simulator", "Project", version.Project, "BuildTime", version.BuildTime, "Release", version.Release, "Commit", version.Commit)
	setupLog.Info("Http", "port ", port)
	setupLog.Info("Directories in in /ca", "directory", caWorkDir)

	files, err := ioutil.ReadDir(caWorkDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		setupLog.Info("Scanning directories", "Directory: ", file.Name(), "is dir", file.IsDir())
	}
	setupLog.Info("Files in /ca", "directory", caWorkDir)

	filesca, err := ioutil.ReadDir(caWorkDir + "/ca")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileca := range filesca {
		setupLog.Info("Scanning files", "File: ", fileca.Name(), "is dir", fileca.IsDir())
	}

	setupLog.Info("Scanning files", "File in directory: ", caWorkDir)
	setupLog.Info("Files in /templates", "directory", caWorkDir)

	filestemplate, err := ioutil.ReadDir(caWorkDir + "/templates")
	if err != nil {
		log.Fatal(err)
	}

	for _, filetemplate := range filestemplate {
		setupLog.Info("Scanning files", "File: ", filetemplate.Name(), "is dir", filetemplate.IsDir())
	}

	certserv, err := certserv.NewCertserv()
	if err != nil {
		fmt.Printf("Cannot initialize NewCertserv() : %s\n", err.Error())
	}
	err = generateServerCertificate(certserv, ips, dns)
	if err != nil {
		fmt.Printf("Cannot generate server certificate generateServerCertificate() : %s\n", err.Error())
	}

	isReady := &atomic.Value{}
	isReady.Store(false)

	go func() {
		setupLog.Info("HandleHealthz", "Readyz probe is negative by default...", "")
		time.Sleep(5 * time.Second) // 5 seconds hardcoded
		isReady.Store(true)
		setupLog.Info("HandleHealthz", "Readyz probe is positive", "")
	}()

	//http.HandleFunc("/", greeting)

	//livenes and readiness
	http.HandleFunc("/healthz", HandleHealthz)
	http.HandleFunc("/readyz", HandleReadyz(isReady))

	//http.HandleFunc("/readyz", HandleReadyz)
	http.HandleFunc("/certnew.cer", certserv.HandleCertnewCer)
	http.HandleFunc("/certnew.p7b", certserv.HandleCertnewP7b)
	http.HandleFunc("/certcarc.asp", certserv.HandleCertcarcAsp)
	http.HandleFunc("/certfnsh.asp", certserv.HandleCertfnshAsp)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), serverPem, serverKey, nil))
}

// Generate certificate for the simulator server TLS
func generateServerCertificate(cs *certserv.Certserv, ips *string, dns *string) error {

	setupLog.Info("Configuration", "workdir ", caWorkDir)
	setupLog.Info("Configuration", "ip ", *ips)
	setupLog.Info("Configuration", "dns ", *dns)
	var ipAddresses []net.IP
	if ips != nil && len(*ips) > 0 {
		for _, ipString := range strings.Split(*ips, ",") {

			ip := net.ParseIP(ipString)
			if ip == nil {
				fmt.Printf("Error parsing ip=%s\n", ipString)
				continue
			}
			ipAddresses = append(ipAddresses, ip)
		}
	}
	var dnsNames []string
	if dns != nil && len(*dns) > 0 {
		dnsNames = strings.Split(*dns, ",")
	}

	organization := []string{"ADCS simulator for cert-manager testing"}

	if len(ipAddresses) == 0 && len(dnsNames) == 0 {
		setupLog.Info("No subjects specified on certificate", "ipAddresses", ipAddresses, "dnsNames", dnsNames)
		return fmt.Errorf("no subjects specified on certificate")

	}
	var commonName string
	if len(dnsNames) > 0 {
		commonName = dnsNames[0]
	} else {
		commonName = ipAddresses[0].String()
	}
	// CSR
	pubKeyAlgo := x509.RSA
	sigAlgo := x509.SHA512WithRSA
	csr := &x509.CertificateRequest{
		Version:            3,
		SignatureAlgorithm: sigAlgo,
		PublicKeyAlgorithm: pubKeyAlgo,
		Subject: pkix.Name{
			Organization: organization,
			CommonName:   commonName,
		},
		DNSNames:    dnsNames,
		IPAddresses: ipAddresses,
		// TODO: work out how best to handle extensions/key usages here
		ExtraExtensions: []pkix.Extension{},
	}
	// Private key
	keySize := 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("error creating x509 key: %s", err.Error())
	}
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	setupLog.Info("Writting", "key file ", serverKey)
	err = os.WriteFile(serverKey, keyBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing key file: %s", err.Error())
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, privateKey)
	if err != nil {
		return fmt.Errorf("error creating x509 certificate request: %s", err.Error())
	}
	err = ioutil.WriteFile(serverCsr, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}), 0644)
	if err != nil {
		return fmt.Errorf("error writing CSR file: %s", err.Error())
	}

	certData, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return fmt.Errorf("error parsing x509 certificate request: %s", err.Error())
	}
	certPem, err := cs.CreateCertificateChainPem(certData)

	if err != nil {
		return fmt.Errorf("error creating x509 certificate: %s", err.Error())
	}
	err = ioutil.WriteFile(serverPem, []byte(certPem), 0644)
	if err != nil {
		return fmt.Errorf("error writing certificate file: %s", err.Error())
	}
	return nil
}

//Based on https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/

// Liveness

func HandleHealthz(w http.ResponseWriter, r *http.Request) {
	setupLog.Info("HandleHealthz", "check", "Ok")
	w.WriteHeader(http.StatusOK)

	return
}

// Readiness
func HandleReadyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		setupLog.Info("HandleReadyz", "check", "Ok")
		w.WriteHeader(http.StatusOK)
	}
}

// https://umesh.dev/posts/how-to-implement-http-basic-auth-in-gogolang
var users = map[string]string{
	"test": "secret",
}

func isAuthorised(username, password string) bool {
	pass, ok := users[username]
	if !ok {
		return false
	}

	return password == pass
}

func greeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "No basic auth present"}`))
		return
	}

	if !isAuthorised(username, password) {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Invalid username or password"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "welcome to golang world!"}`))
	return
}
