package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	//"github.com/kr/pretty"
)

func main() {
	router := mux.NewRouter()

	echoRouter := EchoRouter()
	router.PathPrefix("/echo/").Handler(negroni.New(
		negroni.HandlerFunc(ValidateRequest),
		negroni.Wrap(echoRouter),
	))

	pageRouter := PageRouter()
	router.PathPrefix("/").Handler(negroni.New(
		negroni.Wrap(pageRouter),
	))

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}

func ValidateRequest(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Checking request signature...")

	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		http.Error(rw, "Not Authorized", 401)
		return
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Not Authorized", 401)
		return
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		log.Println("Failed to parse certificate PEM.")
		http.Error(rw, "Not Authorized", 401)
		return
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, "Not Authorized", 401)
		return
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		log.Println("Amazon certificate expired.")
		http.Error(rw, "Not Authorized", 401)
		return
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		log.Println("Amazon certificate invalid.")
		http.Error(rw, "Not Authorized", 401)
		return
	}

	// Verify the key
	publicKey := cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, readerToSHA1(r.Body), encryptedSig)
	if err != nil {
		log.Println("Signature match failed.")
		http.Error(rw, "Not Authorized", 401)
		return
	}

	next(rw, r)
}

func readerToSHA1(input io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(input)

	h := sha1.New()
	io.WriteString(h, buf.String())
	return h.Sum(nil)
}

func readCert(certURL string) ([]byte, error) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, errors.New("Could not download Amazon cert file.")
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("Could not read Amazon cert file.")
	}

	return certContents, nil
}

func verifySignature() bool {
	return false
}

func verifyCert() bool {
	return false
}

func verifyCertURL(path string) bool {
	if !strings.HasSuffix(path, "/echo.api/echo-api-cert.pem") {
		return false
	}

	if !strings.HasPrefix(path, "https://s3.amazonaws.com/echo.api/") && !strings.HasPrefix(path, "https://s3.amazonaws.com:443/echo.api/") {
		return false
	}

	return true
}
