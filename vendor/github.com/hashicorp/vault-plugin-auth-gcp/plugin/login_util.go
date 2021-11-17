package gcpauth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"google.golang.org/api/iamcredentials/v1"
)

func ServiceAccountLoginJwt(iamClient *iamcredentials.Service, exp time.Time, aud, serviceAccount string) (*iamcredentials.SignJwtResponse, error) {
	accountResource := fmt.Sprintf(gcputil.ServiceAccountCredentialsTemplate, serviceAccount)

	payload, err := json.Marshal(map[string]interface{}{
		"sub": serviceAccount,
		"aud": aud,
		"exp": exp.Unix(),
	})
	if err != nil {
		return nil, err
	}
	jwtReq := &iamcredentials.SignJwtRequest{
		Payload: string(payload),
	}
	return iamClient.Projects.ServiceAccounts.SignJwt(accountResource, jwtReq).Do()
}
