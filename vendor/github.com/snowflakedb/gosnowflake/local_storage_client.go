// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type localUtil struct {
}

func (util *localUtil) createClient(_ *execResponseStageInfo, _ bool) (cloudClient, error) {
	return nil, nil
}

func (util *localUtil) uploadOneFileWithRetry(meta *fileMetadata) error {
	var frd *bufio.Reader
	if meta.srcStream != nil {
		b := meta.srcStream
		if meta.realSrcStream != nil {
			b = meta.realSrcStream
		}
		frd = bufio.NewReader(b)
	} else {
		f, err := os.Open(meta.realSrcFileName)
		if err != nil {
			return err
		}
		defer f.Close()
		frd = bufio.NewReader(f)
	}

	user, err := expandUser(meta.stageInfo.Location)
	if err != nil {
		return err
	}
	if !meta.overwrite {
		if _, err := os.Stat(filepath.Join(user, meta.dstFileName)); err == nil {
			meta.dstFileSize = 0
			meta.resStatus = skipped
			return nil
		}
	}
	output, err := os.OpenFile(filepath.Join(user, meta.dstFileName), os.O_CREATE|os.O_WRONLY, readWriteFileMode)
	if err != nil {
		return err
	}
	defer output.Close()
	data := make([]byte, meta.uploadSize)
	for {
		n, err := frd.Read(data)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err = output.Write(data); err != nil {
			return err
		}
	}
	meta.dstFileSize = meta.uploadSize
	meta.resStatus = uploaded
	return nil
}

func (util *localUtil) downloadOneFile(meta *fileMetadata) error {
	srcFileName := meta.srcFileName
	if strings.HasPrefix(meta.srcFileName, fmt.Sprintf("%b", os.PathSeparator)) {
		srcFileName = srcFileName[1:]
	}
	user, err := expandUser(meta.stageInfo.Location)
	if err != nil {
		return err
	}
	fullSrcFileName := path.Join(user, srcFileName)
	user, err = expandUser(meta.localLocation)
	if err != nil {
		return err
	}
	fullDstFileName := path.Join(user, baseName(meta.dstFileName))
	baseDir, err := getDirectory()
	if err != nil {
		return err
	}
	if _, err = os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(fullSrcFileName)
	if err != nil {
		return err
	}
	if err = os.WriteFile(fullDstFileName, data, readWriteFileMode); err != nil {
		return err
	}
	fi, err := os.Stat(fullDstFileName)
	if err != nil {
		return err
	}
	meta.dstFileSize = fi.Size()
	meta.resStatus = downloaded
	return nil
}
