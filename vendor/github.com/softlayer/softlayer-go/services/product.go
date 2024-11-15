/**
 * Copyright 2016-2024 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 * on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// AUTOMATICALLY GENERATED CODE - DO NOT MODIFY

package services

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

// The SoftLayer_Product_Item data type contains general information relating to a single SoftLayer product.
type Product_Item struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductItemService returns an instance of the Product_Item SoftLayer service
func GetProductItemService(sess session.SLSession) Product_Item {
	return Product_Item{Session: sess}
}

func (r Product_Item) Id(id int) Product_Item {
	r.Options.Id = &id
	return r
}

func (r Product_Item) Mask(mask string) Product_Item {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Item) Filter(filter string) Product_Item {
	r.Options.Filter = filter
	return r
}

func (r Product_Item) Limit(limit int) Product_Item {
	r.Options.Limit = &limit
	return r
}

func (r Product_Item) Offset(offset int) Product_Item {
	r.Options.Offset = &offset
	return r
}

// Retrieve
func (r Product_Item) GetActivePresaleEvents() (resp []datatypes.Sales_Presale_Event, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getActivePresaleEvents", nil, &r.Options, &resp)
	return
}

// Retrieve Active usage based prices.
func (r Product_Item) GetActiveUsagePrices() (resp []datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getActiveUsagePrices", nil, &r.Options, &resp)
	return
}

// Retrieve The attribute values for a product item. These are additional properties that give extra information about the product being sold.
func (r Product_Item) GetAttributes() (resp []datatypes.Product_Item_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve Attributes that govern when an item may no longer be available.
func (r Product_Item) GetAvailabilityAttributes() (resp []datatypes.Product_Item_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getAvailabilityAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve An item's special billing type, if applicable.
func (r Product_Item) GetBillingType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getBillingType", nil, &r.Options, &resp)
	return
}

// Retrieve An item's included product item references. Some items have other items included in them that we specifically detail. They are here called Bundled Items. An example is Plesk unlimited. It as a bundled item labeled 'SiteBuilder'. These are the SoftLayer_Product_Item_Bundles objects. See the SoftLayer_Product_Item::bundleItems property for bundle of SoftLayer_Product_Item of objects.
func (r Product_Item) GetBundle() (resp []datatypes.Product_Item_Bundles, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getBundle", nil, &r.Options, &resp)
	return
}

// Retrieve An item's included products. Some items have other items included in them that we specifically detail. They are here called Bundled Items. An example is Plesk unlimited. It as a bundled item labeled 'SiteBuilder'. These are the SoftLayer_Product_Item objects.
func (r Product_Item) GetBundleItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getBundleItems", nil, &r.Options, &resp)
	return
}

// Retrieve When the product capacity is best described as a range, this holds the ceiling of the range.
func (r Product_Item) GetCapacityMaximum() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getCapacityMaximum", nil, &r.Options, &resp)
	return
}

// Retrieve When the product capacity is best described as a range, this holds the floor of the range.
func (r Product_Item) GetCapacityMinimum() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getCapacityMinimum", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates that this product is restricted by a capacity on a related product.
func (r Product_Item) GetCapacityRestrictedProductFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getCapacityRestrictedProductFlag", nil, &r.Options, &resp)
	return
}

// Retrieve An item's associated item categories.
func (r Product_Item) GetCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getCategories", nil, &r.Options, &resp)
	return
}

// Retrieve Some product items have configuration templates which can be used to during provisioning of that product.
func (r Product_Item) GetConfigurationTemplates() (resp []datatypes.Configuration_Template, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getConfigurationTemplates", nil, &r.Options, &resp)
	return
}

// Retrieve An item's conflicts. For example, McAfee LinuxShield cannot be ordered with Windows. It was not meant for that operating system and as such is a conflict.
func (r Product_Item) GetConflicts() (resp []datatypes.Product_Item_Resource_Conflict, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getConflicts", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates that this product is restricted by the number of cores on the compute instance. This is deprecated. Use [[SoftLayer_Product_Item/getCapacityRestrictedProductFlag|getCapacityRestrictedProductFlag]]
func (r Product_Item) GetCoreRestrictedItemFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getCoreRestrictedItemFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Some product items have a downgrade path. This is the first product item in the downgrade path.
func (r Product_Item) GetDowngradeItem() (resp datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getDowngradeItem", nil, &r.Options, &resp)
	return
}

// Retrieve Some product items have a downgrade path. These are those product items.
func (r Product_Item) GetDowngradeItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getDowngradeItems", nil, &r.Options, &resp)
	return
}

// Retrieve An item's category conflicts. For example, 10 Gbps redundant network functionality cannot be ordered with a secondary GPU and as such is a conflict.
func (r Product_Item) GetGlobalCategoryConflicts() (resp []datatypes.Product_Item_Resource_Conflict, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getGlobalCategoryConflicts", nil, &r.Options, &resp)
	return
}

// Retrieve The generic hardware component that this item represents.
func (r Product_Item) GetHardwareGenericComponentModel() (resp datatypes.Hardware_Component_Model_Generic, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getHardwareGenericComponentModel", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetHideFromPortalFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getHideFromPortalFlag", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetIneligibleForAccountDiscountFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getIneligibleForAccountDiscountFlag", nil, &r.Options, &resp)
	return
}

// Retrieve DEPRECATED. An item's inventory status per datacenter.
func (r Product_Item) GetInventory() (resp []datatypes.Product_Package_Inventory, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getInventory", nil, &r.Options, &resp)
	return
}

// Retrieve Flag to indicate the server product is engineered for a multi-server solution. (Deprecated)
func (r Product_Item) GetIsEngineeredServerProduct() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getIsEngineeredServerProduct", nil, &r.Options, &resp)
	return
}

// Retrieve An item's primary item category.
func (r Product_Item) GetItemCategory() (resp datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getItemCategory", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetLocalDiskFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getLocalDiskFlag", nil, &r.Options, &resp)
	return
}

// Retrieve An item's location conflicts. For example, Dual Path network functionality cannot be ordered in WDC and as such is a conflict.
func (r Product_Item) GetLocationConflicts() (resp []datatypes.Product_Item_Resource_Conflict, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getLocationConflicts", nil, &r.Options, &resp)
	return
}

// Retrieve Indicates whether an item is a M.2 disk controller.
func (r Product_Item) GetM2ControllerFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getM2ControllerFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Indicates whether an item is a M.2 drive.
func (r Product_Item) GetM2DriveFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getM2DriveFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The minimum number of bays that support NVMe SSDs.
func (r Product_Item) GetMinimumNvmeBays() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getMinimumNvmeBays", nil, &r.Options, &resp)
	return
}

// Retrieve Indicates whether an item is a NVMe SSD.
func (r Product_Item) GetNvmeDiskFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getNvmeDiskFlag", nil, &r.Options, &resp)
	return
}

// Product Items are the products SoftLayer sells. Items have many properties that help describe what each is for. Each product items holds within it a description, tax rate information, status, and upgrade downgrade path information.
func (r Product_Item) GetObject() (resp datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetObjectStorageClusterGeolocationType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getObjectStorageClusterGeolocationType", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetObjectStorageItemFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getObjectStorageItemFlag", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetObjectStorageServiceClass() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getObjectStorageServiceClass", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of all the SoftLayer_Product_Package(s) in which this item exists.
func (r Product_Item) GetPackages() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPackages", nil, &r.Options, &resp)
	return
}

// Retrieve Indicates whether an item is a PCIe drive.
func (r Product_Item) GetPcieDriveFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPcieDriveFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The number of cores that a processor has.
func (r Product_Item) GetPhysicalCoreCapacity() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPhysicalCoreCapacity", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item) GetPresaleEvents() (resp []datatypes.Sales_Presale_Event, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPresaleEvents", nil, &r.Options, &resp)
	return
}

// Retrieve A product item's prices.
func (r Product_Item) GetPrices() (resp []datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPrices", nil, &r.Options, &resp)
	return
}

// Retrieve The number of private network interfaces provided by a port_speed product.
func (r Product_Item) GetPrivateInterfaceCount() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPrivateInterfaceCount", nil, &r.Options, &resp)
	return
}

// Retrieve The number of public network interfaces provided by a port_speed product.
func (r Product_Item) GetPublicInterfaceCount() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getPublicInterfaceCount", nil, &r.Options, &resp)
	return
}

// Retrieve If an item must be ordered with another item, it will have a requirement item here.
func (r Product_Item) GetRequirements() (resp []datatypes.Product_Item_Requirement, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getRequirements", nil, &r.Options, &resp)
	return
}

// Retrieve An item's rules. This includes the requirements and conflicts to resources that an item has.
func (r Product_Item) GetRules() (resp []datatypes.Product_Item_Rule, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getRules", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer_Software_Description tied to this item. This will only be populated for software items.
func (r Product_Item) GetSoftwareDescription() (resp datatypes.Software_Description, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getSoftwareDescription", nil, &r.Options, &resp)
	return
}

// Retrieve The total number of cores for a speed select server product.
func (r Product_Item) GetSpeedSelectServerCoreCount() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getSpeedSelectServerCoreCount", nil, &r.Options, &resp)
	return
}

// Retrieve Indicates a speed select server item.
func (r Product_Item) GetSpeedSelectServerFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getSpeedSelectServerFlag", nil, &r.Options, &resp)
	return
}

// Retrieve An item's tax category, if applicable.
func (r Product_Item) GetTaxCategory() (resp datatypes.Product_Item_Tax_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getTaxCategory", nil, &r.Options, &resp)
	return
}

// Retrieve Third-party policy assignments for this product.
func (r Product_Item) GetThirdPartyPolicyAssignments() (resp []datatypes.Product_Item_Policy_Assignment, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getThirdPartyPolicyAssignments", nil, &r.Options, &resp)
	return
}

// Retrieve The 3rd party vendor for a support subscription item. (Deprecated)
func (r Product_Item) GetThirdPartySupportVendor() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getThirdPartySupportVendor", nil, &r.Options, &resp)
	return
}

// Retrieve The total number of physical processing cores (excluding virtual cores / hyperthreads) for this server.
func (r Product_Item) GetTotalPhysicalCoreCapacity() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getTotalPhysicalCoreCapacity", nil, &r.Options, &resp)
	return
}

// Retrieve Shows the total number of cores. This is deprecated. Use [[SoftLayer_Product_Item/getCapacity|getCapacity]] for guest_core products and [[SoftLayer_Product_Item/getTotalPhysicalCoreCapacity|getTotalPhysicalCoreCapacity]] for server products
func (r Product_Item) GetTotalPhysicalCoreCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getTotalPhysicalCoreCount", nil, &r.Options, &resp)
	return
}

// Retrieve The total number of processors for this server.
func (r Product_Item) GetTotalProcessorCapacity() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getTotalProcessorCapacity", nil, &r.Options, &resp)
	return
}

// Retrieve Some product items have an upgrade path. This is the next product item in the upgrade path.
func (r Product_Item) GetUpgradeItem() (resp datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getUpgradeItem", nil, &r.Options, &resp)
	return
}

// Retrieve Some product items have an upgrade path. These are those upgrade product items.
func (r Product_Item) GetUpgradeItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item", "getUpgradeItems", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Product_Item_Category data type contains general category information for prices.
type Product_Item_Category struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductItemCategoryService returns an instance of the Product_Item_Category SoftLayer service
func GetProductItemCategoryService(sess session.SLSession) Product_Item_Category {
	return Product_Item_Category{Session: sess}
}

func (r Product_Item_Category) Id(id int) Product_Item_Category {
	r.Options.Id = &id
	return r
}

func (r Product_Item_Category) Mask(mask string) Product_Item_Category {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Item_Category) Filter(filter string) Product_Item_Category {
	r.Options.Filter = filter
	return r
}

func (r Product_Item_Category) Limit(limit int) Product_Item_Category {
	r.Options.Limit = &limit
	return r
}

func (r Product_Item_Category) Offset(offset int) Product_Item_Category {
	r.Options.Offset = &offset
	return r
}

// Returns a list of of active Items in the "Additional Services" package with their active prices for a given product item category and sorts them by price.
func (r Product_Item_Category) GetAdditionalProductsForCategory() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getAdditionalProductsForCategory", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Category) GetBandwidthCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getBandwidthCategories", nil, &r.Options, &resp)
	return
}

// Retrieve The billing items associated with an account that share a category code with an item category's category code.
func (r Product_Item_Category) GetBillingItems() (resp []datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getBillingItems", nil, &r.Options, &resp)
	return
}

// This method returns a collection of computing categories. These categories are also top level items in a service offering.
func (r Product_Item_Category) GetComputingCategories(resetCache *bool) (resp []datatypes.Product_Item_Category, err error) {
	params := []interface{}{
		resetCache,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getComputingCategories", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Category) GetCustomUsageRatesCategories(resetCache *bool) (resp []datatypes.Product_Item_Category, err error) {
	params := []interface{}{
		resetCache,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getCustomUsageRatesCategories", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Category) GetExternalResourceCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getExternalResourceCategories", nil, &r.Options, &resp)
	return
}

// Retrieve This invoice item's "item category group".
func (r Product_Item_Category) GetGroup() (resp datatypes.Product_Item_Category_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getGroup", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of service offering category groups. Each group contains a collection of items associated with this category.
func (r Product_Item_Category) GetGroups() (resp []datatypes.Product_Package_Item_Category_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getGroups", nil, &r.Options, &resp)
	return
}

// Each product item price must be tied to a category for it to be sold. These categories describe how a particular product item is sold. For example, the 250GB hard drive can be sold as disk0, disk1, ... disk11. There are different prices for this product item depending on which category it is. This keeps down the number of products in total.
func (r Product_Item_Category) GetObject() (resp datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Category) GetObjectStorageCategories(resetCache *bool) (resp []datatypes.Product_Item_Category, err error) {
	params := []interface{}{
		resetCache,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getObjectStorageCategories", params, &r.Options, &resp)
	return
}

// Retrieve Any unique options associated with an item category.
func (r Product_Item_Category) GetOrderOptions() (resp []datatypes.Product_Item_Category_Order_Option_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getOrderOptions", nil, &r.Options, &resp)
	return
}

// Retrieve A list of configuration available in this category.'
func (r Product_Item_Category) GetPackageConfigurations() (resp []datatypes.Product_Package_Order_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getPackageConfigurations", nil, &r.Options, &resp)
	return
}

// Retrieve A list of preset configurations this category is used in.'
func (r Product_Item_Category) GetPresetConfigurations() (resp []datatypes.Product_Package_Preset_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getPresetConfigurations", nil, &r.Options, &resp)
	return
}

// Retrieve The question references that are associated with an item category.
func (r Product_Item_Category) GetQuestionReferences() (resp []datatypes.Product_Item_Category_Question_Xref, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getQuestionReferences", nil, &r.Options, &resp)
	return
}

// Retrieve The questions that are associated with an item category.
func (r Product_Item_Category) GetQuestions() (resp []datatypes.Product_Item_Category_Question, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getQuestions", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Category) GetSoftwareCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getSoftwareCategories", nil, &r.Options, &resp)
	return
}

// This method returns a list of subnet categories.
func (r Product_Item_Category) GetSubnetCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getSubnetCategories", nil, &r.Options, &resp)
	return
}

// This method returns a collection of computing categories. These categories are also top level items in a service offering.
func (r Product_Item_Category) GetTopLevelCategories(resetCache *bool) (resp []datatypes.Product_Item_Category, err error) {
	params := []interface{}{
		resetCache,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getTopLevelCategories", params, &r.Options, &resp)
	return
}

// This method returns service product categories that can be canceled via API.  You can use these categories to find the billing items you wish to cancel.
func (r Product_Item_Category) GetValidCancelableServiceItemCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getValidCancelableServiceItemCategories", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Category) GetVlanCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category", "getVlanCategories", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Product_Item_Category_Group data type contains general category group information.
type Product_Item_Category_Group struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductItemCategoryGroupService returns an instance of the Product_Item_Category_Group SoftLayer service
func GetProductItemCategoryGroupService(sess session.SLSession) Product_Item_Category_Group {
	return Product_Item_Category_Group{Session: sess}
}

func (r Product_Item_Category_Group) Id(id int) Product_Item_Category_Group {
	r.Options.Id = &id
	return r
}

func (r Product_Item_Category_Group) Mask(mask string) Product_Item_Category_Group {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Item_Category_Group) Filter(filter string) Product_Item_Category_Group {
	r.Options.Filter = filter
	return r
}

func (r Product_Item_Category_Group) Limit(limit int) Product_Item_Category_Group {
	r.Options.Limit = &limit
	return r
}

func (r Product_Item_Category_Group) Offset(offset int) Product_Item_Category_Group {
	r.Options.Offset = &offset
	return r
}

// Each product item category must be tied to a category group. These category groups describe how a particular product item category is categorized. For example, the disk0, disk1, ... disk11 can be categorized as Server and Attached Services. There are different groups for each of this product item category depending on the function of the item product in the subject category.
func (r Product_Item_Category_Group) GetObject() (resp datatypes.Product_Item_Category_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Category_Group", "getObject", nil, &r.Options, &resp)
	return
}

// Represents the assignment of a policy to a product. The existence of a record means that the associated product is subject to the terms defined in the document content of the policy.
type Product_Item_Policy_Assignment struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductItemPolicyAssignmentService returns an instance of the Product_Item_Policy_Assignment SoftLayer service
func GetProductItemPolicyAssignmentService(sess session.SLSession) Product_Item_Policy_Assignment {
	return Product_Item_Policy_Assignment{Session: sess}
}

func (r Product_Item_Policy_Assignment) Id(id int) Product_Item_Policy_Assignment {
	r.Options.Id = &id
	return r
}

func (r Product_Item_Policy_Assignment) Mask(mask string) Product_Item_Policy_Assignment {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Item_Policy_Assignment) Filter(filter string) Product_Item_Policy_Assignment {
	r.Options.Filter = filter
	return r
}

func (r Product_Item_Policy_Assignment) Limit(limit int) Product_Item_Policy_Assignment {
	r.Options.Limit = &limit
	return r
}

func (r Product_Item_Policy_Assignment) Offset(offset int) Product_Item_Policy_Assignment {
	r.Options.Offset = &offset
	return r
}

// Register the acceptance of the associated policy to product assignment, and link the created record to a Ticket.
func (r Product_Item_Policy_Assignment) AcceptFromTicket(ticketId *int) (resp bool, err error) {
	params := []interface{}{
		ticketId,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Item_Policy_Assignment", "acceptFromTicket", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Policy_Assignment) GetObject() (resp datatypes.Product_Item_Policy_Assignment, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Policy_Assignment", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve the binary contents of the associated PDF policy document.
func (r Product_Item_Policy_Assignment) GetPolicyDocumentContents() (resp []byte, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Policy_Assignment", "getPolicyDocumentContents", nil, &r.Options, &resp)
	return
}

// Retrieve The name of the assigned policy.
func (r Product_Item_Policy_Assignment) GetPolicyName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Policy_Assignment", "getPolicyName", nil, &r.Options, &resp)
	return
}

// Retrieve The [[SoftLayer_Product_Item]] for this policy assignment.
func (r Product_Item_Policy_Assignment) GetProduct() (resp datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Policy_Assignment", "getProduct", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Product_Item_Price data type contains general information relating to a single SoftLayer product item price. You can find out what packages each price is in as well as which category under which this price is sold. All prices are returned in floating point values measured in US Dollars ($USD).
type Product_Item_Price struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductItemPriceService returns an instance of the Product_Item_Price SoftLayer service
func GetProductItemPriceService(sess session.SLSession) Product_Item_Price {
	return Product_Item_Price{Session: sess}
}

func (r Product_Item_Price) Id(id int) Product_Item_Price {
	r.Options.Id = &id
	return r
}

func (r Product_Item_Price) Mask(mask string) Product_Item_Price {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Item_Price) Filter(filter string) Product_Item_Price {
	r.Options.Filter = filter
	return r
}

func (r Product_Item_Price) Limit(limit int) Product_Item_Price {
	r.Options.Limit = &limit
	return r
}

func (r Product_Item_Price) Offset(offset int) Product_Item_Price {
	r.Options.Offset = &offset
	return r
}

// Retrieve The account that the item price is restricted to.
func (r Product_Item_Price) GetAccountRestrictions() (resp []datatypes.Product_Item_Price_Account_Restriction, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getAccountRestrictions", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item_Price) GetAttributes() (resp []datatypes.Product_Item_Price_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve Signifies pricing that is only available on a bare metal reserved capacity order.
func (r Product_Item_Price) GetBareMetalReservedCapacityFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getBareMetalReservedCapacityFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the price is for Big Data OS/Journal disks only. (Deprecated)
func (r Product_Item_Price) GetBigDataOsJournalDiskFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getBigDataOsJournalDiskFlag", nil, &r.Options, &resp)
	return
}

// Retrieve cross reference for bundles
func (r Product_Item_Price) GetBundleReferences() (resp []datatypes.Product_Item_Bundles, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getBundleReferences", nil, &r.Options, &resp)
	return
}

// Retrieve The maximum capacity value for which this price is suitable.
func (r Product_Item_Price) GetCapacityRestrictionMaximum() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getCapacityRestrictionMaximum", nil, &r.Options, &resp)
	return
}

// Retrieve The minimum capacity value for which this price is suitable.
func (r Product_Item_Price) GetCapacityRestrictionMinimum() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getCapacityRestrictionMinimum", nil, &r.Options, &resp)
	return
}

// Retrieve The type of capacity restriction by which this price must abide.
func (r Product_Item_Price) GetCapacityRestrictionType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getCapacityRestrictionType", nil, &r.Options, &resp)
	return
}

// Retrieve All categories which this item is a member.
func (r Product_Item_Price) GetCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getCategories", nil, &r.Options, &resp)
	return
}

// Retrieve Signifies pricing that is only available on a dedicated host virtual server order.
func (r Product_Item_Price) GetDedicatedHostInstanceFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getDedicatedHostInstanceFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether this price defines a software license for its product item.
func (r Product_Item_Price) GetDefinedSoftwareLicenseFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getDefinedSoftwareLicenseFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Eligibility strategy to assess if a customer can order using this price.
func (r Product_Item_Price) GetEligibilityStrategy() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getEligibilityStrategy", nil, &r.Options, &resp)
	return
}

// Retrieve The product item a price is tied to.
func (r Product_Item_Price) GetItem() (resp datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getItem", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Price) GetObject() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item_Price) GetOrderPremiums() (resp []datatypes.Product_Item_Price_Premium, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getOrderPremiums", nil, &r.Options, &resp)
	return
}

// Retrieve cross reference for packages
func (r Product_Item_Price) GetPackageReferences() (resp []datatypes.Product_Package_Item_Prices, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getPackageReferences", nil, &r.Options, &resp)
	return
}

// Retrieve A price's packages under which this item is sold.
func (r Product_Item_Price) GetPackages() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getPackages", nil, &r.Options, &resp)
	return
}

// Retrieve A list of preset configurations this price is used in.'
func (r Product_Item_Price) GetPresetConfigurations() (resp []datatypes.Product_Package_Preset_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getPresetConfigurations", nil, &r.Options, &resp)
	return
}

// Retrieve The type keyname of this price which can be STANDARD, TIERED, or TERM.
func (r Product_Item_Price) GetPriceType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getPriceType", nil, &r.Options, &resp)
	return
}

// Retrieve The pricing location group that this price is applicable for. Prices that have a pricing location group will only be available for ordering with the locations specified on the location group.
func (r Product_Item_Price) GetPricingLocationGroup() (resp datatypes.Location_Group_Pricing, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getPricingLocationGroup", nil, &r.Options, &resp)
	return
}

// Retrieve The number of server cores required to order this item. This is deprecated. Use [[SoftLayer_Product_Item_Price/getCapacityRestrictionMinimum|getCapacityRestrictionMinimum]] and [[SoftLayer_Product_Item_Price/getCapacityRestrictionMaximum|getCapacityRestrictionMaximum]]
func (r Product_Item_Price) GetRequiredCoreCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getRequiredCoreCount", nil, &r.Options, &resp)
	return
}

// Retrieve Signifies pricing that is only available on a reserved capacity virtual server order.
func (r Product_Item_Price) GetReservedCapacityInstanceFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getReservedCapacityInstanceFlag", nil, &r.Options, &resp)
	return
}

// Returns a collection of rate-based [[SoftLayer_Product_Item_Price]] objects associated with the [[SoftLayer_Product_Item]] objects and the [[SoftLayer_Location]] specified. The location is required to get the appropriate rate-based prices because the usage rates may vary from datacenter to datacenter.
func (r Product_Item_Price) GetUsageRatePrices(location *datatypes.Location, items []datatypes.Product_Item) (resp []datatypes.Product_Item_Price, err error) {
	params := []interface{}{
		location,
		items,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price", "getUsageRatePrices", params, &r.Options, &resp)
	return
}

// no documentation yet
type Product_Item_Price_Premium struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductItemPricePremiumService returns an instance of the Product_Item_Price_Premium SoftLayer service
func GetProductItemPricePremiumService(sess session.SLSession) Product_Item_Price_Premium {
	return Product_Item_Price_Premium{Session: sess}
}

func (r Product_Item_Price_Premium) Id(id int) Product_Item_Price_Premium {
	r.Options.Id = &id
	return r
}

func (r Product_Item_Price_Premium) Mask(mask string) Product_Item_Price_Premium {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Item_Price_Premium) Filter(filter string) Product_Item_Price_Premium {
	r.Options.Filter = filter
	return r
}

func (r Product_Item_Price_Premium) Limit(limit int) Product_Item_Price_Premium {
	r.Options.Limit = &limit
	return r
}

func (r Product_Item_Price_Premium) Offset(offset int) Product_Item_Price_Premium {
	r.Options.Offset = &offset
	return r
}

// Retrieve
func (r Product_Item_Price_Premium) GetItemPrice() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price_Premium", "getItemPrice", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item_Price_Premium) GetLocation() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price_Premium", "getLocation", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Item_Price_Premium) GetObject() (resp datatypes.Product_Item_Price_Premium, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price_Premium", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Item_Price_Premium) GetPackage() (resp datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Item_Price_Premium", "getPackage", nil, &r.Options, &resp)
	return
}

// no documentation yet
type Product_Order struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductOrderService returns an instance of the Product_Order SoftLayer service
func GetProductOrderService(sess session.SLSession) Product_Order {
	return Product_Order{Session: sess}
}

func (r Product_Order) Id(id int) Product_Order {
	r.Options.Id = &id
	return r
}

func (r Product_Order) Mask(mask string) Product_Order {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Order) Filter(filter string) Product_Order {
	r.Options.Filter = filter
	return r
}

func (r Product_Order) Limit(limit int) Product_Order {
	r.Options.Limit = &limit
	return r
}

func (r Product_Order) Offset(offset int) Product_Order {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Product_Order) CheckItemAvailability(itemPrices []datatypes.Product_Item_Price, accountId *int, availabilityTypeKeyNames []string) (resp bool, err error) {
	params := []interface{}{
		itemPrices,
		accountId,
		availabilityTypeKeyNames,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "checkItemAvailability", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Order) CheckItemAvailabilityForImageTemplate(imageTemplateId *int, accountId *int, packageId *int, availabilityTypeKeyNames []string) (resp bool, err error) {
	params := []interface{}{
		imageTemplateId,
		accountId,
		packageId,
		availabilityTypeKeyNames,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "checkItemAvailabilityForImageTemplate", params, &r.Options, &resp)
	return
}

// Check order items for conflicts
func (r Product_Order) CheckItemConflicts(itemPrices []datatypes.Product_Item_Price) (resp bool, err error) {
	params := []interface{}{
		itemPrices,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "checkItemConflicts", params, &r.Options, &resp)
	return
}

// This method simply returns a receipt for a previously finalized payment authorization from PayPal. The response matches the response returned from placeOrder when the order was originally placed with PayPal as the payment type.
func (r Product_Order) GetExternalPaymentAuthorizationReceipt(token *string, payerId *string) (resp datatypes.Container_Product_Order_Receipt, err error) {
	params := []interface{}{
		token,
		payerId,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "getExternalPaymentAuthorizationReceipt", params, &r.Options, &resp)
	return
}

// This method is deprecated and always returns nothing.
// Deprecated: This function has been marked as deprecated.
func (r Product_Order) GetNetworks(locationId *int, packageId *int, accountId *int) (resp []datatypes.Container_Product_Order_Network, err error) {
	params := []interface{}{
		locationId,
		packageId,
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "getNetworks", params, &r.Options, &resp)
	return
}

// When the account is on an external reseller brand, this service will provide a SoftLayer_Product_Order with the the pricing adjusted by the external reseller.
func (r Product_Order) GetResellerOrder(orderContainer *datatypes.Container_Product_Order) (resp datatypes.Container_Product_Order, err error) {
	params := []interface{}{
		orderContainer,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "getResellerOrder", params, &r.Options, &resp)
	return
}

// Sometimes taxes cannot be calculated immediately, so we start the calculations and let them run in the background. This method will return the current progress and information related to a specific tax calculation, which allows real-time progress updates on tax calculations.
func (r Product_Order) GetTaxCalculationResult(orderHash *string) (resp datatypes.Container_Tax_Cache, err error) {
	params := []interface{}{
		orderHash,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "getTaxCalculationResult", params, &r.Options, &resp)
	return
}

// Return collections of public and private VLANs that are available during ordering. If a location ID is provided, the resulting VLANs will be limited to that location. If the Virtual Server package id (46) is provided, the VLANs will be narrowed down to those locations that contain routers with the VIRTUAL_IMAGE_STORE data attribute.
//
// For the selectedItems parameter, this is a comma-separated string of category codes and item values. For example:
//
// - `port_speed=10,guest_disk0=LOCAL_DISK`
//
// - `port_speed=100,disk0=SAN_DISK`
//
// - `port_speed=100,private_network_only=1,guest_disk0=LOCAL_DISK`
//
// This parameter is used to narrow the available results down even further. It's not necessary when selecting a VLAN, but it will help avoid errors when attempting to place an order. The only acceptable category codes are:
//
// - `port_speed`
//
// - A disk category, such as `guest_disk0` or `disk0`, with values of either `LOCAL_DISK` or `SAN_DISK`
//
// - `private_network_only`
//
// - `dual_path_network`
//
// For most customers, it's sufficient to only provide the first 2 parameters.
func (r Product_Order) GetVlans(locationId *int, packageId *int, selectedItems *string, vlanIds []int, subnetIds []int, accountId *int, orderContainer *datatypes.Container_Product_Order, hardwareFirewallOrderedFlag *bool) (resp datatypes.Container_Product_Order_Network_Vlans, err error) {
	params := []interface{}{
		locationId,
		packageId,
		selectedItems,
		vlanIds,
		subnetIds,
		accountId,
		orderContainer,
		hardwareFirewallOrderedFlag,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "getVlans", params, &r.Options, &resp)
	return
}

// Use this method to place bare metal server, virtual server and additional service orders with SoftLayer.
// Upon success, your credit card or PayPal account will incur charges for the monthly order total
// (or prorated value if ordered mid billing cycle). If all products on the order are only billed hourly,
// you will be charged on your billing anniversary date, which occurs monthly on the day you ordered your first
// service with SoftLayer. For new customers, you are required to provide billing information when you place an order.
// For existing customers, the credit card on file will be charged. If you're a PayPal customer, a URL will be
// returned from the call to [[SoftLayer_Product_Order/placeOrder]] which is to be used to finish the authorization
// process. This authorization tells PayPal that you indeed want to place an order with SoftLayer.
// From PayPal's web site, you will be redirected back to SoftLayer for your order receipt.
//
// When an order is placed, your order will be in a "pending approval" state. When all internal checks pass,
// your order will be automatically approved. For orders that may need extra attention, a Sales representative
// will review the order and contact you if necessary. Once the order is approved, your server or service will
// be provisioned and available to you shortly thereafter. Depending on the type of server or service ordered,
// provisioning times will vary.
//
// ## Order Containers
//
// When placing API orders, it's important to order your server and services on the appropriate
// [[SoftLayer_Container_Product_Order]]. Failing to provide the correct container may delay your server or service
// from being provisioned in a timely manner. Some common order containers are included below.
//
// **Note:** `SoftLayer_Container_Product_Order_` has been removed from the containers in the table below for readability.
//
// | Product | Order Container | Package Type |
// | ------- | --------------- | ------------ |
// | Bare metal server by CPU | [[SoftLayer_Container_Product_Order_Hardware_Server]] | BARE_METAL_CPU |
// | Bare metal server by core | [[SoftLayer_Container_Product_Order_Hardware_Server]] | BARE_METAL_CORE |
// | Virtual server | [[SoftLayer_Container_Product_Order_Virtual_Guest]] | VIRTUAL_SERVER_INSTANCE |
// | Local & dedicated load balancers | [[SoftLayer_Container_Product_Order_Network_LoadBalancer]] | ADDITIONAL_SERVICES_LOAD_BALANCER |
// | Content delivery network | [[SoftLayer_Container_Product_Order_Network_ContentDelivery_Account]] | ADDITIONAL_SERVICES_CDN |
// | Content delivery network Addon | [[SoftLayer_Container_Product_Order_Network_ContentDelivery_Account_Addon]] | ADDITIONAL_SERVICES_CDN_ADDON |
// | Hardware & software firewalls | [[SoftLayer_Container_Product_Order_Network_Protection_Firewall]] | ADDITIONAL_SERVICES_FIREWALL |
// | Dedicated firewall | [[SoftLayer_Container_Product_Order_Network_Protection_Firewall_Dedicated]] | ADDITIONAL_SERVICES_FIREWALL |
// | Object storage | [[SoftLayer_Container_Product_Order_Network_Storage_Object]] | ADDITIONAL_SERVICES_OBJECT_STORAGE |
// | Object storage (hub) | [[SoftLayer_Container_Product_Order_Network_Storage_Hub]] | ADDITIONAL_SERVICES_OBJECT_STORAGE |
// | Network attached storage | [[SoftLayer_Container_Product_Order_Network_Storage_Nas]] | ADDITIONAL_SERVICES_NETWORK_ATTACHED_STORAGE |
// | Iscsi storage | [[SoftLayer_Container_Product_Order_Network_Storage_Iscsi]] | ADDITIONAL_SERVICES_ISCSI_STORAGE |
// | Evault | [[SoftLayer_Container_Product_Order_Network_Storage_Backup_Evault_Vault]] | ADDITIONAL_SERVICES |
// | Evault Plugin | [[SoftLayer_Container_Product_Order_Network_Storage_Backup_Evault_Plugin]] | ADDITIONAL_SERVICES |
// | Application delivery appliance | [[SoftLayer_Container_Product_Order_Network_Application_Delivery_Controller]] | ADDITIONAL_SERVICES_APPLICATION_DELIVERY_APPLIANCE |
// | Network subnet | [[SoftLayer_Container_Product_Order_Network_Subnet]] | ADDITIONAL_SERVICES |
// | Global IPv4 | [[SoftLayer_Container_Product_Order_Network_Subnet]] | ADDITIONAL_SERVICES_GLOBAL_IP_ADDRESSES |
// | Global IPv6 | [[SoftLayer_Container_Product_Order_Network_Subnet]] | ADDITIONAL_SERVICES_GLOBAL_IP_ADDRESSES |
// | Network VLAN | [[SoftLayer_Container_Product_Order_Network_Vlan]] | ADDITIONAL_SERVICES_NETWORK_VLAN |
// | Portable storage | [[SoftLayer_Container_Product_Order_Virtual_Disk_Image]] | ADDITIONAL_SERVICES_PORTABLE_STORAGE |
// | SSL certificate | [[SoftLayer_Container_Product_Order_Security_Certificate]] | ADDITIONAL_SERVICES_SSL_CERTIFICATE |
// | External authentication | [[SoftLayer_Container_Product_Order_User_Customer_External_Binding]] | ADDITIONAL_SERVICES |
// | Dedicated Host | [[SoftLayer_Container_Product_Order_Virtual_DedicatedHost]] | DEDICATED_HOST |
//
// ## Server example
//
// This example includes a single bare metal server being ordered with monthly billing.
//
// **Warning:** the price ids provided below may be outdated or unavailable, so you will need to determine the
//
// available prices from the bare metal server [[SoftLayer_Product_Package/getAllObjects]], which have a
// [[SoftLayer_Product_Package_Type]] of `BARE_METAL_CPU` or `BARE_METAL_CORE`. You can get a full list of
// package types with [[SoftLayer_Product_Package_Type/getAllObjects]].
//
// ### Bare Metal Ordering
//
// ```xml
// <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://api.service.softlayer.com/soap/v3/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
//
//	<SOAP-ENV:Header>
//	  <ns1:authenticate>
//	    <username>your username</username>
//	    <apiKey>your api key</apiKey>
//	  </ns1:authenticate>
//	</SOAP-ENV:Header>
//	<SOAP-ENV:Body>
//	  <ns1:placeOrder>
//	    <orderData xsi:type="ns1:SoftLayer_Container_Product_Order_Hardware_Server">
//	      <hardware SOAP-ENC:arrayType="ns1:SoftLayer_Hardware[1]" xsi:type="ns1:SoftLayer_HardwareArray">
//	        <item xsi:type="ns1:SoftLayer_Hardware">
//	          <domain xsi:type="xsd:string">example.com</domain>
//	          <hostname xsi:type="xsd:string">server1</hostname>
//	        </item>
//	      </hardware>
//	      <location xsi:type="xsd:string">138124</location>
//	      <packageId xsi:type="xsd:int">142</packageId>
//	      <prices SOAP-ENC:arrayType="ns1:SoftLayer_Product_Item_Price[14]" xsi:type="ns1:SoftLayer_Product_Item_PriceArray">
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">58</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">22337</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">21189</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">876</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">57</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">55</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">21190</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">36381</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">21</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">22013</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">906</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">420</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">418</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">342</id>
//	        </item>
//	      </prices>
//	      <useHourlyPricing xsi:type="xsd:boolean">false</useHourlyPricing>
//	    </orderData>
//	    <saveAsQuote xsi:nil="true" />
//	  </ns1:placeOrder>
//	</SOAP-ENV:Body>
//
// </SOAP-ENV:Envelope>
// ```
//
// ## Virtual server example
//
// This example includes 2 identical virtual servers (except for hostname) being ordered for hourly billing.
// It includes an optional image template id and VLAN data specified on the virtualGuest objects -
// `primaryBackendNetworkComponent` and `primaryNetworkComponent`.
//
// **Warning:** the price ids provided below may be outdated or unavailable, so you will need to determine the
//
// available prices from the virtual server package with [[SoftLayer_Product_Package/getAllObjects]],
// which has a [[SoftLayer_Product_Package_Type]] of `VIRTUAL_SERVER_INSTANCE`.
//
// #### Virtual Ordering
//
// ```xml
// <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://api.service.softlayer.com/soap/v3/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
//
//	<SOAP-ENV:Header>
//	  <ns1:authenticate>
//	    <username>your username</username>
//	    <apiKey>your api key</apiKey>
//	  </ns1:authenticate>
//	</SOAP-ENV:Header>
//	<SOAP-ENV:Body>
//	  <ns1:placeOrder>
//	    <orderData xsi:type="ns1:SoftLayer_Container_Product_Order_Virtual_Guest">
//	      <imageTemplateId xsi:type="xsd:int">13251</imageTemplateId>
//	      <location xsi:type="xsd:string">37473</location>
//	      <packageId xsi:type="xsd:int">46</packageId>
//	      <prices SOAP-ENC:arrayType="ns1:SoftLayer_Product_Item_Price[13]" xsi:type="ns1:SoftLayer_Product_Item_PriceArray">
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">2159</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">55</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">13754</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">1641</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">905</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">1800</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">58</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">21</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">1645</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">272</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">57</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">418</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">420</id>
//	        </item>
//	      </prices>
//	      <quantity xsi:type="xsd:int">2</quantity>
//	      <useHourlyPricing xsi:type="xsd:boolean">true</useHourlyPricing>
//	      <virtualGuests SOAP-ENC:arrayType="ns1:SoftLayer_Virtual_Guest[1]" xsi:type="ns1:SoftLayer_Virtual_GuestArray">
//	        <item xsi:type="ns1:SoftLayer_Virtual_Guest">
//	          <domain xsi:type="xsd:string">example.com</domain>
//	          <hostname xsi:type="xsd:string">server1</hostname>
//	          <primaryBackendNetworkComponent xsi:type="ns1:SoftLayer_Virtual_Guest_Network_Component">
//	            <networkVlan xsi:type="ns1:SoftLayer_Network_Vlan">
//	              <id xsi:type="xsd:int">12345</id>
//	            </networkVlan>
//	          </primaryBackendNetworkComponent>
//	          <primaryNetworkComponent xsi:type="ns1:SoftLayer_Virtual_Guest_Network_Component">
//	            <networkVlan xsi:type="ns1:SoftLayer_Network_Vlan">
//	              <id xsi:type="xsd:int">67890</id>
//	            </networkVlan>
//	          </primaryNetworkComponent>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Virtual_Guest">
//	          <domain xsi:type="xsd:string">example.com</domain>
//	          <hostname xsi:type="xsd:string">server2</hostname>
//	          <primaryBackendNetworkComponent xsi:type="ns1:SoftLayer_Virtual_Guest_Network_Component">
//	            <networkVlan xsi:type="ns1:SoftLayer_Network_Vlan">
//	              <id xsi:type="xsd:int">12345</id>
//	            </networkVlan>
//	          </primaryBackendNetworkComponent>
//	          <primaryNetworkComponent xsi:type="ns1:SoftLayer_Virtual_Guest_Network_Component">
//	            <networkVlan xsi:type="ns1:SoftLayer_Network_Vlan">
//	              <id xsi:type="xsd:int">67890</id>
//	            </networkVlan>
//	          </primaryNetworkComponent>
//	        </item>
//	      </virtualGuests>
//	    </orderData>
//	    <saveAsQuote xsi:nil="true" />
//	  </ns1:placeOrder>
//	</SOAP-ENV:Body>
//
// </SOAP-ENV:Envelope>
// ```
//
// ## VLAN example
//
// **Warning:** the price ids provided below may be outdated or unavailable, so you will need to determine the
//
// available prices from the additional services pacakge with [[SoftLayer_Product_Package/getAllObjects]],
// which has a [[SoftLayer_Product_Package_Type]] of `ADDITIONAL_SERVICES`.
// You can get a full list of [[SoftLayer_Product_Package_Type/getAllObjects|]] to find other available additional
// service packages.<br/><br/>
//
// ### VLAN Ordering
//
// ```xml
// <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://api.service.softlayer.com/soap/v3/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
//
//	<SOAP-ENV:Header>
//	  <ns1:authenticate>
//	    <username>your username</username>
//	    <apiKey>your api key</apiKey>
//	  </ns1:authenticate>
//	</SOAP-ENV:Header>
//	<SOAP-ENV:Body>
//	  <ns1:placeOrder>
//	    <orderData xsi:type="ns1:SoftLayer_Container_Product_Order_Network_Vlan">
//	      <location xsi:type="xsd:string">154820</location>
//	      <packageId xsi:type="xsd:int">0</packageId>
//	      <prices SOAP-ENC:arrayType="ns1:SoftLayer_Product_Item_Price[2]" xsi:type="ns1:SoftLayer_Product_Item_PriceArray">
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">2021</id>
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Product_Item_Price">
//	          <id xsi:type="xsd:int">2018</id>
//	        </item>
//	      </prices>
//	      <useHourlyPricing xsi:type="xsd:boolean">true</useHourlyPricing>
//	    </orderData>
//	    <saveAsQuote xsi:nil="true" />
//	  </ns1:placeOrder>
//	</SOAP-ENV:Body>
//
// </SOAP-ENV:Envelope>
// ```
//
// ## Multiple products example
//
// This example includes a combination of the above examples in a single order. Note that all the configuration
// options for each individual order container are the same as above, except now we encapsulate each one within
// the `orderContainers` property on the base [[SoftLayer_Container_Product_Order]].
//
// **Warning:** not all products are available to be ordered with other products. For example, since
//
// SSL certificates require validation from a 3rd party, the approval process may take days or even weeks,
// and this would not be acceptable when you need your hourly virtual server right now. To better accommodate
// customers, we restrict several products to be ordered individually.
//
// ### Bare metal server + virtual server + VLAN
//
// ```xml
// <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://api.service.softlayer.com/soap/v3/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
//
//	<SOAP-ENV:Header>
//	  <ns1:authenticate>
//	    <username>your username</username>
//	    <apiKey>your api key</apiKey>
//	  </ns1:authenticate>
//	</SOAP-ENV:Header>
//	<SOAP-ENV:Body>
//	  <ns1:placeOrder>
//	    <orderData xsi:type="ns1:SoftLayer_Container_Product_Order">
//	      <orderContainers SOAP-ENC:arrayType="ns1:SoftLayer_Container_Product_Order[3]" xsi:type="ns1:SoftLayer_Container_Product_OrderArray">
//	        <item xsi:type="ns1:SoftLayer_Container_Product_Order_Hardware_Server">
//	          ...
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Container_Product_Order_Virtual_Guest">
//	          ...
//	        </item>
//	        <item xsi:type="ns1:SoftLayer_Container_Product_Order_Network_Vlan">
//	          ...
//	        </item>
//	      </orderContainers>
//	    </orderData>
//	    <saveAsQuote xsi:nil="true" />
//	  </ns1:placeOrder>
//	</SOAP-ENV:Body>
//
// </SOAP-ENV:Envelope>
// ```
func (r Product_Order) PlaceOrder(orderData interface{}, saveAsQuote *bool) (resp datatypes.Container_Product_Order_Receipt, err error) {
	err = datatypes.SetComplexType(orderData)
	if err != nil {
		return
	}
	params := []interface{}{
		orderData,
		saveAsQuote,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "placeOrder", params, &r.Options, &resp)
	return
}

// Use this method for placing server quotes and additional services quotes. The same applies for this as with verifyOrder. Send in the SoftLayer_Container_Product_Order_Hardware_Server for server quotes. After placing the quote, you must go to this URL to finish the order process. After going to this URL, it will direct you back to a SoftLayer webpage that tells us you have finished the process. After this, it will go to sales for final approval.
func (r Product_Order) PlaceQuote(orderData interface{}) (resp datatypes.Container_Product_Order_Receipt, err error) {
	params := []interface{}{
		orderData,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "placeQuote", params, &r.Options, &resp)
	return
}

// This method simply finalizes an authorization from PayPal. It tells SoftLayer that the customer has completed the PayPal process. This is ONLY needed if you, the customer, have your own API into PayPal and wish to automate authorizations from PayPal and our system. For most, this method will not be needed. Once an order is placed using placeOrder() for PayPal customers, a URL is given back to the customer. In it is the token and PayerID. If you want to systematically pay with PayPal, do so then call this method with the token and PayerID.
func (r Product_Order) ProcessExternalPaymentAuthorization(token *string, payerId *string) (resp datatypes.Container_Product_Order, err error) {
	params := []interface{}{
		token,
		payerId,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "processExternalPaymentAuthorization", params, &r.Options, &resp)
	return
}

// Get list of items that are required with the item prices provided
func (r Product_Order) RequiredItems(itemPrices []datatypes.Product_Item_Price) (resp []datatypes.Product_Item, err error) {
	params := []interface{}{
		itemPrices,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "requiredItems", params, &r.Options, &resp)
	return
}

// This service is used to verify that an order meets all the necessary requirements to purchase a server, virtual server or service from SoftLayer. It will verify that the products requested do not conflict. For example, you cannot order a Windows firewall with a Linux operating system. It will also check to make sure you have provided all the products that are required for the [[SoftLayer_Product_Package_Order_Configuration]] associated with the [[SoftLayer_Product_Package]] on each of the [[SoftLayer_Container_Product_Order]] specified.<br/><br/>
//
// This service returns the same container that was provided, but with additional information that can be used for debugging or validation. It will also contain pricing information (prorated if applicable) for each of the products on the order. If an exception occurs during verification, a container with the <code>SoftLayer_Exception_Order</code> exception type will be specified in the result.<br/><br/>
//
// <code>verifyOrder</code> accepts the same [[SoftLayer_Container_Product_Order]] as <code>placeOrder</code>, so see [[SoftLayer_Product_Order/placeOrder]] for more details.
func (r Product_Order) VerifyOrder(orderData interface{}) (resp datatypes.Container_Product_Order, err error) {
	err = datatypes.SetComplexType(orderData)
	if err != nil {
		return
	}
	params := []interface{}{
		orderData,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Order", "verifyOrder", params, &r.Options, &resp)
	return
}

// The SoftLayer_Product_Package data type contains information about packages from which orders can be generated. Packages contain general information regarding what is in them, where they are currently sold, availability, and pricing.
type Product_Package struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductPackageService returns an instance of the Product_Package SoftLayer service
func GetProductPackageService(sess session.SLSession) Product_Package {
	return Product_Package{Session: sess}
}

func (r Product_Package) Id(id int) Product_Package {
	r.Options.Id = &id
	return r
}

func (r Product_Package) Mask(mask string) Product_Package {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Package) Filter(filter string) Product_Package {
	r.Options.Filter = filter
	return r
}

func (r Product_Package) Limit(limit int) Product_Package {
	r.Options.Limit = &limit
	return r
}

func (r Product_Package) Offset(offset int) Product_Package {
	r.Options.Offset = &offset
	return r
}

// Retrieve The preset configurations available only for the authenticated account and this package.
func (r Product_Package) GetAccountRestrictedActivePresets() (resp []datatypes.Product_Package_Preset, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAccountRestrictedActivePresets", nil, &r.Options, &resp)
	return
}

// Retrieve The results from this call are similar to [[SoftLayer_Product_Package/getCategories|getCategories]], but these ONLY include account-restricted prices. Not all accounts have restricted pricing.
func (r Product_Package) GetAccountRestrictedCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAccountRestrictedCategories", nil, &r.Options, &resp)
	return
}

// Retrieve The flag to indicate if there are any restricted prices in a package for the currently-active account.
func (r Product_Package) GetAccountRestrictedPricesFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAccountRestrictedPricesFlag", nil, &r.Options, &resp)
	return
}

// Return a list of Items in the package with their active prices.
func (r Product_Package) GetActiveItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActiveItems", nil, &r.Options, &resp)
	return
}

// <strong>This method is deprecated and should not be used in production code.</strong>
//
// This method will return the [[SoftLayer_Product_Package]] objects from which you can order a bare metal server, virtual server, service (such as CDN or Object Storage) or other software filtered by an attribute type associated with the package. Once you have the package you want to order from, you may query one of various endpoints from that package to get specific information about its products and pricing. See [[SoftLayer_Product_Package/getCategories|getCategories]] or [[SoftLayer_Product_Package/getItems|getItems]] for more information.
// Deprecated: This function has been marked as deprecated.
func (r Product_Package) GetActivePackagesByAttribute(attributeKeyName *string) (resp []datatypes.Product_Package, err error) {
	params := []interface{}{
		attributeKeyName,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActivePackagesByAttribute", params, &r.Options, &resp)
	return
}

// Retrieve The available preset configurations for this package.
func (r Product_Package) GetActivePresets() (resp []datatypes.Product_Package_Preset, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActivePresets", nil, &r.Options, &resp)
	return
}

// [DEPRECATED] This method pulls all the active private hosted cloud packages. This will give you a basic description of the packages that are currently active and from which you can order private hosted cloud configurations.
// Deprecated: This function has been marked as deprecated.
func (r Product_Package) GetActivePrivateHostedCloudPackages() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActivePrivateHostedCloudPackages", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of valid RAM items available for purchase in this package.
func (r Product_Package) GetActiveRamItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActiveRamItems", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of valid server items available for purchase in this package.
func (r Product_Package) GetActiveServerItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActiveServerItems", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of valid software items available for purchase in this package.
func (r Product_Package) GetActiveSoftwareItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActiveSoftwareItems", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of [[SoftLayer_Product_Item_Price]] objects for pay-as-you-go usage.
func (r Product_Package) GetActiveUsagePrices() (resp []datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActiveUsagePrices", nil, &r.Options, &resp)
	return
}

// This method returns a collection of active usage rate [[SoftLayer_Product_Item_Price]] objects for the current package and specified datacenter. Optionally you can retrieve the active usage rate prices for a particular [[SoftLayer_Product_Item_Category]] by specifying a category code as the first parameter. This information is useful so that you can see "pay as you go" rates (if any) for the current package, location and optionally category.
func (r Product_Package) GetActiveUsageRatePrices(locationId *int, categoryCode *string) (resp []datatypes.Product_Item_Price, err error) {
	params := []interface{}{
		locationId,
		categoryCode,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getActiveUsageRatePrices", params, &r.Options, &resp)
	return
}

// Retrieve This flag indicates that the package is an additional service.
func (r Product_Package) GetAdditionalServiceFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAdditionalServiceFlag", nil, &r.Options, &resp)
	return
}

// This method pulls all the active packages. This will give you a basic description of the packages that are currently active
func (r Product_Package) GetAllObjects() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package) GetAttributes() (resp []datatypes.Product_Package_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of valid locations for this package. (Deprecated - Use [[SoftLayer_Product_Package/getRegions|getRegions]])
func (r Product_Package) GetAvailableLocations() (resp []datatypes.Product_Package_Locations, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAvailableLocations", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package) GetAvailablePackagesForImageTemplate(imageTemplate *datatypes.Virtual_Guest_Block_Device_Template_Group) (resp []datatypes.Product_Package, err error) {
	params := []interface{}{
		imageTemplate,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAvailablePackagesForImageTemplate", params, &r.Options, &resp)
	return
}

// Retrieve The maximum number of available disk storage units associated with the servers in a package.
func (r Product_Package) GetAvailableStorageUnits() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getAvailableStorageUnits", nil, &r.Options, &resp)
	return
}

// Retrieve This is a collection of categories ([[SoftLayer_Product_Item_Category]]) associated with a package which can be used for ordering. These categories have several objects prepopulated which are useful when determining the available products for purchase. The categories contain groups ([[SoftLayer_Product_Package_Item_Category_Group]]) that organize the products and prices by similar features. For example, operating systems will be grouped by their manufacturer and virtual server disks will be grouped by their disk type (SAN vs. local). Each group will contain prices ([[SoftLayer_Product_Item_Price]]) which you can use determine the cost of each product. Each price has a product ([[SoftLayer_Product_Item]]) which provides the name and other useful information about the server, service or software you may purchase.
func (r Product_Package) GetCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getCategories", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package) GetCdnItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getCdnItems", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package) GetCloudStorageItems(provider *int) (resp []datatypes.Product_Item, err error) {
	params := []interface{}{
		provider,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getCloudStorageItems", params, &r.Options, &resp)
	return
}

// Retrieve The item categories associated with a package, including information detailing which item categories are required as part of a SoftLayer product order.
func (r Product_Package) GetConfiguration() (resp []datatypes.Product_Package_Order_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getConfiguration", nil, &r.Options, &resp)
	return
}

// Retrieve The default boot category code for the package.
func (r Product_Package) GetDefaultBootCategoryCode() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDefaultBootCategoryCode", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of valid RAM items available for purchase in this package.
func (r Product_Package) GetDefaultRamItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDefaultRamItems", nil, &r.Options, &resp)
	return
}

// Retrieve The node type for a package in a solution deployment.
func (r Product_Package) GetDeploymentNodeType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDeploymentNodeType", nil, &r.Options, &resp)
	return
}

// Retrieve The packages that are allowed in a multi-server solution. (Deprecated)
func (r Product_Package) GetDeploymentPackages() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDeploymentPackages", nil, &r.Options, &resp)
	return
}

// Retrieve The solution deployment type.
func (r Product_Package) GetDeploymentType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDeploymentType", nil, &r.Options, &resp)
	return
}

// Retrieve The package that represents a multi-server solution. (Deprecated)
func (r Product_Package) GetDeployments() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDeployments", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates the package does not allow custom disk partitions.
func (r Product_Package) GetDisallowCustomDiskPartitions() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getDisallowCustomDiskPartitions", nil, &r.Options, &resp)
	return
}

// Retrieve The Softlayer order step is optionally step-based. This returns the first SoftLayer_Product_Package_Order_Step in the step-based order process.
func (r Product_Package) GetFirstOrderStep() (resp datatypes.Product_Package_Order_Step, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getFirstOrderStep", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the package is a specialized network gateway appliance package.
func (r Product_Package) GetGatewayApplianceFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getGatewayApplianceFlag", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates that the package supports GPUs.
func (r Product_Package) GetGpuFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getGpuFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Determines whether the package contains prices that can be ordered hourly.
func (r Product_Package) GetHourlyBillingAvailableFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getHourlyBillingAvailableFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Packages with this flag do not allow monthly orders.
func (r Product_Package) GetHourlyOnlyOrders() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getHourlyOnlyOrders", nil, &r.Options, &resp)
	return
}

// Returns a collection of SoftLayer_Product_Item_Attribute_Type objects.  These item attribute types specifically deal with when an item, SoftLayer_Product_Item, from the product catalog may no longer be available.  The keynames for these attribute types start with 'UNAVAILABLE_AFTER_DATE_*', where the '*' may represent any string.  For example, 'UNAVAILABLE_AFTER_DATE_NEW_ORDERS', signifies that the item is not available for new orders.  There is a catch all attribute type, 'UNAVAILABLE_AFTER_DATE_ALL'.  If an item has one of these availability attributes set, the value should be a valid date in MM/DD/YYYY, indicating the date after which the item will no longer be available.
func (r Product_Package) GetItemAvailabilityTypes() (resp []datatypes.Product_Item_Attribute_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemAvailabilityTypes", nil, &r.Options, &resp)
	return
}

// Retrieve The item-item conflicts associated with a package.
func (r Product_Package) GetItemConflicts() (resp []datatypes.Product_Item_Resource_Conflict, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemConflicts", nil, &r.Options, &resp)
	return
}

// Retrieve The item-location conflicts associated with a package.
func (r Product_Package) GetItemLocationConflicts() (resp []datatypes.Product_Item_Resource_Conflict, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemLocationConflicts", nil, &r.Options, &resp)
	return
}

// Retrieve cross reference for item prices
func (r Product_Package) GetItemPriceReferences() (resp []datatypes.Product_Package_Item_Prices, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemPriceReferences", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of SoftLayer_Product_Item_Prices that are valid for this package.
func (r Product_Package) GetItemPrices() (resp []datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemPrices", nil, &r.Options, &resp)
	return
}

// Return a collection of SoftLayer_Item_Price objects from a collection of SoftLayer_Software_Description
func (r Product_Package) GetItemPricesFromSoftwareDescriptions(softwareDescriptions []datatypes.Software_Description, includeTranslationsFlag *bool, returnAllPricesFlag *bool) (resp []datatypes.Product_Item_Price, err error) {
	params := []interface{}{
		softwareDescriptions,
		includeTranslationsFlag,
		returnAllPricesFlag,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemPricesFromSoftwareDescriptions", params, &r.Options, &resp)
	return
}

// Retrieve A collection of valid items available for purchase in this package.
func (r Product_Package) GetItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItems", nil, &r.Options, &resp)
	return
}

// Return a collection of [[SoftLayer_Product_Item]] objects from a [[SoftLayer_Virtual_Guest_Block_Device_Template_Group]] object
func (r Product_Package) GetItemsFromImageTemplate(imageTemplate *datatypes.Virtual_Guest_Block_Device_Template_Group) (resp []datatypes.Product_Item, err error) {
	params := []interface{}{
		imageTemplate,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getItemsFromImageTemplate", params, &r.Options, &resp)
	return
}

// Retrieve A collection of valid locations for this package. (Deprecated - Use [[SoftLayer_Product_Package/getRegions|getRegions]])
func (r Product_Package) GetLocations() (resp []datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getLocations", nil, &r.Options, &resp)
	return
}

// Retrieve The lowest server [[SoftLayer_Product_Item_Price]] related to this package.
func (r Product_Package) GetLowestServerPrice() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getLowestServerPrice", nil, &r.Options, &resp)
	return
}

// Retrieve The maximum available network speed associated with the package.
func (r Product_Package) GetMaximumPortSpeed() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getMaximumPortSpeed", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package) GetMessageQueueItems() (resp []datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getMessageQueueItems", nil, &r.Options, &resp)
	return
}

// Retrieve The minimum available network speed associated with the package.
func (r Product_Package) GetMinimumPortSpeed() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getMinimumPortSpeed", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates that this is a MongoDB engineered package. (Deprecated)
func (r Product_Package) GetMongoDbEngineeredFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getMongoDbEngineeredFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Services ordered from this package cannot have upgrades or downgrades performed.
func (r Product_Package) GetNoUpgradesFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getNoUpgradesFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the package is not in compliance with EU support.
func (r Product_Package) GetNonEuCompliantFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getNonEuCompliantFlag", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package) GetObject() (resp datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getObject", nil, &r.Options, &resp)
	return
}

// This method will return a collection of [[SoftLayer_Container_Product_Order_Network_Storage_Hub_Datacenter]] objects which contain a datacenter location and all the associated active usage rate prices where object storage is available. This method is really only applicable to the object storage additional service package which has a [[SoftLayer_Product_Package_Type]] of ”'ADDITIONAL_SERVICES_OBJECT_STORAGE”'. This information is useful so that you can see the "pay as you go" rates per datacenter.
func (r Product_Package) GetObjectStorageDatacenters() (resp []datatypes.Container_Product_Order_Network_Storage_Hub_Datacenter, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getObjectStorageDatacenters", nil, &r.Options, &resp)
	return
}

// This method will return a collection of [[SoftLayer_Container_Product_Order_Network_Storage_ObjectStorage_LocationGroup]] objects which contain a location group and all the associated active usage rate prices where object storage is available. This method is really only applicable to the object storage additional service package which has a [[SoftLayer_Product_Package_Type]] of ”'ADDITIONAL_SERVICES_OBJECT_STORAGE”'. This information is useful so that you can see the "pay as you go" rates per location group.
func (r Product_Package) GetObjectStorageLocationGroups() (resp []datatypes.Container_Product_Order_Network_Storage_ObjectStorage_LocationGroup, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getObjectStorageLocationGroups", nil, &r.Options, &resp)
	return
}

// Retrieve The premium price modifiers associated with the [[SoftLayer_Product_Item_Price]] and [[SoftLayer_Location]] objects in a package.
func (r Product_Package) GetOrderPremiums() (resp []datatypes.Product_Item_Price_Premium, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getOrderPremiums", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates if the package may be available in PoP locations in addition to Datacenters.
func (r Product_Package) GetPopLocationAvailabilityFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPopLocationAvailabilityFlag", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates the package is pre-configured. (Deprecated)
func (r Product_Package) GetPreconfiguredFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPreconfiguredFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the package requires the user to define a preset configuration.
func (r Product_Package) GetPresetConfigurationRequiredFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPresetConfigurationRequiredFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the package prevents the user from specifying a Vlan.
func (r Product_Package) GetPreventVlanSelectionFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPreventVlanSelectionFlag", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates the package is for a private hosted cloud deployment. (Deprecated)
func (r Product_Package) GetPrivateHostedCloudPackageFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPrivateHostedCloudPackageFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The server role of the private hosted cloud deployment. (Deprecated)
func (r Product_Package) GetPrivateHostedCloudPackageType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPrivateHostedCloudPackageType", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the package only has access to the private network.
func (r Product_Package) GetPrivateNetworkOnlyFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getPrivateNetworkOnlyFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the package is a specialized mass storage QuantaStor package. (Deprecated)
func (r Product_Package) GetQuantaStorPackageFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getQuantaStorPackageFlag", nil, &r.Options, &resp)
	return
}

// Retrieve This flag indicates the package does not allow different disks with RAID.
func (r Product_Package) GetRaidDiskRestrictionFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getRaidDiskRestrictionFlag", nil, &r.Options, &resp)
	return
}

// Retrieve This flag determines if the package contains a redundant power supply product.
func (r Product_Package) GetRedundantPowerFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getRedundantPowerFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The regional locations that a package is available in.
func (r Product_Package) GetRegions() (resp []datatypes.Location_Region, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getRegions", nil, &r.Options, &resp)
	return
}

// This call is similar to [[SoftLayer_Product_Package/getCategories|getCategories]], except that it does not include account-restricted pricing. Not all accounts have restricted pricing.
func (r Product_Package) GetStandardCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getStandardCategories", nil, &r.Options, &resp)
	return
}

// Retrieve The top level category code for this service offering.
func (r Product_Package) GetTopLevelItemCategoryCode() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getTopLevelItemCategoryCode", nil, &r.Options, &resp)
	return
}

// Retrieve The type of service offering. This property can be used to help filter packages.
func (r Product_Package) GetType() (resp datatypes.Product_Package_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package", "getType", nil, &r.Options, &resp)
	return
}

// Package presets are used to simplify ordering by eliminating the need for price ids when submitting orders.
//
// Orders submitted with a preset id defined will use the prices included in the package preset. Prices submitted on an order with a preset id will replace the prices included in the package preset for that prices category. If the package preset has a fixed configuration flag <em>(fixedConfigurationFlag)</em> set then the prices included in the preset configuration cannot be replaced by prices submitted on the order. The only exception to the fixed configuration flag would be if a price submitted on the order is an account-restricted price for the same product item.
type Product_Package_Preset struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductPackagePresetService returns an instance of the Product_Package_Preset SoftLayer service
func GetProductPackagePresetService(sess session.SLSession) Product_Package_Preset {
	return Product_Package_Preset{Session: sess}
}

func (r Product_Package_Preset) Id(id int) Product_Package_Preset {
	r.Options.Id = &id
	return r
}

func (r Product_Package_Preset) Mask(mask string) Product_Package_Preset {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Package_Preset) Filter(filter string) Product_Package_Preset {
	r.Options.Filter = filter
	return r
}

func (r Product_Package_Preset) Limit(limit int) Product_Package_Preset {
	r.Options.Limit = &limit
	return r
}

func (r Product_Package_Preset) Offset(offset int) Product_Package_Preset {
	r.Options.Offset = &offset
	return r
}

// This method returns all the active package presets.
func (r Product_Package_Preset) GetAllObjects() (resp []datatypes.Product_Package_Preset, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package_Preset) GetAvailableStorageUnits() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getAvailableStorageUnits", nil, &r.Options, &resp)
	return
}

// Retrieve When true this preset is for ordering a Bare Metal Reserved server.
func (r Product_Package_Preset) GetBareMetalReservedFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getBareMetalReservedFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The item categories that are included in this package preset configuration.
func (r Product_Package_Preset) GetCategories() (resp []datatypes.Product_Item_Category, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getCategories", nil, &r.Options, &resp)
	return
}

// Retrieve The compute family this configuration belongs to.
func (r Product_Package_Preset) GetComputeGroup() (resp datatypes.Product_Item_Server_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getComputeGroup", nil, &r.Options, &resp)
	return
}

// Retrieve The preset configuration (category and price).
func (r Product_Package_Preset) GetConfiguration() (resp []datatypes.Product_Package_Preset_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getConfiguration", nil, &r.Options, &resp)
	return
}

// Retrieve When true this preset is only allowed to upgrade/downgrade to other presets in the same compute family.
func (r Product_Package_Preset) GetDisallowedComputeGroupUpgradeFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getDisallowedComputeGroupUpgradeFlag", nil, &r.Options, &resp)
	return
}

// Retrieve A package preset with this flag set will not allow the price's defined in the preset configuration to be overriden during order placement.
func (r Product_Package_Preset) GetFixedConfigurationFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getFixedConfigurationFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The locations this preset configuration is available in. If empty the preset is available in all locations the package is available in.
func (r Product_Package_Preset) GetLocations() (resp []datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getLocations", nil, &r.Options, &resp)
	return
}

// Retrieve The lowest server prices related to this package preset.
func (r Product_Package_Preset) GetLowestPresetServerPrice() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getLowestPresetServerPrice", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package_Preset) GetObject() (resp datatypes.Product_Package_Preset, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The package this preset belongs to.
func (r Product_Package_Preset) GetPackage() (resp datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getPackage", nil, &r.Options, &resp)
	return
}

// Retrieve The item categories associated with a package preset, including information detailing which item categories are required as part of a SoftLayer product order.
func (r Product_Package_Preset) GetPackageConfiguration() (resp []datatypes.Product_Package_Order_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getPackageConfiguration", nil, &r.Options, &resp)
	return
}

// Retrieve The item prices that are included in this package preset configuration.
func (r Product_Package_Preset) GetPrices() (resp []datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getPrices", nil, &r.Options, &resp)
	return
}

// Retrieve Describes how all disks in this preset will be configured.
func (r Product_Package_Preset) GetStorageGroupTemplateArrays() (resp []datatypes.Configuration_Storage_Group_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getStorageGroupTemplateArrays", nil, &r.Options, &resp)
	return
}

// Retrieve The starting hourly price for this configuration. Additional options not defined in the preset may increase the cost.
func (r Product_Package_Preset) GetTotalMinimumHourlyFee() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getTotalMinimumHourlyFee", nil, &r.Options, &resp)
	return
}

// Retrieve The starting monthly price for this configuration. Additional options not defined in the preset may increase the cost.
func (r Product_Package_Preset) GetTotalMinimumRecurringFee() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Preset", "getTotalMinimumRecurringFee", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Product_Package_Server data type contains summarized information for bare metal servers regarding pricing, processor stats, and feature sets.
type Product_Package_Server struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductPackageServerService returns an instance of the Product_Package_Server SoftLayer service
func GetProductPackageServerService(sess session.SLSession) Product_Package_Server {
	return Product_Package_Server{Session: sess}
}

func (r Product_Package_Server) Id(id int) Product_Package_Server {
	r.Options.Id = &id
	return r
}

func (r Product_Package_Server) Mask(mask string) Product_Package_Server {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Package_Server) Filter(filter string) Product_Package_Server {
	r.Options.Filter = filter
	return r
}

func (r Product_Package_Server) Limit(limit int) Product_Package_Server {
	r.Options.Limit = &limit
	return r
}

func (r Product_Package_Server) Offset(offset int) Product_Package_Server {
	r.Options.Offset = &offset
	return r
}

// This method will grab all the package servers.
func (r Product_Package_Server) GetAllObjects() (resp []datatypes.Product_Package_Server, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getAllObjects", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package_Server) GetCatalog() (resp datatypes.Product_Catalog, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getCatalog", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package_Server) GetItem() (resp datatypes.Product_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getItem", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package_Server) GetItemPrice() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getItemPrice", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package_Server) GetObject() (resp datatypes.Product_Package_Server, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package_Server) GetPackage() (resp datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getPackage", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Product_Package_Server) GetPreset() (resp datatypes.Product_Package_Preset, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server", "getPreset", nil, &r.Options, &resp)
	return
}

// The [[SoftLayer_Product_Package_Server_Option]] data type contains various data points associated with package servers that can be used in selection criteria.
type Product_Package_Server_Option struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductPackageServerOptionService returns an instance of the Product_Package_Server_Option SoftLayer service
func GetProductPackageServerOptionService(sess session.SLSession) Product_Package_Server_Option {
	return Product_Package_Server_Option{Session: sess}
}

func (r Product_Package_Server_Option) Id(id int) Product_Package_Server_Option {
	r.Options.Id = &id
	return r
}

func (r Product_Package_Server_Option) Mask(mask string) Product_Package_Server_Option {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Package_Server_Option) Filter(filter string) Product_Package_Server_Option {
	r.Options.Filter = filter
	return r
}

func (r Product_Package_Server_Option) Limit(limit int) Product_Package_Server_Option {
	r.Options.Limit = &limit
	return r
}

func (r Product_Package_Server_Option) Offset(offset int) Product_Package_Server_Option {
	r.Options.Offset = &offset
	return r
}

// This method will grab all the package server options.
func (r Product_Package_Server_Option) GetAllOptions() (resp []datatypes.Product_Package_Server_Option, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server_Option", "getAllOptions", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package_Server_Option) GetObject() (resp datatypes.Product_Package_Server_Option, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server_Option", "getObject", nil, &r.Options, &resp)
	return
}

// This method will grab all the package server options for the specified type.
func (r Product_Package_Server_Option) GetOptions(typ *string) (resp []datatypes.Product_Package_Server_Option, err error) {
	params := []interface{}{
		typ,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Package_Server_Option", "getOptions", params, &r.Options, &resp)
	return
}

// The [[SoftLayer_Product_Package_Type]] object indicates the type for a service offering (package). The type can be used to filter packages. For example, if you are looking for the package representing virtual servers, you can filter on the type's key name of ”'VIRTUAL_SERVER_INSTANCE”'. For bare metal servers by core or CPU, filter on ”'BARE_METAL_CORE”' or ”'BARE_METAL_CPU”', respectively.
type Product_Package_Type struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductPackageTypeService returns an instance of the Product_Package_Type SoftLayer service
func GetProductPackageTypeService(sess session.SLSession) Product_Package_Type {
	return Product_Package_Type{Session: sess}
}

func (r Product_Package_Type) Id(id int) Product_Package_Type {
	r.Options.Id = &id
	return r
}

func (r Product_Package_Type) Mask(mask string) Product_Package_Type {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Package_Type) Filter(filter string) Product_Package_Type {
	r.Options.Filter = filter
	return r
}

func (r Product_Package_Type) Limit(limit int) Product_Package_Type {
	r.Options.Limit = &limit
	return r
}

func (r Product_Package_Type) Offset(offset int) Product_Package_Type {
	r.Options.Offset = &offset
	return r
}

// This method will return all of the available package types.
func (r Product_Package_Type) GetAllObjects() (resp []datatypes.Product_Package_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Type", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Package_Type) GetObject() (resp datatypes.Product_Package_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Type", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve All the packages associated with the given package type.
func (r Product_Package_Type) GetPackages() (resp []datatypes.Product_Package, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Package_Type", "getPackages", nil, &r.Options, &resp)
	return
}

// no documentation yet
type Product_Promotion struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductPromotionService returns an instance of the Product_Promotion SoftLayer service
func GetProductPromotionService(sess session.SLSession) Product_Promotion {
	return Product_Promotion{Session: sess}
}

func (r Product_Promotion) Id(id int) Product_Promotion {
	r.Options.Id = &id
	return r
}

func (r Product_Promotion) Mask(mask string) Product_Promotion {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Promotion) Filter(filter string) Product_Promotion {
	r.Options.Filter = filter
	return r
}

func (r Product_Promotion) Limit(limit int) Product_Promotion {
	r.Options.Limit = &limit
	return r
}

func (r Product_Promotion) Offset(offset int) Product_Promotion {
	r.Options.Offset = &offset
	return r
}

// Retrieves a promotion using its code.
func (r Product_Promotion) FindByPromoCode(code *string) (resp datatypes.Container_Product_Promotion, err error) {
	params := []interface{}{
		code,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Promotion", "findByPromoCode", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Product_Promotion) GetObject() (resp datatypes.Product_Promotion, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Promotion", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Product_Upgrade_Request data type contains general information relating to a hardware, virtual server, or service upgrade. It also relates a [[SoftLayer_Billing_Order]] to a [[SoftLayer_Ticket]].
type Product_Upgrade_Request struct {
	Session session.SLSession
	Options sl.Options
}

// GetProductUpgradeRequestService returns an instance of the Product_Upgrade_Request SoftLayer service
func GetProductUpgradeRequestService(sess session.SLSession) Product_Upgrade_Request {
	return Product_Upgrade_Request{Session: sess}
}

func (r Product_Upgrade_Request) Id(id int) Product_Upgrade_Request {
	r.Options.Id = &id
	return r
}

func (r Product_Upgrade_Request) Mask(mask string) Product_Upgrade_Request {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Product_Upgrade_Request) Filter(filter string) Product_Upgrade_Request {
	r.Options.Filter = filter
	return r
}

func (r Product_Upgrade_Request) Limit(limit int) Product_Upgrade_Request {
	r.Options.Limit = &limit
	return r
}

func (r Product_Upgrade_Request) Offset(offset int) Product_Upgrade_Request {
	r.Options.Offset = &offset
	return r
}

// When a change is made to an upgrade by Sales, this method will approve the changes that were made. A customer must acknowledge the change and approve it so that the upgrade request can proceed.
func (r Product_Upgrade_Request) ApproveChanges() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "approveChanges", nil, &r.Options, &resp)
	return
}

// Retrieve The account that an order belongs to
func (r Product_Upgrade_Request) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve Indicates that the upgrade request has completed or has been cancelled.
func (r Product_Upgrade_Request) GetCompletedFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getCompletedFlag", nil, &r.Options, &resp)
	return
}

// Retrieve This is the invoice associated with the upgrade request. For hourly servers or services, an invoice will not be available.
func (r Product_Upgrade_Request) GetInvoice() (resp datatypes.Billing_Invoice, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getInvoice", nil, &r.Options, &resp)
	return
}

// getObject retrieves a SoftLayer_Product_Upgrade_Request object on your account whose ID corresponds to the ID of the init parameter passed to the SoftLayer_Product_Upgrade_Request service.
func (r Product_Upgrade_Request) GetObject() (resp datatypes.Product_Upgrade_Request, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve An order record associated to the upgrade request
func (r Product_Upgrade_Request) GetOrder() (resp datatypes.Billing_Order, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getOrder", nil, &r.Options, &resp)
	return
}

// Retrieve A server object associated with the upgrade request if any.
func (r Product_Upgrade_Request) GetServer() (resp datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getServer", nil, &r.Options, &resp)
	return
}

// Retrieve The current status of the upgrade request.
func (r Product_Upgrade_Request) GetStatus() (resp datatypes.Product_Upgrade_Request_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getStatus", nil, &r.Options, &resp)
	return
}

// Retrieve The ticket that is used to coordinate the upgrade process.
func (r Product_Upgrade_Request) GetTicket() (resp datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getTicket", nil, &r.Options, &resp)
	return
}

// Retrieve The user that placed the order.
func (r Product_Upgrade_Request) GetUser() (resp datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getUser", nil, &r.Options, &resp)
	return
}

// Retrieve A virtual server object associated with the upgrade request if any.
func (r Product_Upgrade_Request) GetVirtualGuest() (resp datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "getVirtualGuest", nil, &r.Options, &resp)
	return
}

// In case an upgrade cannot be performed, the maintenance window needs to be updated to a future date.
func (r Product_Upgrade_Request) UpdateMaintenanceWindow(maintenanceStartTime *string, maintenanceWindowId *int) (resp bool, err error) {
	params := []interface{}{
		maintenanceStartTime,
		maintenanceWindowId,
	}
	err = r.Session.DoRequest("SoftLayer_Product_Upgrade_Request", "updateMaintenanceWindow", params, &r.Options, &resp)
	return
}
