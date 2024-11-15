package oss

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadFile is multipart file upload.
//
// objectKey    the object name.
// filePath    the local file path to upload.
// partSize    the part size in byte.
// options    the options for uploading object.
//
// error    it's nil if the operation succeeds, otherwise it's an error object.
//
func (bucket Bucket) UploadFile(objectKey, filePath string, partSize int64, options ...Option) error {
	if partSize < MinPartSize || partSize > MaxPartSize {
		return errors.New("oss: part size invalid range (100KB, 5GB]")
	}

	cpConf := getCpConfig(options)
	routines := getRoutines(options)

	if cpConf != nil && cpConf.IsEnable {
		cpFilePath := getUploadCpFilePath(cpConf, filePath, bucket.BucketName, objectKey)
		if cpFilePath != "" {
			return bucket.uploadFileWithCp(objectKey, filePath, partSize, options, cpFilePath, routines)
		}
	}

	return bucket.uploadFile(objectKey, filePath, partSize, options, routines)
}

func getUploadCpFilePath(cpConf *cpConfig, srcFile, destBucket, destObject string) string {
	if cpConf.FilePath == "" && cpConf.DirPath != "" {
		dest := fmt.Sprintf("oss://%v/%v", destBucket, destObject)
		absPath, _ := filepath.Abs(srcFile)
		cpFileName := getCpFileName(absPath, dest, "")
		cpConf.FilePath = cpConf.DirPath + string(os.PathSeparator) + cpFileName
	}
	return cpConf.FilePath
}

// ----- concurrent upload without checkpoint  -----

// getCpConfig gets checkpoint configuration
func getCpConfig(options []Option) *cpConfig {
	cpcOpt, err := FindOption(options, checkpointConfig, nil)
	if err != nil || cpcOpt == nil {
		return nil
	}

	return cpcOpt.(*cpConfig)
}

// getCpFileName return the name of the checkpoint file
func getCpFileName(src, dest, versionId string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(src))
	srcCheckSum := hex.EncodeToString(md5Ctx.Sum(nil))

	md5Ctx.Reset()
	md5Ctx.Write([]byte(dest))
	destCheckSum := hex.EncodeToString(md5Ctx.Sum(nil))

	if versionId == "" {
		return fmt.Sprintf("%v-%v.cp", srcCheckSum, destCheckSum)
	}

	md5Ctx.Reset()
	md5Ctx.Write([]byte(versionId))
	versionCheckSum := hex.EncodeToString(md5Ctx.Sum(nil))
	return fmt.Sprintf("%v-%v-%v.cp", srcCheckSum, destCheckSum, versionCheckSum)
}

// getRoutines gets the routine count. by default it's 1.
func getRoutines(options []Option) int {
	rtnOpt, err := FindOption(options, routineNum, nil)
	if err != nil || rtnOpt == nil {
		return 1
	}

	rs := rtnOpt.(int)
	if rs < 1 {
		rs = 1
	} else if rs > 100 {
		rs = 100
	}

	return rs
}

// getPayer return the payer of the request
func getPayer(options []Option) string {
	payerOpt, err := FindOption(options, HTTPHeaderOssRequester, nil)
	if err != nil || payerOpt == nil {
		return ""
	}
	return payerOpt.(string)
}

// GetProgressListener gets the progress callback
func GetProgressListener(options []Option) ProgressListener {
	isSet, listener, _ := IsOptionSet(options, progressListener)
	if !isSet {
		return nil
	}
	return listener.(ProgressListener)
}

// uploadPartHook is for testing usage
type uploadPartHook func(id int, chunk FileChunk) error

var uploadPartHooker uploadPartHook = defaultUploadPart

func defaultUploadPart(id int, chunk FileChunk) error {
	return nil
}

// workerArg defines worker argument structure
type workerArg struct {
	bucket   *Bucket
	filePath string
	imur     InitiateMultipartUploadResult
	options  []Option
	hook     uploadPartHook
}

// worker is the worker coroutine function
type defaultUploadProgressListener struct {
}

// ProgressChanged no-ops
func (listener *defaultUploadProgressListener) ProgressChanged(event *ProgressEvent) {
}

func worker(id int, arg workerArg, jobs <-chan FileChunk, results chan<- UploadPart, failed chan<- error, die <-chan bool) {
	for chunk := range jobs {
		if err := arg.hook(id, chunk); err != nil {
			failed <- err
			break
		}
		var respHeader http.Header
		p := Progress(&defaultUploadProgressListener{})
		opts := make([]Option, len(arg.options)+2)
		opts = append(opts, arg.options...)

		// use defaultUploadProgressListener
		opts = append(opts, p, GetResponseHeader(&respHeader))

		startT := time.Now().UnixNano() / 1000 / 1000 / 1000
		part, err := arg.bucket.UploadPartFromFile(arg.imur, arg.filePath, chunk.Offset, chunk.Size, chunk.Number, opts...)
		endT := time.Now().UnixNano() / 1000 / 1000 / 1000
		if err != nil {
			arg.bucket.Client.Config.WriteLog(Debug, "upload part error,cost:%d second,part number:%d,request id:%s,error:%s\n", endT-startT, chunk.Number, GetRequestId(respHeader), err.Error())
			failed <- err
			break
		}
		select {
		case <-die:
			return
		default:
		}
		results <- part
	}
}

// scheduler function
func scheduler(jobs chan FileChunk, chunks []FileChunk) {
	for _, chunk := range chunks {
		jobs <- chunk
	}
	close(jobs)
}

func getTotalBytes(chunks []FileChunk) int64 {
	var tb int64
	for _, chunk := range chunks {
		tb += chunk.Size
	}
	return tb
}

// uploadFile is a concurrent upload, without checkpoint
func (bucket Bucket) uploadFile(objectKey, filePath string, partSize int64, options []Option, routines int) error {
	listener := GetProgressListener(options)

	chunks, err := SplitFileByPartSize(filePath, partSize)
	if err != nil {
		return err
	}

	partOptions := ChoiceTransferPartOption(options)
	completeOptions := ChoiceCompletePartOption(options)
	abortOptions := ChoiceAbortPartOption(options)

	// Initialize the multipart upload
	imur, err := bucket.InitiateMultipartUpload(objectKey, options...)
	if err != nil {
		return err
	}

	jobs := make(chan FileChunk, len(chunks))
	results := make(chan UploadPart, len(chunks))
	failed := make(chan error)
	die := make(chan bool)

	var completedBytes int64
	totalBytes := getTotalBytes(chunks)
	event := newProgressEvent(TransferStartedEvent, 0, totalBytes, 0)
	publishProgress(listener, event)

	// Start the worker coroutine
	arg := workerArg{&bucket, filePath, imur, partOptions, uploadPartHooker}
	for w := 1; w <= routines; w++ {
		go worker(w, arg, jobs, results, failed, die)
	}

	// Schedule the jobs
	go scheduler(jobs, chunks)

	// Waiting for the upload finished
	completed := 0
	parts := make([]UploadPart, len(chunks))
	for completed < len(chunks) {
		select {
		case part := <-results:
			completed++
			parts[part.PartNumber-1] = part
			completedBytes += chunks[part.PartNumber-1].Size

			// why RwBytes in ProgressEvent is 0 ?
			// because read or write event has been notified in teeReader.Read()
			event = newProgressEvent(TransferDataEvent, completedBytes, totalBytes, chunks[part.PartNumber-1].Size)
			publishProgress(listener, event)
		case err := <-failed:
			close(die)
			event = newProgressEvent(TransferFailedEvent, completedBytes, totalBytes, 0)
			publishProgress(listener, event)
			bucket.AbortMultipartUpload(imur, abortOptions...)
			return err
		}

		if completed >= len(chunks) {
			break
		}
	}

	event = newProgressEvent(TransferCompletedEvent, completedBytes, totalBytes, 0)
	publishProgress(listener, event)

	// Complete the multpart upload
	_, err = bucket.CompleteMultipartUpload(imur, parts, completeOptions...)
	if err != nil {
		bucket.AbortMultipartUpload(imur, abortOptions...)
		return err
	}
	return nil
}

// ----- concurrent upload with checkpoint  -----
const uploadCpMagic = "FE8BB4EA-B593-4FAC-AD7A-2459A36E2E62"

type uploadCheckpoint struct {
	Magic     string   // Magic
	MD5       string   // Checkpoint file content's MD5
	FilePath  string   // Local file path
	FileStat  cpStat   // File state
	ObjectKey string   // Key
	UploadID  string   // Upload ID
	Parts     []cpPart // All parts of the local file
	CallbackVal string
	CallbackBody *[]byte
}

type cpStat struct {
	Size         int64     // File size
	LastModified time.Time // File's last modified time
	MD5          string    // Local file's MD5
}

type cpPart struct {
	Chunk       FileChunk  // File chunk
	Part        UploadPart // Uploaded part
	IsCompleted bool       // Upload complete flag
}

// isValid checks if the uploaded data is valid---it's valid when the file is not updated and the checkpoint data is valid.
func (cp uploadCheckpoint) isValid(filePath string,options []Option) (bool, error) {

	callbackVal, _ := FindOption(options, HTTPHeaderOssCallback, "")
	if callbackVal != "" && cp.CallbackVal != callbackVal {
		return false, nil
	}
	callbackBody, _ := FindOption(options, responseBody, nil)
	if callbackBody != nil{
		body, _ := json.Marshal(callbackBody)
		if bytes.Equal(*cp.CallbackBody, body) {
			return false, nil
		}
	}
	// Compare the CP's magic number and MD5.
	cpb := cp
	cpb.MD5 = ""
	js, _ := json.Marshal(cpb)
	sum := md5.Sum(js)
	b64 := base64.StdEncoding.EncodeToString(sum[:])

	if cp.Magic != uploadCpMagic || b64 != cp.MD5 {
		return false, nil
	}

	// Make sure if the local file is updated.
	fd, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer fd.Close()

	st, err := fd.Stat()
	if err != nil {
		return false, err
	}

	md, err := calcFileMD5(filePath)
	if err != nil {
		return false, err
	}

	// Compare the file size, file's last modified time and file's MD5
	if cp.FileStat.Size != st.Size() ||
		!cp.FileStat.LastModified.Equal(st.ModTime()) ||
		cp.FileStat.MD5 != md {
		return false, nil
	}

	return true, nil
}

// load loads from the file
func (cp *uploadCheckpoint) load(filePath string) error {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, cp)
	return err
}

// dump dumps to the local file
func (cp *uploadCheckpoint) dump(filePath string) error {
	bcp := *cp

	// Calculate MD5
	bcp.MD5 = ""
	js, err := json.Marshal(bcp)
	if err != nil {
		return err
	}
	sum := md5.Sum(js)
	b64 := base64.StdEncoding.EncodeToString(sum[:])
	bcp.MD5 = b64

	// Serialization
	js, err = json.Marshal(bcp)
	if err != nil {
		return err
	}

	// Dump
	return ioutil.WriteFile(filePath, js, FilePermMode)
}

// updatePart updates the part status
func (cp *uploadCheckpoint) updatePart(part UploadPart) {
	cp.Parts[part.PartNumber-1].Part = part
	cp.Parts[part.PartNumber-1].IsCompleted = true
}

// todoParts returns unfinished parts
func (cp *uploadCheckpoint) todoParts() []FileChunk {
	fcs := []FileChunk{}
	for _, part := range cp.Parts {
		if !part.IsCompleted {
			fcs = append(fcs, part.Chunk)
		}
	}
	return fcs
}

// allParts returns all parts
func (cp *uploadCheckpoint) allParts() []UploadPart {
	ps := []UploadPart{}
	for _, part := range cp.Parts {
		ps = append(ps, part.Part)
	}
	return ps
}

// getCompletedBytes returns completed bytes count
func (cp *uploadCheckpoint) getCompletedBytes() int64 {
	var completedBytes int64
	for _, part := range cp.Parts {
		if part.IsCompleted {
			completedBytes += part.Chunk.Size
		}
	}
	return completedBytes
}

// calcFileMD5 calculates the MD5 for the specified local file
func calcFileMD5(filePath string) (string, error) {
	return "", nil
}

// prepare initializes the multipart upload
func prepare(cp *uploadCheckpoint, objectKey, filePath string, partSize int64, bucket *Bucket, options []Option) error {
	// CP
	cp.Magic = uploadCpMagic
	cp.FilePath = filePath
	cp.ObjectKey = objectKey

	// Local file
	fd, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	st, err := fd.Stat()
	if err != nil {
		return err
	}
	cp.FileStat.Size = st.Size()
	cp.FileStat.LastModified = st.ModTime()
	callbackVal, _ := FindOption(options, HTTPHeaderOssCallback, "")
	cp.CallbackVal = callbackVal.(string)
	callbackBody, _ := FindOption(options, responseBody, nil)
	if callbackBody != nil  {
		body, _ := json.Marshal(callbackBody)
		cp.CallbackBody = &body
	}
	md, err := calcFileMD5(filePath)
	if err != nil {
		return err
	}
	cp.FileStat.MD5 = md

	// Chunks
	parts, err := SplitFileByPartSize(filePath, partSize)
	if err != nil {
		return err
	}

	cp.Parts = make([]cpPart, len(parts))
	for i, part := range parts {
		cp.Parts[i].Chunk = part
		cp.Parts[i].IsCompleted = false
	}

	// Init load
	imur, err := bucket.InitiateMultipartUpload(objectKey, options...)
	if err != nil {
		return err
	}
	cp.UploadID = imur.UploadID

	return nil
}

// complete completes the multipart upload and deletes the local CP files
func complete(cp *uploadCheckpoint, bucket *Bucket, parts []UploadPart, cpFilePath string, options []Option) error {
	imur := InitiateMultipartUploadResult{Bucket: bucket.BucketName,
		Key: cp.ObjectKey, UploadID: cp.UploadID}

	_, err := bucket.CompleteMultipartUpload(imur, parts, options...)
	if err != nil {
		if e, ok := err.(ServiceError);ok && (e.StatusCode == 203 || e.StatusCode == 404) {
			os.Remove(cpFilePath)
		}
		return err
	}
	os.Remove(cpFilePath)
	return err
}

// uploadFileWithCp handles concurrent upload with checkpoint
func (bucket Bucket) uploadFileWithCp(objectKey, filePath string, partSize int64, options []Option, cpFilePath string, routines int) error {
	listener := GetProgressListener(options)

	partOptions := ChoiceTransferPartOption(options)
	completeOptions := ChoiceCompletePartOption(options)

	// Load CP data
	ucp := uploadCheckpoint{}
	err := ucp.load(cpFilePath)
	if err != nil {
		os.Remove(cpFilePath)
	}

	// Load error or the CP data is invalid.
	valid, err := ucp.isValid(filePath,options)
	if err != nil || !valid {
		if err = prepare(&ucp, objectKey, filePath, partSize, &bucket, options); err != nil {
			return err
		}
		os.Remove(cpFilePath)
	}

	chunks := ucp.todoParts()
	imur := InitiateMultipartUploadResult{
		Bucket:   bucket.BucketName,
		Key:      objectKey,
		UploadID: ucp.UploadID}

	jobs := make(chan FileChunk, len(chunks))
	results := make(chan UploadPart, len(chunks))
	failed := make(chan error)
	die := make(chan bool)

	completedBytes := ucp.getCompletedBytes()

	// why RwBytes in ProgressEvent is 0 ?
	// because read or write event has been notified in teeReader.Read()
	event := newProgressEvent(TransferStartedEvent, completedBytes, ucp.FileStat.Size, 0)
	publishProgress(listener, event)

	// Start the workers
	arg := workerArg{&bucket, filePath, imur, partOptions, uploadPartHooker}
	for w := 1; w <= routines; w++ {
		go worker(w, arg, jobs, results, failed, die)
	}

	// Schedule jobs
	go scheduler(jobs, chunks)

	// Waiting for the job finished
	completed := 0
	for completed < len(chunks) {
		select {
		case part := <-results:
			completed++
			ucp.updatePart(part)
			ucp.dump(cpFilePath)
			completedBytes += ucp.Parts[part.PartNumber-1].Chunk.Size
			event = newProgressEvent(TransferDataEvent, completedBytes, ucp.FileStat.Size, ucp.Parts[part.PartNumber-1].Chunk.Size)
			publishProgress(listener, event)
		case err := <-failed:
			close(die)
			event = newProgressEvent(TransferFailedEvent, completedBytes, ucp.FileStat.Size, 0)
			publishProgress(listener, event)
			return err
		}

		if completed >= len(chunks) {
			break
		}
	}

	event = newProgressEvent(TransferCompletedEvent, completedBytes, ucp.FileStat.Size, 0)
	publishProgress(listener, event)

	// Complete the multipart upload
	err = complete(&ucp, &bucket, ucp.allParts(), cpFilePath, completeOptions)
	return err
}
