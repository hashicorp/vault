package pki

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"sort"

	"github.com/hashicorp/vault/sdk/logical"
)

func rebuildIssuersChains(ctx context.Context, s logical.Storage, referenceCert *issuer /* optional */) error {
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
	// To begin, we fetch all known issuers from disk.
	issuers, err := listIssuers(ctx, s)
	if err != nil {
		return fmt.Errorf("unable to list issuers to build chain: %v", err)
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
		return writeIssuer(ctx, s, referenceCert)
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
	sort.SliceStable(issuers, func(i, j int) bool {
		return issuers[i] < issuers[j]
	})

	// We expect each of these maps to be the size of the number of issuers
	// we have (as we're mapping from issuers to other values).
	//
	// The first caches the storage entry for the issuer, the second caches
	// the parsed *x509.Certificate of the issuer itself, and the third and
	// fourth maps that certificate back to the other issuers with that
	// subject (note the keyword _other_: we'll exclude self-loops here) --
	// either via a parent or child relationship.
	issuerIdEntryMap := make(map[issuerId]*issuer, len(issuers))
	issuerIdCertMap := make(map[issuerId]*x509.Certificate, len(issuers))
	issuerIdParentsMap := make(map[issuerId][]issuerId, len(issuers))
	issuerIdChildrenMap := make(map[issuerId][]issuerId, len(issuers))

	// For every known issuer, we map that subject back to the id of issuers
	// containing that subject. This lets us build our issuerId -> parents
	// mapping efficiently. Worst case we'll have a single linear chain where
	// every entry has a distinct subject.
	subjectIssuerIdsMap := make(map[string][]issuerId, len(issuers))

	// First, read every issuer entry from storage. We'll propagate entries
	// to three of the maps here: all but issuerIdParentsMap and
	// issuerIdChildrenMap, which we'll do in a second pass.
	for _, identifier := range issuers {
		var stored *issuer

		// When the reference issuer is provided and matches this identifier,
		// prefer the updated reference copy instead.
		if referenceCert != nil && identifier == referenceCert.ID {
			stored = referenceCert
		} else {
			// Otherwise, fetch it from disk.
			stored, err = fetchIssuerById(ctx, s, identifier)
			if err != nil {
				return fmt.Errorf("unable to fetch issuer %v to build chain: %v", identifier, err)
			}
		}

		if stored == nil || len(stored.Certificate) == 0 {
			return fmt.Errorf("bad issuer while building chain: missing certificate entry: %v", identifier)
		}

		issuerIdEntryMap[identifier] = stored
		cert, err := stored.GetCertificate()
		if err != nil {
			return fmt.Errorf("unable to parse issuer %v to certificate to build chain: %v", identifier, err)
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
	processedIssuers := make(map[issuerId]bool, len(issuers))
	toVisit := make([]issuerId, 0, len(issuers))

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
		var issuer issuerId
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
			includedParentCerts := make(map[string]bool, len(parentCerts)+1)
			includedParentCerts[entry.Certificate] = true
			for _, parentCert := range parentCerts {
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
			msg += fmt.Sprintf("[failed to build chain correctly: unprocessed issuer %v: ok: %v; visited: %v]\n", issuer, ok, visited)
		}
	}
	if len(msg) > 0 {
		return fmt.Errorf(msg)
	}

	// Finally, write all issuers to disk.
	for _, issuer := range issuers {
		entry := issuerIdEntryMap[issuer]

		err := writeIssuer(ctx, s, entry)
		if err != nil {
			return fmt.Errorf("failed to persist issuer (%v) chain to disk: %v", issuer, err)
		}
	}

	// Everything worked \o/
	return nil
}

func addToChainIfNotExisting(includedParentCerts map[string]bool, entry *issuer, certToAdd string) {
	included, ok := includedParentCerts[certToAdd]
	if ok && included {
		return
	}

	entry.CAChain = append(entry.CAChain, certToAdd)
	includedParentCerts[certToAdd] = true
}

func processAnyCliqueOrCycle(
	issuers []issuerId,
	processedIssuers map[issuerId]bool,
	toVisit []issuerId,
	issuerIdEntryMap map[issuerId]*issuer,
	issuerIdCertMap map[issuerId]*x509.Certificate,
	issuerIdParentsMap map[issuerId][]issuerId,
	issuerIdChildrenMap map[issuerId][]issuerId,
	subjectIssuerIdsMap map[string][]issuerId,
) ([]issuerId /* toVisit */, error) {
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
	// discovery of the clique will lead to the cycle.
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

	for _, issuer := range issuers {
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

		var cliques [][]issuerId
		var cycles [][]issuerId
		var closure []issuerId

		var cliquesToProcess []issuerId
		cliquesToProcess = append(cliquesToProcess, issuer)

		for len(cliquesToProcess) > 0 {
			var node issuerId
			node, cliquesToProcess = cliquesToProcess[0], cliquesToProcess[1:]

			// Skip potential clique nodes which have already been processed
			// (either by the topo-sort or by this code).
			if processed, ok := processedIssuers[node]; ok && processed {
				continue
			}
			var skipNode bool = false
			for _, closureNode := range closure {
				if closureNode == node {
					skipNode = true
					break
				}
			}
			if skipNode {
				continue
			}

			// See if this is a node on a clique and find that clique.
			cliqueNodes, err := isOnReissuedClique(issuerIdCertMap, subjectIssuerIdsMap, node)
			if err != nil {
				// Clique is too large.
				return toVisit, err
			}

			// Skip nodes which aren't a clique.
			if len(cliqueNodes) <= 1 {
				continue
			}

			// Add our discovered clique. Note that we avoid duplicate cliques by
			// the skip logic above. Additionally, we know that cliqueNodes must
			// be unique and not duplicated with any existing nodes so we can add
			// all nodes to closure.
			cliques = append(cliques, cliqueNodes)
			closure = append(closure, cliqueNodes...)

			// Try and expand the clique to see if there's common cycles around
			// it.
			foundCycles, err := findCyclesNearClique(processedIssuers, issuerIdChildrenMap, cliqueNodes)
			if err != nil {
				// Cycle is too large.
				return toVisit, err
			}

			// We've already handled the clique, so we can continue early if we haven't
			// found any cycles.
			if len(foundCycles) == 0 {
				continue
			}

			// Assumption: each cycle in foundCycles is in canonical order (see note
			// below about canonical ordering). Deduplicate these against already
			// existing cycles and add them to the closure nodes.
			for _, cycle := range foundCycles {
				cycles = appendCycleIfNotExisting(cycles, cycle)

				// Now, for each cycle node, we need to find all adjacent cliques.
				// We do this by finding each child of the cycle and adding it to
				// the queue. If these nodes aren't cycles, we'll skip them fairly
				// quickly.
				for _, node := range cycle {
					children, ok := issuerIdChildrenMap[node]
					if !ok {
						continue
					}

					for _, child := range children {
						cliquesToProcess = append(cliquesToProcess, child)
					}

					// Finally, check if this node is in the closure; if not,
					// add it.
					inClosure := false
					for _, entry := range closure {
						if entry == node {
							inClosure = true
						}
					}
					if !inClosure {
						closure = append(closure, node)
					}
				}
			}
		}

		// If we lack a closure, go to the next node.
		if len(closure) == 0 {
			continue
		}

		// Ok, we've computed the closure. Now we can build CA nodes and mark
		// everything as processed, growing the toVisit queue in the process.
		// For every node we've found...
		for _, node := range closure {
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
			for cliqueIndex, clique := range cliques {
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

					// Add certs from _this_ clique
					for _, cliqueNode := range clique {
						nodeEntry := issuerIdEntryMap[cliqueNode]
						addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
					}

					// Add certs from all cycles.
					for _, cycle := range cycles {
						for _, cycleNode := range cycle {
							nodeEntry := issuerIdEntryMap[cycleNode]
							addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
						}
					}

					// Add certs from other cliques.
					for otherCliqueIndex, otherClique := range cliques {
						// Skip the present clique...
						if cliqueIndex == otherCliqueIndex {
							continue
						}

						for _, cliqueNode := range otherClique {
							nodeEntry := issuerIdEntryMap[cliqueNode]
							addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
						}
					}

					break
				}
			}

			// Otherwise, it must be part of a cycle.
			for cycleIndex, cycle := range cycles {
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

					// Handle this cycle's order correctly.
					for _, cycleNode := range append(cycle[offsetInCycle:], cycle[0:offsetInCycle]...) {
						nodeEntry := issuerIdEntryMap[cycleNode]
						addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
					}

					// Handle all cliques
					for _, clique := range cliques {
						for _, cliqueNode := range clique {
							nodeEntry := issuerIdEntryMap[cliqueNode]
							addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
						}
					}

					// Finally handle other cycles
					for otherCycleIndex, otherCycle := range cycles {
						if cycleIndex == otherCycleIndex {
							continue
						}

						for _, cycleNode := range otherCycle {
							nodeEntry := issuerIdEntryMap[cycleNode]
							addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
						}
					}

					break
				}
			}

			if !foundNode {
				// Unable to find node; return an error.
				return nil, fmt.Errorf("Unable to find node (%v) in closure (%v) but not in cycles (%v) or cliques (%v)", issuer, closure, cycles, cliques)
			}
		}
	}

	// Otherwise, we must have just have cycles without cliques.
	for _, issuer := range issuers {
		// Skip this node if it is already processed.
		if processed, ok := processedIssuers[issuer]; ok && processed {
			continue
		}

		cycles, err := findAllCyclesWithNode(processedIssuers, issuerIdChildrenMap, issuer, []issuerId{})
		if err != nil {
			// To large of cycle.
			return nil, err
		}

		// If we don't have any cycles, move onto the next node.
		if len(cycles) == 0 {
			continue
		}

		// Finally, for all detected cycles, build the CAChain for nodes in
		// cycles. Since they all share a common parent, they must all contain
		// each other.
		for cycleIndex, cycle := range cycles {
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

				for _, cycleNode := range append(cycle[nodeIndex:], cycle[0:nodeIndex]...) {
					nodeEntry := issuerIdEntryMap[cycleNode]
					addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)
				}

				for otherCycleIndex, otherCycle := range cycles {
					if cycleIndex == otherCycleIndex {
						continue
					}

					for _, cycleNode := range otherCycle {
						nodeEntry := issuerIdEntryMap[cycleNode]
						addToChainIfNotExisting(includedParentCerts, entry, nodeEntry.Certificate)

					}
				}

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

func isOnReissuedClique(
	issuerIdCertMap map[issuerId]*x509.Certificate,
	subjectIssuerIdsMap map[string][]issuerId,
	node issuerId,
) ([]issuerId, error) {
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
	// sufficient. If we index our certificates by subject -> issuerId
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
	var clique []issuerId
	for _, candidate := range candidates {
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

func containsIssuer(collection []issuerId, target issuerId) bool {
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

func appendCycleIfNotExisting(knownCycles [][]issuerId, candidate []issuerId) [][]issuerId {
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

func canonicalizeCycle(cycle []issuerId) []issuerId {
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
	processedIssuers map[issuerId]bool,
	issuerIdChildrenMap map[issuerId][]issuerId,
	cliqueNodes []issuerId,
) ([][]issuerId, error) {
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
	var knownCycles [][]issuerId

	// We know the node has at least one child, since the clique is non-empty.
	for _, child := range issuerIdChildrenMap[cliqueNode] {
		// Skip children that are part of the clique.
		if containsIssuer(excludeNodes, child) {
			continue
		}

		// Find cycles containing this node.
		newCycles, err := findAllCyclesWithNode(processedIssuers, issuerIdChildrenMap, child, excludeNodes)
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

	return knownCycles, nil
}

func findAllCyclesWithNode(
	processedIssuers map[issuerId]bool,
	issuerIdChildrenMap map[issuerId][]issuerId,
	source issuerId,
	exclude []issuerId,
) ([][]issuerId, error) {
	// We wish to find all cycles involving this particular node and report
	// the corresponding paths. This is a full-graph traversal (excluding
	// certain paths) as we're not just checking if a cycle occurred, but
	// instead returning all of cycles with that node.
	//
	// Set some limit on max cycle size.
	maxCycleSize := 8

	// Whether we've visited any given node.
	cycleVisited := make(map[issuerId]bool)

	// Paths to the specified node. Some of these might be cycles.
	pathsTo := make(map[issuerId][][]issuerId)

	// Nodes to visit.
	var visitQueue []issuerId

	// Add the source node to start. In order to set up the paths to a
	// given node, we seed pathsTo with the single path involving just
	// this node
	visitQueue = append(visitQueue, source)
	pathsTo[source] = [][]issuerId{{source}}

	// Begin building paths.
	//
	// Loop invariant:
	//  pathTo[x] contains valid paths to reach this node, from source.
	for len(visitQueue) > 0 {
		var current issuerId
		current, visitQueue = visitQueue[0], visitQueue[1:]

		// If we've already processed this node, we have a cycle. Skip this
		// node for now; we'll build cycles later.
		if processed, ok := cycleVisited[current]; ok && processed {
			continue
		}

		// Mark this node as visited for next time.
		cycleVisited[current] = true

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

			// Since we know that we can visit this node, we should now build
			// all destination paths using this node, from our current node.
			// Since these are all starting at a single path from source,
			// if we have any cycles back to source, we'll find them here.
			for _, path := range pathsTo[current] {
				newPath := append(path, child)
				pathsTo[child] = append(pathsTo[child], newPath)
			}

			// Visit this child next.
			visitQueue = append(visitQueue, child)
		}
	}

	// Ok, we've now exited from our loop. Any cycles would've been detected
	// and their paths recorded in pathsTo. Now we can iterate over these
	// (starting a source), clean them up and validate them.
	var cycles [][]issuerId
	for _, cycle := range pathsTo[source] {
		// Skip the trivial cycle.
		if len(cycle) == 1 {
			continue
		}

		// Validate cycle starts and ends with source.
		if cycle[0] != source {
			return nil, fmt.Errorf("cycle (%v) unexpectedly starts with node %v; expected to start with %v", cycle, cycle[0], source)
		}

		if cycle[len(cycle)-1] != source {
			return nil, fmt.Errorf("cycle (%v) unexpectedly ends with node %v; expected to start with %v", cycle, cycle[len(cycle)-1], source)
		}

		truncatedCycle := cycle[0 : len(cycle)-1]

		if len(truncatedCycle) > maxCycleSize {
			return nil, fmt.Errorf("cycle (%v) exceeds max size: %v > %v", cycle, len(cycle), maxCycleSize)
		}

		cycles = append(cycles, truncatedCycle)
	}

	return cycles, nil
}
