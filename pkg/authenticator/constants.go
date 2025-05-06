package authenticator

var FederatedAuthenticatorIDs = struct {
	Apple      string
	Duo        string
	EmailOTP   string
	Facebook   string
	GitHub     string
	GoogleOIDC string
	Hypr       string
	Iproov     string
	IWAKrb     string
	Microsoft  string
	MSLive     string
	Office365  string
	OIDC       string
	OrgEnt     string
	PwdReset   string
	SAML       string
	SIWE       string
	SMSOTP     string
	Twitter    string
	Yahoo      string
}{
	Apple:      "QXBwbGVPSURDQXV0aGVudGljYXRvcg",
	Duo:        "RHVvQXV0aGVudGljYXRvcg",
	EmailOTP:   "RW1haWxPVFA",
	Facebook:   "RmFjZWJvb2tBdXRoZW50aWNhdG9y",
	GitHub:     "R2l0aHViQXV0aGVudGljYXRvcg",
	GoogleOIDC: "R29vZ2xlT0lEQ0F1dGhlbnRpY2F0b3I",
	Hypr:       "SFlQUkF1dGhlbnRpY2F0b3I",
	Iproov:     "SXByb292QXV0aGVudGljYXRvcg",
	IWAKrb:     "SVdBS2VyYmVyb3NBdXRoZW50aWNhdG9y",
	Microsoft:  "T3BlbklEQ29ubmVjdEF1dGhlbnRpY2F0b3I",
	MSLive:     "TWljcm9zb2Z0V2luZG93c0xpdmVBdXRoZW50aWNhdG9y",
	Office365:  "T2ZmaWNlMzY1QXV0aGVudGljYXRvcg",
	OIDC:       "T3BlbklEQ29ubmVjdEF1dGhlbnRpY2F0b3I",
	OrgEnt:     "T3JnYW5pemF0aW9uQXV0aGVudGljYXRvcg",
	PwdReset:   "cGFzc3dvcmQtcmVzZXQtZW5mb3JjZXI",
	SAML:       "U0FNTFNTT0F1dGhlbnRpY2F0b3I",
	SIWE:       "T3BlbklEQ29ubmVjdEF1dGhlbnRpY2F0b3I",
	SMSOTP:     "U01TT1RQ",
	Twitter:    "VHdpdHRlckF1dGhlbnRpY2F0b3I",
	Yahoo:      "WWFob29PQXV0aDJBdXRoZW50aWNhdG9y",
}

var LocalAuthenticatorIDs = struct {
	ActiveSessionLimitHandler string
	BackupCode                string
	Basic                     string
	EmailOTP                  string
	FIDO                      string
	IdentifierFirst           string
	JWTBasic                  string
	MagicLink                 string
	PassiveSTS                string
	Push                      string
	SMSOTP                    string
	TOTP                      string
	X509Certificate           string
}{
	ActiveSessionLimitHandler: "U2Vzc2lvbkV4ZWN1dG9y",
	BackupCode:                "YmFja3VwLWNvZGUtYXV0aGVudGljYXRvcg",
	Basic:                     "QmFzaWNBdXRoZW50aWNhdG9y",
	EmailOTP:                  "ZW1haWwtb3RwLWF1dGhlbnRpY2F0b3I",
	FIDO:                      "RklET0F1dGhlbnRpY2F0b3I",
	IdentifierFirst:           "SWRlbnRpZmllckV4ZWN1dG9y",
	JWTBasic:                  "SldUQmFzaWNBdXRoZW50aWNhdG9y",
	MagicLink:                 "TWFnaWNMaW5rQXV0aGVudGljYXRvcg",
	PassiveSTS:                "UGFzc2l2ZVNUU0F1dGhlbnRpY2F0b3I",
	Push:                      "cHVzaC1ub3RpZmljYXRpb24tYXV0aGVudGljYXRvcg",
	SMSOTP:                    "c21zLW90cC1hdXRoZW50aWNhdG9y",
	TOTP:                      "dG90cA",
	X509Certificate:           "eDUwOUNlcnRpZmljYXRlQXV0aGVudGljYXRvcg",
}

var SocialAuthenticatorIDs = map[string]struct{}{
	FederatedAuthenticatorIDs.Apple:      {},
	FederatedAuthenticatorIDs.GoogleOIDC: {},
	FederatedAuthenticatorIDs.Facebook:   {},
	FederatedAuthenticatorIDs.Twitter:    {},
	FederatedAuthenticatorIDs.GitHub:     {},
}

var SecondFactorAuthenticatorIDs = map[string]struct{}{
	LocalAuthenticatorIDs.TOTP:         {},
	FederatedAuthenticatorIDs.Iproov:   {},
	FederatedAuthenticatorIDs.Duo:      {},
	FederatedAuthenticatorIDs.PwdReset: {},
}
