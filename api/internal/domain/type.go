package domain

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type SexualPreference string

const (
	PrefHeterosexual SexualPreference = "heterosexual"
	PrefHomosexual   SexualPreference = "homosexual"
	PrefBisexual     SexualPreference = "bisexual"
)

type NotificationType string

const (
	NotifLike    NotificationType = "like"
	NotifView    NotificationType = "view"
	NotifMatch   NotificationType = "match"
	NotifUnlike  NotificationType = "unlike"
	NotifMessage NotificationType = "message"
)

type AuthProvider string

const (
	ProviderLocal    AuthProvider = "local"
	ProviderGoogle   AuthProvider = "google"
	ProviderFacebook AuthProvider = "facebook"
	ProviderApple    AuthProvider = "apple"
)
