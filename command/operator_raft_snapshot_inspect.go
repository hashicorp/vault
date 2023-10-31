// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	protoio "github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftSnapshotInspectCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftSnapshotInspectCommand)(nil)
)

type OperatorRaftSnapshotInspectCommand struct {
	*BaseCommand
	details bool
	depth   int
	filter  string
}

func (c *OperatorRaftSnapshotInspectCommand) Synopsis() string {
	return "Inspects raft snapshot"
}

func (c *OperatorRaftSnapshotInspectCommand) Help() string {
	helpText := `
	Usage: vault operator raft snapshot inspect <snapshot_file>
	
	Inspects a snapshot file.
	
	$ vault operator raft snapshot inspect raft.snap
	
	` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftSnapshotInspectCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "details",
		Target:  &c.details,
		Default: true,
		Usage:   "Provides information about usage for data stored in the snapshot.",
	})

	f.IntVar(&IntVar{
		Name:    "depth",
		Target:  &c.depth,
		Default: 2,
		Usage:   "Can only be used with -details. The key prefix depth used to breakdown KV store data. If set to 0, all keys will be returned. Defaults to 2.",
	})

	f.StringVar(&StringVar{
		Name:    "filter",
		Target:  &c.filter,
		Default: "",
		Usage:   "Can only be used with -details. Limits the key breakdown using this prefix filter.",
	})

	return set
}

func (c *OperatorRaftSnapshotInspectCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftSnapshotInspectCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

type OutputFormat struct {
	Meta         *MetadataInfo
	StatsKV      []typeStats
	TotalCountKV int
	TotalSizeKV  int
}

// SnapshotInfo is used for passing snapshot stat
// information between functions
type SnapshotInfo struct {
	Meta         MetadataInfo
	StatsKV      map[string]typeStats
	TotalCountKV int
	TotalSizeKV  int
}

type MetadataInfo struct {
	ID      string
	Size    int64
	Index   uint64
	Term    uint64
	Version raft.SnapshotVersion
}

type typeStats struct {
	Name  string
	Count int
	Size  int
}

func (c *OperatorRaftSnapshotInspectCommand) Run(args []string) int {
	flags := c.Flags()

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Validate flags
	if c.depth < 0 {
		c.UI.Error("Depth must be equal to or greater than 0")
		return 1
	}

	var file string
	args = c.flags.Args()

	switch len(args) {
	case 0:
		c.UI.Error("Missing FILE argument")
		return 1
	case 1:
		file = args[0]
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error opening snapshot file: %s", err))
		return 1
	}
	defer f.Close()

	// Extract metadata and snapshot info from snapshot file
	var info *SnapshotInfo
	var meta *raft.SnapshotMeta
	info, meta, err = c.Read(hclog.New(nil), f)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading snapshot: %s", err))
		return 1
	}

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error outputting enhanced snapshot data: %s", err))
		return 1
	}

	// Generate structs for the formatter with information we read in
	metaformat := &MetadataInfo{
		ID:      meta.ID,
		Size:    meta.Size,
		Index:   meta.Index,
		Term:    meta.Term,
		Version: meta.Version,
	}

	formattedStatsKV := generateKVStats(*info)

	data := &OutputFormat{
		Meta:         metaformat,
		StatsKV:      formattedStatsKV,
		TotalCountKV: info.TotalCountKV,
		TotalSizeKV:  info.TotalSizeKV,
	}

	if Format(c.UI) != "table" {
		return OutputData(c.UI, data)
	}

	tableData, err := formatTable(data)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(tableData)

	return 0
}

func (c *OperatorRaftSnapshotInspectCommand) kvEnhance(val *pb.StorageEntry, info *SnapshotInfo, read int) {
	if !c.details {
		return
	}

	if val.Key == "" {
		return
	}

	// check for whether a filter is specified. if it is, skip
	// any keys that don't match.
	if len(c.filter) > 0 && !strings.HasPrefix(val.Key, c.filter) {
		return
	}

	split := strings.Split(val.Key, "/")

	// handle the situation where the key is shorter than
	// the specified depth.
	actualDepth := c.depth
	if c.depth == 0 || c.depth > len(split) {
		actualDepth = len(split)
	}

	prefix := strings.Join(split[0:actualDepth], "/")
	kvs := info.StatsKV[prefix]
	if kvs.Name == "" {
		kvs.Name = prefix
	}

	kvs.Count++
	kvs.Size += read
	info.TotalCountKV++
	info.TotalSizeKV += read
	info.StatsKV[prefix] = kvs
}

// Read from snapshot's state.bin and update the SnapshotInfo struct
func (c *OperatorRaftSnapshotInspectCommand) parseState(r io.Reader) (SnapshotInfo, error) {
	info := SnapshotInfo{
		StatsKV: make(map[string]typeStats),
	}

	protoReader := protoio.NewDelimitedReader(r, math.MaxInt32)

	for {
		s := new(pb.StorageEntry)
		if err := protoReader.ReadMsg(s); err != nil {
			if err == io.EOF {
				break
			}
			return info, err
		}
		size := protoReader.GetLastReadSize()
		c.kvEnhance(s, &info, size)
	}

	return info, nil
}

// Read contents of snapshot. Parse metadata and snapshot info
// Also, verify validity of snapshot
func (c *OperatorRaftSnapshotInspectCommand) Read(logger hclog.Logger, in io.Reader) (*SnapshotInfo, *raft.SnapshotMeta, error) {
	// Wrap the reader in a gzip decompressor.
	decomp, err := gzip.NewReader(in)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decompress snapshot: %v", err)
	}

	defer func() {
		if decomp == nil {
			return
		}

		if err := decomp.Close(); err != nil {
			logger.Error("Failed to close snapshot decompressor", "error", err)
		}
	}()

	// Read the archive.
	snapshotInfo, metadata, err := c.read(decomp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read snapshot file: %v", err)
	}

	if err := concludeGzipRead(decomp); err != nil {
		return nil, nil, err
	}

	if err := decomp.Close(); err != nil {
		return nil, nil, err
	}
	decomp = nil
	return snapshotInfo, metadata, nil
}

const (
	TableFormat string = "table"
	JSONFormat  string = "json"
)

func formatTable(info *OutputFormat) (string, error) {
	var b bytes.Buffer
	tw := tabwriter.NewWriter(&b, 8, 8, 6, ' ', 0)

	fmt.Fprintf(tw, " ID\t%s", info.Meta.ID)
	fmt.Fprintf(tw, "\n Size\t%d", info.Meta.Size)
	fmt.Fprintf(tw, "\n Index\t%d", info.Meta.Index)
	fmt.Fprintf(tw, "\n Term\t%d", info.Meta.Term)
	fmt.Fprintf(tw, "\n Version\t%d", info.Meta.Version)
	fmt.Fprintf(tw, "\n")

	if info.StatsKV != nil {
		fmt.Fprintf(tw, "\n")
		fmt.Fprintln(tw, "\n Key Name\tCount\tSize")
		fmt.Fprintf(tw, " %s\t%s\t%s", "----", "----", "----")

		for _, s := range info.StatsKV {
			fmt.Fprintf(tw, "\n %s\t%d\t%s", s.Name, s.Count, ByteSize(uint64(s.Size)))
		}

		fmt.Fprintf(tw, "\n %s\t%s", "----", "----")
		fmt.Fprintf(tw, "\n Total Size\t\t%s", ByteSize(uint64(info.TotalSizeKV)))
	}

	if err := tw.Flush(); err != nil {
		return b.String(), err
	}

	return b.String(), nil
}

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "TB"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "GB"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "KB"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// sortTypeStats sorts the stat slice by count and then
// alphabetically in the case the counts are equal
func sortTypeStats(stats []typeStats) []typeStats {
	// sort alphabetically if size is equal
	sort.Slice(stats, func(i, j int) bool {
		// Sort alphabetically if count is equal
		if stats[i].Count == stats[j].Count {
			return stats[i].Name < stats[j].Name
		}
		return stats[i].Count > stats[j].Count
	})

	return stats
}

// generateKVStats reformats the KV stats to work with
// the output struct that's used to produce the printed
// output the user sees.
func generateKVStats(info SnapshotInfo) []typeStats {
	kvLen := len(info.StatsKV)
	if kvLen > 0 {
		ks := make([]typeStats, 0, kvLen)

		for _, s := range info.StatsKV {
			ks = append(ks, s)
		}

		ks = sortTypeStats(ks)

		return ks
	}

	return nil
}

// hashList manages a list of filenames and their hashes.
type hashList struct {
	hashes map[string]hash.Hash
}

// newHashList returns a new hashList.
func newHashList() *hashList {
	return &hashList{
		hashes: make(map[string]hash.Hash),
	}
}

// Add creates a new hash for the given file.
func (hl *hashList) Add(file string) hash.Hash {
	if existing, ok := hl.hashes[file]; ok {
		return existing
	}

	h := sha256.New()
	hl.hashes[file] = h
	return h
}

// Encode takes the current sum of all the hashes and saves the hash list as a
// SHA256SUMS-style text file.
func (hl *hashList) Encode(w io.Writer) error {
	for file, h := range hl.hashes {
		if _, err := fmt.Fprintf(w, "%x  %s\n", h.Sum([]byte{}), file); err != nil {
			return err
		}
	}
	return nil
}

// DecodeAndVerify reads a SHA256SUMS-style text file and checks the results
// against the current sums for all the hashes.
func (hl *hashList) DecodeAndVerify(r io.Reader) error {
	// Read the file and make sure everything in there has a matching hash.
	seen := make(map[string]struct{})
	s := bufio.NewScanner(r)
	for s.Scan() {
		sha := make([]byte, sha256.Size)
		var file string
		if _, err := fmt.Sscanf(s.Text(), "%x  %s", &sha, &file); err != nil {
			return err
		}

		h, ok := hl.hashes[file]
		if !ok {
			return fmt.Errorf("list missing hash for %q", file)
		}
		if !bytes.Equal(sha, h.Sum([]byte{})) {
			return fmt.Errorf("hash check failed for %q", file)
		}
		seen[file] = struct{}{}
	}
	if err := s.Err(); err != nil {
		return err
	}

	// Make sure everything we had a hash for was seen.
	for file := range hl.hashes {
		if _, ok := seen[file]; !ok {
			return fmt.Errorf("file missing for %q", file)
		}
	}

	return nil
}

// read takes a reader and extracts the snapshot metadata and snapshot
// info. It also checks the integrity of the snapshot data.
func (c *OperatorRaftSnapshotInspectCommand) read(in io.Reader) (*SnapshotInfo, *raft.SnapshotMeta, error) {
	// Start a new tar reader.
	archive := tar.NewReader(in)

	// Create a hash list that we will use to compare with the SHA256SUMS
	// file in the archive.
	hl := newHashList()

	// Populate the hashes for all the files we expect to see. The check at
	// the end will make sure these are all present in the SHA256SUMS file
	// and that the hashes match.
	metaHash := hl.Add("meta.json")
	snapHash := hl.Add("state.bin")

	// Look through the archive for the pieces we care about.
	var shaBuffer bytes.Buffer
	var snapshotInfo SnapshotInfo
	var metadata raft.SnapshotMeta
	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed reading snapshot: %v", err)
		}

		switch hdr.Name {
		case "meta.json":
			// Previously we used json.Decode to decode the archive stream. There are
			// edgecases in which it doesn't read all the bytes from the stream, even
			// though the json object is still being parsed properly. Since we
			// simultaneously feeded everything to metaHash, our hash ended up being
			// different than what we calculated when creating the snapshot. Which in
			// turn made the snapshot verification fail. By explicitly reading the
			// whole thing first we ensure that we calculate the correct hash
			// independent of how json.Decode works internally.
			buf, err := io.ReadAll(io.TeeReader(archive, metaHash))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to read snapshot metadata: %v", err)
			}
			if err := json.Unmarshal(buf, &metadata); err != nil {
				return nil, nil, fmt.Errorf("failed to decode snapshot metadata: %v", err)
			}
		case "state.bin":
			// create reader that writes to snapHash what it reads from archive
			wrappedReader := io.TeeReader(archive, snapHash)
			var err error
			snapshotInfo, err = c.parseState(wrappedReader)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing snapshot state: %v", err)
			}

		case "SHA256SUMS":
			if _, err := io.CopyN(&shaBuffer, archive, 10000); err != nil && err != io.EOF {
				return nil, nil, fmt.Errorf("failed to read snapshot hashes: %v", err)
			}

		case "SHA256SUMS.sealed":
			// Add verification of sealed sum in future
			continue

		default:
			return nil, nil, fmt.Errorf("unexpected file %q in snapshot", hdr.Name)
		}
	}

	// Verify all the hashes.
	if err := hl.DecodeAndVerify(&shaBuffer); err != nil {
		return nil, nil, fmt.Errorf("failed checking integrity of snapshot: %v", err)
	}

	return &snapshotInfo, &metadata, nil
}

// concludeGzipRead should be invoked after you think you've consumed all of
// the data from the gzip stream. It will error if the stream was corrupt.
//
// The docs for gzip.Reader say: "Clients should treat data returned by Read as
// tentative until they receive the io.EOF marking the end of the data."
func concludeGzipRead(decomp *gzip.Reader) error {
	extra, err := io.ReadAll(decomp) // ReadAll consumes the EOF
	if err != nil {
		return err
	}
	if len(extra) != 0 {
		return fmt.Errorf("%d unread uncompressed bytes remain", len(extra))
	}
	return nil
}
