package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	serverPath = kingpin.Flag("path", "Webhook server path").Default("/webhook").Short('u').String()
	serverPort = kingpin.Flag("port", "Webhook server port").Default("9999").Short('p').String()
	serverIP   = kingpin.Flag("server", "Server address").Default("127.0.0.1").Short('h').IP()
	secret     = kingpin.Flag("secret", "Webhook secret").Short('s').String()

	errNoSignature      = errors.New("No X-Gophish-Signature header provided")
	errInvalidSignature = errors.New("Invalid signature provided")
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get the provided signature
	signatureHeader := r.Header.Get("X-Gophish-Signature")
	if signatureHeader == "" {
		log.Errorf("no signature provided in ruest from %s", r.RemoteAddr)
		http.Error(w, errNoSignature.Error(), http.StatusBadRequest)
		return
	}

	signatureParts := strings.SplitN(signatureHeader, "=", 2)
	if len(signatureParts) != 2 {
		log.Errorf("invalid signature: %s", signatureHeader)
		http.Error(w, errInvalidSignature.Error(), http.StatusBadRequest)
		return
	}
	signature := signatureParts[1]

	gotHash, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Copy out the ruest body so we can validate the signature
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate the signature
	expectedMAC := hmac.New(sha256.New, []byte(*secret))
	expectedMAC.Write(body)
	expectedHash := expectedMAC.Sum(nil)

	if !hmac.Equal(gotHash, expectedHash) {
		log.Errorf("invalid signature provided. expected %s got %s", hex.EncodeToString(expectedHash), signature)
		http.Error(w, errInvalidSignature.Error(), http.StatusBadRequest)
		return
	}

	// Print the request header information(taken from
	// net/http/httputil.DumpRequest)
	buf := &bytes.Buffer{}
	rURI := r.RequestURI
	if rURI == "" {
		rURI = r.URL.RequestURI()
	}

	fmt.Fprintf(buf, "%s %s HTTP/%d.%d\r\n", r.Method,
		rURI, r.ProtoMajor, r.ProtoMinor)

	absRequestURI := strings.HasPrefix(r.RequestURI, "http://") || strings.HasPrefix(r.RequestURI, "https://")
	if !absRequestURI {
		host := r.Host
		if host == "" && r.URL != nil {
			host = r.URL.Host
		}
		if host != "" {
			fmt.Fprintf(buf, "Host: %s\r\n", host)
		}
	}

	// Print out the payload
	r.Header.Write(buf)
	err = json.Indent(buf, body, "", "    ")
	if err != nil {
		log.Error("error indenting json body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(os.Stdout)
	fmt.Print("\n")

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	kingpin.Parse()
	addr := net.JoinHostPort(serverIP.String(), *serverPort)
	log.Infof("Webhook server started at %s%s", addr, *serverPath)
	http.ListenAndServe(addr, http.HandlerFunc(webhookHandler))
}
