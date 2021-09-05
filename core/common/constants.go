package common

import (
	"errors"
)

//Application core constants, like errors and default messages
var (
	AlreadyInSyncing               = errors.New("app already trying to sync")
	AlreadyUpdatingHotels          = errors.New("app already trying to update hotels")
	AlreadyUpdatingCities          = errors.New("app already trying to update cities")
	HotelIsUnavailable             = errors.New("هتل در حال حاضر موجود نیست")
	JsonDataIsNotValid             = errors.New("json data is not valid")
	HotelNotFound                  = errors.New("هتل مورد نظر یافت نشد")
	FAQNotFound                    = errors.New("FAQ cannot be found")
	OrderNotFound                  = errors.New("سفارش مورد نظر یافت نشد")
	AmentityCategoryNotFound       = errors.New("دسته بندی مورد نظر یافت نشد")
	CityNotFound                   = errors.New("شهر مورد نظر یافت نشد")
	AtLeastOneAdultNeeded          = errors.New("at least you should enter an adult")
	NationalityUnknown             = errors.New("you should enter at least nationality id or passport number")
	AmenityNotFound                = errors.New("amenity cannot be found")
	BadgeNotFound                  = errors.New("badge cannot be found")
	GetHotelListMaxTryLimitReached = errors.New("maximum try for getting hotels list reached")
	ErrorInConfirmingOrder         = errors.New("error while trying to confirm an order. please check the logs for more information")
	IndraOrderNotFound             = errors.New("indra order cannot be found")
	CannotGetHotelDataForSync      = errors.New("cannot get hotel data for sync")
	AccConsumerNum                 = 1
	AccConsumerSize                = 1
	ProviderRateLimitProblem       = errors.New("provider rate limit constrain problem detected")
	ProviderUnknownProblem         = errors.New("مشکلی در انجام درخواست شما رخ داده است. لطفا دوباره تلاش کنید")
	ProviderGateWayTimeOutProblem  = errors.New("provider gateway timeout(504) problem detected")
	HotelReserveForbidden          = errors.New("رزرو این هتل فقط در محیط پروداکشن امکان پذیر می باشد")
	RoomIsNotAvailable             = errors.New("رزرو این اتاق امکال پذیر نیست")
	DatesNotMatchError             = errors.New("تاریخ های انتخابی مغایرت دارد")

	HotelType_Hotel          = "hotel"
	HotelType_HotelApartment = "hotelapartment"
)

const MaxAgeAsAChild = 12
const NonRefundableDefaultMessage = "امکان لغو رزرو وجود ندارد"

var HotelAmenities = map[int]string{
	1027: "استخر",
	1048: "صندوق امانات",
	1056: "میز بیلیارد",
	1062: "سالن بدنسازی",
}

const SmsNotificationMessage = "حساب کاربری هتل علی بابا موجودی کافی برای ثبت سفارش ندارد. لطفا بررسی کنید.\n حساب کاربری : {account} \n موجودی فعلی : {currentBalance} \n مبلغ حد موجودی برای گزارش : {balanceLimit}"

var HotelTypes = map[string]int{
	HotelType_Hotel:          204,
	HotelType_HotelApartment: 219,
}

var HotelTypesString = map[int]string{
	204: HotelType_Hotel,
	219: HotelType_HotelApartment,
}

const (
	RefundStatus_PaymentFinalized = "PaymentFinalized"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

const MaxUint64 = ^uint64(0)
const MinUint64 = 0
const MaxInt64 = int64(MaxUint64 >> 1)
const MinInt64 = -MaxInt64 - 1

const ABChannelName = "JABAMA"
