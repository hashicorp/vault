// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// CreateIdentityProviderDetails The representation of CreateIdentityProviderDetails
type CreateIdentityProviderDetails interface {

	// The OCID of your tenancy.
	GetCompartmentId() *string

	// The name you assign to the `IdentityProvider` during creation.
	// The name must be unique across all `IdentityProvider` objects in the
	// tenancy and cannot be changed.
	GetName() *string

	// The description you assign to the `IdentityProvider` during creation.
	// Does not have to be unique, and it's changeable.
	GetDescription() *string

	// The identity provider service or product.
	// Supported identity providers are Oracle Identity Cloud Service (IDCS) and Microsoft
	// Active Directory Federation Services (ADFS).
	// Example: `IDCS`
	GetProductType() CreateIdentityProviderDetailsProductTypeEnum

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	GetFreeformTags() map[string]string

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	GetDefinedTags() map[string]map[string]interface{}
}

type createidentityproviderdetails struct {
	JsonData      []byte
	CompartmentId *string                                      `mandatory:"true" json:"compartmentId"`
	Name          *string                                      `mandatory:"true" json:"name"`
	Description   *string                                      `mandatory:"true" json:"description"`
	ProductType   CreateIdentityProviderDetailsProductTypeEnum `mandatory:"true" json:"productType"`
	FreeformTags  map[string]string                            `mandatory:"false" json:"freeformTags"`
	DefinedTags   map[string]map[string]interface{}            `mandatory:"false" json:"definedTags"`
	Protocol      string                                       `json:"protocol"`
}

// UnmarshalJSON unmarshals json
func (m *createidentityproviderdetails) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalercreateidentityproviderdetails createidentityproviderdetails
	s := struct {
		Model Unmarshalercreateidentityproviderdetails
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.CompartmentId = s.Model.CompartmentId
	m.Name = s.Model.Name
	m.Description = s.Model.Description
	m.ProductType = s.Model.ProductType
	m.FreeformTags = s.Model.FreeformTags
	m.DefinedTags = s.Model.DefinedTags
	m.Protocol = s.Model.Protocol

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *createidentityproviderdetails) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Protocol {
	case "SAML2":
		mm := CreateSaml2IdentityProviderDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetCompartmentId returns CompartmentId
func (m createidentityproviderdetails) GetCompartmentId() *string {
	return m.CompartmentId
}

//GetName returns Name
func (m createidentityproviderdetails) GetName() *string {
	return m.Name
}

//GetDescription returns Description
func (m createidentityproviderdetails) GetDescription() *string {
	return m.Description
}

//GetProductType returns ProductType
func (m createidentityproviderdetails) GetProductType() CreateIdentityProviderDetailsProductTypeEnum {
	return m.ProductType
}

//GetFreeformTags returns FreeformTags
func (m createidentityproviderdetails) GetFreeformTags() map[string]string {
	return m.FreeformTags
}

//GetDefinedTags returns DefinedTags
func (m createidentityproviderdetails) GetDefinedTags() map[string]map[string]interface{} {
	return m.DefinedTags
}

func (m createidentityproviderdetails) String() string {
	return common.PointerString(m)
}

// CreateIdentityProviderDetailsProductTypeEnum Enum with underlying type: string
type CreateIdentityProviderDetailsProductTypeEnum string

// Set of constants representing the allowable values for CreateIdentityProviderDetailsProductTypeEnum
const (
	CreateIdentityProviderDetailsProductTypeIdcs CreateIdentityProviderDetailsProductTypeEnum = "IDCS"
	CreateIdentityProviderDetailsProductTypeAdfs CreateIdentityProviderDetailsProductTypeEnum = "ADFS"
)

var mappingCreateIdentityProviderDetailsProductType = map[string]CreateIdentityProviderDetailsProductTypeEnum{
	"IDCS": CreateIdentityProviderDetailsProductTypeIdcs,
	"ADFS": CreateIdentityProviderDetailsProductTypeAdfs,
}

// GetCreateIdentityProviderDetailsProductTypeEnumValues Enumerates the set of values for CreateIdentityProviderDetailsProductTypeEnum
func GetCreateIdentityProviderDetailsProductTypeEnumValues() []CreateIdentityProviderDetailsProductTypeEnum {
	values := make([]CreateIdentityProviderDetailsProductTypeEnum, 0)
	for _, v := range mappingCreateIdentityProviderDetailsProductType {
		values = append(values, v)
	}
	return values
}
