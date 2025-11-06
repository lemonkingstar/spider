package prequest

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	InsecureSkipVerify bool
	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string
	// the password to decrypt the certificate
	Password string
}

func ClientTslConfNoVerity() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
	}
}

func ClientTslConfVerityServer(caFile string) (*tls.Config, error) {
	caPool, err := loadCa(caFile)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{
		RootCAs: caPool,
	}

	return conf, nil
}

func ClientTLSConfVerity(caFile, certFile, keyFile, passwd string) (*tls.Config, error) {
	caPool, err := loadCa(caFile)
	if err != nil {
		return nil, err
	}

	cert, err := loadCertificates(certFile, keyFile, passwd)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            caPool,
		Certificates:       []tls.Certificate{*cert},
	}

	return conf, nil
}

func ServerTslConf(caFile, certFile, keyFile, passwd string) (*tls.Config, error) {
	if "" == caFile {
		return ServerTslConfVerity(certFile, keyFile, passwd)
	}

	return ServerTslConfVerityClient(caFile, certFile, keyFile, passwd)
}

func ServerTslConfVerity(certFile, keyFile, passwd string) (*tls.Config, error) {

	cert, err := loadCertificates(certFile, keyFile, passwd)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{*cert},
	}

	return conf, nil
}

func ServerTslConfVerityClient(caFile, certFile, keyFile, passwd string) (*tls.Config, error) {
	caPool, err := loadCa(caFile)
	if err != nil {
		return nil, err
	}

	cert, err := loadCertificates(certFile, keyFile, passwd)
	if err != nil {
		return nil, err
	}

	conf := &tls.Config{
		ClientCAs:    caPool,
		Certificates: []tls.Certificate{*cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return conf, nil
}

func loadCa(caFile string) (*x509.CertPool, error) {
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(ca); ok != true {
		return nil, fmt.Errorf("append ca cert failed")
	}

	return caPool, nil
}

func loadCertificates(certFile, keyFile, passwd string) (*tls.Certificate, error) {
	//key file
	priKey, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	if "" != passwd {
		priPem, _ := pem.Decode(priKey)
		if priPem == nil {
			return nil, fmt.Errorf("decode private key failed")
		}

		priDecrPem, decErr := x509.DecryptPEMBlock(priPem, []byte(passwd))
		if decErr != nil {
			return nil, decErr
		}

		priKey = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: priDecrPem,
		})
	}

	//certificate
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(certData, priKey)
	if err != nil {
		return nil, err
	}

	return &tlsCert, nil
}
