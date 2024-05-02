// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package puller

import (
	"context"
	"os"
	"slices"
	"testing"
)

func TestEnsureDownloaded(t *testing.T) {
	const (
		command = "vault-plugin-auth-jwt"
		version = "v0.17.0"
	)

	dir := t.TempDir()
	in := DownloadPluginInput{
		Directory: dir,
		Command:   command,
		Version:   version,
	}
	selected, sha256, err := EnsurePluginDownloaded(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if sha256 == nil {
		t.Fatal("expected sha256 sum")
	}
	in.SHA256Sum = sha256
	_, err = os.Lstat(in.targetFile())
	if err != nil {
		t.Fatal(err)
	}
	if selected != version {
		t.Fatalf("expected %s, got %s", version, selected)
	}

	selected, sha256, err = EnsurePluginDownloaded(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(in.SHA256Sum, sha256) {
		t.Fatalf("expected %x, got %x", in.SHA256Sum, sha256)
	}
	if selected != version {
		t.Fatalf("expected %s, got %s", version, selected)
	}
}

func TestEnsureDownloadedLatestVersion(t *testing.T) {
	const (
		command = "vault-plugin-auth-jwt"
	)

	dir := t.TempDir()
	in := DownloadPluginInput{
		Directory: dir,
		Command:   command,
	}
	selected, sha256, err := EnsurePluginDownloaded(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if sha256 == nil {
		t.Fatal("expected sha256 sum")
	}
	in.SHA256Sum = sha256
	in.Version = selected
	_, err = os.Lstat(in.targetFile())
	if err != nil {
		t.Fatal(err)
	}
	if selected != "0.20.3" {
		t.Fatalf("expected %s, got %s", "0.20.3", selected)
	}
}
