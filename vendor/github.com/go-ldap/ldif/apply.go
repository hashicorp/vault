package ldif

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

// Apply sends the LDIF entries to the server and does the changes as
// given by the entries.
//
// All *ldap.Entry are converted to an *ldap.AddRequest.
//
// By default, it returns on the first error. To continue with applying the
// LDIF, set the continueOnErr argument to true - in this case the errors
// are logged with log.Printf()
func (l *LDIF) Apply(conn ldap.Client, continueOnErr bool) error {
	for _, entry := range l.Entries {
		switch {
		case entry.Entry != nil:
			add := ldap.NewAddRequest(entry.Entry.DN, entry.Add.Controls)
			for _, attr := range entry.Entry.Attributes {
				add.Attribute(attr.Name, attr.Values)
			}
			entry.Add = add
			fallthrough
		case entry.Add != nil:
			if err := conn.Add(entry.Add); err != nil {
				if continueOnErr {
					log.Printf("ERROR: Failed to add %s: %s", entry.Add.DN, err)
					continue
				}
				return fmt.Errorf("failed to add %s: %s", entry.Add.DN, err)
			}

		case entry.Del != nil:
			if err := conn.Del(entry.Del); err != nil {
				if continueOnErr {
					log.Printf("ERROR: Failed to delete %s: %s", entry.Del.DN, err)
					continue
				}
				return fmt.Errorf("failed to delete %s: %s", entry.Del.DN, err)
			}

		case entry.Modify != nil:
			if err := conn.Modify(entry.Modify); err != nil {
				if continueOnErr {
					log.Printf("ERROR: Failed to modify %s: %s", entry.Modify.DN, err)
					continue
				}
				return fmt.Errorf("failed to modify %s: %s", entry.Modify.DN, err)
			}
		}
	}
	return nil
}
