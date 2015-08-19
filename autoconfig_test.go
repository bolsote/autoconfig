package main

import "testing"

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

func TestLookupEmpty(t *testing.T) {
	domain := Domain{"marshland.ovh", ClientConfig{}}
	proto, service := "udp", "imaps"

	_, _, err := domain.lookup(service, proto)

	if err == nil {
		t.Errorf("Service %q://%q should have errored", proto, service)
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
