package cert

// Cert object describes the configuration for bmcert
type Cert struct {
    SkipVerify bool
    Verbose    bool
    OverWrite  bool
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
	IssuingCa string       `json:"issuing_ca"`
	PrivateKey string 	    `json:"private_key"`
	PrivateKeyType string `json:"private_key_type"`
	SerialNumber string    `json:"serial_number"`
}

type certificateRequest struct {
	CommonName string `json:"common_name"`
	AltNames   string `json:"alt_names"`
	IPSans     string `json:"ip_sans"`
	URISans    string `json:"uri_sans"`
}

type apiResponse struct {
	RequestID  string      `json:"request_id"`
	LeaseID    string      `json:"lease_id"`    // usually null
	Renewable  bool        `json:"renewable"`
	LeaseDuration float32  `json:"lease_duration"`
	Data SignedCertificate `json:"data"`
	WrapInfo   string      `json:"wrap_info"`   // usually null
	Warnings   string      `json:"warnings"`    // usually null
	Auth       string      `json:"auth"`        // usually null
}
