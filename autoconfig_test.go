package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLookupValid(t *testing.T) {
	domain := Domain{"marshland.ovh", ClientConfig{}}
	proto, service := "tcp", "imaps"

	got_addr, got_port, err := domain.lookup(service, proto)
	want_addr, want_port := "hermes.marshland.ovh", uint16(993)

	if err != nil {
		t.Errorf("Service %q://%q returned error: %v", proto, service, err)
	}
	if got_addr != want_addr {
		t.Errorf("Service %q://%q returned address %q, expected %q", proto, service, got_addr, want_addr)
	}
	if got_port != want_port {
		t.Errorf("Service %q://%q returned port %d, expected %d", proto, service, got_port, want_port)
	}
}

func TestLookupError(t *testing.T) {
	domain := Domain{"marshland.ovh", ClientConfig{}}
	proto, service := "udp", "imaps"

	_, _, err := domain.lookup(service, proto)

	if err == nil {
		t.Errorf("Service %q://%q should have errored", proto, service)
	}
}

func TestConfigError(t *testing.T) {
	domain := Domain{"marshland.co.uk", ClientConfig{}}
	_, err := domain.GenerateXml()

	if err == nil {
		t.Error("Incoming server parsing should have errored")
	}
}

func TestConfigIncoming(t *testing.T) {
	domain := Domain{"marshland.ovh", ClientConfig{}}
	domain.GenerateXml()

	got := domain.config.Providers[0].IncomingServers[0]
	want := IncomingServer{}
	want.Type = "imap"
	want.Hostname = "hermes.marshland.ovh"
	want.Port = 993
	want.SocketType = "SSL"
	want.Authentication = "password-cleartext"
	want.Username = "%EMAILLOCALPART%"

	if got != want {
		t.Error("Incoming server doesn't match expected value")
	}
}

func TestConfigOutgoing(t *testing.T) {
	domain := Domain{"marshland.ovh", ClientConfig{}}
	domain.GenerateXml()

	got := domain.config.Providers[0].OutgoingServers[0]
	want := OutgoingServer{}
	want.Type = "smtp"
	want.Hostname = "hermes.marshland.ovh"
	want.Port = 465
	want.SocketType = "SSL"
	want.Authentication = "password-cleartext"
	want.Username = "%EMAILLOCALPART%"

	if got != want {
		t.Error("Incoming server doesn't match expected value")
	}
}

func TestHTTServer(t *testing.T) {
	d := &Domain{"marshland.ovh", ClientConfig{}}

	s := httptest.NewServer(http.HandlerFunc(d.HttpHandler))
	defer s.Close()

	r, err := http.Get(s.URL)
	if err != nil {
		t.Error("HTTP server failed")
	}

	config, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		t.Error("HTTP response reading failed")
	}

	c := ClientConfig{}
	if xml.Unmarshal([]byte(config), &c) != nil {
		t.Error("HTTP response unmarshalling failed")
	}

	wantIn := IncomingServer{}
	wantIn.Type = "imap"
	wantIn.Hostname = "hermes.marshland.ovh"
	wantIn.Port = 993
	wantIn.SocketType = "SSL"
	wantIn.Authentication = "password-cleartext"
	wantIn.Username = "%EMAILLOCALPART%"

	wantOut := OutgoingServer{}
	wantOut.Type = "smtp"
	wantOut.Hostname = "hermes.marshland.ovh"
	wantOut.Port = 465
	wantOut.SocketType = "SSL"
	wantOut.Authentication = "password-cleartext"
	wantOut.Username = "%EMAILLOCALPART%"

	if c.Providers[0].IncomingServers[0] != wantIn {
		t.Error("Incoming server doesn't match expected value")
	}
	if c.Providers[0].OutgoingServers[0] != wantOut {
		t.Error("Outgoing server doesn't match expected value")
	}
}

func TestHTTPServerError(t *testing.T) {
	d := &Domain{"marshland.co.uk", ClientConfig{}}

	s := httptest.NewServer(http.HandlerFunc(d.HttpHandler))
	defer s.Close()

	r, _ := http.Get(s.URL)
	if r.StatusCode != 500 {
		t.Error("HTTP server should have errored")
	}
}
