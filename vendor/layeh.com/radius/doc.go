// Package radius provides a RADIUS client and server.
//
// Attributes
//
// The following tables list the attributes automatically registered in the
// Builtin dictionary. Each row contains the attributes' name, type (number),
// and Go data type.
//
// The following attributes are defined by RFC 2865:
//
//  User-Name                 1   string
//  User-Password             2   string
//  CHAP-Password             3   []byte
//  NAS-IP-Address            4   net.IP
//  NAS-Port                  5   uint32
//  Service-Type              6   uint32
//  Framed-Protocol           7   uint32
//  Framed-IP-Address         8   net.IP
//  Framed-IP-Netmask         9   net.IP
//  Framed-Routing            10  uint32
//  Filter-Id                 11  string
//  Framed-MTU                12  uint32
//  Framed-Compression        13  uint32
//  Login-IP-Host             14  net.IP
//  Login-Service             15  uint32
//  Login-TCP-Port            16  uint32
//  Reply-Message             18  string
//  Callback-Number           19  []byte
//  Callback-Id               20  []byte
//  Framed-Route              22  string
//  Framed-IPX-Network        23  net.IP
//  State                     24  []byte
//  Class                     25  []byte
//  Vendor-Specific           26  VendorSpecific
//  Session-Timeout           27  uint32
//  Idle-Timeout              28  uint32
//  Termination-Action        29  uint32
//  Called-Station-Id         30  []byte
//  Calling-Station-Id        31  []byte
//  NAS-Identifier            32  []byte
//  Proxy-State               33  []byte
//  Login-LAT-Service         34  []byte
//  Login-LAT-Node            35  []byte
//  Login-LAT-Group           36  []byte
//  Framed-AppleTalk-Link     37  uint32
//  Framed-AppleTalk-Network  38  uint32
//  Framed-AppleTalk-Zone     39  []byte
//  CHAP-Challenge            60  []byte
//  NAS-Port-Type             61  uint32
//  Port-Limit                62  uint32
//  Login-LAT-Port            63  []byte
//
// The following attributes are defined by RFC 2866:
//
//  Acct-Status-Type       40  uint32
//  Acct-Delay-Time        41  uint32
//  Acct-Input-Octets      42  uint32
//  Acct-Output-Octets     43  uint32
//  Acct-Session-Id        44  string
//  Acct-Authentic         45  uint32
//  Acct-Session-Time      46  uint32
//  Acct-Input-Packets     47  uint32
//  Acct-Output-Packets    48  uint32
//  Acct-Terminate-Cause   49  uint32
//  Acct-Multi-Session-Id  50  string
//  Acct-Link-Count        51  uint32
package radius // import "layeh.com/radius"
