// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

//lint:file-ignore U1000 Ignore all unused code

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabriel-vasile/mimetype"
)

type (
	cloudType   string
	commandType string
)

const (
	fileProtocol              = "file://"
	dataSizeThreshold int64   = 64 * 1024 * 1024
	injectWaitPut             = 0
	isWindows                 = runtime.GOOS == "windows"
	mb                float64 = 1024.0 * 1024.0
)

const (
	uploadCommand   commandType = "UPLOAD"
	downloadCommand commandType = "DOWNLOAD"
	unknownCommand  commandType = "UNKNOWN"

	putRegexp string = `(?i)^(?:/\*.*\*/\s*)*put\s+`
	getRegexp string = `(?i)^(?:/\*.*\*/\s*)*get\s+`
)

const (
	s3Client    cloudType = "S3"
	azureClient cloudType = "AZURE"
	gcsClient   cloudType = "GCS"
	local       cloudType = "LOCAL_FS"
)

type resultStatus int

const (
	errStatus resultStatus = iota
	uploaded
	downloaded
	skipped
	renewToken
	renewPresignedURL
	notFoundFile
	needRetry
	needRetryWithLowerConcurrency
)

func (rs resultStatus) String() string {
	return [...]string{"ERROR", "UPLOADED", "DOWNLOADED", "SKIPPED",
		"RENEW_TOKEN", "RENEW_PRESIGNED_URL", "NOT_FOUND_FILE", "NEED_RETRY",
		"NEED_RETRY_WITH_LOWER_CONCURRENCY"}[rs]
}

func (rs resultStatus) isSet() bool {
	return uploaded <= rs && rs <= needRetryWithLowerConcurrency
}

// SnowflakeFileTransferOptions enables users to specify options regarding
// files transfers such as PUT/GET
type SnowflakeFileTransferOptions struct {
	showProgressBar    bool
	RaisePutGetError   bool
	MultiPartThreshold int64

	/* streaming PUT */
	compressSourceFromStream bool

	/* streaming GET */
	GetFileToStream bool

	/* PUT */
	putCallback             *snowflakeProgressPercentage
	putAzureCallback        *snowflakeProgressPercentage
	putCallbackOutputStream *io.Writer

	/* GET */
	getCallback             *snowflakeProgressPercentage
	getAzureCallback        *snowflakeProgressPercentage
	getCallbackOutputStream *io.Writer
}

type snowflakeFileTransferAgent struct {
	ctx                         context.Context
	sc                          *snowflakeConn
	data                        *execResponseData
	command                     string
	commandType                 commandType
	stageLocationType           cloudType
	fileMetadata                []*fileMetadata
	encryptionMaterial          []*snowflakeFileEncryption
	stageInfo                   *execResponseStageInfo
	results                     []*fileMetadata
	sourceStream                *bytes.Buffer
	srcLocations                []string
	autoCompress                bool
	srcCompression              string
	parallel                    int64
	overwrite                   bool
	srcFiles                    []string
	localLocation               string
	srcFileToEncryptionMaterial map[string]*snowflakeFileEncryption
	useAccelerateEndpoint       bool
	presignedURLs               []string
	options                     *SnowflakeFileTransferOptions
	streamBuffer                *bytes.Buffer
}

func (sfa *snowflakeFileTransferAgent) execute() error {
	var err error
	if err = sfa.parseCommand(); err != nil {
		return err
	}
	if err = sfa.initFileMetadata(); err != nil {
		return err
	}

	if sfa.commandType == uploadCommand {
		if err = sfa.processFileCompressionType(); err != nil {
			return err
		}
	}

	if err = sfa.transferAccelerateConfig(); err != nil {
		return err
	}

	if sfa.commandType == downloadCommand {
		if _, err = os.Stat(sfa.localLocation); os.IsNotExist(err) {
			if err = os.MkdirAll(sfa.localLocation, os.ModePerm); err != nil {
				return err
			}
		}
	}

	if sfa.stageLocationType == local {
		if _, err = os.Stat(sfa.stageInfo.Location); os.IsNotExist(err) {
			if err = os.MkdirAll(sfa.stageInfo.Location, os.ModePerm); err != nil {
				return err
			}
		}
	}

	if err = sfa.updateFileMetadataWithPresignedURL(); err != nil {
		return err
	}

	smallFileMetas := make([]*fileMetadata, 0)
	largeFileMetas := make([]*fileMetadata, 0)

	for _, meta := range sfa.fileMetadata {
		meta.overwrite = sfa.overwrite
		meta.sfa = sfa
		meta.options = sfa.options
		if sfa.stageLocationType != local {
			sizeThreshold := sfa.options.MultiPartThreshold
			meta.options.MultiPartThreshold = sizeThreshold
			if meta.srcFileSize > sizeThreshold && sfa.commandType == uploadCommand {
				meta.parallel = sfa.parallel
				largeFileMetas = append(largeFileMetas, meta)
			} else {
				meta.parallel = 1
				smallFileMetas = append(smallFileMetas, meta)
			}
		} else {
			meta.parallel = 1
			smallFileMetas = append(smallFileMetas, meta)
		}
	}

	if sfa.commandType == uploadCommand {
		if err = sfa.upload(largeFileMetas, smallFileMetas); err != nil {
			return err
		}
	} else {
		if err = sfa.download(smallFileMetas); err != nil {
			return err
		}
	}

	return nil
}

func (sfa *snowflakeFileTransferAgent) parseCommand() error {
	var err error
	if sfa.data.Command != "" {
		sfa.commandType = commandType(sfa.data.Command)
	} else {
		sfa.commandType = unknownCommand
	}

	sfa.initEncryptionMaterial()
	if len(sfa.data.SrcLocations) == 0 {
		return (&SnowflakeError{
			Number:   ErrInvalidStageLocation,
			SQLState: sfa.data.SQLState,
			QueryID:  sfa.data.QueryID,
			Message:  "failed to parse location",
		}).exceptionTelemetry(sfa.sc)
	}
	sfa.srcLocations = sfa.data.SrcLocations

	if sfa.commandType == uploadCommand {
		if sfa.sourceStream != nil {
			sfa.srcFiles = sfa.srcLocations // streaming PUT
		} else {
			sfa.srcFiles, err = sfa.expandFilenames(sfa.srcLocations)
			if err != nil {
				return err
			}
		}
		sfa.autoCompress = sfa.data.AutoCompress
		sfa.srcCompression = strings.ToLower(sfa.data.SourceCompression)
	} else {
		sfa.srcFiles = sfa.srcLocations
		sfa.srcFileToEncryptionMaterial = make(map[string]*snowflakeFileEncryption)
		if len(sfa.data.SrcLocations) == len(sfa.encryptionMaterial) {
			for i, srcFile := range sfa.srcFiles {
				sfa.srcFileToEncryptionMaterial[srcFile] = sfa.encryptionMaterial[i]
			}
		} else if len(sfa.encryptionMaterial) != 0 {
			return (&SnowflakeError{
				Number:      ErrInternalNotMatchEncryptMaterial,
				SQLState:    sfa.data.SQLState,
				QueryID:     sfa.data.QueryID,
				Message:     errMsgInternalNotMatchEncryptMaterial,
				MessageArgs: []interface{}{len(sfa.data.SrcLocations), len(sfa.encryptionMaterial)},
			}).exceptionTelemetry(sfa.sc)
		}

		sfa.localLocation, err = expandUser(sfa.data.LocalLocation)
		if err != nil {
			return err
		}
		if fi, err := os.Stat(sfa.localLocation); err != nil || !fi.IsDir() {
			return (&SnowflakeError{
				Number:      ErrLocalPathNotDirectory,
				SQLState:    sfa.data.SQLState,
				QueryID:     sfa.data.QueryID,
				Message:     errMsgLocalPathNotDirectory,
				MessageArgs: []interface{}{sfa.localLocation},
			}).exceptionTelemetry(sfa.sc)
		}
	}

	sfa.parallel = 1
	if sfa.data.Parallel != 0 {
		sfa.parallel = sfa.data.Parallel
	}
	sfa.overwrite = sfa.data.Overwrite
	sfa.stageLocationType = cloudType(strings.ToUpper(sfa.data.StageInfo.LocationType))
	sfa.stageInfo = &sfa.data.StageInfo
	sfa.presignedURLs = make([]string, 0)
	if len(sfa.data.PresignedURLs) != 0 {
		sfa.presignedURLs = sfa.data.PresignedURLs
	}

	if sfa.getStorageClient(sfa.stageLocationType) == nil {
		return (&SnowflakeError{
			Number:      ErrInvalidStageFs,
			SQLState:    sfa.data.SQLState,
			QueryID:     sfa.data.QueryID,
			Message:     errMsgInvalidStageFs,
			MessageArgs: []interface{}{sfa.stageLocationType},
		}).exceptionTelemetry(sfa.sc)
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) initEncryptionMaterial() {
	sfa.encryptionMaterial = make([]*snowflakeFileEncryption, 0)
	wrapper := sfa.data.EncryptionMaterial

	if sfa.commandType == uploadCommand {
		if wrapper.QueryID != "" {
			sfa.encryptionMaterial = append(sfa.encryptionMaterial, &wrapper.snowflakeFileEncryption)
		}
	} else {
		for _, encmat := range wrapper.EncryptionMaterials {
			if encmat.QueryID != "" {
				sfa.encryptionMaterial = append(sfa.encryptionMaterial, &encmat)
			}
		}
	}
}

func (sfa *snowflakeFileTransferAgent) expandFilenames(locations []string) ([]string, error) {
	canonicalLocations := make([]string, 0)
	for _, fileName := range locations {
		if sfa.commandType == uploadCommand {
			var err error
			fileName, err = expandUser(fileName)
			if err != nil {
				return []string{}, err
			}
			if !filepath.IsAbs(fileName) {
				cwd, err := getDirectory()
				if err != nil {
					return []string{}, err
				}
				fileName = filepath.Join(cwd, fileName)
			}
			if isWindows && len(fileName) > 2 && fileName[0] == '/' && fileName[2] == ':' {
				// Windows path: /C:/data/file1.txt where it starts with slash
				// followed by a drive letter and colon.
				fileName = fileName[1:]
			}
			files, err := filepath.Glob(fileName)
			if err != nil {
				return []string{}, err
			}
			canonicalLocations = append(canonicalLocations, files...)
		} else {
			canonicalLocations = append(canonicalLocations, fileName)
		}
	}
	return canonicalLocations, nil
}

func (sfa *snowflakeFileTransferAgent) initFileMetadata() error {
	sfa.fileMetadata = []*fileMetadata{}
	if sfa.commandType == uploadCommand {
		if len(sfa.srcFiles) == 0 {
			fileName := sfa.data.SrcLocations
			return (&SnowflakeError{
				Number:      ErrFileNotExists,
				SQLState:    sfa.data.SQLState,
				QueryID:     sfa.data.QueryID,
				Message:     errMsgFileNotExists,
				MessageArgs: []interface{}{fileName},
			}).exceptionTelemetry(sfa.sc)
		}
		if sfa.sourceStream != nil {
			fileName := sfa.srcFiles[0]
			srcFileSize := int64(sfa.sourceStream.Len())
			sfa.fileMetadata = append(sfa.fileMetadata, &fileMetadata{
				name:              baseName(fileName),
				srcFileName:       fileName,
				srcStream:         sfa.sourceStream,
				srcFileSize:       srcFileSize,
				stageLocationType: sfa.stageLocationType,
				stageInfo:         sfa.stageInfo,
			})
		} else {
			for i, fileName := range sfa.srcFiles {
				fi, err := os.Stat(fileName)
				if os.IsNotExist(err) {
					return (&SnowflakeError{
						Number:      ErrFileNotExists,
						SQLState:    sfa.data.SQLState,
						QueryID:     sfa.data.QueryID,
						Message:     errMsgFileNotExists,
						MessageArgs: []interface{}{fileName},
					}).exceptionTelemetry(sfa.sc)
				} else if fi.IsDir() {
					return (&SnowflakeError{
						Number:      ErrFileNotExists,
						SQLState:    sfa.data.SQLState,
						QueryID:     sfa.data.QueryID,
						Message:     errMsgFileNotExists,
						MessageArgs: []interface{}{fileName},
					}).exceptionTelemetry(sfa.sc)
				}
				sfa.fileMetadata = append(sfa.fileMetadata, &fileMetadata{
					name:              baseName(fileName),
					srcFileName:       fileName,
					srcFileSize:       fi.Size(),
					stageLocationType: sfa.stageLocationType,
					stageInfo:         sfa.stageInfo,
				})
				if len(sfa.encryptionMaterial) > 0 {
					sfa.fileMetadata[i].encryptionMaterial = sfa.encryptionMaterial[0]
				}
			}
		}

		if len(sfa.encryptionMaterial) > 0 {
			for _, meta := range sfa.fileMetadata {
				meta.encryptionMaterial = sfa.encryptionMaterial[0]
			}
		}
	} else if sfa.commandType == downloadCommand {
		for _, fileName := range sfa.srcFiles {
			if len(fileName) > 0 {
				firstPathSep := strings.Index(fileName, "/")
				dstFileName := fileName
				if firstPathSep >= 0 {
					dstFileName = fileName[firstPathSep+1:]
				}
				sfa.fileMetadata = append(sfa.fileMetadata, &fileMetadata{
					name:              baseName(fileName),
					srcFileName:       fileName,
					dstFileName:       dstFileName,
					dstStream:         new(bytes.Buffer),
					stageLocationType: sfa.stageLocationType,
					stageInfo:         sfa.stageInfo,
					localLocation:     sfa.localLocation,
				})
			}
		}
		// TODO is this necessary?
		for _, meta := range sfa.fileMetadata {
			fileName := meta.srcFileName
			if val, ok := sfa.srcFileToEncryptionMaterial[fileName]; ok {
				meta.encryptionMaterial = val
			}
		}
	}

	return nil
}

func (sfa *snowflakeFileTransferAgent) processFileCompressionType() error {
	var userSpecifiedSourceCompression *compressionType
	var autoDetect bool
	if sfa.srcCompression == "auto_detect" {
		autoDetect = true
	} else if sfa.srcCompression == "none" {
		autoDetect = false
	} else {
		userSpecifiedSourceCompression = lookupByMimeSubType(sfa.srcCompression)
		if userSpecifiedSourceCompression == nil || !userSpecifiedSourceCompression.isSupported {
			return (&SnowflakeError{
				Number:      ErrCompressionNotSupported,
				SQLState:    sfa.data.SQLState,
				QueryID:     sfa.data.QueryID,
				Message:     errMsgFeatureNotSupported,
				MessageArgs: []interface{}{userSpecifiedSourceCompression},
			}).exceptionTelemetry(sfa.sc)
		}
		autoDetect = false
	}

	gzipCompression := compressionTypes["GZIP"]
	for _, meta := range sfa.fileMetadata {
		fileName := meta.srcFileName
		var currentFileCompressionType *compressionType
		if autoDetect {
			currentFileCompressionType = lookupByExtension(filepath.Ext(fileName))
			if currentFileCompressionType == nil {
				var mtype *mimetype.MIME
				var err error
				if meta.srcStream != nil {
					r := getReaderFromBuffer(&meta.srcStream)
					mtype, err = mimetype.DetectReader(r)
					if err != nil {
						return err
					}
					io.ReadAll(r) // flush out tee buffer
				} else {
					mtype, err = mimetype.DetectFile(fileName)
					if err != nil {
						return err
					}
				}
				currentFileCompressionType = lookupByExtension(mtype.Extension())
			}

			if currentFileCompressionType != nil && !currentFileCompressionType.isSupported {
				return (&SnowflakeError{
					Number:      ErrCompressionNotSupported,
					SQLState:    sfa.data.SQLState,
					QueryID:     sfa.data.QueryID,
					Message:     errMsgFeatureNotSupported,
					MessageArgs: []interface{}{userSpecifiedSourceCompression},
				}).exceptionTelemetry(sfa.sc)
			}
		} else {
			currentFileCompressionType = userSpecifiedSourceCompression
		}

		if currentFileCompressionType != nil {
			meta.srcCompressionType = currentFileCompressionType
			if currentFileCompressionType.isSupported {
				meta.dstCompressionType = currentFileCompressionType
				meta.requireCompress = false
				meta.dstFileName = meta.name
			} else {
				return (&SnowflakeError{
					Number:      ErrCompressionNotSupported,
					SQLState:    sfa.data.SQLState,
					QueryID:     sfa.data.QueryID,
					Message:     errMsgFeatureNotSupported,
					MessageArgs: []interface{}{userSpecifiedSourceCompression},
				}).exceptionTelemetry(sfa.sc)
			}
		} else {
			meta.requireCompress = sfa.autoCompress
			meta.srcCompressionType = nil
			if sfa.autoCompress {
				dstFileName := meta.name + compressionTypes["GZIP"].fileExtension
				meta.dstFileName = dstFileName
				meta.dstCompressionType = gzipCompression
			} else {
				meta.dstFileName = meta.name
				meta.dstCompressionType = nil
			}
		}
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) updateFileMetadataWithPresignedURL() error {
	// presigned URL only applies to GCS
	if sfa.stageLocationType == gcsClient {
		if sfa.commandType == uploadCommand {
			filePathToBeReplaced := sfa.getLocalFilePathFromCommand(sfa.command)
			for _, meta := range sfa.fileMetadata {
				filePathToBeReplacedWith := strings.TrimRight(filePathToBeReplaced, meta.dstFileName) + meta.dstFileName
				commandWithSingleFile := strings.ReplaceAll(sfa.command, filePathToBeReplaced, filePathToBeReplacedWith)
				req := execRequest{
					SQLText: commandWithSingleFile,
				}
				headers := getHeaders()
				headers[httpHeaderAccept] = headerContentTypeApplicationJSON
				jsonBody, err := json.Marshal(req)
				if err != nil {
					return err
				}
				data, err := sfa.sc.rest.FuncPostQuery(
					sfa.sc.ctx,
					sfa.sc.rest,
					&url.Values{},
					headers,
					jsonBody,
					sfa.sc.rest.RequestTimeout,
					getOrGenerateRequestIDFromContext(sfa.sc.ctx),
					sfa.sc.cfg)
				if err != nil {
					return err
				}

				if data.Data.StageInfo != (execResponseStageInfo{}) {
					meta.stageInfo = &data.Data.StageInfo
					meta.presignedURL = nil
					if meta.stageInfo.PresignedURL != "" {
						meta.presignedURL, err = url.Parse(meta.stageInfo.PresignedURL)
						if err != nil {
							return err
						}
					}
				}
			}
		} else if sfa.commandType == downloadCommand {
			for i, meta := range sfa.fileMetadata {
				if len(sfa.presignedURLs) > 0 {
					var err error
					meta.presignedURL, err = url.Parse(sfa.presignedURLs[i])
					if err != nil {
						return err
					}
				} else {
					meta.presignedURL = nil
				}
			}
		} else {
			return (&SnowflakeError{
				Number:      ErrCommandNotRecognized,
				SQLState:    sfa.data.SQLState,
				QueryID:     sfa.data.QueryID,
				Message:     errMsgCommandNotRecognized,
				MessageArgs: []interface{}{sfa.commandType},
			}).exceptionTelemetry(sfa.sc)
		}
	}
	return nil
}

type s3BucketAccelerateConfigGetter interface {
	GetBucketAccelerateConfiguration(ctx context.Context, params *s3.GetBucketAccelerateConfigurationInput, optFns ...func(*s3.Options)) (*s3.GetBucketAccelerateConfigurationOutput, error)
}

type s3ClientCreator interface {
	extractBucketNameAndPath(location string) (*s3Location, error)
	createClient(info *execResponseStageInfo, useAccelerateEndpoint bool) (cloudClient, error)
}

func (sfa *snowflakeFileTransferAgent) transferAccelerateConfigWithUtil(s3Util s3ClientCreator) error {
	s3Loc, err := s3Util.extractBucketNameAndPath(sfa.stageInfo.Location)
	if err != nil {
		return err
	}
	s3Cli, err := s3Util.createClient(sfa.stageInfo, false)
	if err != nil {
		return err
	}
	client, ok := s3Cli.(s3BucketAccelerateConfigGetter)
	if !ok {
		return (&SnowflakeError{
			Number:   ErrFailedToConvertToS3Client,
			SQLState: sfa.data.SQLState,
			QueryID:  sfa.data.QueryID,
			Message:  errMsgFailedToConvertToS3Client,
		}).exceptionTelemetry(sfa.sc)
	}
	ret, err := client.GetBucketAccelerateConfiguration(context.Background(), &s3.GetBucketAccelerateConfigurationInput{
		Bucket: &s3Loc.bucketName,
	})
	sfa.useAccelerateEndpoint = ret != nil && ret.Status == "Enabled"
	if err != nil {
		logger.WithContext(sfa.sc.ctx).Warnln("An error occurred when getting accelerate config:", err)
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) transferAccelerateConfig() error {
	if sfa.stageLocationType == s3Client {
		s3Util := new(snowflakeS3Client)
		return sfa.transferAccelerateConfigWithUtil(s3Util)
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) getLocalFilePathFromCommand(command string) string {
	if len(command) == 0 || !strings.Contains(command, fileProtocol) {
		return ""
	}
	if !regexp.MustCompile(putRegexp).Match([]byte(command)) {
		return ""
	}

	filePathBeginIdx := strings.Index(command, fileProtocol)
	isFilePathQuoted := command[filePathBeginIdx-1] == '\''
	filePathBeginIdx += len(fileProtocol)
	var filePathEndIdx int
	filePath := ""

	if isFilePathQuoted {
		filePathEndIdx = filePathBeginIdx + strings.Index(command[filePathBeginIdx:], "'")
		if filePathEndIdx > filePathBeginIdx {
			filePath = command[filePathBeginIdx:filePathEndIdx]
		}
	} else {
		indexList := make([]int, 0)
		delims := []rune{' ', '\n', ';'}
		for _, delim := range delims {
			index := strings.Index(command[filePathBeginIdx:], string(delim))
			if index != -1 {
				indexList = append(indexList, index)
			}
		}
		filePathEndIdx = -1
		if getMin(indexList) != -1 {
			filePathEndIdx = filePathBeginIdx + getMin(indexList)
		}
		if filePathEndIdx > filePathBeginIdx {
			filePath = command[filePathBeginIdx:filePathEndIdx]
		} else {
			filePath = command[filePathBeginIdx:]
		}
	}
	return filePath
}

func (sfa *snowflakeFileTransferAgent) upload(
	largeFileMetadata []*fileMetadata,
	smallFileMetadata []*fileMetadata) error {
	client, err := sfa.getStorageClient(sfa.stageLocationType).
		createClient(sfa.stageInfo, sfa.useAccelerateEndpoint)
	if err != nil {
		return err
	}
	for _, meta := range smallFileMetadata {
		meta.client = client
	}
	for _, meta := range largeFileMetadata {
		meta.client = client
	}

	if len(smallFileMetadata) > 0 {
		logger.WithContext(sfa.sc.ctx).Infof("uploading %v small files", len(smallFileMetadata))
		if err = sfa.uploadFilesParallel(smallFileMetadata); err != nil {
			return err
		}
	}
	if len(largeFileMetadata) > 0 {
		logger.WithContext(sfa.sc.ctx).Infof("uploading %v large files", len(largeFileMetadata))
		if err = sfa.uploadFilesSequential(largeFileMetadata); err != nil {
			return err
		}
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) download(
	fileMetadata []*fileMetadata) error {
	client, err := sfa.getStorageClient(sfa.stageLocationType).
		createClient(sfa.stageInfo, sfa.useAccelerateEndpoint)
	if err != nil {
		return err
	}
	for _, meta := range fileMetadata {
		meta.client = client
	}

	logger.WithContext(sfa.sc.ctx).Infof("downloading %v files", len(fileMetadata))
	if err = sfa.downloadFilesParallel(fileMetadata); err != nil {
		return err
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) uploadFilesParallel(fileMetas []*fileMetadata) error {
	idx := 0
	fileMetaLen := len(fileMetas)
	var err error
	for idx < fileMetaLen {
		endOfIdx := intMin(fileMetaLen, idx+int(sfa.parallel))
		targetMeta := fileMetas[idx:endOfIdx]
		for {
			var wg sync.WaitGroup
			results := make([]*fileMetadata, len(targetMeta))
			errors := make([]error, len(targetMeta))
			for i, meta := range targetMeta {
				wg.Add(1)
				go func(k int, m *fileMetadata) {
					defer wg.Done()
					results[k], errors[k] = sfa.uploadOneFile(m)
				}(i, meta)
			}
			wg.Wait()

			// append errors with no result associated to separate array
			var errorMessages []string
			for i, result := range results {
				if result == nil {
					if errors[i] == nil {
						errorMessages = append(errorMessages, "unknown error")
					} else {
						errorMessages = append(errorMessages, errors[i].Error())
					}
				}
			}
			if errorMessages != nil {
				// sort the error messages to be more deterministic as the goroutines may finish in different order each time
				sort.Strings(errorMessages)
				return fmt.Errorf("errors during file upload:\n%v", strings.Join(errorMessages, "\n"))
			}

			retryMeta := make([]*fileMetadata, 0)
			for i, result := range results {
				result.errorDetails = errors[i]
				if result.resStatus == renewToken || result.resStatus == renewPresignedURL {
					retryMeta = append(retryMeta, result)
				} else {
					sfa.results = append(sfa.results, result)
				}
			}
			if len(retryMeta) == 0 {
				break
			}

			needRenewToken := false
			for _, result := range retryMeta {
				if result.resStatus == renewToken {
					needRenewToken = true
				}
			}

			if needRenewToken {
				client, err := sfa.renewExpiredClient()
				if err != nil {
					return err
				}
				for _, result := range retryMeta {
					result.client = client
				}
				if endOfIdx < fileMetaLen {
					for i := idx + int(sfa.parallel); i < fileMetaLen; i++ {
						fileMetas[i].client = client
					}
				}
			}

			for _, result := range retryMeta {
				if result.resStatus == renewPresignedURL {
					sfa.updateFileMetadataWithPresignedURL()
					break
				}
			}
			targetMeta = retryMeta
		}
		if endOfIdx == fileMetaLen {
			break
		}
		idx += int(sfa.parallel)
	}
	return err
}

func (sfa *snowflakeFileTransferAgent) uploadFilesSequential(fileMetas []*fileMetadata) error {
	idx := 0
	fileMetaLen := len(fileMetas)
	for idx < fileMetaLen {
		res, err := sfa.uploadOneFile(fileMetas[idx])
		if err != nil {
			return err
		}

		if res.resStatus == renewToken {
			client, err := sfa.renewExpiredClient()
			if err != nil {
				return err
			}
			for i := idx; i < fileMetaLen; i++ {
				fileMetas[i].client = client
			}
			continue
		} else if res.resStatus == renewPresignedURL {
			sfa.updateFileMetadataWithPresignedURL()
			continue
		}

		sfa.results = append(sfa.results, res)
		idx++
		if injectWaitPut > 0 {
			time.Sleep(injectWaitPut)
		}
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) uploadOneFile(meta *fileMetadata) (*fileMetadata, error) {
	meta.realSrcFileName = meta.srcFileName
	tmpDir, err := os.MkdirTemp(sfa.sc.cfg.TmpDirPath, "")
	if err != nil {
		return nil, err
	}
	meta.tmpDir = tmpDir
	defer os.RemoveAll(tmpDir) // cleanup

	fileUtil := new(snowflakeFileUtil)
	if meta.requireCompress {
		if meta.srcStream != nil {
			meta.realSrcStream, _, err = fileUtil.compressFileWithGzipFromStream(&meta.srcStream)
		} else {
			meta.realSrcFileName, _, err = fileUtil.compressFileWithGzip(meta.srcFileName, tmpDir)
		}
		if err != nil {
			return nil, err
		}
	}

	if meta.srcStream != nil {
		if meta.realSrcStream != nil {
			meta.sha256Digest, meta.uploadSize, err = fileUtil.getDigestAndSizeForStream(&meta.realSrcStream)
		} else {
			meta.sha256Digest, meta.uploadSize, err = fileUtil.getDigestAndSizeForStream(&meta.srcStream)
		}
	} else {
		meta.sha256Digest, meta.uploadSize, err = fileUtil.getDigestAndSizeForFile(meta.realSrcFileName)
	}
	if err != nil {
		return meta, err
	}

	client := sfa.getStorageClient(sfa.stageLocationType)
	if err = client.uploadOneFileWithRetry(meta); err != nil {
		return meta, err
	}
	return meta, nil
}

func (sfa *snowflakeFileTransferAgent) downloadFilesParallel(fileMetas []*fileMetadata) error {
	idx := 0
	fileMetaLen := len(fileMetas)
	var err error
	for idx < fileMetaLen {
		endOfIdx := intMin(fileMetaLen, idx+int(sfa.parallel))
		targetMeta := fileMetas[idx:endOfIdx]
		for {
			var wg sync.WaitGroup
			results := make([]*fileMetadata, len(targetMeta))
			errors := make([]error, len(targetMeta))
			for i, meta := range targetMeta {
				wg.Add(1)
				go func(k int, m *fileMetadata) {
					defer wg.Done()
					results[k], errors[k] = sfa.downloadOneFile(m)
				}(i, meta)
			}
			wg.Wait()

			retryMeta := make([]*fileMetadata, 0)
			for i, result := range results {
				result.errorDetails = errors[i]
				if result.resStatus == renewToken || result.resStatus == renewPresignedURL {
					retryMeta = append(retryMeta, result)
				} else {
					sfa.results = append(sfa.results, result)
				}
			}
			if len(retryMeta) == 0 {
				break
			}
			logger.WithContext(sfa.sc.ctx).Infof("%v retries found", len(retryMeta))

			needRenewToken := false
			for _, result := range retryMeta {
				if result.resStatus == renewToken {
					needRenewToken = true
				}
				logger.WithContext(sfa.sc.ctx).Infof(
					"retying download file %v with status %v",
					result.name, result.resStatus)
			}

			if needRenewToken {
				client, err := sfa.renewExpiredClient()
				if err != nil {
					return err
				}
				for _, result := range retryMeta {
					result.client = client
				}
				if endOfIdx < fileMetaLen {
					for i := idx + int(sfa.parallel); i < fileMetaLen; i++ {
						fileMetas[i].client = client
					}
				}
			}

			for _, result := range retryMeta {
				if result.resStatus == renewPresignedURL {
					sfa.updateFileMetadataWithPresignedURL()
					break
				}
			}
			targetMeta = retryMeta
		}
		if endOfIdx == fileMetaLen {
			break
		}
		idx += int(sfa.parallel)
	}
	return err
}

func (sfa *snowflakeFileTransferAgent) downloadOneFile(meta *fileMetadata) (*fileMetadata, error) {
	tmpDir, err := os.MkdirTemp(sfa.sc.cfg.TmpDirPath, "")
	if err != nil {
		return nil, err
	}
	meta.tmpDir = tmpDir
	defer os.RemoveAll(tmpDir) // cleanup
	client := sfa.getStorageClient(sfa.stageLocationType)
	if err = client.downloadOneFile(meta); err != nil {
		meta.dstFileSize = -1
		if !meta.resStatus.isSet() {
			meta.resStatus = errStatus
		}
		meta.errorDetails = fmt.Errorf(err.Error() + ", file=" + meta.dstFileName)
		return meta, err
	}
	return meta, nil
}

func (sfa *snowflakeFileTransferAgent) getStorageClient(stageLocationType cloudType) storageUtil {
	if stageLocationType == local {
		return &localUtil{}
	} else if stageLocationType == s3Client || stageLocationType == azureClient || stageLocationType == gcsClient {
		return &remoteStorageUtil{}
	}
	return nil
}

func (sfa *snowflakeFileTransferAgent) renewExpiredClient() (cloudClient, error) {
	data, err := sfa.sc.exec(
		sfa.sc.ctx,
		sfa.command,
		false,
		false,
		false,
		[]driver.NamedValue{})
	if err != nil {
		return nil, err
	}
	storageClient := sfa.getStorageClient(sfa.stageLocationType)
	return storageClient.createClient(&data.Data.StageInfo, sfa.useAccelerateEndpoint)
}

func (sfa *snowflakeFileTransferAgent) result() (*execResponse, error) {
	// inherit old response data
	data := sfa.data
	rowset := make([]fileTransferResultType, 0)
	if sfa.commandType == uploadCommand {
		if len(sfa.results) > 0 {
			for _, meta := range sfa.results {
				var srcCompressionType, dstCompressionType *compressionType
				if meta.srcCompressionType != nil {
					srcCompressionType = meta.srcCompressionType
				} else {
					srcCompressionType = &compressionType{
						name: "NONE",
					}
				}
				if meta.dstCompressionType != nil {
					dstCompressionType = meta.dstCompressionType
				} else {
					dstCompressionType = &compressionType{
						name: "NONE",
					}
				}
				errorDetails := meta.errorDetails
				srcFileSize := meta.srcFileSize
				dstFileSize := meta.dstFileSize
				if sfa.options.RaisePutGetError && errorDetails != nil {
					return nil, (&SnowflakeError{
						Number:   ErrFailedToUploadToStage,
						SQLState: sfa.data.SQLState,
						QueryID:  sfa.data.QueryID,
						Message:  errorDetails.Error(),
					}).exceptionTelemetry(sfa.sc)
				}
				rowset = append(rowset, fileTransferResultType{
					meta.name,
					meta.srcFileName,
					meta.dstFileName,
					srcFileSize,
					dstFileSize,
					srcCompressionType,
					dstCompressionType,
					meta.resStatus,
					meta.errorDetails,
				})
			}
			sort.Slice(rowset, func(i, j int) bool {
				return rowset[i].srcFileName < rowset[j].srcFileName
			})
			ccrs := make([][]*string, 0, len(rowset))
			for _, rs := range rowset {
				srcFileSize := fmt.Sprintf("%v", rs.srcFileSize)
				dstFileSize := fmt.Sprintf("%v", rs.dstFileSize)
				resStatus := rs.resStatus.String()
				errorStr := ""
				if rs.errorDetails != nil {
					errorStr = rs.errorDetails.Error()
				}
				ccrs = append(ccrs, []*string{
					&rs.srcFileName,
					&rs.dstFileName,
					&srcFileSize,
					&dstFileSize,
					&rs.srcCompressionType.name,
					&rs.dstCompressionType.name,
					&resStatus,
					&errorStr,
				})
			}
			data.RowSet = ccrs
			cc := make([]chunkRowType, len(ccrs))
			populateJSONRowSet(cc, ccrs)
			data.QueryResultFormat = "json"
			rt := []execResponseRowType{
				{Name: "source", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "target", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "source_size", ByteLength: 64, Length: 64, Type: "FIXED", Scale: 0, Nullable: false},
				{Name: "target_size", ByteLength: 64, Length: 64, Type: "FIXED", Scale: 0, Nullable: false},
				{Name: "source_compression", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "target_compression", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "status", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "message", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
			}
			data.RowType = rt
			return &execResponse{Data: *data, Success: true}, nil
		}
	} else { // DOWNLOAD
		if len(sfa.results) > 0 {
			for _, meta := range sfa.results {
				dstFileSize := meta.dstFileSize
				errorDetails := meta.errorDetails
				if sfa.options.RaisePutGetError && errorDetails != nil {
					return nil, (&SnowflakeError{
						Number:   ErrFailedToDownloadFromStage,
						SQLState: sfa.data.SQLState,
						QueryID:  sfa.data.QueryID,
						Message:  errorDetails.Error(),
					}).exceptionTelemetry(sfa.sc)
				}

				rowset = append(rowset, fileTransferResultType{
					"", "", meta.dstFileName, 0, dstFileSize,
					nil, nil, meta.resStatus, meta.errorDetails,
				})
			}
			sort.Slice(rowset, func(i, j int) bool {
				return rowset[i].srcFileName < rowset[j].srcFileName
			})
			ccrs := make([][]*string, 0, len(rowset))
			for _, rs := range rowset {
				dstFileSize := fmt.Sprintf("%v", rs.dstFileSize)
				resStatus := rs.resStatus.String()
				errorStr := ""
				if rs.errorDetails != nil {
					errorStr = rs.errorDetails.Error()
				}
				ccrs = append(ccrs, []*string{
					&rs.dstFileName,
					&dstFileSize,
					&resStatus,
					&errorStr,
				})
			}
			data.RowSet = ccrs
			cc := make([]chunkRowType, len(ccrs))
			populateJSONRowSet(cc, ccrs)
			data.QueryResultFormat = "json"
			rt := []execResponseRowType{
				{Name: "file", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "size", ByteLength: 64, Length: 64, Type: "FIXED", Scale: 0, Nullable: false},
				{Name: "status", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
				{Name: "message", ByteLength: 10000, Length: 10000, Type: "TEXT", Scale: 0, Nullable: false},
			}
			data.RowType = rt
			return &execResponse{Data: *data, Success: true}, nil
		}
	}
	return nil, (&SnowflakeError{
		Number:   ErrNotImplemented,
		SQLState: sfa.data.SQLState,
		QueryID:  sfa.data.QueryID,
		Message:  errMsgNotImplemented,
	}).exceptionTelemetry(sfa.sc)
}

func isFileTransfer(query string) bool {
	putRe := regexp.MustCompile(putRegexp)
	getRe := regexp.MustCompile(getRegexp)
	return putRe.Match([]byte(query)) || getRe.Match([]byte(query))
}

type snowflakeProgressPercentage struct {
	filename        string
	fileSize        float64
	outputStream    *io.Writer
	showProgressBar bool
	seenSoFar       int64
	done            bool
	startTime       time.Time
}

func (spp *snowflakeProgressPercentage) call(bytesAmount int64) {
	if spp.outputStream != nil {
		spp.seenSoFar += bytesAmount
		percentage := spp.percent(spp.seenSoFar, spp.fileSize)
		if !spp.done {
			spp.done = spp.updateProgress(spp.filename, spp.startTime, spp.fileSize, percentage, spp.outputStream, spp.showProgressBar)
		}
	}
}

func (spp *snowflakeProgressPercentage) percent(seenSoFar int64, size float64) float64 {
	if float64(seenSoFar) >= size || size <= 0 {
		return 1.0
	}
	return float64(seenSoFar) / size
}

func (spp *snowflakeProgressPercentage) updateProgress(filename string, startTime time.Time, totalSize float64, progress float64, outputStream *io.Writer, showProgressBar bool) bool {
	barLength := 10
	totalSize /= mb
	status := ""
	elapsedTime := time.Since(startTime)

	var throughput float64
	if elapsedTime != 0.0 {
		throughput = totalSize / elapsedTime.Seconds()
	}

	if progress < 0 {
		progress = 0
		status = "Halt...\r\n"
	}
	if progress >= 1 {
		status = fmt.Sprintf("Done (%.3fs, %.2fMB/s)", elapsedTime.Seconds(), throughput)
	}
	if status == "" && showProgressBar {
		status = fmt.Sprintf("(%.3fsm %.2fMB/s)", elapsedTime.Seconds(), throughput)
	}
	if status != "" {
		block := int(math.Round(float64(barLength) * progress))
		text := fmt.Sprintf("\r%v(%.2fMB): [%v] %.2f%% %v ", filename, totalSize, strings.Repeat("#", block)+strings.Repeat("-", barLength-block), progress*100, status)
		(*outputStream).Write([]byte(text))
	}
	return progress == 1.0
}
