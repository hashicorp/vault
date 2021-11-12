package configutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/hcl"

	wrapping "github.com/hashicorp/go-kms-wrapping"

	"github.com/hashicorp/hcl/hcl/ast"
)

var (
	// Allow override within the ent side of things.
	entValidateKmsLibrary = defaultValidateKmsLibrary
	nameRegexp            = regexp.MustCompile("^" + framework.GenericNameRegex("validate") + "$")
)

// KMSLibrary is a per-server configuration that will be further augmented with managed key configuration to
// build up a KMS wrapper type to access HSMs
type KMSLibrary struct {
	FoundKeys []string `hcl:",decodedFields"`
	Type      string   `hcl:"-"`
	Name      string   `hcl:"name"`
	Library   string   `hcl:"library"`
}

func (k *KMSLibrary) GoString() string {
	return fmt.Sprintf("*%#v", *k)
}

func defaultValidateKmsLibrary(kms *KMSLibrary) error {
	switch kms.Type {
	case wrapping.PKCS11:
		return fmt.Errorf("KMS type 'pkcs11' requires the Vault Enterprise HSM binary")

	default:
		return fmt.Errorf("unknown KMS type %q", kms.Type)
	}
}

func parseKmsLibraries(result *SharedConfig, list *ast.ObjectList) error {
	result.KmsLibraries = make(map[string]*KMSLibrary, len(list.Items))

	for _, item := range list.Items {
		library, err := decodeItem(item)
		if err != nil {
			return err
		}

		if err := validate(library); err != nil {
			return err
		}

		if _, ok := result.KmsLibraries[library.Name]; ok {
			return fmt.Errorf("duplicated kms_library configuration sections with name %s", library.Name)
		}

		result.KmsLibraries[library.Name] = library
	}
	return nil
}

func decodeItem(item *ast.ObjectItem) (*KMSLibrary, error) {
	library := &KMSLibrary{}
	if err := hcl.DecodeObject(&library, item.Val); err != nil {
		return nil, multierror.Prefix(err, "kms_library")
	}

	if len(item.Keys) != 1 {
		return nil, errors.New("kms_library section was missing a type")
	}

	library.Type = strings.ToLower(item.Keys[0].Token.Value().(string))
	library.Name = strings.ToLower(library.Name)

	return library, nil
}

func validate(obj *KMSLibrary) error {
	if obj.Library == "" {
		return fmt.Errorf("library key can not be blank within kms_library type: %s", obj.Type)
	}

	if obj.Name == "" {
		return fmt.Errorf("name key can not be blank within kms_library type: %s", obj.Name)
	}

	if !nameRegexp.MatchString(obj.Name) {
		return fmt.Errorf("value ('%s') for name field contained invalid characters", obj.Name)
	}

	if err := entValidateKmsLibrary(obj); err != nil {
		return err
	}

	return nil
}
