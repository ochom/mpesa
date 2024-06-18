package utils

import (
	"testing"

	"github.com/ochom/mpesa/src/domain"
)

var soapXML = `
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
		  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		  <soapenv:Body>
		    <ns1:C2BPaymentConfirmationRequest xmlns:ns1="http://cps.huawei.com/cpsinterface/c2bpayment">
		      <TransType>Pay Bill</TransType>
		      <TransID>SFI9TYPX9J</TransID>
		      <TransTime>20240618082632</TransTime>
		      <TransAmount>270.00</TransAmount>
		      <BusinessShortCode>290028</BusinessShortCode>
		      <BillRefNumber>254715842888</BillRefNumber>
		      <OrgAccountBalance>106915.00</OrgAccountBalance>
		      <MSISDN>207790dd287fab583ffbe503bda29278861b4c063d891b0ef931e21ef83988e7</MSISDN>
		      <KYCInfo>
		        <KYCName>[Personal Details][First Name]</KYCName>
		        <KYCValue>Chebet</KYCValue>
		      </KYCInfo>
		      <KYCInfo>
		        <KYCName>[Personal Details][Middle Name]</KYCName>
		        <KYCValue/>
		      </KYCInfo>
		      <KYCInfo>
		        <KYCName>[Personal Details][Last Name]</KYCName>
		        <KYCValue/>
		      </KYCInfo>
		    </ns1:C2BPaymentConfirmationRequest>
		  </soapenv:Body>
		</soapenv:Envelope>
	`

func TestParseXml(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "parse xml",
			args: args{
				data: soapXML,
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseXml[domain.SoapPaymentConfirmationRequest](tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("ParseXml() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}
