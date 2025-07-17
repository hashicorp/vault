// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"crypto/x509"
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/helper/errutil"
)

func prettyIssuer(issuerIdEntryMap map[issuing.IssuerID]*issuing.IssuerEntry, issuer issuing.IssuerID) string {
	if entry, ok := issuerIdEntryMap[issuer]; ok && len(entry.Name) > 0 {
		return "[id:" + string(issuer) + "/name:" + entry.Name + "]"
	}

	return "[" + string(issuer) + "]"
}

func (sc *storageContext) rebuildIssuersChains(referenceCert *issuing.IssuerEntry /* optional */) error {
	// This function rebuilds the CAChain field of all known issuers. This
	// function should usually be invoked when a new issuer is added to the
	// pool of issuers.
	//
	// In addition to the context and storage, we take an optional
	// referenceCert parameter -- an issuer certificate that we should write
	// to storage once done, but which might not be persisted yet (either due
	// to new values on it or due to it not yet existing in the list). This is
	// helpful when calling e.g., importIssuer(...) (from storage.go), to allow
	// the newly imported issuer to have its CAChain field computed, but
	// without writing and re-reading it from storage (potentially failing in
	// the process if chain building failed).
	//
	// Our contract guarantees that, if referenceCert is provided, we'll write
	// it to storage. Further, we guarantee that (given the issuers haven't
	// changed), the results will be stable on multiple calls to rebuild the
	// chain.
	//
	// Note that at no point in time do we fetch the private keys associated
	// with any issuers. It is sufficient to merely look at the issuers
	// themselves.
	//
	// To begin, we fetch all known issuers from disk.
	issuers, err := sc.listIssuers()
	if err != nil {
		return fmt.Errorf("unable to list issuers to build chain: %w", err)
	}

	// Fast path: no issuers means we can set the reference cert's value, if
	// provided, to itself.
	if len(issuers) == 0 {
		if referenceCert == nil {
			// Nothing to do; no reference cert was provided.
			return nil
		}

		// Otherwise, the only entry in the chain (that we know about) is the
		// certificate itself.
		referenceCert.CAChain = []string{referenceCert.Certificate}
		return sc.writeIssuer(referenceCert)
	}

	// Our provided reference cert might not be in the list of issuers. In
	// that case, add it manually.
	if referenceCert != nil {
		missing := true
		for _, issuer := range issuers {
			if issuer == referenceCert.ID {
				missing = false
				break
			}
		}

		if missing {
			issuers = append(issuers, referenceCert.ID)
		}
	}

	// Now call a stable sorting algorithm here. We want to ensure the results
	// are the same across multiple calls to rebuildIssuersChains with the same
	// input data.
	//
	// Note: while we want to ensure referenceCert is written last (because it
	// is the user-facing action), we need to balance this with always having
	// a stable chain order, regardless of which certificate was chosen as the
	// reference cert. (E.g., for a given collection of unchanging certificates,
	// if we repeatedly set+unset a manual chain, triggering rebuilds, we should
	// always have the same chain after each unset). Thus, delay the write of
	// the referenceCert below when persisting -- but keep the sort AFTER the
	// referenceCert was added to the list, not before.
	//
	// (Otherwise, if this is called with one existing issuer and one new
	//  reference cert, and the reference cert sorts before the existing
	//  issuer, we will sort this list and have persisted the new issuer
	//  first, and may fail on the subsequent write to the existing issuer.
	//  Alternatively, if we don't sort the issuers in this order and there's
	//  a parallel chain (where cert A is a child of both B and C, with
	//  C.ID < B.ID and C was passed in as the yet unwritten referenceCert),
	//  then we'll create a chain with order A -> B -> C on initial write (as
	//  A and B come from disk) but A -> C -> B on subsequent writes (when all
	//  certs come from disk). Thus the sort must be done after adding in the
	//  referenceCert, thus sorting it consistently, but its write must be
	//  singled out to occur last.)
	sort.SliceStable(issuers, func(i, j int) bool {
		return issuers[i] > issuers[j]
	})

	// We expect each of these maps to be the size of the number of issuers
	// we have (as we're mapping from issuers to other values).
	//
	// The first caches the storage entry for the issuer, the second caches
	// the parsed *x509.Certificate of the issuer itself, and the third and
	// fourth maps that certificate back to the other issuers with that
	// subject (note the keyword _other_: we'll exclude self-loops here) --
	// either via a parent or child relationship.
	issuerIdEntryMap := make(map[issuing.IssuerID]*issuing.IssuerEntry, len(issuers))
	issuerIdCertMap := make(map[issuing.IssuerID]*x509.Certificate, len(issuers))
	issuerIdParentsMap := make(map[issuing.IssuerID][]issuing.IssuerID, len(issuers))
	issuerIdChildrenMap := make(map[issuing.IssuerID][]issuing.IssuerID, len(issuers))

	// For every known issuer, we map that subject back to the id of issuers
	// containing that subject. This lets us build our IssuerID -> parents
	// mapping efficiently. Worst case we'll have a single linear chain where
	// every entry has a distinct subject.
	subjectIssuerIdsMap := make(map[string][]issuing.IssuerID, len(issuers))

	// First, read every issuer entry from storage. We'll propagate entries
	// to three of the maps here: all but issuerIdParentsMap and
	// issuerIdChildrenMap, which we'll do in a second pass.
	for _, identifier := range issuers {
		var stored *issuing.IssuerEntry

		// When the reference issuer is provided and matches this identifier,
		// prefer the updated reference copy instead.
		if referenceCert != nil && identifier == referenceCert.ID {
			stored = referenceCert
		} else {
			// Otherwise, fetch it from disk.
			stored, err = sc.fetchIssuerById(identifier)
			if err != nil {
				return fmt.Errorf("unable to fetch issuer %v to build chain: %w", identifier, err)
			}
		}

		if stored == nil || len(stored.Certificate) == 0 {
			return fmt.Errorf("bad issuer while building chain: missing certificate entry: %v", identifier)
		}

		issuerIdEntryMap[identifier] = stored
		cert, err := stored.GetCertificate()
		if err != nil {
			return fmt.Errorf("unable to parse issuer %v to certificate to build chain: %w", identifier, err)
		}

		issuerIdCertMap[identifier] = cert
		subjectIssuerIdsMap[string(cert.RawSubject)] = append(subjectIssuerIdsMap[string(cert.RawSubject)], identifier)
	}

	// Now that we have the subj->issuer map built, we can build the parent
	// and child mappings. We iterate over all issuers and build it one step
	// at a time.
	//
	// This is worst case O(n^2) because all of the issuers could have the
	// same name and be self-signed certs with different keys. That makes the
	// chain building (below) fast as they've all got empty parents/children
	// maps.
	//
	// Note that the order of iteration is stable. Why? We've built
	// subjectIssuerIdsMap from the (above) sorted issuers by appending the
	// next entry to the present list; since they're already sorted, that
	// lookup will also be sorted. Thus, each of these iterations are also
	// in sorted order, so the resulting map entries (of ids) are also sorted.
	// Thus, the graph structure is in sorted order and thus the toposort
	// below will be stable.
	for _, child := range issuers {
		// Fetch the certificate as we'll need it later.
		childCert := issuerIdCertMap[child]

		parentSubject := string(issuerIdCertMap[child].RawIssuer)
		parentCerts, ok := subjectIssuerIdsMap[parentSubject]
		if !ok {
			// When the issuer isn't known to Vault, the lookup by the issuer
			// will be empty. This most commonly occurs when intermediates are
			// directly added (via intermediate/set-signed) without providing
			// the root.
			continue
		}

		// Now, iterate over all possible parents and assign the child/parent
		// relationship.
		for _, parent := range parentCerts {
			// Skip self-references to the exact same certificate.
			if child == parent {
				continue
			}

			// While we could use Subject/Authority Key Identifier (SKI/AKI)
			// as a heuristic for whether or not this relationship is valid,
			// this is insufficient as otherwise valid CA certificates could
			// elide this information. That means its best to actually validate
			// the signature (e.g., call child.CheckSignatureFrom(parent))
			// instead.
			parentCert := issuerIdCertMap[parent]
			if err := childCert.CheckSignatureFrom(parentCert); err != nil {
				// We cannot return an error here as it could be that this
				// signature is entirely valid -- but just for a different
				// key. Instead, skip adding the parent->child and
				// child->parent link.
				continue
			}

			// Otherwise, we can append it to the map, allowing us to walk the
			// issuer->parent mapping.
			issuerIdParentsMap[child] = append(issuerIdParentsMap[child], parent)

			// Also cross-add the child relationship step at the same time.
			issuerIdChildrenMap[parent] = append(issuerIdChildrenMap[parent], child)
		}
	}

	// Finally, we consult RFC 8446 Section 4.4.2 for creating an algorithm for
	// building the chain:
	//
	// > ... The sender's certificate MUST come in the first
	// > CertificateEntry in the list.  Each following certificate SHOULD
	// > directly certify the one immediately preceding it.  Because
	// > certificate validation requires that trust anchors be distributed
	// > independently, a certificate that specifies a trust anchor MAY be
	// > omitted from the chain, provided that supported peers are known to
	// > possess any omitted certificates.
	// >
	// > Note: Prior to TLS 1.3, "certificate_list" ordering required each
	// > certificate to certify the one immediately preceding it; however,
	// > some implementations allowed some flexibility.  Servers sometimes
	// > send both a current and deprecated intermediate for transitional
	// > purposes, and others are simply configured incorrectly, but these
	// > cases can nonetheless be validated properly.  For maximum
	// > compatibility, all implementations SHOULD be prepared to handle
	// > potentially extraneous certificates and arbitrary orderings from any
	// > TLS version, with the exception of the end-entity certificate which
	// > MUST be first.
	//
	// So, we take this to mean we should build chains via DFS: each issuer is
	// explored until an empty parent pointer (i.e., self-loop) is reached and
	// then the last most recently seen duplicate parent link is then explored.
	//
	// However, we don't actually need to do a DFS (per issuer) here. We can
	// simply invert the (pseudo-)directed graph, i.e., topologically sort it.
	// Some number of certs (roots without cross-signing) lack parent issuers.
	// These are already "done" from the PoV of chain building. We can thus
	// iterating through the parent mapping to find entries without parents to
	// start the sort. After processing, we can add all children and visit them
	// if all parents have been processed.
	//
	// Note though, that while topographical sorting is equivalent to the DFS,
	// we have to take care to make it a pseudo-DAG. This means handling the
	// most common 2-star (2-clique) sub-graphs of reissued certificates,
	// manually building their chain prior to starting the topographical sort.
	//
	// This thus runs in O(|V| + |E|) -> O(n^2) in the number of issuers.
	processedIssuers := make(map[issuing.IssuerID]bool, len(issuers))
	toVisit := make([]issuing.IssuerID, 0, len(issuers))

	// Handle any explicitly constructed certificate chains. Here, we don't
	// validate much what the user provides; if they provide since-deleted
	// refs, skip them; if they duplicate entries, add them multiple times.
	// The other chain building logic will be able to deduplicate them when
	// used as parents to other certificates.
	for _, candidate := range issuers {
		entry := issuerIdEntryMap[candidate]
		if len(entry.ManualChain) == 0 {
			continue
		}

		entry.CAChain = nil
		for _, parentId := range entry.ManualChain {
			parentEntry := issuerIdEntryMap[parentId]
			if parentEntry == nil {
				continue
			}

			entry.CAChain = append(entry.CAChain, parentEntry.Certificate)
		}

		// Mark this node as processed and add its children.
		processedIssuers[candidate] = true
		children, ok := issuerIdChildrenMap[candidate]
		if !ok {
			continue
		}

		toVisit = append(toVisit, children...)
	}

	// Setup the toVisit queue.
	for _, candidate := range issuers {
		parentCerts, ok := issuerIdParentsMap[candidate]
		if ok && len(parentCerts) > 0 {
			// Assumption: no self-loops in the parent mapping, so if there's
			// a non-empty parent mapping it means we can skip this node as
			// it can't be processed yet.
			continue
		}

		// Because this candidate has no known parent issuers; update the
		// list.
		toVisit = append(toVisit, candidate)
	}

	// If the queue is empty (and we know we have issuers), trigger the
	// clique/cycle detection logic so we aren't starved for nodes.
	if len(toVisit) == 0 {
		toVisit, err = processAnyCliqueOrCycle(issuers, processedIssuers, toVisit, issuerIdEntryMap, issuerIdCertMap, issuerIdParentsMap, issuerIdChildrenMap, subjectIssuerIdsMap)
		if err != nil {
			return err
		}
	}

	// Now actually build the CAChain entries... Use a safety mechanism to
	// ensure we don't accidentally infinite-loop (if we introduce a bug).
	maxVisitCount := len(issuers)*len(issuers)*len(issuers) + 100
	for len(toVisit) > 0 && maxVisitCount >= 0 {
		var issuer issuing.IssuerID
		issuer, toVisit = toVisit[0], toVisit[1:]

		// If (and only if) we're presently starved for next nodes to visit,
		// attempt to resolve cliques and cycles again to fix that. This is
		// because all-cycles cycle detection is at least as costly as
		// traversing the entire graph a couple of times.
		//
		// Additionally, we do this immediately after popping a node from the
		// queue as we wish to ensure we never become starved for nodes.
		if len(toVisit) == 0 {
			toVisit, err = processAnyCliqueOrCycle(issuers, processedIssuers, toVisit, issuerIdEntryMap, issuerIdCertMap, issuerIdParentsMap, issuerIdChildrenMap, subjectIssuerIdsMap)
			if err != nil {
				return err
			}
		}

		// Self-loops and cross-signing might lead to this node already being
		// processed; skip it on the second pass.
		if processed, ok := processedIssuers[issuer]; ok && processed {
			continue
		}

		// Check our parent certs now; if they are all processed, we can
		// process this node. Otherwise, we'll re-add this to the queue
		// when the last parent is processed (and we re-add its children).
		parentCerts, ok := issuerIdParentsMap[issuer]
		if ok && len(parentCerts) > 0 {
			// For each parent, validate that we've processed it.
			mustSkip := false
			for _, parentCert := range parentCerts {
				if processed, ok := processedIssuers[parentCert]; !ok || !processed {
					mustSkip = true
					break
				}
			}

			if mustSkip {
				// Skip this node for now, we'll come back to it later.
				continue
			}
		}

		// Now we can build the chain. Start with the current cert...
		entry := issuerIdEntryMap[issuer]
		entry.CAChain = []string{entry.Certificate}

		// ...and add all parents into it. Note that we have to tell if
		// that parent was already visited or not.
		if ok && len(parentCerts) > 0 {
			// Split children into two categories: roots and intermediates.
			// When building a straight-line chain, we want to prefer the
			// root (thus, ending the verification) to any cross-signed
			// intermediates. If a root is cross-signed, we'll include it's
			// cross-signed cert in _its_ chain, thus ignoring our duplicate
			// parent here.
			//
			// Why? When you step from the present node ("issuer") onto one
			// of its parents, if you step onto a root, it is a no-op: you
			// can still visit all of the neighbors (because any neighbors,
			// if they exist, must be cross-signed alternative paths).
			// However, if you directly step onto the cross-signed, now you're
			// taken in an alternative direction (via its chain), and must
			// revisit any roots later.
			var roots []issuing.IssuerID
			var intermediates []issuing.IssuerID
			for _, parentCertId := range parentCerts {
				if bytes.Equal(issuerIdCertMap[parentCertId].RawSubject, issuerIdCertMap[parentCertId].RawIssuer) {
					roots = append(roots, parentCertId)
				} else {
					intermediates = append(intermediates, parentCertId)
				}
			}

			if len(parentCerts) > 1024*1024*1024 {
				return errutil.InternalError{Err: fmt.Sprintf("error building certificate chain, %d is too many parent certs",
					len(parentCerts))}
			}
			includedParentCerts := make(map[string]bool, len(parentCerts)+1)
			includedParentCerts[entry.Certificate] = true
			for _, parentCert := range append(roots, intermediates...) {
				// See discussion of the algorithm above as to why this is
				// in the correct order. However, note that we do need to
				// exclude duplicate certs, hence the map above.
				//
				// Assumption: issuerIdEntryMap and issuerIdParentsMap is well
				// constructed.
				parent := issuerIdEntryMap[parentCert]
				for _, parentChainCert := range parent.CAChain {
					addToChainIfNotExisting(includedParentCerts, entry, parentChainCert)
				}
			}
		}

		// Now, mark this node as processed and go and visit all of its
		// children.
		processedIssuers[issuer] = true

		childrenCerts, ok := issuerIdChildrenMap[issuer]
		if ok && len(childrenCerts) > 0 {
			toVisit = append(toVisit, childrenCerts...)
		}
	}

	// Assumption: no nodes left unprocessed. They should've either been
	// reached through the parent->child addition or they should've been
	// self-loops.
	var msg string
	for _, issuer := range issuers {
		if visited, ok := processedIssuers[issuer]; !ok || !visited {
			pretty := prettyIssuer(issuerIdEntryMap, issuer)
			msg += fmt.Sprintf("[failed to build chain correctly: unprocessed issuer %v: ok: %v; visited: %v]\n", pretty, ok, visited)
		}
	}
	if len(msg) > 0 {
		return errors.New(msg)
	}

	// Finally, write all issuers to disk.
	//
	// See the note above when sorting issuers for why we delay persisting
	// the referenceCert, if it was provided.
	for _, issuer := range issuers {
		entry := issuerIdEntryMap[issuer]

		if referenceCert != nil && issuer == referenceCert.ID {
			continue
		}

		err := sc.writeIssuer(entry)
		if err != nil {
			pretty := prettyIssuer(issuerIdEntryMap, issuer)
			return fmt.Errorf("failed to persist issuer (%v) chain to disk: %w", pretty, err)
		}
	}
	if referenceCert != nil {
		err := sc.writeIssuer(issuerIdEntryMap[referenceCert.ID])
		if err != nil {
			pretty := prettyIssuer(issuerIdEntryMap, referenceCert.ID)
			return fmt.Errorf("failed to persist issuer (%v) chain to disk: %w", pretty, err)
		}
	}

	// Everything worked \o/
	return nil
}

func addToChainIfNotExisting(includedParentCerts map[string]bool, entry *issuing.IssuerEntry, certToAdd string) {
	included, ok := includedParentCerts[certToAdd]
	if ok && included {
		return
	}

	entry.CAChain = append(entry.CAChain, certToAdd)
	includedParentCerts[certToAdd] = true
}

func processAnyCliqueOrCycle(
	issuers []issuing.IssuerID,
	processedIssuers map[issuing.IssuerID]bool,
	toVisit []issuing.IssuerID,
	issuerIdEntryMap map[issuing.IssuerID]*issuing.IssuerEntry,
	issuerIdCertMap map[issuing.IssuerID]*x509.Certificate,
	issuerIdParentsMap map[issuing.IssuerID][]issuing.IssuerID,
	issuerIdChildrenMap map[issuing.IssuerID][]issuing.IssuerID,
	subjectIssuerIdsMap map[string][]issuing.IssuerID,
) ([]issuing.IssuerID /* toVisit */, error) {
	// Topological sort really only works on directed acyclic graphs (DAGs).
	// But a pool of arbitrary (issuer) certificates are actually neither!
	// This pool could contain both cliques and cycles. Because this could
	// block chain construction, we need to handle these cases.
	//
	// Within the helper for rebuildIssuersChains, we realize that we might
	// have certain pathological cases where cliques and cycles might _mix_.
	// This warrants handling them outside of the topo-sort code, effectively
	// acting as a node-collapsing technique (turning many nodes into one).
	// In reality, we just special-case this and handle the processing of
	// these nodes manually, fixing their CAChain value and then skipping
	// them.
	//
	// Since clique detection is (in this case) cheap (at worst O(n) on the
	// size of the graph), we favor it over the cycle detection logic. The
	// order (in the case of mixed cliques+cycles) doesn't matter, as the
	// discovery of the clique will lead to the cycle. We additionally find
	// all (unprocessed) cliques first, so our cycle detection code can avoid
	// falling into cliques.
	//
	// We need to be able to handle cliques adjacent to cycles. This is
	// necessary because a cross-signed cert (with same subject and key as
	// the clique, but different issuer) could be part of a cycle; this cycle
	// loop forms a parent chain (that topo-sort can't resolve) -- AND the
	// clique itself mixes with this, so resolving one or the other isn't
	// sufficient (as the reissued clique plus the cross-signed cert
	// effectively acts as a single node in the cycle). Oh, and there might
	// be multiple cycles. :-)
	//
	// We also might just have cycles, separately from reissued cliques.
	//
	// The nice thing about both cliques and cycles is that, as long as you
	// deduplicate your certs, all issuers in the collection (including the
	// mixed collection) have the same chain entries, just in different
	// orders (preferring the cycle and appending the remaining clique
	// entries afterwards).

	// To begin, cache all cliques that we know about.
	allCliques, issuerIdCliqueMap, allCliqueNodes, err := findAllCliques(processedIssuers, issuerIdCertMap, subjectIssuerIdsMap, issuers)
	if err != nil {
		// Found a clique that is too large; exit with an error.
		return nil, err
	}

	for _, issuer := range issuers {
		// Skip anything that's already been processed.
		if processed, ok := processedIssuers[issuer]; ok && processed {
			continue
		}

		// This first branch is finding cliques. However, finding a clique is
		// not sufficient as discussed above -- we also need to find any
		// incident cycle as this cycle is a parent and child to the clique,
		// which means the cycle nodes _must_ include the clique _and_ the
		// clique must include the cycle (in the CA Chain computation).
		// However, its not sufficient to just do one and then the other:
		// we need the closure of all cliques (and their incident cycles).
		// Finally -- it isn't enough to consider this chain in isolation
		// either. We need to consider _all_ parents and ensure they've been
		// processed before processing this closure.
		var cliques [][]issuing.IssuerID
		var cycles [][]issuing.IssuerID
		closure := make(map[issuing.IssuerID]bool)

		var cliquesToProcess []issuing.IssuerID
		cliquesToProcess = append(cliquesToProcess, issuer)

		for len(cliquesToProcess) > 0 {
			var node issuing.IssuerID
			node, cliquesToProcess = cliquesToProcess[0], cliquesToProcess[1:]

			// Skip potential clique nodes which have already been processed
			// (either by the topo-sort or by this clique-finding code).
			if processed, ok := processedIssuers[node]; ok && processed {
				continue
			}
			if nodeInClosure, ok := closure[node]; ok && nodeInClosure {
				continue
			}

			// Check if we have a clique for this node from our computed
			// collection of cliques.
			cliqueId, ok := issuerIdCliqueMap[node]
			if !ok {
				continue
			}
			cliqueNodes := allCliques[cliqueId]

			// Add our discovered clique. Note that we avoid duplicate cliques by
			// the skip logic above. Additionally, we know that cliqueNodes must
			// be unique and not duplicated with any existing nodes so we can add
			// all nodes to closure.
			cliques = append(cliques, cliqueNodes)
			for _, node := range cliqueNodes {
				closure[node] = true
			}

			// Try and expand the clique to see if there's common cycles around
			// it. We exclude _all_ clique nodes from the expansion path, because
			// it will unnecessarily bloat the detected cycles AND we know that
			// we'll find them again from the neighborhood search.
			//
			// Additionally, note that, detection of cycles should be independent
			// of cliques: cliques form under reissuance, and cycles form via
			// cross-signing chains; the latter ensures that any cliques can be
			// strictly bypassed from cycles (but the chain construction later
			// ensures we pull in the cliques into the cycles).
			foundCycles, err := findCyclesNearClique(processedIssuers, issuerIdCertMap, issuerIdChildrenMap, allCliqueNodes)
			if err != nil {
				// Cycle is too large.
				return toVisit, err
			}

			// Assumption: each cycle in foundCycles is in canonical order (see note
			// below about canonical ordering). Deduplicate these against already
			// existing cycles and add them to the closure nodes.
			for _, cycle := range foundCycles {
				cycles = appendCycleIfNotExisting(cycles, cycle)

				// Now, for each cycle node, we need to find all adjacent cliques.
				// We do this by finding each child of the cycle and adding it to
				// the queue. If these nodes aren't on cliques, we'll skip them
				// fairly quickly since the cliques were pre-computed.
				for _, cycleNode := range cycle {
					children, ok := issuerIdChildrenMap[cycleNode]
					if !ok {
						continue
					}

					cliquesToProcess = append(cliquesToProcess, children...)

					// While we're here, add this cycle node to the closure.
					closure[cycleNode] = true
				}
			}
		}

		// Before we begin, we need to compute the _parents_ of the nodes in
		// these cliques and cycles and ensure they've all been processed (if
		// they're not already part of the closure).
		parents, ok := computeParentsFromClosure(processedIssuers, issuerIdParentsMap, closure)
		if !ok {
			// At least one parent wasn't processed; skip this cliques and
			// cycles group for now until they have all been processed.
			continue
		}

		// Ok, we've computed the closure. Now we can build CA nodes and mark
		// everything as processed, growing the toVisit queue in the process.
		// For every node we've found...
		for node := range closure {
			// Skip anything that's already been processed.
			if processed, ok := processedIssuers[node]; ok && processed {
				continue
			}

			// Before we begin, mark this node as processed (so we can continue
			// later) and add children to toVisit.
			processedIssuers[node] = true
			childrenCerts, ok := issuerIdChildrenMap[node]
			if ok && len(childrenCerts) > 0 {
				toVisit = append(toVisit, childrenCerts...)
			}

			// It can either be part of a clique or a cycle. We wish to add
			// the nodes of whatever grouping
			foundNode := false
			for _, clique := range cliques {
				inClique := false
				for _, cliqueNode := range clique {
					if cliqueNode == node {
						inClique = true
						break
					}
				}

				if inClique {
					foundNode = true

					// Compute this node's CAChain. Note order doesn't matter
					// (within the clique), but we'll preserve the relative
					// order of associated cycles.
					entry := issuerIdEntryMap[node]
					entry.CAChain = []string{entry.Certificate}

					includedParentCerts := make(map[string]bool, len(closure)+1)
					includedParentCerts[entry.Certificate] = true

					// First add nodes from this clique, then all cycles, and then
					// all other cliques.
					addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, clique)
					addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, cycles...)
					addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, cliques...)
					addParentChainsToEntry(issuerIdEntryMap, includedParentCerts, entry, parents)

					break
				}
			}

			// Otherwise, it must be part of a cycle.
			for _, cycle := range cycles {
				inCycle := false
				offsetInCycle := 0
				for index, cycleNode := range cycle {
					if cycleNode == node {
						inCycle = true
						offsetInCycle = index
						break
					}
				}

				if inCycle {
					foundNode = true

					// Compute this node's CAChain. Note that order within cycles
					// matters, but we'll preserve the relative order.
					entry := issuerIdEntryMap[node]
					entry.CAChain = []string{entry.Certificate}

					includedParentCerts := make(map[string]bool, len(closure)+1)
					includedParentCerts[entry.Certificate] = true

					// First add nodes from this cycle, then all cliques, then all
					// other cycles, and finally from parents.
					orderedCycle := append(cycle[offsetInCycle:], cycle[0:offsetInCycle]...)
					addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, orderedCycle)
					addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, cliques...)
					addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, cycles...)
					addParentChainsToEntry(issuerIdEntryMap, includedParentCerts, entry, parents)

					break
				}
			}

			if !foundNode {
				// Unable to find node; return an error. This shouldn't happen
				// generally.
				pretty := prettyIssuer(issuerIdEntryMap, issuer)
				return nil, fmt.Errorf("unable to find node (%v) in closure (%v) but not in cycles (%v) or cliques (%v)", pretty, closure, cycles, cliques)
			}
		}
	}

	// We might also have cycles without having associated cliques. We assume
	// that any cliques (if they existed and were relevant for the remaining
	// cycles) were processed at this point. However, we might still have
	// unprocessed cliques (and related cycles) at this point _if_ an
	// unrelated cycle is the parent to that clique+cycle group.
	for _, issuer := range issuers {
		// Skip this node if it is already processed.
		if processed, ok := processedIssuers[issuer]; ok && processed {
			continue
		}

		// Cliques should've been processed by now, if they were necessary
		// for processable cycles, so ignore them from here to avoid
		// bloating our search paths.
		cycles, err := findAllCyclesWithNode(processedIssuers, issuerIdCertMap, issuerIdChildrenMap, issuer, allCliqueNodes)
		if err != nil {
			// To large of cycle.
			return nil, err
		}

		closure := make(map[issuing.IssuerID]bool)
		for _, cycle := range cycles {
			for _, node := range cycle {
				closure[node] = true
			}
		}

		// Before we begin, we need to compute the _parents_ of the nodes in
		// these cycles and ensure they've all been processed (if they're not
		// part of the closure).
		parents, ok := computeParentsFromClosure(processedIssuers, issuerIdParentsMap, closure)
		if !ok {
			// At least one parent wasn't processed; skip this cycle
			// group for now until they have all been processed.
			continue
		}

		// Finally, for all detected cycles, build the CAChain for nodes in
		// cycles. Since they all share a common parent, they must all contain
		// each other.
		for _, cycle := range cycles {
			// For each node in each cycle
			for nodeIndex, node := range cycle {
				// If the node is processed already, skip it.
				if processed, ok := processedIssuers[node]; ok && processed {
					continue
				}

				// Otherwise, build its CAChain.
				entry := issuerIdEntryMap[node]
				entry.CAChain = []string{entry.Certificate}

				// No indication as to size of chain here
				includedParentCerts := make(map[string]bool)
				includedParentCerts[entry.Certificate] = true

				// First add nodes from this cycle, then all other cycles, and
				// finally from parents.
				orderedCycle := append(cycle[nodeIndex:], cycle[0:nodeIndex]...)
				addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, orderedCycle)
				addNodeCertsToEntry(issuerIdEntryMap, issuerIdChildrenMap, includedParentCerts, entry, cycles...)
				addParentChainsToEntry(issuerIdEntryMap, includedParentCerts, entry, parents)

				// Finally, mark the node as processed and add the remaining
				// children to toVisit.
				processedIssuers[node] = true
				childrenCerts, ok := issuerIdChildrenMap[node]
				if ok && len(childrenCerts) > 0 {
					toVisit = append(toVisit, childrenCerts...)
				}
			}
		}
	}

	return toVisit, nil
}

func findAllCliques(
	processedIssuers map[issuing.IssuerID]bool,
	issuerIdCertMap map[issuing.IssuerID]*x509.Certificate,
	subjectIssuerIdsMap map[string][]issuing.IssuerID,
	issuers []issuing.IssuerID,
) ([][]issuing.IssuerID, map[issuing.IssuerID]int, []issuing.IssuerID, error) {
	var allCliques [][]issuing.IssuerID
	issuerIdCliqueMap := make(map[issuing.IssuerID]int)
	var allCliqueNodes []issuing.IssuerID

	for _, node := range issuers {
		// Check if the node has already been visited...
		if processed, ok := processedIssuers[node]; ok && processed {
			// ...if so it might have had a manually constructed chain; skip
			// it for clique detection.
			continue
		}
		if _, ok := issuerIdCliqueMap[node]; ok {
			// ...if so it must be on another clique; skip the clique finding
			// so we don't get duplicated cliques.
			continue
		}

		// See if this is a node on a clique and find that clique.
		cliqueNodes, err := isOnReissuedClique(processedIssuers, issuerIdCertMap, subjectIssuerIdsMap, node)
		if err != nil {
			// Clique is too large.
			return nil, nil, nil, err
		}

		// Skip nodes which really aren't a clique.
		if len(cliqueNodes) <= 1 {
			continue
		}

		// Add this clique and update the mapping. A given node can only be in one
		// clique.
		cliqueId := len(allCliques)
		allCliques = append(allCliques, cliqueNodes)
		allCliqueNodes = append(allCliqueNodes, cliqueNodes...)
		for _, cliqueNode := range cliqueNodes {
			issuerIdCliqueMap[cliqueNode] = cliqueId
		}
	}

	return allCliques, issuerIdCliqueMap, allCliqueNodes, nil
}

func isOnReissuedClique(
	processedIssuers map[issuing.IssuerID]bool,
	issuerIdCertMap map[issuing.IssuerID]*x509.Certificate,
	subjectIssuerIdsMap map[string][]issuing.IssuerID,
	node issuing.IssuerID,
) ([]issuing.IssuerID, error) {
	// Finding max cliques in arbitrary graphs is a nearly pathological
	// problem, usually left to the realm of SAT solvers and NP-Complete
	// theoretical.
	//
	// We're not dealing with arbitrary graphs though. We're dealing with
	// a highly regular, highly structured constructed graph.
	//
	// Reissued cliques form in certificate chains when two conditions hold:
	//
	// 1. The Subject of the certificate matches the Issuer.
	// 2. The underlying public key is the same, resulting in the signature
	//    validating for any pair of certs.
	//
	// This follows from the definition of a reissued certificate (same key
	// material, subject, and issuer but with a different serial number and
	// a different validity period). The structure means that the graph is
	// highly regular: given a partial or self-clique, if any candidate node
	// can satisfy this relation with any node of the existing clique, it must
	// mean it must form a larger clique and satisfy this relationship with
	// all other nodes in the existing clique.
	//
	// (Aside: this is not the only type of clique, but it is the only type
	//  of 3+ node clique. A 2-star is emitted from certain graphs, but we
	//  chose to handle that case in the cycle detection code rather than
	//  under this reissued clique detection code).
	//
	// What does this mean for our algorithm? A simple greedy search is
	// sufficient. If we index our certificates by subject -> IssuerID
	// (and cache its value across calls, which we've already done for
	// building the parent/child relationship), we can find all other issuers
	// with the same public key and subject as the existing node fairly
	// easily.
	//
	// However, we should also set some reasonable bounds on clique size.
	// Let's limit it to 6 nodes.
	maxCliqueSize := 6

	// Per assumptions of how we've built the graph, these map lookups should
	// both exist.
	cert := issuerIdCertMap[node]
	subject := string(cert.RawSubject)
	issuer := string(cert.RawIssuer)
	candidates := subjectIssuerIdsMap[subject]

	// If the given node doesn't have the same subject and issuer, it isn't
	// a valid clique node.
	if subject != issuer {
		return nil, nil
	}

	// We have two choices here for validating that the two keys are the same:
	// perform a cheap ASN.1 encoding comparison of the public keys, which
	// _should_ be the same but may not be, or perform a more costly (but
	// which should definitely be correct) signature verification. We prefer
	// cheap and call it good enough.
	spki := cert.RawSubjectPublicKeyInfo

	// We know candidates has everything satisfying _half_ of the first
	// condition (the subject half), so validate they match the other half
	// (the issuer half) and the second condition. For node (which is
	// included in candidates), the condition should vacuously hold.
	var clique []issuing.IssuerID
	for _, candidate := range candidates {
		// Skip already processed nodes, even if they could be clique
		// candidates. We'll treat them as any other (already processed)
		// external parent in that scenario.
		if processed, ok := processedIssuers[candidate]; ok && processed {
			continue
		}

		candidateCert := issuerIdCertMap[candidate]
		hasRightKey := bytes.Equal(candidateCert.RawSubjectPublicKeyInfo, spki)
		hasMatchingIssuer := string(candidateCert.RawIssuer) == issuer

		if hasRightKey && hasMatchingIssuer {
			clique = append(clique, candidate)
		}
	}

	// Clique is invalid if it contains zero or one nodes.
	if len(clique) <= 1 {
		return nil, nil
	}

	// Validate it is within the acceptable clique size.
	if len(clique) > maxCliqueSize {
		return clique, fmt.Errorf("error building issuer chains: excessively reissued certificate: %v entries", len(clique))
	}

	// Must be a valid clique.
	return clique, nil
}

func containsIssuer(collection []issuing.IssuerID, target issuing.IssuerID) bool {
	if len(collection) == 0 {
		return false
	}

	for _, needle := range collection {
		if needle == target {
			return true
		}
	}

	return false
}

func appendCycleIfNotExisting(knownCycles [][]issuing.IssuerID, candidate []issuing.IssuerID) [][]issuing.IssuerID {
	// There's two ways to do cycle detection: canonicalize the cycles,
	// rewriting them to have the least (or max) element first or just
	// brute force the detection.
	//
	// Canonicalizing them is faster and easier to write (just compare
	// canonical forms) so do that instead.
	canonicalized := canonicalizeCycle(candidate)

	found := false
	for _, existing := range knownCycles {
		if len(existing) != len(canonicalized) {
			continue
		}

		equivalent := true
		for index, node := range canonicalized {
			if node != existing[index] {
				equivalent = false
				break
			}
		}

		if equivalent {
			found = true
			break
		}
	}

	if !found {
		return append(knownCycles, canonicalized)
	}

	return knownCycles
}

func canonicalizeCycle(cycle []issuing.IssuerID) []issuing.IssuerID {
	// Find the minimum value and put it at the head, keeping the relative
	// ordering the same.
	minIndex := 0
	for index, entry := range cycle {
		if entry < cycle[minIndex] {
			minIndex = index
		}
	}

	ret := append(cycle[minIndex:], cycle[0:minIndex]...)
	if len(ret) != len(cycle) {
		panic("ABORT")
	}

	return ret
}

func findCyclesNearClique(
	processedIssuers map[issuing.IssuerID]bool,
	issuerIdCertMap map[issuing.IssuerID]*x509.Certificate,
	issuerIdChildrenMap map[issuing.IssuerID][]issuing.IssuerID,
	cliqueNodes []issuing.IssuerID,
) ([][]issuing.IssuerID, error) {
	// When we have a reissued clique, we need to find all cycles next to it.
	// Presumably, because they all have non-empty parents, they should not
	// have been visited yet. We further know that (because we're exploring
	// the children path), any processed check would be unnecessary as all
	// children shouldn't have been processed yet (since their parents aren't
	// either).
	//
	// So, we can explore each of the children of any one clique node and
	// find all cycles using that node, until we come back to the starting
	// node, excluding the clique and other cycles.
	cliqueNode := cliqueNodes[0]

	// Copy the clique nodes as excluded nodes; we'll avoid exploring cycles
	// which have parents that have been already explored.
	excludeNodes := cliqueNodes[:]
	var knownCycles [][]issuing.IssuerID

	// We know the node has at least one child, since the clique is non-empty.
	for _, child := range issuerIdChildrenMap[cliqueNode] {
		// Skip children that are part of the clique.
		if containsIssuer(excludeNodes, child) {
			continue
		}

		// Find cycles containing this node.
		newCycles, err := findAllCyclesWithNode(processedIssuers, issuerIdCertMap, issuerIdChildrenMap, child, excludeNodes)
		if err != nil {
			// Found too large of a cycle
			return nil, err
		}

		// Add all cycles into the known cycles list.
		for _, cycle := range newCycles {
			knownCycles = appendCycleIfNotExisting(knownCycles, cycle)
		}

		// Exclude only the current child. Adding everything in the cycles
		// results might prevent discovery of other valid cycles.
		excludeNodes = append(excludeNodes, child)
	}

	// Sort cycles from longest->shortest.
	sort.SliceStable(knownCycles, func(i, j int) bool {
		return len(knownCycles[i]) < len(knownCycles[j])
	})

	return knownCycles, nil
}

func findAllCyclesWithNode(
	processedIssuers map[issuing.IssuerID]bool,
	issuerIdCertMap map[issuing.IssuerID]*x509.Certificate,
	issuerIdChildrenMap map[issuing.IssuerID][]issuing.IssuerID,
	source issuing.IssuerID,
	exclude []issuing.IssuerID,
) ([][]issuing.IssuerID, error) {
	// We wish to find all cycles involving this particular node and report
	// the corresponding paths. This is a full-graph traversal (excluding
	// certain paths) as we're not just checking if a cycle occurred, but
	// instead returning all of cycles with that node.
	//
	// Set some limit on max cycle size.
	maxCycleSize := 8

	// Whether we've visited any given node.
	cycleVisited := make(map[issuing.IssuerID]bool)
	visitCounts := make(map[issuing.IssuerID]int)
	parentCounts := make(map[issuing.IssuerID]map[issuing.IssuerID]bool)

	// Paths to the specified node. Some of these might be cycles.
	pathsTo := make(map[issuing.IssuerID][][]issuing.IssuerID)

	// Nodes to visit.
	var visitQueue []issuing.IssuerID

	// Add the source node to start. In order to set up the paths to a
	// given node, we seed pathsTo with the single path involving just
	// this node
	visitQueue = append(visitQueue, source)
	pathsTo[source] = [][]issuing.IssuerID{{source}}

	// Begin building paths.
	//
	// Loop invariant:
	//  pathTo[x] contains valid paths to reach this node, from source.
	for len(visitQueue) > 0 {
		var current issuing.IssuerID
		current, visitQueue = visitQueue[0], visitQueue[1:]

		// If we've already processed this node, we have a cycle. Skip this
		// node for now; we'll build cycles later.
		if processed, ok := cycleVisited[current]; ok && processed {
			continue
		}

		// Mark this node as visited for next time.
		cycleVisited[current] = true
		if _, ok := visitCounts[current]; !ok {
			visitCounts[current] = 0
		}
		visitCounts[current] += 1

		// For every child of this node...
		children, ok := issuerIdChildrenMap[current]
		if !ok {
			// Node has no children, nothing else we can do.
			continue
		}

		for _, child := range children {
			// Ensure we can visit this child; exclude processedIssuers and
			// exclude lists.
			if childProcessed, ok := processedIssuers[child]; ok && childProcessed {
				continue
			}

			skipNode := false
			for _, excluded := range exclude {
				if excluded == child {
					skipNode = true
					break
				}
			}

			if skipNode {
				continue
			}

			// Track this parent->child relationship to know when to exit.
			setOfParents, ok := parentCounts[child]
			if !ok {
				setOfParents = make(map[issuing.IssuerID]bool)
				parentCounts[child] = setOfParents
			}
			_, existingParent := setOfParents[current]
			setOfParents[current] = true

			// Since we know that we can visit this node, we should now build
			// all destination paths using this node, from our current node.
			//
			// Since these are all starting at a single path from source,
			// if we have any cycles back to source, we'll find them here.
			//
			// Only add this if it is a net-new path that doesn't repeat
			// (either internally -- indicating an internal cycle -- or
			//  externally with an existing path).
			addedPath := false
			if _, ok := pathsTo[child]; !ok {
				pathsTo[child] = make([][]issuing.IssuerID, 0)
			}

			for _, path := range pathsTo[current] {
				if child != source {
					// We only care about source->source cycles. If this
					// cycles, but isn't a source->source cycle, don't add
					// this path.
					foundSelf := false
					for _, node := range path {
						if child == node {
							foundSelf = true
							break
						}
					}
					if foundSelf {
						// Skip this path.
						continue
					}
				}

				if len(path) > 1024*1024*1024 {
					return nil, errutil.InternalError{Err: fmt.Sprintf("Error updating certificate path: path of length %d is too long", len(path))}
				}
				// Make sure to deep copy the path.
				newPath := make([]issuing.IssuerID, 0, len(path)+1)
				newPath = append(newPath, path...)
				newPath = append(newPath, child)

				isSamePath := false
				for _, childPath := range pathsTo[child] {
					if len(childPath) != len(newPath) {
						continue
					}

					isSamePath = true
					for index, node := range childPath {
						if newPath[index] != node {
							isSamePath = false
							break
						}
					}

					if isSamePath {
						break
					}
				}

				if !isSamePath {
					pathsTo[child] = append(pathsTo[child], newPath)
					addedPath = true
				}
			}

			// Add this child as a candidate to visit next.
			visitQueue = append(visitQueue, child)

			// If there's a new parent or we found a new path, then we should
			// revisit this child, to update _its_ children and see if there's
			// another new path. Eventually the paths will stabilize and we'll
			// end up with no new parents or paths.
			if !existingParent || addedPath {
				cycleVisited[child] = false
			}
		}
	}

	// Ok, we've now exited from our loop. Any cycles would've been detected
	// and their paths recorded in pathsTo. Now we can iterate over these
	// (starting a source), clean them up and validate them.
	var cycles [][]issuing.IssuerID
	for _, cycle := range pathsTo[source] {
		// Skip the trivial cycle.
		if len(cycle) == 1 && cycle[0] == source {
			continue
		}

		// Validate cycle starts and ends with source.
		if cycle[0] != source {
			return nil, fmt.Errorf("cycle (%v) unexpectedly starts with node %v; expected to start with %v", cycle, cycle[0], source)
		}

		// If the cycle doesn't start/end with the source,
		// skip it.
		if cycle[len(cycle)-1] != source {
			continue
		}

		truncatedCycle := cycle[0 : len(cycle)-1]
		if len(truncatedCycle) >= maxCycleSize {
			return nil, fmt.Errorf("cycle (%v) exceeds max size: %v > %v", cycle, len(cycle), maxCycleSize)
		}

		// Now one last thing: our cycle was built via parent->child
		// traversal, but we want child->parent ordered cycles. So,
		// just reverse it.
		reversed := reversedCycle(truncatedCycle)
		cycles = appendCycleIfNotExisting(cycles, reversed)
	}

	// Sort cycles from longest->shortest.
	sort.SliceStable(cycles, func(i, j int) bool {
		return len(cycles[i]) > len(cycles[j])
	})

	return cycles, nil
}

func reversedCycle(cycle []issuing.IssuerID) []issuing.IssuerID {
	var result []issuing.IssuerID
	for index := len(cycle) - 1; index >= 0; index-- {
		result = append(result, cycle[index])
	}

	return result
}

func computeParentsFromClosure(
	processedIssuers map[issuing.IssuerID]bool,
	issuerIdParentsMap map[issuing.IssuerID][]issuing.IssuerID,
	closure map[issuing.IssuerID]bool,
) (map[issuing.IssuerID]bool, bool) {
	parents := make(map[issuing.IssuerID]bool)
	for node := range closure {
		nodeParents, ok := issuerIdParentsMap[node]
		if !ok {
			continue
		}

		for _, parent := range nodeParents {
			if nodeInClosure, ok := closure[parent]; ok && nodeInClosure {
				continue
			}

			parents[parent] = true
			if processed, ok := processedIssuers[parent]; ok && processed {
				continue
			}

			return nil, false
		}
	}

	return parents, true
}

func addNodeCertsToEntry(
	issuerIdEntryMap map[issuing.IssuerID]*issuing.IssuerEntry,
	issuerIdChildrenMap map[issuing.IssuerID][]issuing.IssuerID,
	includedParentCerts map[string]bool,
	entry *issuing.IssuerEntry,
	issuersCollection ...[]issuing.IssuerID,
) {
	for _, collection := range issuersCollection {
		// Find a starting point into this collection such that it verifies
		// something in the existing collection.
		offset := 0
		for index, issuer := range collection {
			children, ok := issuerIdChildrenMap[issuer]
			if !ok {
				continue
			}

			foundChild := false
			for _, child := range children {
				childEntry := issuerIdEntryMap[child]
				if inChain, ok := includedParentCerts[childEntry.Certificate]; ok && inChain {
					foundChild = true
					break
				}
			}

			if foundChild {
				offset = index
				break
			}
		}

		// Assumption: collection is in child -> parent order. For cliques,
		// this is trivially true because everyone can validate each other,
		// but for cycles we have to ensure that in findAllCyclesWithNode.
		// This allows us to build the chain in the correct order.
		for _, issuer := range append(collection[offset:], collection[0:offset]...) {
			nodeEntry := issuerIdEntryMap[issuer]
			addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
		}
	}
}

func addParentChainsToEntry(
	issuerIdEntryMap map[issuing.IssuerID]*issuing.IssuerEntry,
	includedParentCerts map[string]bool,
	entry *issuing.IssuerEntry,
	parents map[issuing.IssuerID]bool,
) {
	for parent := range parents {
		nodeEntry := issuerIdEntryMap[parent]
		for _, cert := range nodeEntry.CAChain {
			addToChainIfNotExisting(includedParentCerts, entry, cert)
		}
	}
}
