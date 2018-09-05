package gcpauth

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"google.golang.org/api/iam/v1"
	"time"
)

func ServiceAccountLoginJwt(
	iamClient *iam.Service, exp time.Time, aud, project, serviceAccount string) (*iam.SignJwtResponse, error) {
	accountResource := fmt.Sprintf(gcputil.ServiceAccountTemplate, project, serviceAccount)

	payload, err := json.Marshal(map[string]interface{}{
		"sub": serviceAccount,
		"aud": aud,
		"exp": exp.Unix(),
	})
	if err != nil {
		return nil, err
	}
	jwtReq := &iam.SignJwtRequest{
		Payload: string(payload),
	}
	return iamClient.Projects.ServiceAccounts.SignJwt(accountResource, jwtReq).Do()
}
