// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20190118

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AlgorithmInfo struct {

	// 算法的标识
	KeyUsage *string `json:"KeyUsage,omitempty" name:"KeyUsage"`

	// 算法的名称
	Algorithm *string `json:"Algorithm,omitempty" name:"Algorithm"`
}

type AsymmetricRsaDecryptRequest struct {
	*tchttp.BaseRequest

	// CMK的唯一标识
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 使用PublicKey加密的密文，Base64编码
	Ciphertext *string `json:"Ciphertext,omitempty" name:"Ciphertext"`

	// 在使用公钥加密时对应的算法：当前支持RSAES_PKCS1_V1_5、RSAES_OAEP_SHA_1、RSAES_OAEP_SHA_256
	Algorithm *string `json:"Algorithm,omitempty" name:"Algorithm"`
}

func (r *AsymmetricRsaDecryptRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AsymmetricRsaDecryptRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AsymmetricRsaDecryptResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的唯一标识
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 解密后的明文，base64编码
		Plaintext *string `json:"Plaintext,omitempty" name:"Plaintext"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *AsymmetricRsaDecryptResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AsymmetricRsaDecryptResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AsymmetricSm2DecryptRequest struct {
	*tchttp.BaseRequest

	// CMK的唯一标识
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 使用PublicKey加密的密文，Base64编码。密文长度不能超过256字节。
	Ciphertext *string `json:"Ciphertext,omitempty" name:"Ciphertext"`
}

func (r *AsymmetricSm2DecryptRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AsymmetricSm2DecryptRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AsymmetricSm2DecryptResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的唯一标识
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 解密后的明文，base64编码
		Plaintext *string `json:"Plaintext,omitempty" name:"Plaintext"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *AsymmetricSm2DecryptResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AsymmetricSm2DecryptResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CancelKeyDeletionRequest struct {
	*tchttp.BaseRequest

	// 需要被取消删除的CMK的唯一标志
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *CancelKeyDeletionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CancelKeyDeletionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CancelKeyDeletionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一标志被取消删除的CMK。
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *CancelKeyDeletionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CancelKeyDeletionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateKeyRequest struct {
	*tchttp.BaseRequest

	// 作为密钥更容易辨识，更容易被人看懂的别名， 不可为空，1-60个字母数字 - _ 的组合，首字符必须为字母或者数字。以 kms- 作为前缀的用于云产品使用，Alias 不可重复。
	Alias *string `json:"Alias,omitempty" name:"Alias"`

	// CMK 的描述，最大1024字节
	Description *string `json:"Description,omitempty" name:"Description"`

	// 指定key的用途，默认为  "ENCRYPT_DECRYPT" 表示创建对称加解密密钥，其它支持用途 “ASYMMETRIC_DECRYPT_RSA_2048” 表示创建用于加解密的RSA2048非对称密钥，“ASYMMETRIC_DECRYPT_SM2” 表示创建用于加解密的SM2非对称密钥
	KeyUsage *string `json:"KeyUsage,omitempty" name:"KeyUsage"`

	// 指定key类型，默认为1，1表示默认类型，由KMS创建CMK密钥，2 表示EXTERNAL 类型，该类型需要用户导入密钥材料，参考 GetParametersForImport 和 ImportKeyMaterial 接口
	Type *uint64 `json:"Type,omitempty" name:"Type"`
}

func (r *CreateKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的全局唯一标识符
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 作为密钥更容易辨识，更容易被人看懂的别名
		Alias *string `json:"Alias,omitempty" name:"Alias"`

		// 密钥创建时间，unix时间戳
		CreateTime *uint64 `json:"CreateTime,omitempty" name:"CreateTime"`

		// CMK的描述
	// 注意：此字段可能返回 null，表示取不到有效值。
		Description *string `json:"Description,omitempty" name:"Description"`

		// CMK的状态
		KeyState *string `json:"KeyState,omitempty" name:"KeyState"`

		// CMK的用途
		KeyUsage *string `json:"KeyUsage,omitempty" name:"KeyUsage"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DecryptRequest struct {
	*tchttp.BaseRequest

	// 待解密的密文数据
	CiphertextBlob *string `json:"CiphertextBlob,omitempty" name:"CiphertextBlob"`

	// key/value对的json字符串，如果Encrypt指定了该参数，则在调用Decrypt API时需要提供同样的参数，最大支持1024字符
	EncryptionContext *string `json:"EncryptionContext,omitempty" name:"EncryptionContext"`
}

func (r *DecryptRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DecryptRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DecryptResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的全局唯一标识
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 解密后的明文。该字段是base64编码的，为了得到原始明文，调用方需要进行base64解码
		Plaintext *string `json:"Plaintext,omitempty" name:"Plaintext"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DecryptResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DecryptResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteImportedKeyMaterialRequest struct {
	*tchttp.BaseRequest

	// 指定需要删除密钥材料的EXTERNAL CMK。
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *DeleteImportedKeyMaterialRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteImportedKeyMaterialRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteImportedKeyMaterialResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteImportedKeyMaterialResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteImportedKeyMaterialResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeKeyRequest struct {
	*tchttp.BaseRequest

	// CMK全局唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *DescribeKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 密钥属性信息
	// 注意：此字段可能返回 null，表示取不到有效值。
		KeyMetadata *KeyMetadata `json:"KeyMetadata,omitempty" name:"KeyMetadata"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeKeysRequest struct {
	*tchttp.BaseRequest

	// 查询CMK的ID列表，批量查询一次最多支持100个KeyId
	KeyIds []*string `json:"KeyIds,omitempty" name:"KeyIds" list`
}

func (r *DescribeKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeKeysRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返回的属性信息列表
	// 注意：此字段可能返回 null，表示取不到有效值。
		KeyMetadatas []*KeyMetadata `json:"KeyMetadatas,omitempty" name:"KeyMetadatas" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeKeysResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableKeyRequest struct {
	*tchttp.BaseRequest

	// CMK唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *DisableKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DisableKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableKeyRotationRequest struct {
	*tchttp.BaseRequest

	// CMK唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *DisableKeyRotationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableKeyRotationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableKeyRotationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DisableKeyRotationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableKeyRotationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableKeysRequest struct {
	*tchttp.BaseRequest

	// 需要批量禁用的CMK Id 列表，CMK数量最大支持100
	KeyIds []*string `json:"KeyIds,omitempty" name:"KeyIds" list`
}

func (r *DisableKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableKeysRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DisableKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *DisableKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DisableKeysResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableKeyRequest struct {
	*tchttp.BaseRequest

	// CMK唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *EnableKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *EnableKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableKeyRotationRequest struct {
	*tchttp.BaseRequest

	// CMK唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *EnableKeyRotationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableKeyRotationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableKeyRotationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *EnableKeyRotationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableKeyRotationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableKeysRequest struct {
	*tchttp.BaseRequest

	// 需要批量启用的CMK Id 列表， CMK数量最大支持100
	KeyIds []*string `json:"KeyIds,omitempty" name:"KeyIds" list`
}

func (r *EnableKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableKeysRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EnableKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *EnableKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EnableKeysResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EncryptRequest struct {
	*tchttp.BaseRequest

	// 调用CreateKey生成的CMK全局唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 被加密的明文数据，该字段必须使用base64编码，原文最大长度支持4K
	Plaintext *string `json:"Plaintext,omitempty" name:"Plaintext"`

	// key/value对的json字符串，如果指定了该参数，则在调用Decrypt API时需要提供同样的参数，最大支持1024个字符
	EncryptionContext *string `json:"EncryptionContext,omitempty" name:"EncryptionContext"`
}

func (r *EncryptRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EncryptRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EncryptResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 加密后经过base64编码的密文
		CiphertextBlob *string `json:"CiphertextBlob,omitempty" name:"CiphertextBlob"`

		// 加密使用的CMK的全局唯一标识
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *EncryptResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EncryptResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GenerateDataKeyRequest struct {
	*tchttp.BaseRequest

	// CMK全局唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 指定生成Datakey的加密算法以及Datakey大小，AES_128或者AES_256。KeySpec 和 NumberOfBytes 必须指定一个
	KeySpec *string `json:"KeySpec,omitempty" name:"KeySpec"`

	// 生成的DataKey的长度，同时指定NumberOfBytes和KeySpec时，以NumberOfBytes为准。最小值为1， 最大值为1024。KeySpec 和 NumberOfBytes 必须指定一个
	NumberOfBytes *uint64 `json:"NumberOfBytes,omitempty" name:"NumberOfBytes"`

	// key/value对的json字符串，如果使用该字段，则返回的DataKey在解密时需要填入相同的字符串
	EncryptionContext *string `json:"EncryptionContext,omitempty" name:"EncryptionContext"`
}

func (r *GenerateDataKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GenerateDataKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GenerateDataKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的全局唯一标识
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 生成的DataKey的明文，该明文使用base64编码，用户需要使用base64解码得到明文
		Plaintext *string `json:"Plaintext,omitempty" name:"Plaintext"`

		// DataKey加密后经过base64编码的密文，用户需要自行保存密文
		CiphertextBlob *string `json:"CiphertextBlob,omitempty" name:"CiphertextBlob"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *GenerateDataKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GenerateDataKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GenerateRandomRequest struct {
	*tchttp.BaseRequest

	// 生成的随机数的长度。最小值为1， 最大值为1024。
	NumberOfBytes *uint64 `json:"NumberOfBytes,omitempty" name:"NumberOfBytes"`
}

func (r *GenerateRandomRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GenerateRandomRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GenerateRandomResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 生成的随机数的明文，该明文使用base64编码，用户需要使用base64解码得到明文。
		Plaintext *string `json:"Plaintext,omitempty" name:"Plaintext"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *GenerateRandomResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GenerateRandomResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetKeyRotationStatusRequest struct {
	*tchttp.BaseRequest

	// CMK唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *GetKeyRotationStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetKeyRotationStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetKeyRotationStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 密钥轮换是否开启
		KeyRotationEnabled *bool `json:"KeyRotationEnabled,omitempty" name:"KeyRotationEnabled"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetKeyRotationStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetKeyRotationStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetParametersForImportRequest struct {
	*tchttp.BaseRequest

	// CMK的唯一标识，获取密钥参数的CMK必须是EXTERNAL类型，即在CreateKey时指定Type=2 类型的CMK。
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 指定加密密钥材料的算法，目前支持RSAES_PKCS1_V1_5、RSAES_OAEP_SHA_1、RSAES_OAEP_SHA_256
	WrappingAlgorithm *string `json:"WrappingAlgorithm,omitempty" name:"WrappingAlgorithm"`

	// 指定加密密钥材料的类型，目前只支持RSA_2048
	WrappingKeySpec *string `json:"WrappingKeySpec,omitempty" name:"WrappingKeySpec"`
}

func (r *GetParametersForImportRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetParametersForImportRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetParametersForImportResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的唯一标识，用于指定目标导入密钥材料的CMK。
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 导入密钥材料需要的token，用于作为 ImportKeyMaterial 的参数。
		ImportToken *string `json:"ImportToken,omitempty" name:"ImportToken"`

		// 用于加密密钥材料的RSA公钥，base64编码。使用PublicKey base64解码后的公钥将导入密钥进行加密后作为 ImportKeyMaterial 的参数。
		PublicKey *string `json:"PublicKey,omitempty" name:"PublicKey"`

		// 该导出token和公钥的有效期，超过该时间后无法导入，需要重新调用GetParametersForImport获取。
		ParametersValidTo *uint64 `json:"ParametersValidTo,omitempty" name:"ParametersValidTo"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetParametersForImportResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetParametersForImportResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetPublicKeyRequest struct {
	*tchttp.BaseRequest

	// CMK的唯一标识。
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *GetPublicKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetPublicKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetPublicKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的唯一标识。
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 经过base64编码的公钥内容。
		PublicKey *string `json:"PublicKey,omitempty" name:"PublicKey"`

		// PEM格式的公钥内容。
		PublicKeyPem *string `json:"PublicKeyPem,omitempty" name:"PublicKeyPem"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetPublicKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetPublicKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetServiceStatusRequest struct {
	*tchttp.BaseRequest
}

func (r *GetServiceStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetServiceStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetServiceStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// KMS服务是否开通， true 表示已开通
		ServiceEnabled *bool `json:"ServiceEnabled,omitempty" name:"ServiceEnabled"`

		// 服务不可用类型： 0-未购买，1-正常， 2-欠费停服， 3-资源释放
	// 注意：此字段可能返回 null，表示取不到有效值。
		InvalidType *int64 `json:"InvalidType,omitempty" name:"InvalidType"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetServiceStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetServiceStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ImportKeyMaterialRequest struct {
	*tchttp.BaseRequest

	// 使用GetParametersForImport 返回的PublicKey加密后的密钥材料base64编码。对于国密版本region的KMS，导入的密钥材料长度要求为 128 bit，FIPS版本region的KMS， 导入的密钥材料长度要求为 256 bit。
	EncryptedKeyMaterial *string `json:"EncryptedKeyMaterial,omitempty" name:"EncryptedKeyMaterial"`

	// 通过调用GetParametersForImport获得的导入令牌。
	ImportToken *string `json:"ImportToken,omitempty" name:"ImportToken"`

	// 指定导入密钥材料的CMK，需要和GetParametersForImport 指定的CMK相同。
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 密钥材料过期时间 unix 时间戳，不指定或者 0 表示密钥材料不会过期，若指定过期时间，需要大于当前时间点，最大支持 2147443200。
	ValidTo *uint64 `json:"ValidTo,omitempty" name:"ValidTo"`
}

func (r *ImportKeyMaterialRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ImportKeyMaterialRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ImportKeyMaterialResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ImportKeyMaterialResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ImportKeyMaterialResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Key struct {

	// CMK的全局唯一标识。
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

type KeyMetadata struct {

	// CMK的全局唯一标识
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 作为密钥更容易辨识，更容易被人看懂的别名
	Alias *string `json:"Alias,omitempty" name:"Alias"`

	// 密钥创建时间
	CreateTime *uint64 `json:"CreateTime,omitempty" name:"CreateTime"`

	// CMK的描述
	Description *string `json:"Description,omitempty" name:"Description"`

	// CMK的状态， 取值为：Enabled | Disabled | PendingDelete | PendingImport
	KeyState *string `json:"KeyState,omitempty" name:"KeyState"`

	// CMK用途，取值为: ENCRYPT_DECRYPT | ASYMMETRIC_DECRYPT_RSA_2048 | ASYMMETRIC_DECRYPT_SM2
	KeyUsage *string `json:"KeyUsage,omitempty" name:"KeyUsage"`

	// CMK类型，2 表示符合FIPS标准，4表示符合国密标准
	Type *int64 `json:"Type,omitempty" name:"Type"`

	// 创建者
	CreatorUin *uint64 `json:"CreatorUin,omitempty" name:"CreatorUin"`

	// 是否开启了密钥轮换功能
	KeyRotationEnabled *bool `json:"KeyRotationEnabled,omitempty" name:"KeyRotationEnabled"`

	// CMK的创建者，用户创建的为 user，授权各云产品自动创建的为对应的产品名
	Owner *string `json:"Owner,omitempty" name:"Owner"`

	// 在密钥轮换开启状态下，下次轮换的时间
	NextRotateTime *uint64 `json:"NextRotateTime,omitempty" name:"NextRotateTime"`

	// 计划删除的时间
	// 注意：此字段可能返回 null，表示取不到有效值。
	DeletionDate *uint64 `json:"DeletionDate,omitempty" name:"DeletionDate"`

	// CMK 密钥材料类型，由KMS创建的为： TENCENT_KMS， 由用户导入的类型为：EXTERNAL
	// 注意：此字段可能返回 null，表示取不到有效值。
	Origin *string `json:"Origin,omitempty" name:"Origin"`

	// 在Origin为  EXTERNAL 时有效，表示密钥材料的有效日期， 0 表示不过期
	// 注意：此字段可能返回 null，表示取不到有效值。
	ValidTo *uint64 `json:"ValidTo,omitempty" name:"ValidTo"`
}

type ListAlgorithmsRequest struct {
	*tchttp.BaseRequest
}

func (r *ListAlgorithmsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListAlgorithmsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListAlgorithmsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 本地区支持的对称加密算法
		SymmetricAlgorithms []*AlgorithmInfo `json:"SymmetricAlgorithms,omitempty" name:"SymmetricAlgorithms" list`

		// 本地区支持的非对称加密算法
		AsymmetricAlgorithms []*AlgorithmInfo `json:"AsymmetricAlgorithms,omitempty" name:"AsymmetricAlgorithms" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ListAlgorithmsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListAlgorithmsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListKeyDetailRequest struct {
	*tchttp.BaseRequest

	// 含义跟 SQL 查询的 Offset 一致，表示本次获取从按一定顺序排列数组的第 Offset 个元素开始，缺省为0
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 含义跟 SQL 查询的 Limit 一致，表示本次最多获取 Limit 个元素。缺省值为10，最大值为200
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 根据创建者角色筛选，默认 0 表示用户自己创建的cmk， 1 表示授权其它云产品自动创建的cmk
	Role *uint64 `json:"Role,omitempty" name:"Role"`

	// 根据CMK创建时间排序， 0 表示按照降序排序，1表示按照升序排序
	OrderType *uint64 `json:"OrderType,omitempty" name:"OrderType"`

	// 根据CMK状态筛选， 0表示全部CMK， 1 表示仅查询Enabled CMK， 2 表示仅查询Disabled CMK，3 表示查询PendingDelete 状态的CMK(处于计划删除状态的Key)，4 表示查询 PendingImport 状态的CMK
	KeyState *uint64 `json:"KeyState,omitempty" name:"KeyState"`

	// 根据KeyId或者Alias进行模糊匹配查询
	SearchKeyAlias *string `json:"SearchKeyAlias,omitempty" name:"SearchKeyAlias"`

	// 根据CMK类型筛选， "TENCENT_KMS" 表示筛选密钥材料由KMS创建的CMK， "EXTERNAL" 表示筛选密钥材料需要用户导入的 EXTERNAL类型CMK，"ALL" 或者不设置表示两种类型都查询，大小写敏感。
	Origin *string `json:"Origin,omitempty" name:"Origin"`

	// 根据CMK的KeyUsage筛选，为空表示筛选全部，可使用的参数为：ENCRYPT_DECRYPT 或 ASYMMETRIC_DECRYPT_RSA_2048 或 ASYMMETRIC_DECRYPT_SM2
	KeyUsage *string `json:"KeyUsage,omitempty" name:"KeyUsage"`
}

func (r *ListKeyDetailRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListKeyDetailRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListKeyDetailResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK的总数量
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 返回的属性信息列表。
	// 注意：此字段可能返回 null，表示取不到有效值。
		KeyMetadatas []*KeyMetadata `json:"KeyMetadatas,omitempty" name:"KeyMetadatas" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ListKeyDetailResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListKeyDetailResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListKeysRequest struct {
	*tchttp.BaseRequest

	// 含义跟 SQL 查询的 Offset 一致，表示本次获取从按一定顺序排列数组的第 Offset 个元素开始，缺省为0
	Offset *uint64 `json:"Offset,omitempty" name:"Offset"`

	// 含义跟 SQL 查询的 Limit 一致，表示本次获最多获取 Limit 个元素。缺省值为10，最大值为200
	Limit *uint64 `json:"Limit,omitempty" name:"Limit"`

	// 根据创建者角色筛选，默认 0 表示用户自己创建的cmk， 1 表示授权其它云产品自动创建的cmk
	Role *uint64 `json:"Role,omitempty" name:"Role"`
}

func (r *ListKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListKeysRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// CMK列表数组
	// 注意：此字段可能返回 null，表示取不到有效值。
		Keys []*Key `json:"Keys,omitempty" name:"Keys" list`

		// CMK的总数量
		TotalCount *uint64 `json:"TotalCount,omitempty" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ListKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListKeysResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ReEncryptRequest struct {
	*tchttp.BaseRequest

	// 需要重新加密的密文
	CiphertextBlob *string `json:"CiphertextBlob,omitempty" name:"CiphertextBlob"`

	// 重新加密使用的CMK，如果为空，则使用密文原有的CMK重新加密（若密钥没有轮换则密文不会刷新）
	DestinationKeyId *string `json:"DestinationKeyId,omitempty" name:"DestinationKeyId"`

	// CiphertextBlob 密文加密时使用的key/value对的json字符串。如果加密时未使用，则为空
	SourceEncryptionContext *string `json:"SourceEncryptionContext,omitempty" name:"SourceEncryptionContext"`

	// 重新加密使用的key/value对的json字符串，如果使用该字段，则返回的新密文在解密时需要填入相同的字符串
	DestinationEncryptionContext *string `json:"DestinationEncryptionContext,omitempty" name:"DestinationEncryptionContext"`
}

func (r *ReEncryptRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ReEncryptRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ReEncryptResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 重新加密后的密文
		CiphertextBlob *string `json:"CiphertextBlob,omitempty" name:"CiphertextBlob"`

		// 重新加密使用的CMK
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 重新加密前密文使用的CMK
		SourceKeyId *string `json:"SourceKeyId,omitempty" name:"SourceKeyId"`

		// true表示密文已经重新加密。同一个CMK进行重加密，在密钥没有发生轮换的情况下不会进行实际重新加密操作，返回原密文
		ReEncrypted *bool `json:"ReEncrypted,omitempty" name:"ReEncrypted"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ReEncryptResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ReEncryptResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ScheduleKeyDeletionRequest struct {
	*tchttp.BaseRequest

	// CMK的唯一标志
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// 计划删除时间区间[7,30]
	PendingWindowInDays *uint64 `json:"PendingWindowInDays,omitempty" name:"PendingWindowInDays"`
}

func (r *ScheduleKeyDeletionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ScheduleKeyDeletionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ScheduleKeyDeletionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 计划删除执行时间
		DeletionDate *uint64 `json:"DeletionDate,omitempty" name:"DeletionDate"`

		// 唯一标志被计划删除的CMK
		KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *ScheduleKeyDeletionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ScheduleKeyDeletionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateAliasRequest struct {
	*tchttp.BaseRequest

	// 新的别名，1-60个字符或数字的组合
	Alias *string `json:"Alias,omitempty" name:"Alias"`

	// CMK的全局唯一标识符
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *UpdateAliasRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateAliasRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateAliasResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateAliasResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateAliasResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateKeyDescriptionRequest struct {
	*tchttp.BaseRequest

	// 新的描述信息，最大支持1024字节
	Description *string `json:"Description,omitempty" name:"Description"`

	// 需要修改描述信息的CMK ID
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`
}

func (r *UpdateKeyDescriptionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateKeyDescriptionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateKeyDescriptionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateKeyDescriptionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateKeyDescriptionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
