// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-gometrics/api"
)

// UpdateCheck determines if the check needs to be updated (new metrics, tags, etc.)
func (cm *CheckManager) UpdateCheck(newMetrics map[string]*api.CheckBundleMetric) {
	// only if check manager is enabled
	if !cm.enabled {
		return
	}

	// only if checkBundle has been populated
	if cm.checkBundle == nil {
		return
	}

	// only if there is *something* to update
	if !cm.forceCheckUpdate && len(newMetrics) == 0 && len(cm.metricTags) == 0 {
		return
	}

	// refresh check bundle (in case there were changes made by other apps or in UI)
	checkBundle, err := cm.apih.FetchCheckBundleByCID(api.CIDType(cm.checkBundle.CID))
	if err != nil {
		cm.Log.Printf("[ERROR] unable to fetch up-to-date check bundle %v", err)
		return
	}
	cm.cbmu.Lock()
	cm.checkBundle = checkBundle
	cm.cbmu.Unlock()

	cm.addNewMetrics(newMetrics)

	if len(cm.metricTags) > 0 {
		// note: if a tag has been added (queued) for a metric which never gets sent
		//       the tags will be discarded. (setting tags does not *create* metrics.)
		for metricName, metricTags := range cm.metricTags {
			for metricIdx, metric := range cm.checkBundle.Metrics {
				if metric.Name == metricName {
					cm.checkBundle.Metrics[metricIdx].Tags = metricTags
					break
				}
			}
			cm.mtmu.Lock()
			delete(cm.metricTags, metricName)
			cm.mtmu.Unlock()
		}
		cm.forceCheckUpdate = true
	}

	if cm.forceCheckUpdate {
		newCheckBundle, err := cm.apih.UpdateCheckBundle(cm.checkBundle)
		if err != nil {
			cm.Log.Printf("[ERROR] updating check bundle %v", err)
			return
		}

		cm.forceCheckUpdate = false
		cm.cbmu.Lock()
		cm.checkBundle = newCheckBundle
		cm.cbmu.Unlock()
		cm.inventoryMetrics()
	}

}

// Initialize CirconusMetrics instance. Attempt to find a check otherwise create one.
// use cases:
//
// check [bundle] by submission url
// check [bundle] by *check* id (note, not check_bundle id)
// check [bundle] by search
// create check [bundle]
func (cm *CheckManager) initializeTrapURL() error {
	if cm.trapURL != "" {
		return nil
	}

	cm.trapmu.Lock()
	defer cm.trapmu.Unlock()

	// special case short-circuit: just send to a url, no check management
	// up to user to ensure that if url is https that it will work (e.g. not self-signed)
	if cm.checkSubmissionURL != "" {
		if !cm.enabled {
			cm.trapURL = cm.checkSubmissionURL
			cm.trapLastUpdate = time.Now()
			return nil
		}
	}

	if !cm.enabled {
		return errors.New("unable to initialize trap, check manager is disabled")
	}

	var err error
	var check *api.Check
	var checkBundle *api.CheckBundle
	var broker *api.Broker

	if cm.checkSubmissionURL != "" {
		check, err = cm.apih.FetchCheckBySubmissionURL(cm.checkSubmissionURL)
		if err != nil {
			return err
		}
		if !check.Active {
			return fmt.Errorf("[ERROR] Check ID %v is not active", check.CID)
		}
		// extract check id from check object returned from looking up using submission url
		// set m.CheckId to the id
		// set m.SubmissionUrl to "" to prevent trying to search on it going forward
		// use case: if the broker is changed in the UI metrics would stop flowing
		// unless the new submission url can be fetched with the API (which is no
		// longer possible using the original submission url)
		var id int
		id, err = strconv.Atoi(strings.Replace(check.CID, "/check/", "", -1))
		if err == nil {
			cm.checkID = api.IDType(id)
			cm.checkSubmissionURL = ""
		} else {
			cm.Log.Printf(
				"[WARN] SubmissionUrl check to Check ID: unable to convert %s to int %q\n",
				check.CID, err)
		}
	} else if cm.checkID > 0 {
		check, err = cm.apih.FetchCheckByID(cm.checkID)
		if err != nil {
			return err
		}
		if !check.Active {
			return fmt.Errorf("[ERROR] Check ID %v is not active", check.CID)
		}
	} else {
		if checkBundle == nil {
			// old search (instanceid as check.target)
			searchCriteria := fmt.Sprintf(
				"(active:1)(type:\"%s\")(host:\"%s\")(tags:%s)", cm.checkType, cm.checkTarget, strings.Join(cm.checkSearchTag, ","))
			checkBundle, err = cm.checkBundleSearch(searchCriteria, map[string]string{})
			if err != nil {
				return err
			}
		}

		if checkBundle == nil {
			// new search (check.target != instanceid, instanceid encoded in notes field)
			searchCriteria := fmt.Sprintf(
				"(active:1)(type:\"%s\")(tags:%s)", cm.checkType, strings.Join(cm.checkSearchTag, ","))
			filterCriteria := map[string]string{"f_notes": cm.getNotes()}
			checkBundle, err = cm.checkBundleSearch(searchCriteria, filterCriteria)
			if err != nil {
				return err
			}
		}

		if checkBundle == nil {
			// err==nil && checkBundle==nil is "no check bundles matched"
			// an error *should* be returned for any other invalid scenario
			checkBundle, broker, err = cm.createNewCheck()
			if err != nil {
				return err
			}
		}
	}

	if checkBundle == nil {
		if check != nil {
			checkBundle, err = cm.apih.FetchCheckBundleByCID(api.CIDType(check.CheckBundleCID))
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("[ERROR] Unable to retrieve, find, or create check")
		}
	}

	if broker == nil {
		broker, err = cm.apih.FetchBrokerByCID(api.CIDType(checkBundle.Brokers[0]))
		if err != nil {
			return err
		}
	}

	// retain to facilitate metric management (adding new metrics specifically)
	cm.checkBundle = checkBundle
	cm.inventoryMetrics()

	// determine the trap url to which metrics should be PUT
	if checkBundle.Type == "httptrap" {
		cm.trapURL = api.URLType(checkBundle.Config.SubmissionURL)
	} else {
		// build a submission_url for non-httptrap checks out of mtev_reverse url
		if len(checkBundle.ReverseConnectURLs) == 0 {
			return fmt.Errorf("%s is not an HTTPTRAP check and no reverse connection urls found", checkBundle.Checks[0])
		}
		mtevURL := checkBundle.ReverseConnectURLs[0]
		mtevURL = strings.Replace(mtevURL, "mtev_reverse", "https", 1)
		mtevURL = strings.Replace(mtevURL, "check", "module/httptrap", 1)
		cm.trapURL = api.URLType(fmt.Sprintf("%s/%s", mtevURL, checkBundle.Config.ReverseSecret))
	}

	// used when sending as "ServerName" get around certs not having IP SANS
	// (cert created with server name as CN but IP used in trap url)
	cn, err := cm.getBrokerCN(broker, cm.trapURL)
	if err != nil {
		return err
	}
	cm.trapCN = BrokerCNType(cn)

	cm.trapLastUpdate = time.Now()

	return nil
}

// Search for a check bundle given a predetermined set of criteria
func (cm *CheckManager) checkBundleSearch(criteria string, filter map[string]string) (*api.CheckBundle, error) {
	checkBundles, err := cm.apih.CheckBundleFilterSearch(api.SearchQueryType(criteria), filter)
	if err != nil {
		return nil, err
	}

	if len(checkBundles) == 0 {
		return nil, nil // trigger creation of a new check
	}

	numActive := 0
	checkID := -1

	for idx, check := range checkBundles {
		if check.Status == statusActive {
			numActive++
			checkID = idx
		}
	}

	if numActive > 1 {
		return nil, fmt.Errorf("[ERROR] multiple check bundles match criteria %s", criteria)
	}

	return &checkBundles[checkID], nil
}

// Create a new check to receive metrics
func (cm *CheckManager) createNewCheck() (*api.CheckBundle, *api.Broker, error) {
	checkSecret := string(cm.checkSecret)
	if checkSecret == "" {
		secret, err := cm.makeSecret()
		if err != nil {
			secret = "myS3cr3t"
		}
		checkSecret = secret
	}

	broker, err := cm.getBroker()
	if err != nil {
		return nil, nil, err
	}

	config := &api.CheckBundle{
		Brokers:     []string{broker.CID},
		Config:      api.CheckBundleConfig{AsyncMetrics: true, Secret: checkSecret},
		DisplayName: string(cm.checkDisplayName),
		Metrics:     []api.CheckBundleMetric{},
		MetricLimit: 0,
		Notes:       cm.getNotes(),
		Period:      60,
		Status:      statusActive,
		Tags:        append(cm.checkSearchTag, cm.checkTags...),
		Target:      string(cm.checkTarget),
		Timeout:     10,
		Type:        string(cm.checkType),
	}

	checkBundle, err := cm.apih.CreateCheckBundle(config)
	if err != nil {
		return nil, nil, err
	}

	return checkBundle, broker, nil
}

// Create a dynamic secret to use with a new check
func (cm *CheckManager) makeSecret() (string, error) {
	hash := sha256.New()
	x := make([]byte, 2048)
	if _, err := rand.Read(x); err != nil {
		return "", err
	}
	hash.Write(x)
	return hex.EncodeToString(hash.Sum(nil))[0:16], nil
}

func (cm *CheckManager) getNotes() string {
	return fmt.Sprintf("cgm_instanceid|%s", cm.checkInstanceID)
}
