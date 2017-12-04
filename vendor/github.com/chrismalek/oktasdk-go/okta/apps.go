package okta

import (
	"fmt"
	"net/url"
	"time"
)

type AppsService service

// AppFilterOptions is used to generate a "Filter" to search for different Apps
// The values here coorelate to API Search paramgters on the group API
type AppFilterOptions struct {
	NextURL       *url.URL `url:"-"`
	GetAllPages   bool     `url:"-"`
	NumberOfPages int      `url:"-"`
	Limit         int      `url:"limit,omitempty"`
}

type App struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Label         string    `json:"label"`
	Status        string    `json:"status"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Created       time.Time `json:"created"`
	Accessibility struct {
		SelfService      bool        `json:"selfService"`
		ErrorRedirectURL interface{} `json:"errorRedirectUrl"`
		LoginRedirectURL interface{} `json:"loginRedirectUrl"`
	} `json:"accessibility"`
	Visibility struct {
		AutoSubmitToolbar bool `json:"autoSubmitToolbar"`
		Hide              struct {
			IOS bool `json:"iOS"`
			Web bool `json:"web"`
		} `json:"hide"`
		AppLinks struct {
			TestorgoneCustomsaml20App1Link bool `json:"testorgone_customsaml20app_1_link"`
		} `json:"appLinks"`
	} `json:"visibility"`
	Features    []interface{} `json:"features"`
	SignOnMode  string        `json:"signOnMode"`
	Credentials struct {
		UserNameTemplate struct {
			Template string `json:"template"`
			Type     string `json:"type"`
		} `json:"userNameTemplate"`
		Signing struct {
		} `json:"signing"`
	} `json:"credentials"`
	Settings struct {
		App struct {
		} `json:"app"`
		Notifications struct {
			Vpn struct {
				Network struct {
					Connection string `json:"connection"`
				} `json:"network"`
				Message interface{} `json:"message"`
				HelpURL interface{} `json:"helpUrl"`
			} `json:"vpn"`
		} `json:"notifications"`
		SignOn struct {
			DefaultRelayState     string        `json:"defaultRelayState"`
			SsoAcsURL             string        `json:"ssoAcsUrl"`
			IdpIssuer             string        `json:"idpIssuer"`
			Audience              string        `json:"audience"`
			Recipient             string        `json:"recipient"`
			Destination           string        `json:"destination"`
			SubjectNameIDTemplate string        `json:"subjectNameIdTemplate"`
			SubjectNameIDFormat   string        `json:"subjectNameIdFormat"`
			ResponseSigned        bool          `json:"responseSigned"`
			AssertionSigned       bool          `json:"assertionSigned"`
			SignatureAlgorithm    string        `json:"signatureAlgorithm"`
			DigestAlgorithm       string        `json:"digestAlgorithm"`
			HonorForceAuthn       bool          `json:"honorForceAuthn"`
			AuthnContextClassRef  string        `json:"authnContextClassRef"`
			SpIssuer              interface{}   `json:"spIssuer"`
			RequestCompressed     bool          `json:"requestCompressed"`
			AttributeStatements   []interface{} `json:"attributeStatements"`
		} `json:"signOn"`
	} `json:"settings"`
	Links struct {
		Logo []struct {
			Name string `json:"name"`
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"logo"`
		AppLinks []struct {
			Name string `json:"name"`
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"appLinks"`
		Help struct {
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"help"`
		Users struct {
			Href string `json:"href"`
		} `json:"users"`
		Deactivate struct {
			Href string `json:"href"`
		} `json:"deactivate"`
		Groups struct {
			Href string `json:"href"`
		} `json:"groups"`
		Metadata struct {
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"metadata"`
	} `json:"_links"`
}

func (a App) String() string {
	// return Stringify(g)
	return fmt.Sprintf("App:(ID: {%v} - Name: {%v})\n", a.ID, a.Name)
}

// GetByID gets a group from OKTA by the Gropu ID. An error is returned if the group is not found
func (a *AppsService) GetByID(appID string) (*App, *Response, error) {

	u := fmt.Sprintf("apps/%v", appID)
	req, err := a.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	app := new(App)

	resp, err := a.client.Do(req, app)

	if err != nil {
		return nil, resp, err
	}

	return app, resp, err
}

type AppUser struct {
	ID              string     `json:"id"`
	ExternalID      string     `json:"externalId"`
	Created         time.Time  `json:"created"`
	LastUpdated     time.Time  `json:"lastUpdated"`
	Scope           string     `json:"scope"`
	Status          string     `json:"status"`
	StatusChanged   *time.Time `json:"statusChanged"`
	PasswordChanged *time.Time `json:"passwordChanged"`
	SyncState       string     `json:"syncState"`
	LastSync        *time.Time `json:"lastSync"`
	Credentials     struct {
		UserName string `json:"userName"`
		Password struct {
		} `json:"password"`
	} `json:"credentials"`
	Profile struct {
		SecondEmail      interface{} `json:"secondEmail"`
		LastName         string      `json:"lastName"`
		MobilePhone      interface{} `json:"mobilePhone"`
		Email            string      `json:"email"`
		SalesforceGroups []string    `json:"salesforceGroups"`
		Role             string      `json:"role"`
		FirstName        string      `json:"firstName"`
		Profile          string      `json:"profile"`
	} `json:"profile"`
	Links struct {
		App struct {
			Href string `json:"href"`
		} `json:"app"`
		User struct {
			Href string `json:"href"`
		} `json:"user"`
	} `json:"_links"`
}

// GetUsers returns the members in an App
//   Pass in an optional AppFilterOptions struct to filter the results
//   The Users in the app are returned
func (a *AppsService) GetUsers(appID string, opt *AppFilterOptions) (appUsers []AppUser, resp *Response, err error) {

	pagesRetreived := 0
	var u string
	if opt.NextURL != nil {
		u = opt.NextURL.String()
	} else {
		u = fmt.Sprintf("apps/%v/users", appID)

		if opt.Limit == 0 {
			opt.Limit = defaultLimit
		}

		u, _ = addOptions(u, opt)
	}

	req, err := a.client.NewRequest("GET", u, nil)

	if err != nil {
		fmt.Printf("____ERROR HERE\n")
		return nil, nil, err
	}
	resp, err = a.client.Do(req, &appUsers)

	if err != nil {
		fmt.Printf("____ERROR HERE 2\n")
		return nil, resp, err
	}

	pagesRetreived++

	if (opt.NumberOfPages > 0 && pagesRetreived < opt.NumberOfPages) || opt.GetAllPages {

		for {

			if pagesRetreived == opt.NumberOfPages {
				break
			}
			if resp.NextURL != nil {

				var userPage []AppUser
				pageOpts := new(AppFilterOptions)
				pageOpts.NextURL = resp.NextURL
				pageOpts.Limit = opt.Limit
				pageOpts.NumberOfPages = 1

				userPage, resp, err = a.GetUsers(appID, pageOpts)

				if err != nil {
					return appUsers, resp, err
				} else {
					appUsers = append(appUsers, userPage...)
					pagesRetreived++
				}
			} else {
				break
			}

		}
	}

	return appUsers, resp, err
}
