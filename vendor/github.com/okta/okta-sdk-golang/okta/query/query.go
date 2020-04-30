/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package query

import (
	"net/url"
	"strconv"
)

type Params struct {
	Q                    string `json:"q,omitempty"`
	After                string `json:"after,omitempty"`
	Limit                int64  `json:"limit,omitempty"`
	Filter               string `json:"filter,omitempty"`
	Expand               string `json:"expand,omitempty"`
	IncludeNonDeleted    *bool  `json:"includeNonDeleted,omitempty"`
	Activate             *bool  `json:"activate,omitempty"`
	TargetAid            string `json:"targetAid,omitempty"`
	QueryScope           string `json:"query_scope,omitempty"`
	SendEmail            *bool  `json:"sendEmail,omitempty"`
	RemoveUsers          *bool  `json:"removeUsers,omitempty"`
	ManagedBy            string `json:"managedBy,omitempty"`
	Until                string `json:"until,omitempty"`
	Since                string `json:"since,omitempty"`
	SortOrder            string `json:"sortOrder,omitempty"`
	Type                 string `json:"type,omitempty"`
	Status               string `json:"status,omitempty"`
	Format               string `json:"format,omitempty"`
	Search               string `json:"search,omitempty"`
	Provider             string `json:"provider,omitempty"`
	NextLogin            string `json:"nextLogin,omitempty"`
	Strict               *bool  `json:"strict,omitempty"`
	ShowAll              *bool  `json:"showAll,omitempty"`
	UpdatePhone          *bool  `json:"updatePhone,omitempty"`
	TemplateId           string `json:"templateId,omitempty"`
	TokenLifetimeSeconds int64  `json:"tokenLifetimeSeconds,omitempty"`
	TempPassword         *bool  `json:"tempPassword,omitempty"`
	OauthTokens          *bool  `json:"oauthTokens,omitempty"`
}

func NewQueryParams(paramOpt ...ParamOptions) *Params {
	p := &Params{}

	for _, par := range paramOpt {
		par(p)
	}

	return p
}

type ParamOptions func(*Params)

func WithQ(queryQ string) ParamOptions {
	return func(p *Params) {
		p.Q = queryQ
	}
}
func WithAfter(queryAfter string) ParamOptions {
	return func(p *Params) {
		p.After = queryAfter
	}
}
func WithLimit(queryLimit int64) ParamOptions {
	return func(p *Params) {
		p.Limit = queryLimit
	}
}
func WithFilter(queryFilter string) ParamOptions {
	return func(p *Params) {
		p.Filter = queryFilter
	}
}
func WithExpand(queryExpand string) ParamOptions {
	return func(p *Params) {
		p.Expand = queryExpand
	}
}
func WithIncludeNonDeleted(queryIncludeNonDeleted bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryIncludeNonDeleted
		p.IncludeNonDeleted = b
	}
}
func WithActivate(queryActivate bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryActivate
		p.Activate = b
	}
}
func WithTargetAid(queryTargetAid string) ParamOptions {
	return func(p *Params) {
		p.TargetAid = queryTargetAid
	}
}
func WithQueryScope(queryQueryScope string) ParamOptions {
	return func(p *Params) {
		p.QueryScope = queryQueryScope
	}
}
func WithSendEmail(querySendEmail bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = querySendEmail
		p.SendEmail = b
	}
}
func WithRemoveUsers(queryRemoveUsers bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryRemoveUsers
		p.RemoveUsers = b
	}
}
func WithManagedBy(queryManagedBy string) ParamOptions {
	return func(p *Params) {
		p.ManagedBy = queryManagedBy
	}
}
func WithUntil(queryUntil string) ParamOptions {
	return func(p *Params) {
		p.Until = queryUntil
	}
}
func WithSince(querySince string) ParamOptions {
	return func(p *Params) {
		p.Since = querySince
	}
}
func WithSortOrder(querySortOrder string) ParamOptions {
	return func(p *Params) {
		p.SortOrder = querySortOrder
	}
}
func WithType(queryType string) ParamOptions {
	return func(p *Params) {
		p.Type = queryType
	}
}
func WithStatus(queryStatus string) ParamOptions {
	return func(p *Params) {
		p.Status = queryStatus
	}
}
func WithFormat(queryFormat string) ParamOptions {
	return func(p *Params) {
		p.Format = queryFormat
	}
}
func WithSearch(querySearch string) ParamOptions {
	return func(p *Params) {
		p.Search = querySearch
	}
}
func WithProvider(queryProvider string) ParamOptions {
	return func(p *Params) {
		p.Provider = queryProvider
	}
}
func WithNextLogin(queryNextLogin string) ParamOptions {
	return func(p *Params) {
		p.NextLogin = queryNextLogin
	}
}
func WithStrict(queryStrict bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryStrict
		p.Strict = b
	}
}
func WithShowAll(queryShowAll bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryShowAll
		p.ShowAll = b
	}
}
func WithUpdatePhone(queryUpdatePhone bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryUpdatePhone
		p.UpdatePhone = b
	}
}
func WithTemplateId(queryTemplateId string) ParamOptions {
	return func(p *Params) {
		p.TemplateId = queryTemplateId
	}
}
func WithTokenLifetimeSeconds(queryTokenLifetimeSeconds int64) ParamOptions {
	return func(p *Params) {
		p.TokenLifetimeSeconds = queryTokenLifetimeSeconds
	}
}
func WithTempPassword(queryTempPassword bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryTempPassword
		p.TempPassword = b
	}
}
func WithOauthTokens(queryOauthTokens bool) ParamOptions {
	return func(p *Params) {
		b := new(bool)
		*b = queryOauthTokens
		p.OauthTokens = b
	}
}

func (p *Params) String() string {
	qs := url.Values{}

	if p.Q != "" {
		qs.Add(`q`, p.Q)
	}
	if p.After != "" {
		qs.Add(`after`, p.After)
	}
	if p.Limit != 0 {
		qs.Add(`limit`, strconv.FormatInt(p.Limit, 10))
	}
	if p.Filter != "" {
		qs.Add(`filter`, p.Filter)
	}
	if p.Expand != "" {
		qs.Add(`expand`, p.Expand)
	}
	if p.IncludeNonDeleted != nil {
		qs.Add(`includeNonDeleted`, strconv.FormatBool(*p.IncludeNonDeleted))
	}
	if p.Activate != nil {
		qs.Add(`activate`, strconv.FormatBool(*p.Activate))
	}
	if p.TargetAid != "" {
		qs.Add(`targetAid`, p.TargetAid)
	}
	if p.QueryScope != "" {
		qs.Add(`query_scope`, p.QueryScope)
	}
	if p.SendEmail != nil {
		qs.Add(`sendEmail`, strconv.FormatBool(*p.SendEmail))
	}
	if p.RemoveUsers != nil {
		qs.Add(`removeUsers`, strconv.FormatBool(*p.RemoveUsers))
	}
	if p.ManagedBy != "" {
		qs.Add(`managedBy`, p.ManagedBy)
	}
	if p.Until != "" {
		qs.Add(`until`, p.Until)
	}
	if p.Since != "" {
		qs.Add(`since`, p.Since)
	}
	if p.SortOrder != "" {
		qs.Add(`sortOrder`, p.SortOrder)
	}
	if p.Type != "" {
		qs.Add(`type`, p.Type)
	}
	if p.Status != "" {
		qs.Add(`status`, p.Status)
	}
	if p.Format != "" {
		qs.Add(`format`, p.Format)
	}
	if p.Search != "" {
		qs.Add(`search`, p.Search)
	}
	if p.Provider != "" {
		qs.Add(`provider`, p.Provider)
	}
	if p.NextLogin != "" {
		qs.Add(`nextLogin`, p.NextLogin)
	}
	if p.Strict != nil {
		qs.Add(`strict`, strconv.FormatBool(*p.Strict))
	}
	if p.ShowAll != nil {
		qs.Add(`showAll`, strconv.FormatBool(*p.ShowAll))
	}
	if p.UpdatePhone != nil {
		qs.Add(`updatePhone`, strconv.FormatBool(*p.UpdatePhone))
	}
	if p.TemplateId != "" {
		qs.Add(`templateId`, p.TemplateId)
	}
	if p.TokenLifetimeSeconds != 0 {
		qs.Add(`tokenLifetimeSeconds`, strconv.FormatInt(p.TokenLifetimeSeconds, 10))
	}
	if p.TempPassword != nil {
		qs.Add(`tempPassword`, strconv.FormatBool(*p.TempPassword))
	}
	if p.OauthTokens != nil {
		qs.Add(`oauthTokens`, strconv.FormatBool(*p.OauthTokens))
	}

	if len(qs) != 0 {
		return "?" + qs.Encode()
	}
	return ""
}
