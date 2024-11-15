/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auth

import (
	"fmt"
	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
)

type Credential interface {
}

func ToCredentialsProvider(credential Credential) (provider credentials.CredentialsProvider, err error) {
	if credential == nil {
		provider = credentials.NewDefaultCredentialsProvider()
		return
	}

	switch instance := credential.(type) {
	case *credentials.AccessKeyCredential:
		{
			provider = credentials.NewStaticAKCredentialsProvider(instance.AccessKeyId, instance.AccessKeySecret)
			return
		}
	case *credentials.StsTokenCredential:
		{
			provider = credentials.NewStaticSTSCredentialsProvider(instance.AccessKeyId, instance.AccessKeySecret, instance.AccessKeyStsToken)
			return
		}
	case *credentials.BearerTokenCredential:
		{
			provider = credentials.NewBearerTokenCredentialsProvider(instance.BearerToken)
			return
		}
	case *credentials.RamRoleArnCredential:
		{
			preProvider := credentials.NewStaticAKCredentialsProvider(instance.AccessKeyId, instance.AccessKeySecret)
			provider, err = credentials.NewRAMRoleARNCredentialsProvider(
				preProvider,
				instance.RoleArn,
				instance.RoleSessionName,
				instance.RoleSessionExpiration,
				instance.Policy,
				instance.StsRegion,
				instance.ExternalId)
			return
		}
	case *credentials.RsaKeyPairCredential:
		{
			provider, err = credentials.NewRSAKeyPairCredentialsProvider(instance.PublicKeyId, instance.PrivateKey, instance.SessionExpiration)
			return
		}
	case *credentials.EcsRamRoleCredential:
		{
			provider = credentials.NewECSRAMRoleCredentialsProvider(instance.RoleName)
			return
		}
	case credentials.CredentialsProvider:
		{
			provider = instance
			return
		}
	default:
		message := fmt.Sprintf(errors.UnsupportedCredentialErrorMessage, reflect.TypeOf(credential))
		err = errors.NewClientError(errors.UnsupportedCredentialErrorCode, message, nil)
	}
	return
}
