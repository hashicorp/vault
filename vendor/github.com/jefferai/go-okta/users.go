package okta

import (
	"time"
)

type User struct {
	ID              string     `json:"id"`
	Status          string     `json:"status"`
	Created         *time.Time `json:"created"`
	Activated       *time.Time `json:"activated"`
	StatusChanged   *time.Time `json:"statusChanged"`
	LastLogin       *time.Time `json:"lastLogin"`
	LastUpdated     *time.Time `json:"lastUpdated"`
	PasswordChanged *time.Time `json:"passwordChanged"`
	Profile         struct {
		Login             string `json:"login"`
		FirstName         string `json:"firstName"`
		LastName          string `json:"lastName"`
		NickName          string `json:"nickName"`
		DisplayName       string `json:"displayName"`
		Email             string `json:"email"`
		SecondEmail       string `json:"secondEmail"`
		ProfileURL        string `json:"profileUrl"`
		PreferredLanguage string `json:"preferredLanguage"`
		UserType          string `json:"userType"`
		Organization      string `json:"organization"`
		Title             string `json:"title"`
		Division          string `json:"division"`
		Department        string `json:"department"`
		CostCenter        string `json:"costCenter"`
		EmployeeNumber    string `json:"employeeNumber"`
		MobilePhone       string `json:"mobilePhone"`
		PrimaryPhone      string `json:"primaryPhone"`
		StreetAddress     string `json:"streetAddress"`
		City              string `json:"city"`
		State             string `json:"state"`
		ZipCode           string `json:"zipCode"`
		CountryCode       string `json:"countryCode"`
	} `json:"profile"`
	Credentials struct {
		Password struct {
		} `json:"password"`
		RecoveryQuestion struct {
			Question string `json:"question"`
		} `json:"recovery_question"`
		Provider struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"provider"`
	} `json:"credentials"`
	Links struct {
		ResetPassword struct {
			Href string `json:"href"`
		} `json:"resetPassword"`
		ResetFactors struct {
			Href string `json:"href"`
		} `json:"resetFactors"`
		ExpirePassword struct {
			Href string `json:"href"`
		} `json:"expirePassword"`
		ForgotPassword struct {
			Href string `json:"href"`
		} `json:"forgotPassword"`
		ChangeRecoveryQuestion struct {
			Href string `json:"href"`
		} `json:"changeRecoveryQuestion"`
		Deactivate struct {
			Href string `json:"href"`
		} `json:"deactivate"`
		ChangePassword struct {
			Href string `json:"href"`
		} `json:"changePassword"`
	} `json:"_links"`
}

type Groups []struct {
	ID      string `json:"id"`
	Profile struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"profile"`
}
