package domain

import "encoding/xml"

// MpesaExpressRequest the payload required to initiate an mpesa stk push
type MpesaExpressRequest struct {
	Amount           string `json:"amount"`
	PhoneNumber      string `json:"phone_number"`
	AccountReference string `json:"account_reference"`
	CallbackUrl      string `json:"callback_url"`
}

// MpesaExpressCallback the response from an mpesa stk push
type MpesaExpressCallback struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDescription string `json:"ResultDescription"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string `json:"Name"`
					Value any    `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// ValidationRequest store mpesa stk-push requests and respective responses
type ValidationRequest struct {
	TransactionType   string
	TransID           string
	TransTime         string
	TransAmount       string
	BusinessShortCode string
	BillRefNumber     string
	InvoiceNumber     string
	OrgAccountBalance string
	ThirdPartyTransID string
	MSISDN            string
	FirstName         string
	MiddleName        string
	LastName          string
}

// SoapPaymentConfirmationRequest represents the structure of the XML content
type SoapPaymentConfirmationRequest struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    Body     `xml:"Body"`
}

// Body represents the body of the SOAP envelope containing the payment confirmation request
type Body struct {
	C2BPaymentConfirmationRequest C2BPaymentConfirmationRequest `xml:"http://cps.huawei.com/cpsinterface/c2bpayment C2BPaymentConfirmationRequest"`
}

// C2BPaymentConfirmationRequest represents the structure of the payment confirmation request
type C2BPaymentConfirmationRequest struct {
	TransType         string    `xml:"TransType"`
	TransID           string    `xml:"TransID"`
	TransTime         string    `xml:"TransTime"`
	TransAmount       string    `xml:"TransAmount"`
	BusinessShortCode string    `xml:"BusinessShortCode"`
	BillRefNumber     string    `xml:"BillRefNumber"`
	OrgAccountBalance string    `xml:"OrgAccountBalance"`
	MSISDN            string    `xml:"MSISDN"`
	KYCInfo           []KYCInfo `xml:"KYCInfo"`
}

// KYCInfo represents the KYCInfo element within SoapPaymentConfirmationRequest
type KYCInfo struct {
	KYCName  string `xml:"KYCName"`
	KYCValue string `xml:"KYCValue"`
}
