package roles

import (
	"fmt"

	"github.com/hashicorp/vault/helper/activedirectory"
)

func getServiceAccountByName(adClient *activedirectory.Client, serviceAccountName string) (*activedirectory.Entry, error) {

	// TODO can this be in a serviceAccountNameFilter() util method in the higher package?
	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}

	entries, err := adClient.Search(filters)
	if err != nil {
		return nil, err
	}

	if len(entries) <= 0 {
		return nil, fmt.Errorf("service account of %s must already exist in active directory, searches are case sensitive", serviceAccountName)
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("expected one matching service account, but received %s", entries)
	}
	return entries[0], nil
}
