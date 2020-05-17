package ebay

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_IsNewItemFromTSV(t *testing.T) {

	var tsv string = "v1|110194763041|0	10 Colt Firearms Pins, Patch Cca Iacp Shot Nra Condition	http://i.ebayimg.com/00/s/ODA4WDgwNA==/z/MFsAAOSwDuJW1zYF/$_1.JPG?set_id=880000500F	Collectibles:Historical Memorabilia:Other Historical Memorabilia	208	Availability	therampantcolt	2210.0	99.94		Unspecified		94740293	1000	New	42.5	USD	1234	SELLER_DEFINED_VARIATIONS		99	US		TOP_RATED	AVAILABLE	FALSE		MORE_THAN	10	TRUE	30	DAY	EXCHANGE	BUYER	10		PAYPAL	SHIP_HOME	COUNTRY:US|CN;REGION:WORLD_WIDE|ASIA		94740293	190198459558	Apple	MQAM2LL/A	Q29sb3I=:TXVsdGkgQ29sb3I=;TWF0ZXJpYWw=:Q2xvdGggJiBWaW55bA==;U2hhZGU=:TXVsdGkgQ29sb3I=;U2l6ZQ==:MTMiIHggNyIgeCA0IiArIDYiIGhhbmRsZQ==;U3R5bGU=:Q29zbWV0aWMgU2hhdmluZyBCYWc=;VHlwZQ==:VHJhdmVsIEJhZw==	http://i.ebayimg.com/00/s/MTM3OFg5MTI=/z/vqUAAMXQ82FRGUmB/$_57.JPG?set_id=880000500F|http://i.ebayimg.com/00/s/MzYwWDQ1OA==/z/3nQAAOxyqUpQ6ZJH/$_57.JPG?set_id=8800005007														FedEx	FedEx 2Day	EXPEDITED	8.99	FIXED	1.99	1	110194763041"

	expItem := Item{
		ID:                            "v1|110194763041|0",
		Title:                         "10 Colt Firearms Pins, Patch Cca Iacp Shot Nra Condition",
		ImageURL:                      "http://i.ebayimg.com/00/s/ODA4WDgwNA==/z/MFsAAOSwDuJW1zYF/$_1.JPG?set_id=880000500F",
		Category:                      "Collectibles:Historical Memorabilia:Other Historical Memorabilia",
		CategoryID:                    "208",
		BuyingOptions:                 "Availability", // TODO: To check this: it is not an expected value
		SellerUsername:                "therampantcolt",
		SellerFeedbackPercentage:      "2210.0",
		SellerFeedbackScore:           "99.94",
		GTIN:                          "",
		Brand:                         "Unspecified",
		MPN:                           "",
		EPID:                          "94740293",
		ConditionID:                   "1000",
		Condition:                     "New",
		PriceValue:                    "42.5",
		PriceCurrency:                 "USD",
		PrimaryItemGroupID:            "1234",
		PrimaryItemGroupType:          "SELLER_DEFINED_VARIATIONS",
		EndDate:                       "",
		SellerItemRevision:            "99",
		LocationCountry:               "US",
		LocalizedAspects:              "",
		SellerTrustLevel:              "TOP_RATED",
		Availability:                  "AVAILABLE",
		ImageAlteringProhibited:       "FALSE",
		EstimatedAvailableQuantity:    "",
		AvailabilityThresholdType:     "MORE_THAN",
		AvailabilityThreshold:         "10",
		ReturnsAccepted:               "TRUE",
		ReturnPeriodValue:             "30",
		ReturnPeriodUnit:              "DAY",
		RefundMethod:                  "EXCHANGE",
		ReturnMethod:                  "BUYER",
		ReturnShippingCostPayer:       "10",
		AcceptedPaymentMethods:        "PAYPAL",
		DeliveryOptions:               "SHIP_HOME",
		ShipToIncludedRegions:         "COUNTRY:US|CN;REGION:WORLD_WIDE|ASIA",
		ShipToExcludedRegions:         "",
		InferredEPID:                  "94740293",
		InferredGTIN:                  "190198459558",
		InferredBrand:                 "Apple",
		InferredMPN:                   "MQAM2LL/A",
		InferredLocalizedAspects:      "Q29sb3I=:TXVsdGkgQ29sb3I=;TWF0ZXJpYWw=:Q2xvdGggJiBWaW55bA==;U2hhZGU=:TXVsdGkgQ29sb3I=;U2l6ZQ==:MTMiIHggNyIgeCA0IiArIDYiIGhhbmRsZQ==;U3R5bGU=:Q29zbWV0aWMgU2hhdmluZyBCYWc=;VHlwZQ==:VHJhdmVsIEJhZw==",
		AdditionalImages:              "http://i.ebayimg.com/00/s/MTM3OFg5MTI=/z/vqUAAMXQ82FRGUmB/$_57.JPG?set_id=880000500F|http://i.ebayimg.com/00/s/MzYwWDQ1OA==/z/3nQAAOxyqUpQ6ZJH/$_57.JPG?set_id=8800005007",
		OriginalPriceValue:            "",
		OriginalPriceCurrency:         "",
		DiscountAmount:                "",
		DiscountPercentage:            "",
		EnergyEfficiencyClass:         "",
		QualifiedPrograms:             "",
		LotSize:                       "",
		LengthUnitOfMeasure:           "",
		PackageWidth:                  "",
		PackageHeight:                 "",
		PackageLength:                 "",
		WeightUnitOfMeasure:           "",
		PackageWeight:                 "",
		ShippingCarrierCode:           "FedEx",
		ShippingServiceCode:           "FedEx 2Day",
		ShippingType:                  "EXPEDITED",
		ShippingCost:                  "8.99",
		ShippingCostType:              "FIXED",
		AdditionalShippingCostPerUnit: "1.99",
		QuantityUsedForEstimate:       "1",
		UnitPrice:                     "110194763041",
		UnitPricingMeasure:            "",
		LegacyItemID:                  "",
		Alerts:                        "",
	}

	item := NewItemFromTSV(tsv)
	assert.DeepEqual(t, expItem, *item)
}
