package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
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
	Server
}
type OutgoingServer struct {
	Server
}

type Provider struct {
	Id               string `xml:"id,attr"`
	Domain           string `xml:"domain"`
	DisplayName      string `xml:"displayName"`
	DisplayShortName string `xml:"displayShortName"`

	IncomingServers []IncomingServer `xml:"incomingServer"`
	OutgoingServers []OutgoingServer `xml:"outgoingServer"`
}

type ClientConfig struct {
	XMLName xml.Name `xml:"clientConfig"`

	Version   string     `xml:"version,attr"`
	Providers []Provider `xml:"emailProvider"`
}

type Domain struct {
	domain string
	config ClientConfig
}

// Lookup the given service, protocol pair in the domain SRV records.
func (d *Domain) lookup(service, proto string) (string, uint16, error) {
	_, addresses, err := net.LookupSRV(service, proto, d.domain)

	if err != nil {
		return "", 0, err
	}
	if len(addresses) == 0 {
		return "", 0, errors.New("No SRV records available for the given domain")
	}

	return strings.Trim(addresses[0].Target, "."), addresses[0].Port, nil
}

// Generate an autoconfig XML document based on the information obtained from
// querying the domain SRV records.
func (d *Domain) generate_xml() ([]byte, error) {
	// Incoming server.
	address_in, port_in, err := d.lookup("imaps", "tcp")
	if err != nil {
		return nil, err
	}

	incoming := IncomingServer{}
	incoming.Type = "imap"
	incoming.Hostname = address_in
	incoming.Port = port_in
	incoming.SocketType = "SSL"
	incoming.Authentication = "password-cleartext"
	incoming.Username = "%EMAILLOCALPART%"

	// Outgoing server.
	address_out, port_out, err := d.lookup("submission", "tcp")
	if err != nil {
		return nil, err
	}

	outgoing := OutgoingServer{}
	outgoing.Type = "smtp"
	outgoing.Hostname = address_out
	outgoing.Port = port_out
	outgoing.SocketType = "SSL"
	outgoing.Authentication = "password-cleartext"
	outgoing.Username = "%EMAILLOCALPART%"

	// Final data mangling.
	d.config = ClientConfig{
		Version: "1.1",
		Providers: []Provider{
			Provider{
				Id:               d.domain,
				Domain:           d.domain,
				DisplayName:      d.domain,
				DisplayShortName: d.domain,
				IncomingServers:  []IncomingServer{incoming},
				OutgoingServers:  []OutgoingServer{outgoing},
			},
		},
	}

	xmlconfig, err := xml.Marshal(&d.config)
	if err != nil {
		return nil, err
	}

	return xmlconfig, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	domain := &Domain{"marshland.ovh", ClientConfig{}}
	xmlconfig, err := domain.generate_xml()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, xml.Header, string(xmlconfig))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9090", nil)
}
