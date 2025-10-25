package models

type MileageLogBilling struct {
	Log MileageLog
	TotalTripCost USD
	TotalMemberBillings USD
	MemberBills map[int]MemberMileageLogBilling
}

type PerMemberBilling struct {

}

type MemberMileageLogBilling struct {
	Member                Member
	RegularTripsCost      USD
	LongDistanceTripsCost USD
}

type BillingMethod interface {
	Name() string
	TripCost(multiplier float64, UseSecondaryRate bool) USD
}

type SimplePerMileBilling struct {
	BillingName string
	BasePerMile USD
}

func (b *SimplePerMileBilling) Name() string {
	return b.BillingName
}

func (b *SimplePerMileBilling) TripCost(multiplier float64, UseSecondaryRate bool) USD {
	// if trip is less than 1 mile, trip = 1 mile
	if multiplier == 0.0 {
		return b.BasePerMile.Multiply(1.0)
	}
	return b.BasePerMile.Multiply(multiplier)
}

type TruckBilling struct {
	BillingName      string
	BasePerMile      USD
	SecondaryPerMile USD
	MinimumFee       USD
}

func (b *TruckBilling) Name() string {
	return b.BillingName
}

func (b *TruckBilling) TripCost(multiplier float64, UseSecondaryRate bool) USD {
	rate := b.BasePerMile

	if UseSecondaryRate {
		rate = b.SecondaryPerMile
	}

	cost := rate.Multiply(multiplier)
	// if trip is less than 1 mile, trip = 1 mile
	if multiplier == 0.0 {
		cost = rate.Multiply(1.0)
	}

	if cost.Float64() < b.MinimumFee.Float64() {
		return b.MinimumFee
	} else {
		return cost
	}
}

type LongDistanceBilling struct {
	BillingName   string
	SingleDayRate USD
	MultiDayRate  USD
}

func (b *LongDistanceBilling) Name() string {
	return b.BillingName
}

func (b *LongDistanceBilling) TripCost(multiplier float64, UseSecondaryRate bool) USD {
	if multiplier == 1 {
		return b.SingleDayRate
	} else {
		return b.MultiDayRate.Multiply(multiplier)
	}
}
