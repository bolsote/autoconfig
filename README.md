autoconfig.go
=============

[![Build Status (Unix/Mac)](https://travis-ci.org/bolsote/autoconfig.go.svg?branch=master)](https://travis-ci.org/bolsote/autoconfig.go)
[![Build Status (Win)](https://ci.appveyor.com/api/projects/status/1mo5wb97ac7da7bw?svg=true)](https://ci.appveyor.com/project/bolsote/autoconfig-go)
[![Test Coverate](https://codecov.io/github/bolsote/autoconfig.go/coverage.svg?branch=master)](https://codecov.io/github/bolsote/autoconfig.go?branch=master)
[![Code Quality](https://codeclimate.com/github/bolsote/autoconfig.go/badges/gpa.svg)](https://codeclimate.com/github/bolsote/autoconfig.go)

This small Go application translates DNS SRV records into an [Autoconfiguration file](https://developer.mozilla.org/en-US/docs/Mozilla/Thunderbird/Autoconfiguration/FileFormat/HowTo), compatible with a number of email clients. It works for my setup, which means a couple of assumptions:

 * Everything goes through SSL.
 * Password are sent in clear text; if you use SSL, anything else is just overkill.
 
The code is easy to change if you want something else, though. I just won't support it.

In case your SRV records have multiple entries, the code will just keep the one with the highest priority. A nice thing to have would be multiple server entries. The [Autoconfiguration file definition](https://wiki.mozilla.org/Thunderbird:Autoconfiguration:ConfigFileFormat) states that, in there are multiple ``incomingServer`` or ``outgoingServer`` entries, they must appear in order of priority. Since ``net.LookupSRV()`` returns all found SRV records [sorted by priority](http://golang.org/pkg/net/#LookupSRV), just returning all of them in ``lookup()`` and iterating over them (preserving their order) to generate multiple ``*Server`` entries would suffice. Since I don't use such a feature, it's not implemented.

For the sake of completeness, let's include a list of possible features. I will get to implement them if I ever use them, or if I feel like playing with this:

 * Multiple servers.
 * Automatically guess ``SocketType`` from the port number. It would not work for unconventional setups, but if you have one, I'm sure you can modify the code yourself.
 * Try to open a session to determine the ``Authentication`` type. It would also help in trying to guess what goes in ``Username``, although it would produce failed login attemps (careful with software like fail2ban).
 * Maybe, in order to make all this more useful, make some of the ``*Server`` fields configurable (such as ``Username`` or ``Authentication``), e.g., through a configuration file.
