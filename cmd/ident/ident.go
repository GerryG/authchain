package main

import (
	"net"
	"net/url"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	testCert()
}

func testCert() {

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{"ACME Corp"},
			Country:       []string{"USA"},
			Province:      []string{"Warner"},
			Locality:      []string{"Loony Tunes"},
			StreetAddress: []string{"1001 Main St"},
			PostalCode:    []string{"12300"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{ 1, 2, 3, 4, 5,  },
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,

	}
	u, _ :=	url.Parse("http://acme.com]")
	mycert := &x509.Certificate{
		SerialNumber: big.NewInt(9999),
		Issuer: pkix.Name{
			Organization:  []string{"ACME Corp"},
			Country:       []string{"USA"},
			Province:      []string{"Warner"},
			Locality:      []string{"Loony Tunes"},
			StreetAddress: []string{"1001 Main St"},
			PostalCode:    []string{"12300"},
		},
		Subject: pkix.Name{
			Organization:  []string{"Test Org Corp."},
			Country:       []string{"USA"},
			Province:      []string{"Tester"},
			Locality:      []string{"Test City"},
			StreetAddress: []string{"1001 Main St"},
			PostalCode:    []string{"12399"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		DNSNames:	[]string{"acme.com"},
		IPAddresses:	[]net.IP{net.ParseIP("192.168.0.0")},
		URIs:		[]*url.URL{u},
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 1024)

	mypriv, _ := rsa.GenerateKey(rand.Reader, 1024)

	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, &priv.PublicKey, priv)

	mycrt_b, err := x509.CreateCertificate(rand.Reader, mycert, ca, &mypriv.PublicKey, priv)

	// Write ca cert (pub key)
	certOut, err := os.Create("testdata/cacert.pem")
	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b})
	certOut.Close()
	log.Print("written cacert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("testdata/cakey.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("written cakey.pem\n")

	// Write test cert (pub key)
	certOut, err = os.Create("testdata/testcert.pem")
	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: mycrt_b})
	certOut.Close()
	log.Print("written testcert.pem\n")

	// Private key
	keyOut, err = os.OpenFile("testdata/testkey.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(mypriv)})
	keyOut.Close()
	log.Print("written testkey.pem\n")
}

