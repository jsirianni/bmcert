package cert

// CertConfig object describes the configuration for
// bmcert
type CertConfig struct {
    SkipVerify bool
    Verbose    bool
    Hostname     string
    OutputDir    string
    OutputFormat string
    Password     string
    AltNames     string
    IPsans       string
    URISans      string

	// certificate object is left unexported as it should
	// not be accessable outside of this library
	certificateReq certificateRequest
}

// SignedCertificate object describes a certificate signed
// by Vault
type SignedCertificate struct {
	Certificate string      `json:"certificate"`
	Issuing_ca string       `json:"issuing_ca"`
	Private_key string 	    `json:"private_key"`
	Private_key_type string `json:"private_key_type"`
	Serial_number string    `json:"serial_number"`
}

type certificateRequest struct {
	Common_name string `json:"common_name"`
	Alt_names   string `json:"alt_names"`
	Ip_sans     string `json:"ip_sans"`
	Uri_sans    string `json:"uri_sans"`
}

type apiResponse struct {
	Request_id string      `json:"request_id"`
	Lease_id   string      `json:"lease_id"`    // usually null
	Renewable  bool        `json:"renewable"`
	Lease_duration float32 `json:"lease_duration"`
	Data SignedCertificate `json:"data"`
	Wrap_info  string      `json:"wrap_info"`   // usually null
	Warnings   string      `json:"warnings"`    // usually null
	Auth       string      `json:"auth"`        // usually null
}
