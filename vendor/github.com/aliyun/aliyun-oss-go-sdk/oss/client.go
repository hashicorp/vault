// Package oss implements functions for access oss service.
// It has two main struct Client and Bucket.
package oss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// Client SDK's entry point. It's for bucket related options such as create/delete/set bucket (such as set/get ACL/lifecycle/referer/logging/website).
// Object related operations are done by Bucket class.
// Users use oss.New to create Client instance.
//
type (
	// Client OSS client
	Client struct {
		Config     *Config      // OSS client configuration
		Conn       *Conn        // Send HTTP request
		HTTPClient *http.Client //http.Client to use - if nil will make its own
	}

	// ClientOption client option such as UseCname, Timeout, SecurityToken.
	ClientOption func(*Client)
)

// New creates a new client.
//
// endpoint    the OSS datacenter endpoint such as http://oss-cn-hangzhou.aliyuncs.com .
// accessKeyId    access key Id.
// accessKeySecret    access key secret.
//
// Client    creates the new client instance, the returned value is valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func New(endpoint, accessKeyID, accessKeySecret string, options ...ClientOption) (*Client, error) {
	// Configuration
	config := getDefaultOssConfig()
	config.Endpoint = endpoint
	config.AccessKeyID = accessKeyID
	config.AccessKeySecret = accessKeySecret

	// URL parse
	url := &urlMaker{}

	// HTTP connect
	conn := &Conn{config: config, url: url}

	// OSS client
	client := &Client{
		Config: config,
		Conn:   conn,
	}

	// Client options parse
	for _, option := range options {
		option(client)
	}

	err := url.InitExt(config.Endpoint, config.IsCname, config.IsUseProxy, config.IsPathStyle)
	if err != nil {
		return nil, err
	}

	if config.AuthVersion != AuthV1 && config.AuthVersion != AuthV2 && config.AuthVersion != AuthV4 {
		return nil, fmt.Errorf("Init client Error, invalid Auth version: %v", config.AuthVersion)
	}

	// Create HTTP connection
	err = conn.init(config, url, client.HTTPClient)

	return client, err
}

// SetRegion set region for client
//
// region    the region, such as cn-hangzhou
func (client *Client) SetRegion(region string) {
	client.Config.Region = region
}

// SetCloudBoxId set CloudBoxId for client
//
// cloudBoxId    the id of cloudBox
func (client *Client) SetCloudBoxId(cloudBoxId string) {
	client.Config.CloudBoxId = cloudBoxId
}

// SetProduct set Product type for client
//
// Product    product type
func (client *Client) SetProduct(product string) {
	client.Config.Product = product
}

// Bucket gets the bucket instance.
//
// bucketName    the bucket name.
// Bucket    the bucket object, when error is nil.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) Bucket(bucketName string) (*Bucket, error) {
	err := CheckBucketName(bucketName)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		client,
		bucketName,
	}, nil
}

// CreateBucket creates a bucket.
//
// bucketName    the bucket name, it's globably unique and immutable. The bucket name can only consist of lowercase letters, numbers and dash ('-').
//               It must start with lowercase letter or number and the length can only be between 3 and 255.
// options    options for creating the bucket, with optional ACL. The ACL could be ACLPrivate, ACLPublicRead, and ACLPublicReadWrite. By default it's ACLPrivate.
//            It could also be specified with StorageClass option, which supports StorageStandard, StorageIA(infrequent access), StorageArchive.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) CreateBucket(bucketName string, options ...Option) error {
	headers := make(map[string]string)
	handleOptions(headers, options)

	buffer := new(bytes.Buffer)

	var cbConfig createBucketConfiguration
	cbConfig.StorageClass = StorageStandard

	isStorageSet, valStroage, _ := IsOptionSet(options, storageClass)
	isRedundancySet, valRedundancy, _ := IsOptionSet(options, redundancyType)
	isObjectHashFuncSet, valHashFunc, _ := IsOptionSet(options, objectHashFunc)
	if isStorageSet {
		cbConfig.StorageClass = valStroage.(StorageClassType)
	}

	if isRedundancySet {
		cbConfig.DataRedundancyType = valRedundancy.(DataRedundancyType)
	}

	if isObjectHashFuncSet {
		cbConfig.ObjectHashFunction = valHashFunc.(ObjecthashFuncType)
	}

	bs, err := xml.Marshal(cbConfig)
	if err != nil {
		return err
	}
	buffer.Write(bs)
	contentType := http.DetectContentType(buffer.Bytes())
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// create bucket xml
func (client Client) CreateBucketXml(bucketName string, xmlBody string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// ListBuckets lists buckets of the current account under the given endpoint, with optional filters.
//
// options    specifies the filters such as Prefix, Marker and MaxKeys. Prefix is the bucket name's prefix filter.
//            And marker makes sure the returned buckets' name are greater than it in lexicographic order.
//            Maxkeys limits the max keys to return, and by default it's 100 and up to 1000.
//            For the common usage scenario, please check out list_bucket.go in the sample.
// ListBucketsResponse    the response object if error is nil.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) ListBuckets(options ...Option) (ListBucketsResult, error) {
	var out ListBucketsResult

	params, err := GetRawParams(options)
	if err != nil {
		return out, err
	}

	resp, err := client.do("GET", "", params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// ListCloudBoxes lists cloud boxes of the current account under the given endpoint, with optional filters.
//
// options    specifies the filters such as Prefix, Marker and MaxKeys. Prefix is the bucket name's prefix filter.
//            And marker makes sure the returned buckets' name are greater than it in lexicographic order.
//            Maxkeys limits the max keys to return, and by default it's 100 and up to 1000.
//            For the common usage scenario, please check out list_bucket.go in the sample.
// ListBucketsResponse    the response object if error is nil.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) ListCloudBoxes(options ...Option) (ListCloudBoxResult, error) {
	var out ListCloudBoxResult

	params, err := GetRawParams(options)
	if err != nil {
		return out, err
	}

	params["cloudboxes"] = nil

	resp, err := client.do("GET", "", params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// IsBucketExist checks if the bucket exists
//
// bucketName    the bucket name.
//
// bool    true if it exists, and it's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) IsBucketExist(bucketName string) (bool, error) {
	listRes, err := client.ListBuckets(Prefix(bucketName), MaxKeys(1))
	if err != nil {
		return false, err
	}

	if len(listRes.Buckets) == 1 && listRes.Buckets[0].Name == bucketName {
		return true, nil
	}
	return false, nil
}

// DeleteBucket deletes the bucket. Only empty bucket can be deleted (no object and parts).
//
// bucketName    the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucket(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// GetBucketLocation gets the bucket location.
//
// Checks out the following link for more information :
// https://www.alibabacloud.com/help/en/object-storage-service/latest/getbucketlocation
//
// bucketName    the bucket name
//
// string    bucket's datacenter location
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketLocation(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["location"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var LocationConstraint string
	err = xmlUnmarshal(resp.Body, &LocationConstraint)
	return LocationConstraint, err
}

// SetBucketACL sets bucket's ACL.
//
// bucketName    the bucket name
// bucketAcl    the bucket ACL: ACLPrivate, ACLPublicRead and ACLPublicReadWrite.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketACL(bucketName string, bucketACL ACLType, options ...Option) error {
	headers := map[string]string{HTTPHeaderOssACL: string(bucketACL)}
	params := map[string]interface{}{}
	params["acl"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketACL gets the bucket ACL.
//
// bucketName    the bucket name.
//
// GetBucketAclResponse    the result object, and it's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketACL(bucketName string, options ...Option) (GetBucketACLResult, error) {
	var out GetBucketACLResult
	params := map[string]interface{}{}
	params["acl"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// SetBucketLifecycle sets the bucket's lifecycle.
//
// For more information, checks out following link:
// https://www.alibabacloud.com/help/en/object-storage-service/latest/putbucketlifecycle
//
// bucketName    the bucket name.
// rules    the lifecycle rules. There're two kind of rules: absolute time expiration and relative time expiration in days and day/month/year respectively.
//          Check out sample/bucket_lifecycle.go for more details.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketLifecycle(bucketName string, rules []LifecycleRule, options ...Option) error {
	err := verifyLifecycleRules(rules)
	if err != nil {
		return err
	}
	lifecycleCfg := LifecycleConfiguration{Rules: rules}
	bs, err := xml.Marshal(lifecycleCfg)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["lifecycle"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// SetBucketLifecycleXml sets the bucket's lifecycle rule from xml config
func (client Client) SetBucketLifecycleXml(bucketName string, xmlBody string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["lifecycle"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// DeleteBucketLifecycle deletes the bucket's lifecycle.
//
//
// bucketName    the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketLifecycle(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["lifecycle"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// GetBucketLifecycle gets the bucket's lifecycle settings.
//
// bucketName    the bucket name.
//
// GetBucketLifecycleResponse    the result object upon successful request. It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketLifecycle(bucketName string, options ...Option) (GetBucketLifecycleResult, error) {
	var out GetBucketLifecycleResult
	params := map[string]interface{}{}
	params["lifecycle"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)

	// NonVersionTransition is not suggested to use
	// to keep compatible
	for k, rule := range out.Rules {
		if len(rule.NonVersionTransitions) > 0 {
			out.Rules[k].NonVersionTransition = &(out.Rules[k].NonVersionTransitions[0])
		}
	}
	return out, err
}

func (client Client) GetBucketLifecycleXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["lifecycle"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// SetBucketReferer sets the bucket's referer whitelist and the flag if allowing empty referrer.
//
// To avoid stealing link on OSS data, OSS supports the HTTP referrer header. A whitelist referrer could be set either by API or web console, as well as
// the allowing empty referrer flag. Note that this applies to requests from web browser only.
// For example, for a bucket os-example and its referrer http://www.aliyun.com, all requests from this URL could access the bucket.
// For more information, please check out this link :
// https://www.alibabacloud.com/help/en/object-storage-service/latest/putbucketreferer
//
// bucketName    the bucket name.
// referrers    the referrer white list. A bucket could have a referrer list and each referrer supports one '*' and multiple '?' as wildcards.
//             The sample could be found in sample/bucket_referer.go
// allowEmptyReferer    the flag of allowing empty referrer. By default it's true.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketReferer(bucketName string, referrers []string, allowEmptyReferer bool, options ...Option) error {
	rxml := RefererXML{}
	rxml.AllowEmptyReferer = allowEmptyReferer
	if referrers == nil {
		rxml.RefererList = append(rxml.RefererList, "")
	} else {
		for _, referrer := range referrers {
			rxml.RefererList = append(rxml.RefererList, referrer)
		}
	}

	bs, err := xml.Marshal(rxml)
	if err != nil {
		return err
	}

	return client.PutBucketRefererXml(bucketName, string(bs), options...)
}

// SetBucketRefererV2 gets the bucket's referer white list.
//
// setBucketReferer   SetBucketReferer bucket referer config in struct format.
//
// GetBucketRefererResponse    the result object upon successful request. It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketRefererV2(bucketName string, setBucketReferer RefererXML, options ...Option) error {
	bs, err := xml.Marshal(setBucketReferer)
	if err != nil {
		return err
	}
	return client.PutBucketRefererXml(bucketName, string(bs), options...)
}

// PutBucketRefererXml set bucket's style
// bucketName    the bucket name.
// xmlData		 the style in xml format
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketRefererXml(bucketName, xmlData string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlData))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["referer"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketReferer gets the bucket's referrer white list.
// bucketName    the bucket name.
// GetBucketRefererResult  the result object upon successful request. It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketReferer(bucketName string, options ...Option) (GetBucketRefererResult, error) {
	var out GetBucketRefererResult
	body, err := client.GetBucketRefererXml(bucketName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// GetBucketRefererXml gets the bucket's referrer white list.
// bucketName    the bucket name.
// GetBucketRefererResponse the bucket referer config result in xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketRefererXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["referer"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

// SetBucketLogging sets the bucket logging settings.
//
// OSS could automatically store the access log. Only the bucket owner could enable the logging.
// Once enabled, OSS would save all the access log into hourly log files in a specified bucket.
// For more information, please check out https://www.alibabacloud.com/help/en/object-storage-service/latest/putbucketlogging
//
// bucketName    bucket name to enable the log.
// targetBucket    the target bucket name to store the log files.
// targetPrefix    the log files' prefix.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketLogging(bucketName, targetBucket, targetPrefix string,
	isEnable bool, options ...Option) error {
	var err error
	var bs []byte
	if isEnable {
		lxml := LoggingXML{}
		lxml.LoggingEnabled.TargetBucket = targetBucket
		lxml.LoggingEnabled.TargetPrefix = targetPrefix
		bs, err = xml.Marshal(lxml)
	} else {
		lxml := loggingXMLEmpty{}
		bs, err = xml.Marshal(lxml)
	}

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["logging"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// DeleteBucketLogging deletes the logging configuration to disable the logging on the bucket.
//
// bucketName    the bucket name to disable the logging.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketLogging(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["logging"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// GetBucketLogging gets the bucket's logging settings
//
// bucketName    the bucket name
// GetBucketLoggingResponse    the result object upon successful request. It's only valid when error is nil.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketLogging(bucketName string, options ...Option) (GetBucketLoggingResult, error) {
	var out GetBucketLoggingResult
	params := map[string]interface{}{}
	params["logging"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// SetBucketWebsite sets the bucket's static website's index and error page.
//
// OSS supports static web site hosting for the bucket data. When the bucket is enabled with that, you can access the file in the bucket like the way to access a static website.
// For more information, please check out: https://www.alibabacloud.com/help/en/object-storage-service/latest/putbucketwebsite
//
// bucketName    the bucket name to enable static web site.
// indexDocument    index page.
// errorDocument    error page.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketWebsite(bucketName, indexDocument, errorDocument string, options ...Option) error {
	wxml := WebsiteXML{}
	wxml.IndexDocument.Suffix = indexDocument
	wxml.ErrorDocument.Key = errorDocument

	bs, err := xml.Marshal(wxml)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["website"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// SetBucketWebsiteDetail sets the bucket's static website's detail
//
// OSS supports static web site hosting for the bucket data. When the bucket is enabled with that, you can access the file in the bucket like the way to access a static website.
// For more information, please check out: https://www.alibabacloud.com/help/en/object-storage-service/latest/putbucketwebsite
//
// bucketName the bucket name to enable static web site.
//
// wxml the website's detail
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketWebsiteDetail(bucketName string, wxml WebsiteXML, options ...Option) error {
	bs, err := xml.Marshal(wxml)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["website"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// SetBucketWebsiteXml sets the bucket's static website's rule
//
// OSS supports static web site hosting for the bucket data. When the bucket is enabled with that, you can access the file in the bucket like the way to access a static website.
// For more information, please check out: https://www.alibabacloud.com/help/en/object-storage-service/latest/putbucketwebsite
//
// bucketName the bucket name to enable static web site.
//
// wxml the website's detail
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketWebsiteXml(bucketName string, webXml string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(webXml))

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["website"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// DeleteBucketWebsite deletes the bucket's static web site settings.
//
// bucketName    the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketWebsite(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["website"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// OpenMetaQuery Enables the metadata management feature for a bucket.
//
// bucketName    the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) OpenMetaQuery(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["metaQuery"] = nil
	params["comp"] = "add"
	resp, err := client.do("POST", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetMetaQueryStatus Queries the information about the metadata index library of a bucket.
//
// bucketName    the bucket name
//
// GetMetaQueryStatusResult    the result object upon successful request. It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetMetaQueryStatus(bucketName string, options ...Option) (GetMetaQueryStatusResult, error) {
	var out GetMetaQueryStatusResult
	params := map[string]interface{}{}
	params["metaQuery"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// DoMetaQuery Queries the objects that meet specified conditions and lists the information about objects based on specified fields and sorting methods.
//
// bucketName   the bucket name
//
// metaQuery    the option of query
//
// DoMetaQueryResult   the result object upon successful request. It's only valid when error is nil.
// error it's nil if no error, otherwise it's an error object.
//
func (client Client) DoMetaQuery(bucketName string, metaQuery MetaQuery, options ...Option) (DoMetaQueryResult, error) {
	var out DoMetaQueryResult
	bs, err := xml.Marshal(metaQuery)
	if err != nil {
		return out, err
	}
	out, err = client.DoMetaQueryXml(bucketName, string(bs), options...)
	return out, err
}

// DoMetaQueryXml Queries the objects that meet specified conditions and lists the information about objects based on specified fields and sorting methods.
//
// bucketName   the bucket name
//
// metaQuery    the option of query
//
// DoMetaQueryResult   the result object upon successful request. It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DoMetaQueryXml(bucketName string, metaQueryXml string, options ...Option) (DoMetaQueryResult, error) {
	var out DoMetaQueryResult
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(metaQueryXml))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["metaQuery"] = nil
	params["comp"] = "query"
	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// CloseMetaQuery Disables the metadata management feature for a bucket.
//
// bucketName    the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) CloseMetaQuery(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["metaQuery"] = nil
	params["comp"] = "delete"
	resp, err := client.do("POST", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketWebsite gets the bucket's default page (index page) and the error page.
//
// bucketName    the bucket name
//
// GetBucketWebsiteResponse    the result object upon successful request. It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketWebsite(bucketName string, options ...Option) (GetBucketWebsiteResult, error) {
	var out GetBucketWebsiteResult
	params := map[string]interface{}{}
	params["website"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// GetBucketWebsiteXml gets the bucket's website config xml config.
//
// bucketName    the bucket name
//
// string   the bucket's xml config, It's only valid when error is nil.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketWebsiteXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["website"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	out := string(body)
	return out, err
}

// SetBucketCORS sets the bucket's CORS rules
//
// For more information, please check out https://help.aliyun.com/document_detail/oss/user_guide/security_management/cors.html
//
// bucketName    the bucket name
// corsRules    the CORS rules to set. The related sample code is in sample/bucket_cors.go.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketCORS(bucketName string, corsRules []CORSRule, options ...Option) error {
	corsxml := CORSXML{}
	for _, v := range corsRules {
		cr := CORSRule{}
		cr.AllowedMethod = v.AllowedMethod
		cr.AllowedOrigin = v.AllowedOrigin
		cr.AllowedHeader = v.AllowedHeader
		cr.ExposeHeader = v.ExposeHeader
		cr.MaxAgeSeconds = v.MaxAgeSeconds
		corsxml.CORSRules = append(corsxml.CORSRules, cr)
	}

	bs, err := xml.Marshal(corsxml)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["cors"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// SetBucketCORSV2 sets the bucket's CORS rules
//
// bucketName    the bucket name
// putBucketCORS    the CORS rules to set.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketCORSV2(bucketName string, putBucketCORS PutBucketCORS, options ...Option) error {
	bs, err := xml.Marshal(putBucketCORS)
	if err != nil {
		return err
	}
	err = client.SetBucketCORSXml(bucketName, string(bs), options...)
	return err
}

func (client Client) SetBucketCORSXml(bucketName string, xmlBody string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["cors"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// DeleteBucketCORS deletes the bucket's static website settings.
//
// bucketName    the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketCORS(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["cors"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// GetBucketCORS gets the bucket's CORS settings.
//
// bucketName    the bucket name.
// GetBucketCORSResult    the result object upon successful request. It's only valid when error is nil.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketCORS(bucketName string, options ...Option) (GetBucketCORSResult, error) {
	var out GetBucketCORSResult
	params := map[string]interface{}{}
	params["cors"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

func (client Client) GetBucketCORSXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["cors"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// GetBucketInfo gets the bucket information.
//
// bucketName    the bucket name.
// GetBucketInfoResult    the result object upon successful request. It's only valid when error is nil.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketInfo(bucketName string, options ...Option) (GetBucketInfoResult, error) {
	var out GetBucketInfoResult
	params := map[string]interface{}{}
	params["bucketInfo"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)

	// convert None to ""
	if err == nil {
		if out.BucketInfo.SseRule.KMSMasterKeyID == "None" {
			out.BucketInfo.SseRule.KMSMasterKeyID = ""
		}

		if out.BucketInfo.SseRule.SSEAlgorithm == "None" {
			out.BucketInfo.SseRule.SSEAlgorithm = ""
		}

		if out.BucketInfo.SseRule.KMSDataEncryption == "None" {
			out.BucketInfo.SseRule.KMSDataEncryption = ""
		}
	}
	return out, err
}

// SetBucketVersioning set bucket versioning:Enabled、Suspended
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) SetBucketVersioning(bucketName string, versioningConfig VersioningConfig, options ...Option) error {
	var err error
	var bs []byte
	bs, err = xml.Marshal(versioningConfig)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["versioning"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketVersioning get bucket versioning status:Enabled、Suspended
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketVersioning(bucketName string, options ...Option) (GetBucketVersioningResult, error) {
	var out GetBucketVersioningResult
	params := map[string]interface{}{}
	params["versioning"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)

	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// SetBucketEncryption set bucket encryption config
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) SetBucketEncryption(bucketName string, encryptionRule ServerEncryptionRule, options ...Option) error {
	var err error
	var bs []byte
	bs, err = xml.Marshal(encryptionRule)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["encryption"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketEncryption get bucket encryption
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketEncryption(bucketName string, options ...Option) (GetBucketEncryptionResult, error) {
	var out GetBucketEncryptionResult
	params := map[string]interface{}{}
	params["encryption"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)

	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// DeleteBucketEncryption delete bucket encryption config
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error bucket
func (client Client) DeleteBucketEncryption(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["encryption"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

//
// SetBucketTagging add tagging to bucket
// bucketName  name of bucket
// tagging    tagging to be added
// error        nil if success, otherwise error
func (client Client) SetBucketTagging(bucketName string, tagging Tagging, options ...Option) error {
	var err error
	var bs []byte
	bs, err = xml.Marshal(tagging)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["tagging"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketTagging get tagging of the bucket
// bucketName  name of bucket
// error      nil if success, otherwise error
func (client Client) GetBucketTagging(bucketName string, options ...Option) (GetBucketTaggingResult, error) {
	var out GetBucketTaggingResult
	params := map[string]interface{}{}
	params["tagging"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

//
// DeleteBucketTagging delete bucket tagging
// bucketName  name of bucket
// error      nil if success, otherwise error
//
func (client Client) DeleteBucketTagging(bucketName string, options ...Option) error {
	key, _ := FindOption(options, "tagging", nil)
	params := map[string]interface{}{}
	if key == nil {
		params["tagging"] = nil
	} else {
		params["tagging"] = key.(string)
	}

	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// GetBucketStat get bucket stat
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketStat(bucketName string, options ...Option) (GetBucketStatResult, error) {
	var out GetBucketStatResult
	params := map[string]interface{}{}
	params["stat"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// GetBucketPolicy API operation for Object Storage Service.
//
// Get the policy from the bucket.
//
// bucketName 	 the bucket name.
//
// string		 return the bucket's policy, and it's only valid when error is nil.
//
// error   		 it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketPolicy(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["policy"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	out := string(body)
	return out, err
}

// SetBucketPolicy API operation for Object Storage Service.
//
// Set the policy from the bucket.
//
// bucketName the bucket name.
//
// policy the bucket policy.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketPolicy(bucketName string, policy string, options ...Option) error {
	params := map[string]interface{}{}
	params["policy"] = nil

	buffer := strings.NewReader(policy)

	resp, err := client.do("PUT", bucketName, params, nil, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// DeleteBucketPolicy API operation for Object Storage Service.
//
// Deletes the policy from the bucket.
//
// bucketName the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketPolicy(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["policy"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// SetBucketRequestPayment API operation for Object Storage Service.
//
// Set the requestPayment of bucket
//
// bucketName the bucket name.
//
// paymentConfig the payment configuration
//
// error it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketRequestPayment(bucketName string, paymentConfig RequestPaymentConfiguration, options ...Option) error {
	params := map[string]interface{}{}
	params["requestPayment"] = nil

	var bs []byte
	bs, err := xml.Marshal(paymentConfig)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketRequestPayment API operation for Object Storage Service.
//
// Get bucket requestPayment
//
// bucketName the bucket name.
//
// RequestPaymentConfiguration the payment configuration
//
// error it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketRequestPayment(bucketName string, options ...Option) (RequestPaymentConfiguration, error) {
	var out RequestPaymentConfiguration
	params := map[string]interface{}{}
	params["requestPayment"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// GetUserQoSInfo API operation for Object Storage Service.
//
// Get user qos.
//
// UserQoSConfiguration the User Qos and range Information.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetUserQoSInfo(options ...Option) (UserQoSConfiguration, error) {
	var out UserQoSConfiguration
	params := map[string]interface{}{}
	params["qosInfo"] = nil

	resp, err := client.do("GET", "", params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// SetBucketQoSInfo API operation for Object Storage Service.
//
// Set Bucket Qos information.
//
// bucketName the bucket name.
//
// qosConf the qos configuration.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketQoSInfo(bucketName string, qosConf BucketQoSConfiguration, options ...Option) error {
	params := map[string]interface{}{}
	params["qosInfo"] = nil

	var bs []byte
	bs, err := xml.Marshal(qosConf)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentTpye := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentTpye

	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketQosInfo API operation for Object Storage Service.
//
// Get Bucket Qos information.
//
// bucketName the bucket name.
//
// BucketQoSConfiguration the  return qos configuration.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketQosInfo(bucketName string, options ...Option) (BucketQoSConfiguration, error) {
	var out BucketQoSConfiguration
	params := map[string]interface{}{}
	params["qosInfo"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// DeleteBucketQosInfo API operation for Object Storage Service.
//
// Delete Bucket QoS information.
//
// bucketName the bucket name.
//
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketQosInfo(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["qosInfo"] = nil

	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// SetBucketInventory API operation for Object Storage Service
//
// Set the Bucket inventory.
//
// bucketName the bucket name.
//
// inventoryConfig the inventory configuration.
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) SetBucketInventory(bucketName string, inventoryConfig InventoryConfiguration, options ...Option) error {
	params := map[string]interface{}{}
	params["inventoryId"] = inventoryConfig.Id
	params["inventory"] = nil

	var bs []byte
	bs, err := xml.Marshal(inventoryConfig)

	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// SetBucketInventoryXml API operation for Object Storage Service
//
// Set the Bucket inventory
//
// bucketName the bucket name.
//
// xmlBody the inventory configuration.
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) SetBucketInventoryXml(bucketName string, xmlBody string, options ...Option) error {
	var inventoryConfig InventoryConfiguration
	err := xml.Unmarshal([]byte(xmlBody), &inventoryConfig)
	if err != nil {
		return err
	}

	if inventoryConfig.Id == "" {
		return fmt.Errorf("inventory id is empty in xml")
	}

	params := map[string]interface{}{}
	params["inventoryId"] = inventoryConfig.Id
	params["inventory"] = nil

	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketInventory API operation for Object Storage Service
//
// Get the Bucket inventory.
//
// bucketName tht bucket name.
//
// strInventoryId the inventory id.
//
// InventoryConfiguration the inventory configuration.
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) GetBucketInventory(bucketName string, strInventoryId string, options ...Option) (InventoryConfiguration, error) {
	var out InventoryConfiguration
	params := map[string]interface{}{}
	params["inventory"] = nil
	params["inventoryId"] = strInventoryId

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// GetBucketInventoryXml API operation for Object Storage Service
//
// Get the Bucket inventory.
//
// bucketName tht bucket name.
//
// strInventoryId the inventory id.
//
// InventoryConfiguration the inventory configuration.
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) GetBucketInventoryXml(bucketName string, strInventoryId string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["inventory"] = nil
	params["inventoryId"] = strInventoryId

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// ListBucketInventory API operation for Object Storage Service
//
// List the Bucket inventory.
//
// bucketName tht bucket name.
//
// continuationToken the users token.
//
// ListInventoryConfigurationsResult list all inventory configuration by .
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) ListBucketInventory(bucketName, continuationToken string, options ...Option) (ListInventoryConfigurationsResult, error) {
	var out ListInventoryConfigurationsResult
	params := map[string]interface{}{}
	params["inventory"] = nil
	if continuationToken == "" {
		params["continuation-token"] = nil
	} else {
		params["continuation-token"] = continuationToken
	}

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// ListBucketInventoryXml API operation for Object Storage Service
//
// List the Bucket inventory.
//
// bucketName tht bucket name.
//
// continuationToken the users token.
//
// ListInventoryConfigurationsResult list all inventory configuration by .
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) ListBucketInventoryXml(bucketName, continuationToken string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["inventory"] = nil
	if continuationToken == "" {
		params["continuation-token"] = nil
	} else {
		params["continuation-token"] = continuationToken
	}

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// DeleteBucketInventory API operation for Object Storage Service.
//
// Delete Bucket inventory information.
//
// bucketName tht bucket name.
//
// strInventoryId the inventory id.
//
// error    it's nil if no error, otherwise it's an error.
//
func (client Client) DeleteBucketInventory(bucketName, strInventoryId string, options ...Option) error {
	params := map[string]interface{}{}
	params["inventory"] = nil
	params["inventoryId"] = strInventoryId

	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// SetBucketAsyncTask API operation for set async fetch task
//
// bucketName tht bucket name.
//
// asynConf  configruation
//
// error  it's nil if success, otherwise it's an error.
func (client Client) SetBucketAsyncTask(bucketName string, asynConf AsyncFetchTaskConfiguration, options ...Option) (AsyncFetchTaskResult, error) {
	var out AsyncFetchTaskResult
	params := map[string]interface{}{}
	params["asyncFetch"] = nil

	var bs []byte
	bs, err := xml.Marshal(asynConf)

	if err != nil {
		return out, err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)

	if err != nil {
		return out, err
	}

	defer resp.Body.Close()
	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// GetBucketAsyncTask API operation for set async fetch task
//
// bucketName tht bucket name.
//
// taskid  returned by SetBucketAsyncTask
//
// error  it's nil if success, otherwise it's an error.
func (client Client) GetBucketAsyncTask(bucketName string, taskID string, options ...Option) (AsynFetchTaskInfo, error) {
	var out AsynFetchTaskInfo
	params := map[string]interface{}{}
	params["asyncFetch"] = nil

	headers := make(map[string]string)
	headers[HTTPHeaderOssTaskID] = taskID
	resp, err := client.do("GET", bucketName, params, headers, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// InitiateBucketWorm creates bucket worm Configuration
// bucketName the bucket name.
// retentionDays the retention period in days
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) InitiateBucketWorm(bucketName string, retentionDays int, options ...Option) (string, error) {
	var initiateWormConf InitiateWormConfiguration
	initiateWormConf.RetentionPeriodInDays = retentionDays

	var respHeader http.Header
	isOptSet, _, _ := IsOptionSet(options, responseHeader)
	if !isOptSet {
		options = append(options, GetResponseHeader(&respHeader))
	}

	bs, err := xml.Marshal(initiateWormConf)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["worm"] = nil

	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respOpt, _ := FindOption(options, responseHeader, nil)
	wormID := ""
	err = CheckRespCode(resp.StatusCode, []int{http.StatusOK})
	if err == nil && respOpt != nil {
		wormID = (respOpt.(*http.Header)).Get("x-oss-worm-id")
	}
	return wormID, err
}

// AbortBucketWorm delete bucket worm Configuration
// bucketName the bucket name.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) AbortBucketWorm(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["worm"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// CompleteBucketWorm complete bucket worm Configuration
// bucketName the bucket name.
// wormID the worm id
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) CompleteBucketWorm(bucketName string, wormID string, options ...Option) error {
	params := map[string]interface{}{}
	params["wormId"] = wormID
	resp, err := client.do("POST", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// ExtendBucketWorm exetend bucket worm Configuration
// bucketName the bucket name.
// retentionDays the retention period in days
// wormID the worm id
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) ExtendBucketWorm(bucketName string, retentionDays int, wormID string, options ...Option) error {
	var extendWormConf ExtendWormConfiguration
	extendWormConf.RetentionPeriodInDays = retentionDays

	bs, err := xml.Marshal(extendWormConf)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["wormId"] = wormID
	params["wormExtend"] = nil

	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketWorm get bucket worm Configuration
// bucketName the bucket name.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketWorm(bucketName string, options ...Option) (WormConfiguration, error) {
	var out WormConfiguration
	params := map[string]interface{}{}
	params["worm"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// SetBucketTransferAcc set bucket transfer acceleration configuration
// bucketName the bucket name.
// accConf bucket transfer acceleration configuration
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) SetBucketTransferAcc(bucketName string, accConf TransferAccConfiguration, options ...Option) error {
	bs, err := xml.Marshal(accConf)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := make(map[string]string)
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["transferAcceleration"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketTransferAcc get bucket transfer acceleration configuration
// bucketName the bucket name.
// accConf bucket transfer acceleration configuration
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketTransferAcc(bucketName string, options ...Option) (TransferAccConfiguration, error) {
	var out TransferAccConfiguration
	params := map[string]interface{}{}
	params["transferAcceleration"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// DeleteBucketTransferAcc delete bucket transfer acceleration configuration
// bucketName the bucket name.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketTransferAcc(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["transferAcceleration"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// PutBucketReplication put bucket replication configuration
// bucketName    the bucket name.
// xmlBody    the replication configuration.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) PutBucketReplication(bucketName string, xmlBody string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["replication"] = nil
	params["comp"] = "add"
	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// PutBucketRTC put bucket replication rtc
// bucketName    the bucket name.
// rtc the bucket rtc config.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) PutBucketRTC(bucketName string, rtc PutBucketRTC, options ...Option) error {
	bs, err := xml.Marshal(rtc)
	if err != nil {
		return err
	}
	err = client.PutBucketRTCXml(bucketName, string(bs), options...)
	return err
}

// PutBucketRTCXml put bucket rtc configuration
// bucketName    the bucket name.
// xmlBody    the rtc configuration in xml format.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) PutBucketRTCXml(bucketName string, xmlBody string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["rtc"] = nil
	resp, err := client.do("PUT", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketReplication get bucket replication configuration
// bucketName    the bucket name.
// string    the replication configuration.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketReplication(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["replication"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// DeleteBucketReplication delete bucket replication configuration
// bucketName    the bucket name.
// ruleId    the ID of the replication configuration.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) DeleteBucketReplication(bucketName string, ruleId string, options ...Option) error {
	replicationxml := ReplicationXML{}
	replicationxml.ID = ruleId

	bs, err := xml.Marshal(replicationxml)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	params := map[string]interface{}{}
	params["replication"] = nil
	params["comp"] = "delete"
	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketReplicationLocation get the locations of the target bucket that can be copied to
// bucketName    the bucket name.
// string    the locations of the target bucket that can be copied to.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketReplicationLocation(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["replicationLocation"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// GetBucketReplicationProgress get the replication progress of bucket
// bucketName    the bucket name.
// ruleId    the ID of the replication configuration.
// string    the replication progress of bucket.
// error    it's nil if no error, otherwise it's an error object.
//
func (client Client) GetBucketReplicationProgress(bucketName string, ruleId string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["replicationProgress"] = nil
	if ruleId != "" {
		params["rule-id"] = ruleId
	}

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// GetBucketAccessMonitor get bucket's access monitor config
// bucketName    the bucket name.
// GetBucketAccessMonitorResult  the access monitor configuration result of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketAccessMonitor(bucketName string, options ...Option) (GetBucketAccessMonitorResult, error) {
	var out GetBucketAccessMonitorResult
	body, err := client.GetBucketAccessMonitorXml(bucketName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// GetBucketAccessMonitorXml get bucket's access monitor config
// bucketName    the bucket name.
// string  the access monitor configuration result of bucket xml foramt.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketAccessMonitorXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["accessmonitor"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// PutBucketAccessMonitor get bucket's access monitor config
// bucketName    the bucket name.
// accessMonitor the access monitor configuration of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketAccessMonitor(bucketName string, accessMonitor PutBucketAccessMonitor, options ...Option) error {
	bs, err := xml.Marshal(accessMonitor)
	if err != nil {
		return err
	}
	err = client.PutBucketAccessMonitorXml(bucketName, string(bs), options...)
	return err
}

// PutBucketAccessMonitorXml get bucket's access monitor config
// bucketName    the bucket name.
// xmlData		 the access monitor configuration in xml foramt
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketAccessMonitorXml(bucketName string, xmlData string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlData))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType
	params := map[string]interface{}{}
	params["accessmonitor"] = nil
	resp, err := client.do("PUT", bucketName, params, nil, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// ListBucketCname list bucket's binding cname
// bucketName    the bucket name.
// string    the xml configuration of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) ListBucketCname(bucketName string, options ...Option) (ListBucketCnameResult, error) {
	var out ListBucketCnameResult
	body, err := client.GetBucketCname(bucketName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// GetBucketCname get bucket's binding cname
// bucketName    the bucket name.
// string    the xml configuration of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketCname(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["cname"] = nil

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// CreateBucketCnameToken create a token for the cname.
// bucketName    the bucket name.
// cname    a custom domain name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) CreateBucketCnameToken(bucketName string, cname string, options ...Option) (CreateBucketCnameTokenResult, error) {
	var out CreateBucketCnameTokenResult
	params := map[string]interface{}{}
	params["cname"] = nil
	params["comp"] = "token"

	rxml := CnameConfigurationXML{}
	rxml.Domain = cname

	bs, err := xml.Marshal(rxml)
	if err != nil {
		return out, err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// GetBucketCnameToken get a token for the cname
// bucketName    the bucket name.
// cname    a custom domain name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketCnameToken(bucketName string, cname string, options ...Option) (GetBucketCnameTokenResult, error) {
	var out GetBucketCnameTokenResult
	params := map[string]interface{}{}
	params["cname"] = cname
	params["comp"] = "token"

	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	err = xmlUnmarshal(resp.Body, &out)
	return out, err
}

// PutBucketCnameXml map a custom domain name to a bucket
// bucketName    the bucket name.
// xmlBody the cname configuration in xml foramt
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketCnameXml(bucketName string, xmlBody string, options ...Option) error {
	params := map[string]interface{}{}
	params["cname"] = nil
	params["comp"] = "add"

	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlBody))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// PutBucketCname map a custom domain name to a bucket
// bucketName    the bucket name.
// cname    a custom domain name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketCname(bucketName string, cname string, options ...Option) error {
	rxml := CnameConfigurationXML{}
	rxml.Domain = cname
	bs, err := xml.Marshal(rxml)
	if err != nil {
		return err
	}
	return client.PutBucketCnameXml(bucketName, string(bs), options...)
}

// PutBucketCnameWithCertificate map a custom domain name to a bucket
// bucketName    the bucket name.
// PutBucketCname    the bucket cname config in struct format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketCnameWithCertificate(bucketName string, putBucketCname PutBucketCname, options ...Option) error {
	bs, err := xml.Marshal(putBucketCname)
	if err != nil {
		return err
	}
	return client.PutBucketCnameXml(bucketName, string(bs), options...)
}

// DeleteBucketCname remove the mapping of the custom domain name from a bucket.
// bucketName    the bucket name.
// cname    a custom domain name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) DeleteBucketCname(bucketName string, cname string, options ...Option) error {
	params := map[string]interface{}{}
	params["cname"] = nil
	params["comp"] = "delete"

	rxml := CnameConfigurationXML{}
	rxml.Domain = cname

	bs, err := xml.Marshal(rxml)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType

	resp, err := client.do("POST", bucketName, params, headers, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// PutBucketResourceGroup set bucket's resource group
// bucketName    the bucket name.
// resourceGroup the resource group configuration of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketResourceGroup(bucketName string, resourceGroup PutBucketResourceGroup, options ...Option) error {
	bs, err := xml.Marshal(resourceGroup)
	if err != nil {
		return err
	}
	err = client.PutBucketResourceGroupXml(bucketName, string(bs), options...)
	return err
}

// PutBucketResourceGroupXml set bucket's resource group
// bucketName    the bucket name.
// xmlData		 the resource group in xml format
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketResourceGroupXml(bucketName string, xmlData string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlData))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType
	params := map[string]interface{}{}
	params["resourceGroup"] = nil
	resp, err := client.do("PUT", bucketName, params, nil, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketResourceGroup get bucket's resource group
// bucketName    the bucket name.
// GetBucketResourceGroupResult  the resource group configuration result of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketResourceGroup(bucketName string, options ...Option) (GetBucketResourceGroupResult, error) {
	var out GetBucketResourceGroupResult
	body, err := client.GetBucketResourceGroupXml(bucketName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// GetBucketResourceGroupXml get bucket's resource group
// bucketName    the bucket name.
// string  the resource group result of bucket xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketResourceGroupXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["resourceGroup"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// PutBucketStyle set bucket's style
// bucketName    the bucket name.
// styleContent the style content.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketStyle(bucketName, styleName string, styleContent string, options ...Option) error {
	bs := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Style><Content>%s</Content></Style>", styleContent)
	err := client.PutBucketStyleXml(bucketName, styleName, bs, options...)
	return err
}

// PutBucketStyleXml set bucket's style
// bucketName    the bucket name.
// styleName the style name.
// xmlData		 the style in xml format
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketStyleXml(bucketName, styleName, xmlData string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlData))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType
	params := map[string]interface{}{}
	params["style"] = nil
	params["styleName"] = styleName
	resp, err := client.do("PUT", bucketName, params, nil, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketStyle get bucket's style
// bucketName    the bucket name.
// styleName the bucket style name.
// GetBucketStyleResult  the style result of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketStyle(bucketName, styleName string, options ...Option) (GetBucketStyleResult, error) {
	var out GetBucketStyleResult
	body, err := client.GetBucketStyleXml(bucketName, styleName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// GetBucketStyleXml get bucket's style
// bucketName    the bucket name.
// styleName the bucket style name.
// string  the style result of bucket in xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketStyleXml(bucketName, styleName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["style"] = nil
	params["styleName"] = styleName
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// ListBucketStyle get bucket's styles
// bucketName    the bucket name.
// GetBucketListStyleResult  the list style result of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) ListBucketStyle(bucketName string, options ...Option) (GetBucketListStyleResult, error) {
	var out GetBucketListStyleResult
	body, err := client.ListBucketStyleXml(bucketName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// ListBucketStyleXml get bucket's list style
// bucketName    the bucket name.
// string  the style result of bucket in xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) ListBucketStyleXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["style"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// DeleteBucketStyle delete bucket's style
// bucketName    the bucket name.
// styleName the bucket style name.
// string  the style result of bucket in xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) DeleteBucketStyle(bucketName, styleName string, options ...Option) error {
	params := map[string]interface{}{}
	params["style"] = bucketName
	params["styleName"] = styleName
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// PutBucketResponseHeader set bucket response header
// bucketName    the bucket name.
// xmlData		 the resource group in xml format
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketResponseHeader(bucketName string, responseHeader PutBucketResponseHeader, options ...Option) error {
	bs, err := xml.Marshal(responseHeader)
	if err != nil {
		return err
	}
	err = client.PutBucketResponseHeaderXml(bucketName, string(bs), options...)
	return err
}

// PutBucketResponseHeaderXml set bucket response header
// bucketName    the bucket name.
// xmlData		 the bucket response header in xml format
// error    it's nil if no error, otherwise it's an error object.
func (client Client) PutBucketResponseHeaderXml(bucketName, xmlData string, options ...Option) error {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(xmlData))
	contentType := http.DetectContentType(buffer.Bytes())
	headers := map[string]string{}
	headers[HTTPHeaderContentType] = contentType
	params := map[string]interface{}{}
	params["responseHeader"] = nil
	resp, err := client.do("PUT", bucketName, params, nil, buffer, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusOK})
}

// GetBucketResponseHeader get bucket's response header.
// bucketName    the bucket name.
// GetBucketResponseHeaderResult  the response header result of bucket.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketResponseHeader(bucketName string, options ...Option) (GetBucketResponseHeaderResult, error) {
	var out GetBucketResponseHeaderResult
	body, err := client.GetBucketResponseHeaderXml(bucketName, options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// GetBucketResponseHeaderXml get bucket's resource group
// bucketName    the bucket name.
// string  the response header result of bucket xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) GetBucketResponseHeaderXml(bucketName string, options ...Option) (string, error) {
	params := map[string]interface{}{}
	params["responseHeader"] = nil
	resp, err := client.do("GET", bucketName, params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// DeleteBucketResponseHeader delete response header from a bucket.
// bucketName    the bucket name.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) DeleteBucketResponseHeader(bucketName string, options ...Option) error {
	params := map[string]interface{}{}
	params["responseHeader"] = nil
	resp, err := client.do("DELETE", bucketName, params, nil, nil, options...)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return CheckRespCode(resp.StatusCode, []int{http.StatusNoContent})
}

// DescribeRegions get describe regions
// GetDescribeRegionsResult  the  result of bucket in xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) DescribeRegions(options ...Option) (DescribeRegionsResult, error) {
	var out DescribeRegionsResult
	body, err := client.DescribeRegionsXml(options...)
	if err != nil {
		return out, err
	}
	err = xmlUnmarshal(strings.NewReader(body), &out)
	return out, err
}

// DescribeRegionsXml get describe regions
// string  the style result of bucket in xml format.
// error    it's nil if no error, otherwise it's an error object.
func (client Client) DescribeRegionsXml(options ...Option) (string, error) {
	params, err := GetRawParams(options)
	if err != nil {
		return "", err
	}
	if params["regions"] == nil {
		params["regions"] = nil
	}
	resp, err := client.do("GET", "", params, nil, nil, options...)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	out := string(body)
	return out, err
}

// LimitUploadSpeed set upload bandwidth limit speed,default is 0,unlimited
// upSpeed KB/s, 0 is unlimited,default is 0
// error it's nil if success, otherwise failure
func (client Client) LimitUploadSpeed(upSpeed int) error {
	if client.Config == nil {
		return fmt.Errorf("client config is nil")
	}
	return client.Config.LimitUploadSpeed(upSpeed)
}

// LimitDownloadSpeed set download bandwidth limit speed,default is 0,unlimited
// downSpeed KB/s, 0 is unlimited,default is 0
// error it's nil if success, otherwise failure
func (client Client) LimitDownloadSpeed(downSpeed int) error {
	if client.Config == nil {
		return fmt.Errorf("client config is nil")
	}
	return client.Config.LimitDownloadSpeed(downSpeed)
}

// UseCname sets the flag of using CName. By default it's false.
//
// isUseCname    true: the endpoint has the CName, false: the endpoint does not have cname. Default is false.
//
func UseCname(isUseCname bool) ClientOption {
	return func(client *Client) {
		client.Config.IsCname = isUseCname
	}
}

// ForcePathStyle sets the flag of using Path Style. By default it's false.
//
// isPathStyle    true: the endpoint has the Path Style, false: the endpoint does not have Path Style. Default is false.
//
func ForcePathStyle(isPathStyle bool) ClientOption {
	return func(client *Client) {
		client.Config.IsPathStyle = isPathStyle
	}
}

// Timeout sets the HTTP timeout in seconds.
//
// connectTimeoutSec    HTTP timeout in seconds. Default is 10 seconds. 0 means infinite (not recommended)
// readWriteTimeout    HTTP read or write's timeout in seconds. Default is 20 seconds. 0 means infinite.
//
func Timeout(connectTimeoutSec, readWriteTimeout int64) ClientOption {
	return func(client *Client) {
		client.Config.HTTPTimeout.ConnectTimeout =
			time.Second * time.Duration(connectTimeoutSec)
		client.Config.HTTPTimeout.ReadWriteTimeout =
			time.Second * time.Duration(readWriteTimeout)
		client.Config.HTTPTimeout.HeaderTimeout =
			time.Second * time.Duration(readWriteTimeout)
		client.Config.HTTPTimeout.IdleConnTimeout =
			time.Second * time.Duration(readWriteTimeout)
		client.Config.HTTPTimeout.LongTimeout =
			time.Second * time.Duration(readWriteTimeout*10)
	}
}

// MaxConns sets the HTTP max connections for a client.
//
// maxIdleConns    controls the maximum number of idle (keep-alive) connections across all hosts. Default is 100.
// maxIdleConnsPerHost    controls the maximum idle (keep-alive) connections to keep per-host. Default is 100.
// maxConnsPerHost    limits the total number of connections per host. Default is no limit.
//
func MaxConns(maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int) ClientOption {
	return func(client *Client) {
		client.Config.HTTPMaxConns.MaxIdleConns = maxIdleConns
		client.Config.HTTPMaxConns.MaxIdleConnsPerHost = maxIdleConnsPerHost
		client.Config.HTTPMaxConns.MaxConnsPerHost = maxConnsPerHost
	}
}

// SecurityToken sets the temporary user's SecurityToken.
//
// token    STS token
//
func SecurityToken(token string) ClientOption {
	return func(client *Client) {
		client.Config.SecurityToken = strings.TrimSpace(token)
	}
}

// EnableMD5 enables MD5 validation.
//
// isEnableMD5    true: enable MD5 validation; false: disable MD5 validation.
//
func EnableMD5(isEnableMD5 bool) ClientOption {
	return func(client *Client) {
		client.Config.IsEnableMD5 = isEnableMD5
	}
}

// MD5ThresholdCalcInMemory sets the memory usage threshold for computing the MD5, default is 16MB.
//
// threshold    the memory threshold in bytes. When the uploaded content is more than 16MB, the temp file is used for computing the MD5.
//
func MD5ThresholdCalcInMemory(threshold int64) ClientOption {
	return func(client *Client) {
		client.Config.MD5Threshold = threshold
	}
}

// EnableCRC enables the CRC checksum. Default is true.
//
// isEnableCRC    true: enable CRC checksum; false: disable the CRC checksum.
//
func EnableCRC(isEnableCRC bool) ClientOption {
	return func(client *Client) {
		client.Config.IsEnableCRC = isEnableCRC
	}
}

// UserAgent specifies UserAgent. The default is aliyun-sdk-go/1.2.0 (windows/-/amd64;go1.5.2).
//
// userAgent    the user agent string.
//
func UserAgent(userAgent string) ClientOption {
	return func(client *Client) {
		client.Config.UserAgent = userAgent
		client.Config.UserSetUa = true
	}
}

// Proxy sets the proxy (optional). The default is not using proxy.
//
// proxyHost    the proxy host in the format "host:port". For example, proxy.com:80 .
//
func Proxy(proxyHost string) ClientOption {
	return func(client *Client) {
		client.Config.IsUseProxy = true
		client.Config.ProxyHost = proxyHost
	}
}

// AuthProxy sets the proxy information with user name and password.
//
// proxyHost    the proxy host in the format "host:port". For example, proxy.com:80 .
// proxyUser    the proxy user name.
// proxyPassword    the proxy password.
//
func AuthProxy(proxyHost, proxyUser, proxyPassword string) ClientOption {
	return func(client *Client) {
		client.Config.IsUseProxy = true
		client.Config.ProxyHost = proxyHost
		client.Config.IsAuthProxy = true
		client.Config.ProxyUser = proxyUser
		client.Config.ProxyPassword = proxyPassword
	}
}

//
// HTTPClient sets the http.Client in use to the one passed in
//
func HTTPClient(HTTPClient *http.Client) ClientOption {
	return func(client *Client) {
		client.HTTPClient = HTTPClient
	}
}

//
// SetLogLevel sets the oss sdk log level
//
func SetLogLevel(LogLevel int) ClientOption {
	return func(client *Client) {
		client.Config.LogLevel = LogLevel
	}
}

//
// SetLogger sets the oss sdk logger
//
func SetLogger(Logger *log.Logger) ClientOption {
	return func(client *Client) {
		client.Config.Logger = Logger
	}
}

// SetCredentialsProvider sets function for get the user's ak
func SetCredentialsProvider(provider CredentialsProvider) ClientOption {
	return func(client *Client) {
		client.Config.CredentialsProvider = provider
	}
}

// SetLocalAddr sets function for local addr
func SetLocalAddr(localAddr net.Addr) ClientOption {
	return func(client *Client) {
		client.Config.LocalAddr = localAddr
	}
}

// AuthVersion  sets auth version: v1 or v2 signature which oss_server needed
func AuthVersion(authVersion AuthVersionType) ClientOption {
	return func(client *Client) {
		client.Config.AuthVersion = authVersion
	}
}

// AdditionalHeaders sets special http headers needed to be signed
func AdditionalHeaders(headers []string) ClientOption {
	return func(client *Client) {
		client.Config.AdditionalHeaders = headers
	}
}

// RedirectEnabled only effective from go1.7 onward,RedirectEnabled set http redirect enabled or not
func RedirectEnabled(enabled bool) ClientOption {
	return func(client *Client) {
		client.Config.RedirectEnabled = enabled
	}
}

// InsecureSkipVerify skip verifying tls certificate file
func InsecureSkipVerify(enabled bool) ClientOption {
	return func(client *Client) {
		client.Config.InsecureSkipVerify = enabled
	}
}

// Region  set region
func Region(region string) ClientOption {
	return func(client *Client) {
		client.Config.Region = region
	}
}

// CloudBoxId  set cloudBox id
func CloudBoxId(cloudBoxId string) ClientOption {
	return func(client *Client) {
		client.Config.CloudBoxId = cloudBoxId
	}
}

// Product  set product type
func Product(product string) ClientOption {
	return func(client *Client) {
		client.Config.Product = product
	}
}

// VerifyObjectStrict  sets the flag of verifying object name strictly.
func VerifyObjectStrict(enable bool) ClientOption {
	return func(client *Client) {
		client.Config.VerifyObjectStrict = enable
	}
}

// Private
func (client Client) do(method, bucketName string, params map[string]interface{},
	headers map[string]string, data io.Reader, options ...Option) (*Response, error) {
	err := CheckBucketName(bucketName)
	if len(bucketName) > 0 && err != nil {
		return nil, err
	}

	// option headers
	addHeaders := make(map[string]string)
	err = handleOptions(addHeaders, options)
	if err != nil {
		return nil, err
	}

	// merge header
	if headers == nil {
		headers = make(map[string]string)
	}

	for k, v := range addHeaders {
		if _, ok := headers[k]; !ok {
			headers[k] = v
		}
	}

	resp, err := client.Conn.Do(method, bucketName, "", params, headers, data, 0, nil)

	// get response header
	respHeader, _ := FindOption(options, responseHeader, nil)
	if respHeader != nil {
		pRespHeader := respHeader.(*http.Header)
		if resp != nil {
			*pRespHeader = resp.Headers
		}
	}

	return resp, err
}
