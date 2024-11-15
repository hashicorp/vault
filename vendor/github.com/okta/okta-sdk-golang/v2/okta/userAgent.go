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

package okta

import "runtime"

type UserAgent struct {
	goVersion string

	osName string

	osVersion string

	config *config
}

func NewUserAgent(config *config) UserAgent {
	ua := UserAgent{}
	ua.config = config
	ua.goVersion = runtime.Version()
	ua.osName = runtime.GOOS
	ua.osVersion = runtime.GOARCH

	return ua
}

func (ua UserAgent) String() string {
	userAgentString := "okta-sdk-golang/" + Version + " "
	userAgentString += "golang/" + ua.goVersion + " "
	userAgentString += ua.osName + "/" + ua.osVersion

	if ua.config.UserAgentExtra != "" {
		userAgentString += " " + ua.config.UserAgentExtra
	}

	return userAgentString
}
