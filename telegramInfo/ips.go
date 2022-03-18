package telegramInfo

import (
	"net"
)

// Map for DataCenter IPs
var dcIPs = map[int8]net.IP{
	1: net.ParseIP("149.154.175.50"),
	2: net.ParseIP("149.154.167.50"),
	3: net.ParseIP("149.154.175.100"),
	4: net.ParseIP("149.154.167.91"),
	5: net.ParseIP("91.108.56.100"),
}

func GetIPFromDC(dcID int8) (net.IP, error) {
	if data, ok := dcIPs[dcID]; ok {
		return data, nil
	} else {
		return nil, ErrWrongDCID
	}
}
