package ebay

import (
	"strings"
)

const (
	indexID                            = 0
	indexTitle                         = 1
	indexImageURL                      = 2
	indexCategory                      = 3
	indexCategoryID                    = 4
	indexBuyingOptions                 = 5
	indexSellerUsername                = 6
	indexSellerFeedbackPercentage      = 7
	indexSellerFeedbackScore           = 8
	indexGTIN                          = 9
	indexBrand                         = 10
	indexMPN                           = 11
	indexEPID                          = 12
	indexConditionID                   = 13
	indexCondition                     = 14
	indexPriceValue                    = 15
	indexPriceCurrency                 = 16
	indexPrimaryItemGroupID            = 17
	indexPrimaryItemGroupType          = 18
	indexEndDate                       = 19
	indexSellerItemRevision            = 20
	indexLocationCountry               = 21
	indexLocalizedAspects              = 22
	indexSellerTrustLevel              = 23
	indexAvailability                  = 24
	indexImageAlteringProhibited       = 25
	indexEstimatedAvailableQuantity    = 26
	indexAvailabilityThresholdType     = 27
	indexAvailabilityThreshold         = 28
	indexReturnsAccepted               = 29
	indexReturnPeriodValue             = 30
	indexReturnPeriodUnit              = 31
	indexRefundMethod                  = 32
	indexReturnMethod                  = 33
	indexReturnShippingCostPayer       = 34
	indexRestockingFeePercentage       = 35
	indexAcceptedPaymentMethods        = 36
	indexDeliveryOptions               = 37
	indexShipToIncludedRegions         = 38
	indexShipToExcludedRegions         = 39
	indexInferredEPID                  = 40
	indexInferredGTIN                  = 41
	indexInferredBrand                 = 42
	indexInferredMPN                   = 43
	indexInferredLocalizedAspects      = 44
	indexAdditionalImages              = 45
	indexOriginalPriceValue            = 46
	indexOriginalPriceCurrency         = 47
	indexDiscountAmount                = 48
	indexDiscountPercentage            = 49
	indexEnergyEfficiencyClass         = 50
	indexQualifiedPrograms             = 51
	indexLotSize                       = 52
	indexLengthUnitOfMeasure           = 53
	indexPackageWidth                  = 54
	indexPackageHeight                 = 55
	indexPackageLength                 = 56
	indexWeightUnitOfMeasure           = 57
	indexPackageWeight                 = 58
	indexShippingCarrierCode           = 59
	indexShippingServiceCode           = 60
	indexShippingType                  = 61
	indexShippingCost                  = 62
	indexShippingCostType              = 63
	indexAdditionalShippingCostPerUnit = 64
	indexQuantityUsedForEstimate       = 65
	indexUnitPrice                     = 66
	indexUnitPricingMeasure            = 67
	indexLegacyItemID                  = 68
	indexAlerts                        = 69
)

// Item represents an eBay Listing Item from the Feed.
// Note that no manipulation is done here: the values are extracted directly from the Feed file.
// Check https://developer.ebay.com/api-docs/buy/feed/resources/item/methods/getItemFeed for details.
type Item struct {
	ID                            string
	Title                         string
	ImageURL                      string
	Category                      string
	CategoryID                    string
	BuyingOptions                 string
	SellerUsername                string
	SellerFeedbackPercentage      string
	SellerFeedbackScore           string
	GTIN                          string
	Brand                         string
	MPN                           string
	EPID                          string
	ConditionID                   string
	Condition                     string
	PriceValue                    string
	PriceCurrency                 string
	PrimaryItemGroupID            string
	PrimaryItemGroupType          string
	EndDate                       string
	SellerItemRevision            string
	LocationCountry               string
	LocalizedAspects              string
	SellerTrustLevel              string
	Availability                  string
	ImageAlteringProhibited       string
	EstimatedAvailableQuantity    string
	AvailabilityThresholdType     string
	AvailabilityThreshold         string
	ReturnsAccepted               string
	ReturnPeriodValue             string
	ReturnPeriodUnit              string
	RefundMethod                  string
	ReturnMethod                  string
	ReturnShippingCostPayer       string
	RestockingFeePercentage       string
	AcceptedPaymentMethods        string
	DeliveryOptions               string
	ShipToIncludedRegions         string
	ShipToExcludedRegions         string
	InferredEPID                  string
	InferredGTIN                  string
	InferredBrand                 string
	InferredMPN                   string
	InferredLocalizedAspects      string
	AdditionalImages              string
	OriginalPriceValue            string
	OriginalPriceCurrency         string
	DiscountAmount                string
	DiscountPercentage            string
	EnergyEfficiencyClass         string
	QualifiedPrograms             string
	LotSize                       string
	LengthUnitOfMeasure           string
	PackageWidth                  string
	PackageHeight                 string
	PackageLength                 string
	WeightUnitOfMeasure           string
	PackageWeight                 string
	ShippingCarrierCode           string
	ShippingServiceCode           string
	ShippingType                  string
	ShippingCost                  string
	ShippingCostType              string
	AdditionalShippingCostPerUnit string
	QuantityUsedForEstimate       string
	UnitPrice                     string
	UnitPricingMeasure            string
	LegacyItemID                  string
	Alerts                        string
}

// NewItemFromTSV creates a new Item from its TSV definition as given in Feed file.
// The tsv string is a row from the TSV file.
func NewItemFromTSV(tsv string) *Item {
	values := strings.SplitAfter(tsv, "\t")

	item := &Item{
		ID:                            getStringValue(indexID, values),
		Title:                         getStringValue(indexTitle, values),
		ImageURL:                      getStringValue(indexImageURL, values),
		Category:                      getStringValue(indexCategory, values),
		CategoryID:                    getStringValue(indexCategoryID, values),
		BuyingOptions:                 getStringValue(indexBuyingOptions, values),
		SellerUsername:                getStringValue(indexSellerUsername, values),
		SellerFeedbackPercentage:      getStringValue(indexSellerFeedbackPercentage, values),
		SellerFeedbackScore:           getStringValue(indexSellerFeedbackScore, values),
		GTIN:                          getStringValue(indexGTIN, values),
		Brand:                         getStringValue(indexBrand, values),
		MPN:                           getStringValue(indexMPN, values),
		EPID:                          getStringValue(indexEPID, values),
		ConditionID:                   getStringValue(indexConditionID, values),
		Condition:                     getStringValue(indexCondition, values),
		PriceValue:                    getStringValue(indexPriceValue, values),
		PriceCurrency:                 getStringValue(indexPriceCurrency, values),
		PrimaryItemGroupID:            getStringValue(indexPrimaryItemGroupID, values),
		PrimaryItemGroupType:          getStringValue(indexPrimaryItemGroupType, values),
		EndDate:                       getStringValue(indexEndDate, values),
		SellerItemRevision:            getStringValue(indexSellerItemRevision, values),
		LocationCountry:               getStringValue(indexLocationCountry, values),
		LocalizedAspects:              getStringValue(indexLocalizedAspects, values),
		SellerTrustLevel:              getStringValue(indexSellerTrustLevel, values),
		Availability:                  getStringValue(indexAvailability, values),
		ImageAlteringProhibited:       getStringValue(indexImageAlteringProhibited, values),
		EstimatedAvailableQuantity:    getStringValue(indexEstimatedAvailableQuantity, values),
		AvailabilityThresholdType:     getStringValue(indexAvailabilityThresholdType, values),
		AvailabilityThreshold:         getStringValue(indexAvailabilityThreshold, values),
		ReturnsAccepted:               getStringValue(indexReturnsAccepted, values),
		ReturnPeriodValue:             getStringValue(indexReturnPeriodValue, values),
		ReturnPeriodUnit:              getStringValue(indexReturnPeriodUnit, values),
		RefundMethod:                  getStringValue(indexRefundMethod, values),
		ReturnMethod:                  getStringValue(indexReturnMethod, values),
		ReturnShippingCostPayer:       getStringValue(indexReturnShippingCostPayer, values),
		RestockingFeePercentage:       getStringValue(indexRestockingFeePercentage, values),
		AcceptedPaymentMethods:        getStringValue(indexAcceptedPaymentMethods, values),
		DeliveryOptions:               getStringValue(indexDeliveryOptions, values),
		ShipToIncludedRegions:         getStringValue(indexShipToIncludedRegions, values),
		ShipToExcludedRegions:         getStringValue(indexShipToExcludedRegions, values),
		InferredEPID:                  getStringValue(indexInferredEPID, values),
		InferredGTIN:                  getStringValue(indexInferredGTIN, values),
		InferredBrand:                 getStringValue(indexInferredBrand, values),
		InferredMPN:                   getStringValue(indexInferredMPN, values),
		InferredLocalizedAspects:      getStringValue(indexInferredLocalizedAspects, values),
		AdditionalImages:              getStringValue(indexAdditionalImages, values),
		OriginalPriceValue:            getStringValue(indexOriginalPriceValue, values),
		OriginalPriceCurrency:         getStringValue(indexOriginalPriceCurrency, values),
		DiscountAmount:                getStringValue(indexDiscountAmount, values),
		DiscountPercentage:            getStringValue(indexDiscountPercentage, values),
		EnergyEfficiencyClass:         getStringValue(indexEnergyEfficiencyClass, values),
		QualifiedPrograms:             getStringValue(indexQualifiedPrograms, values),
		LotSize:                       getStringValue(indexLotSize, values),
		LengthUnitOfMeasure:           getStringValue(indexLengthUnitOfMeasure, values),
		PackageWidth:                  getStringValue(indexPackageWidth, values),
		PackageHeight:                 getStringValue(indexPackageHeight, values),
		PackageLength:                 getStringValue(indexPackageLength, values),
		WeightUnitOfMeasure:           getStringValue(indexWeightUnitOfMeasure, values),
		PackageWeight:                 getStringValue(indexPackageWeight, values),
		ShippingCarrierCode:           getStringValue(indexShippingCarrierCode, values),
		ShippingServiceCode:           getStringValue(indexShippingServiceCode, values),
		ShippingType:                  getStringValue(indexShippingType, values),
		ShippingCost:                  getStringValue(indexShippingCost, values),
		ShippingCostType:              getStringValue(indexShippingCostType, values),
		AdditionalShippingCostPerUnit: getStringValue(indexAdditionalShippingCostPerUnit, values),
		QuantityUsedForEstimate:       getStringValue(indexQuantityUsedForEstimate, values),
		UnitPrice:                     getStringValue(indexUnitPrice, values),
		UnitPricingMeasure:            getStringValue(indexUnitPricingMeasure, values),
		LegacyItemID:                  getStringValue(indexLegacyItemID, values),
		Alerts:                        getStringValue(indexAlerts, values),
	}

	return item
}

func getStringValue(index int, values []string) string {
	if index < len(values) {
		return strings.TrimSpace(values[index])
	}

	return ""
}
