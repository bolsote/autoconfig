autoconfig.go
=============

This small Go application translates DNS SRV records into an [Autoconfiguration file](https://developer.mozilla.org/en-US/docs/Mozilla/Thunderbird/Autoconfiguration/FileFormat/HowTo), compatible with a number of email clients. It works for my setup, which means a couple of assumptions:

 * Everything goes through SSL.
 * Password are sent in clear text; if you use SSL, anything else is just overkill.
 
The code is easy to change if you want something else, though. I just won't support it.

In case your SRV records have multiple entries, the code will just keep the one with the highest priority. A nice thing to have would be multiple server entries. The [Autoconfiguration file definition](https://wiki.mozilla.org/Thunderbird:Autoconfiguration:ConfigFileFormat) states that, in there are multiple ``incomingServer`` or ``outgoingServer`` entries, they must appear in order of priority. A simple way to implement this would be to sort the records by priority, and, keeping that order, generate multiple ``*Server`` entries. Since I don't use such a feature, it's not implemented.