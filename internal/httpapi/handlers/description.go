package handlers

import (
	"log"
	"net/http"
)

const descriptionXML = `<?xml version="1.0" encoding="UTF-8" ?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
	<specVersion>
		<major>1</major>
		<minor>0</minor>
	</specVersion>

	<URLBase>http://192.168.30.104:80/</URLBase>

	<device>
		<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>

		<friendlyName>Home Assistant Bridge (192.168.30.104)</friendlyName>

		<manufacturer>Royal Philips Electronics</manufacturer>
		<manufacturerURL>http://www.philips.com</manufacturerURL>

		<modelDescription>Philips hue Personal Wireless Lighting</modelDescription>
		<modelName>Philips hue bridge 2015</modelName>
		<modelNumber>BSB002</modelNumber>
		<modelURL>http://www.meethue.com</modelURL>

		<serialNumber>001788FFFE23BFC2</serialNumber>

		<UDN>uuid:2f402f80-da50-11e1-9b23-001788255acc</UDN>
	</device>
</root>`

func GetDescription(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf(
			"HTTP request: %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		w.Header().Set(
			"Content-Type",
			"text/xml; charset=\"utf-8\"",
		)

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(descriptionXML))
	}
}
