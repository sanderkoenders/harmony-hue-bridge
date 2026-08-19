package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sanderkoenders/harmony-hue-bridge/internal/bridge"
)

const descriptionXmlTemplate = `
	<?xml version="1.0" encoding="UTF-8" ?>
	<root xmlns="urn:schemas-upnp-org:device-1-0">
		<specVersion>
			<major>1</major>
			<minor>0</minor>
		</specVersion>

		<URLBase>http://%s:%d/</URLBase>

		<device>
			<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>

			<friendlyName>Home Assistant Bridge (%s)</friendlyName>

			<manufacturer>Royal Philips Electronics</manufacturer>
			<manufacturerURL>http://www.philips.com</manufacturerURL>

			<modelDescription>Philips hue Personal Wireless Lighting</modelDescription>
			<modelName>Philips hue bridge 2015</modelName>
			<modelNumber>BSB002</modelNumber>
			<modelURL>http://www.meethue.com</modelURL>

			<serialNumber>%s</serialNumber>

			<UDN>uuid:%s</UDN>
		</device>
	</root>`

func parseDescriptionXML(bridge *bridge.Bridge) string {
	return fmt.Sprintf(descriptionXmlTemplate, bridge.IpAddr, bridge.Port, bridge.IpAddr, bridge.Username, bridge.UUID)
}

func HandleDescription(logger *log.Logger, bridge *bridge.Bridge) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("HTTP %s %s from %s", r.Method, r.URL.RequestURI(), r.RemoteAddr)

		w.Header().Set(
			"Content-Type",
			"text/xml; charset=\"utf-8\"",
		)

		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(parseDescriptionXML(bridge)))
	}
}
