package geoip

import (
	"go.dtapp.net/gostring"
)

var licenseKey = "" // 许可证密钥

func GetGeoLite2AsnDownloadUrl() string {
	return gostring.Replace("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=YOUR_LICENSE_KEY&suffix=tar.gz", "YOUR_LICENSE_KEY", licenseKey)
}

//func GetGeoLite2AsnCsvDownloadUrl() string {
//	return gostring.Replace("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN-CSV&license_key=YOUR_LICENSE_KEY&suffix=zip", "YOUR_LICENSE_KEY", licenseKey)
//}

func GetGeoLite2CityDownloadUrl() string {
	return gostring.Replace("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=YOUR_LICENSE_KEY&suffix=tar.gz", "YOUR_LICENSE_KEY", licenseKey)
}

//func GetGeoLite2CityCsvDownloadUrl() string {
//	return gostring.Replace("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City-CSV&license_key=YOUR_LICENSE_KEY&suffix=zip", "YOUR_LICENSE_KEY", licenseKey)
//}

func GetGeoLite2CountryDownloadUrl() string {
	return gostring.Replace("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country&license_key=YOUR_LICENSE_KEY&suffix=tar.gz", "YOUR_LICENSE_KEY", licenseKey)
}

//func GetGeoLite2CountryCsvDownloadUrl() string {
//	return gostring.Replace("https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-Country-CSV&license_key=YOUR_LICENSE_KEY&suffix=zip", "YOUR_LICENSE_KEY", licenseKey)
//}
