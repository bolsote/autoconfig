package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	Type           string `xml:"type,attr"`
	Hostname       string `xml:"hostname"`
	Port           uint16 `xml:"port"`
	SocketType     string `xml:"socketType"`
	Authentication string `xml:"authentication"`
	Username       string `xml:"username"`
}
type IncomingServer struct {
	XMLName struct{} `xml:"incomingServer"`

	Server
}
type OutgoingServer struct {
	XMLName struct{} `xml:"outgoingServer"`

	Server
}

type Provider struct {
	XMLName struct{} `xml:"emailProvider"`

	Id               string `xml:"id,attr"`
	Domain           string `xml:"domain"`
	DisplayName      string `xml:"displayName"`
	DisplayShortName string `xml:"displayShortName"`

	IncomingServers []IncomingServer
	OutgoingServers []OutgoingServer
}

type ClientConfig struct {
	XMLName struct{} `xml:"clientConfig"`

	Version   string `xml:"version,attr"`
	Providers []Provider
}

// Lookup the given service, protocol pair in the domain SRV records.
func lookup(service, proto, domain string) (string, uint16, error) {
	_, addresses, err := net.LookupSRV(service, proto, domain)

	if err != nil {
		return "", 0, err
	}
	if len(addresses) == 0 {
		return "", 0, errors.New("No SRV records available for the given domain")
	}

	return addresses[0].Target, addresses[0].Port, nil
}

// Generate an autoconfig XML document based on the information obtained from
// querying the domain SRV records.
func generate_xml(domain string) ([]byte, error) {
	// Incoming server.
	address_in, port_in, err := lookup("imaps", "tcp", domain)
	if err != nil {
		return nil, err
	}

	incoming := IncomingServer{}
	incoming.Type = "imap"
	incoming.Hostname = address_in
	incoming.Port = port_in
	incoming.SocketType = "STARTTLS"
	incoming.Authentication = "password-encrypted"
	incoming.Username = "%EMAILADDRESS%"

	// Outgoing server.
	address_out, port_out, err := lookup("submission", "tcp", domain)
	if err != nil {
		return nil, err
	}

	outgoing := OutgoingServer{}
	outgoing.Type = "smtp"
	outgoing.Hostname = address_out
	outgoing.Port = port_out
	outgoing.SocketType = "STARTTLS"
	outgoing.Authentication = "password-encrypted"
	outgoing.Username = "%EMAILADDRESS%"

	// Final data mangling.
	config := ClientConfig{
		Version: "1.1",
		Providers: []Provider{
			Provider{
				Id:               domain,
				Domain:           domain,
				DisplayName:      domain,
				DisplayShortName: domain,
				IncomingServers:  []IncomingServer{incoming},
				OutgoingServers:  []OutgoingServer{outgoing},
			},
		},
	}

	xmlconfig, err := xml.MarshalIndent(&config, "", "\t")
	if err != nil {
		return nil, err
	}

	return xmlconfig, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	xmlconfig, err := generate_xml("marshland.ovh")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, xml.Header, string(xmlconfig))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9000", nil)
}
