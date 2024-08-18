package main

import (
	"fmt"

	"github.com/matthxwpavin/ticketing/prettyjson"
)

type ConsultRequest struct {
	ConsultRequestId   int64   `json:"-"` // Don't appear on api information, used by 'LoadPreviousConsultRequest()' set for CreateConsultRequest is know this payload come from RetryCreate() service.
	SymptomSummary     string  `json:"symptomSummary"`
	Affected           int32   `json:"affected"`
	AffectedUnit       string  `json:"affectedUnit"`
	SymptomFileIds     []int64 `json:"symptomFileIds"`
	InsurancePolicyIds []int64 `json:"insurancePolicyIds"`
	PromotionCouponIds []int64 `json:"promotionCouponIds"`
	PaymentInfoId      int32   `json:"paymentInfoId"`
	ProviderType       int32   `json:"providerType"`
	ProviderId         string  `json:"providerId"`
	SpecialtyId        int32   `json:"specialtyId"`
	AppointmentId      *int32  `json:"appointmentId"`
	// Additional required for approving consult request
	AddressId         int32   `json:"addressId"`
	AttachmentFileIds []int32 `json:"attachmentFileIds"`
	UserId            *int64  `json:"userId"`
	Origin            string  `json:"origin"`
	CountryIso        string  `json:"-"`

	InsuranceCompanyBenefitId *int32 `json:"-" query:"insuranceCompanyBenefitId"`
}

func main() {
	fmt.Println(prettyjson.Stringify(&ConsultRequest{}))
}
