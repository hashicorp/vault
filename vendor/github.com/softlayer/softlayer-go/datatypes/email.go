/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * AUTOMATICALLY GENERATED CODE - DO NOT MODIFY
 */

package datatypes

// no documentation yet
type Email_Subscription struct {
	Entity

	// Brief description of the purpose of the email.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Enabled *bool `json:"enabled,omitempty" xmlrpc:"enabled,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Email template name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Email_Subscription_Group struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Email subscription group name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of all email subscriptions associated with this group.
	SubscriptionCount *uint `json:"subscriptionCount,omitempty" xmlrpc:"subscriptionCount,omitempty"`

	// All email subscriptions associated with this group.
	Subscriptions []Email_Subscription `json:"subscriptions,omitempty" xmlrpc:"subscriptions,omitempty"`
}

// no documentation yet
type Email_Subscription_Suppression_User struct {
	Entity

	// no documentation yet
	Subscription *Email_Subscription `json:"subscription,omitempty" xmlrpc:"subscription,omitempty"`
}
