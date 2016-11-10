// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"github.com/circonus-labs/circonus-gometrics/api"
)

// IsMetricActive checks whether a given metric name is currently active(enabled)
func (cm *CheckManager) IsMetricActive(name string) bool {
	active, _ := cm.availableMetrics[name]
	return active
}

// ActivateMetric determines if a given metric should be activated
func (cm *CheckManager) ActivateMetric(name string) bool {
	active, exists := cm.availableMetrics[name]

	if !exists {
		return true
	}

	if !active && cm.forceMetricActivation {
		return true
	}

	return false
}

// AddMetricTags updates check bundle metrics with tags
func (cm *CheckManager) AddMetricTags(metricName string, tags []string, appendTags bool) bool {
	tagsUpdated := false

	if len(tags) == 0 {
		return tagsUpdated
	}

	if _, exists := cm.metricTags[metricName]; !exists {
		foundMetric := false

		for _, metric := range cm.checkBundle.Metrics {
			if metric.Name == metricName {
				foundMetric = true
				cm.metricTags[metricName] = metric.Tags
				break
			}
		}

		if !foundMetric {
			cm.metricTags[metricName] = []string{}
		}
	}

	action := "no new"
	if appendTags {
		numNewTags := countNewTags(cm.metricTags[metricName], tags)
		if numNewTags > 0 {
			action = "Added"
			cm.metricTags[metricName] = append(cm.metricTags[metricName], tags...)
			tagsUpdated = true
		}
	} else {
		action = "Set"
		cm.metricTags[metricName] = tags
		tagsUpdated = true
	}

	if cm.Debug {
		cm.Log.Printf("[DEBUG] %s metric tag(s) %s %v\n", action, metricName, tags)
	}

	return tagsUpdated
}

// addNewMetrics updates a check bundle with new metrics
func (cm *CheckManager) addNewMetrics(newMetrics map[string]*api.CheckBundleMetric) bool {
	updatedCheckBundle := false

	if cm.checkBundle == nil || len(newMetrics) == 0 {
		return updatedCheckBundle
	}

	cm.cbmu.Lock()
	defer cm.cbmu.Unlock()

	numCurrMetrics := len(cm.checkBundle.Metrics)
	numNewMetrics := len(newMetrics)

	if numCurrMetrics+numNewMetrics >= cap(cm.checkBundle.Metrics) {
		nm := make([]api.CheckBundleMetric, numCurrMetrics+numNewMetrics)
		copy(nm, cm.checkBundle.Metrics)
		cm.checkBundle.Metrics = nm
	}

	cm.checkBundle.Metrics = cm.checkBundle.Metrics[0 : numCurrMetrics+numNewMetrics]

	i := 0
	for _, metric := range newMetrics {
		cm.checkBundle.Metrics[numCurrMetrics+i] = *metric
		i++
		updatedCheckBundle = true
	}

	if updatedCheckBundle {
		cm.forceCheckUpdate = true
	}

	return updatedCheckBundle
}

// inventoryMetrics creates list of active metrics in check bundle
func (cm *CheckManager) inventoryMetrics() {
	availableMetrics := make(map[string]bool)
	for _, metric := range cm.checkBundle.Metrics {
		availableMetrics[metric.Name] = metric.Status == "active"
	}
	cm.availableMetrics = availableMetrics
}

// countNewTags returns a count of new tags which do not exist in the current list of tags
func countNewTags(currTags []string, newTags []string) int {
	if len(newTags) == 0 {
		return 0
	}

	if len(currTags) == 0 {
		return len(newTags)
	}

	newTagCount := 0

	for _, newTag := range newTags {
		found := false
		for _, currTag := range currTags {
			if newTag == currTag {
				found = true
				break
			}
		}
		if !found {
			newTagCount++
		}
	}

	return newTagCount
}
