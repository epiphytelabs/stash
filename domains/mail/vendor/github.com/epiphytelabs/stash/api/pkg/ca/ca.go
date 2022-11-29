package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/pkg/errors"
)

type CA struct {
	cert tls.Certificate
}

func New(domain string) (*CA, error) {
	cert, err := generateCACertificate(domain)
	if err != nil {
		return nil, err
	}

	c := &CA{
		cert: *cert,
	}

	return c, nil
}

func LoadFiles(pub, key string) (*CA, error) {
	cert, err := tls.LoadX509KeyPair(pub, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c := &CA{
		cert: cert,
	}

	return c, nil
}

func Blocks(c tls.Certificate) ([]byte, []byte, error) {
	if len(c.Certificate) < 1 {
		return nil, nil, errors.New("invalid certificate")
	}

	pk, ok := c.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("invalid private key")
	}

	p := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Certificate[0]})
	k := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	return p, k, nil
}

func (c *CA) Certificate() (*x509.Certificate, error) {
	cc, err := x509.ParseCertificate(c.cert.Certificate[0])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cc, nil
}

func (c *CA) Pool() *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.cert.Certificate[0]}))
	return pool
}

func (c *CA) GenerateClient(id string) (*tls.Certificate, error) {
	rkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cpub, err := x509.ParseCertificate(c.cert.Certificate[0])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   id,
			Organization: []string{"ca"},
		},
		Issuer:    cpub.Subject,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(2 * 365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageEmailProtection,
		},
		BasicConstraintsValid: true,
	}

	data, err := x509.CreateCertificate(rand.Reader, &template, cpub, &rkey.PublicKey, c.cert.PrivateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: data})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rkey)})

	cert, err := tls.X509KeyPair(pub, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cert, nil
}

func (c *CA) GenerateServer(host string, alts []string) (*tls.Certificate, error) {
	rkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cpub, err := x509.ParseCertificate(c.cert.Certificate[0])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dns := []string{host, fmt.Sprintf("*.%s", host)}

	for _, alt := range alts {
		dns = append(dns, alt, fmt.Sprintf("*.%s", alt))
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   host,
			Organization: []string{"ca"},
		},
		Issuer:    cpub.Subject,
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(2 * 365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageCodeSigning,
		},
		BasicConstraintsValid: true,
		DNSNames:              dns,
	}

	data, err := x509.CreateCertificate(rand.Reader, &template, cpub, &rkey.PublicKey, c.cert.PrivateKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: data})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rkey)})

	cert, err := tls.X509KeyPair(pub, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cert, nil
}

func (c *CA) WriteFiles(pub, key string) error {
	p, k, err := Blocks(c.cert)
	if err != nil {
		return err
	}

	if err := os.WriteFile(pub, p, 0644); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(key, k, 0600); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func generateCACertificate(domain string) (*tls.Certificate, error) {
	rkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	template := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{domain},
		SerialNumber:          serial,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(20 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		Subject: pkix.Name{
			CommonName: domain,
		},
	}

	data, err := x509.CreateCertificate(rand.Reader, &template, &template, &rkey.PublicKey, rkey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: data})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rkey)})

	cert, err := tls.X509KeyPair(pub, key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cert, nil
}
