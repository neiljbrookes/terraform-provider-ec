// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package deploymentresource

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/terraform-provider-ec/ec/internal/util"
)

// flattenKibanaResources takes the kibana resource models and returns them flattened.
func flattenKibanaResources(in []*models.KibanaResourceInfo, name string) []interface{} {
	var result = make([]interface{}, 0, len(in))
	for _, res := range in {
		var m = make(map[string]interface{})
		if util.IsCurrentKibanaPlanEmpty(res) || isKibanaResourceStopped(res) {
			continue
		}

		if res.RefID != nil && *res.RefID != "" {
			m["ref_id"] = *res.RefID
		}

		if res.Info.ClusterID != nil && *res.Info.ClusterID != "" {
			m["resource_id"] = *res.Info.ClusterID
		}

		var plan = res.Info.PlanInfo.Current.Plan
		if plan.Kibana != nil {
			m["version"] = plan.Kibana.Version
		}

		if res.Region != nil {
			m["region"] = *res.Region
		}

		if topology := flattenKibanaTopology(plan); len(topology) > 0 {
			m["topology"] = topology
		}

		if res.ElasticsearchClusterRefID != nil {
			m["elasticsearch_cluster_ref_id"] = *res.ElasticsearchClusterRefID
		}

		for k, v := range util.FlattenClusterEndpoint(res.Info.Metadata) {
			m[k] = v
		}

		if c := flattenKibanaConfig(plan.Kibana); len(c) > 0 {
			m["config"] = c
		}

		result = append(result, m)
	}

	return result
}

func flattenKibanaTopology(plan *models.KibanaClusterPlan) []interface{} {
	var result = make([]interface{}, 0, len(plan.ClusterTopology))
	for _, topology := range plan.ClusterTopology {
		var m = make(map[string]interface{})
		if topology.Size == nil || topology.Size.Value == nil || *topology.Size.Value == 0 {
			continue
		}

		if topology.InstanceConfigurationID != "" {
			m["instance_configuration_id"] = topology.InstanceConfigurationID
		}

		if topology.Size != nil {
			m["size"] = util.MemoryToState(*topology.Size.Value)
			m["size_resource"] = *topology.Size.Resource

		}

		m["zone_count"] = topology.ZoneCount

		result = append(result, m)
	}

	return result
}

func flattenKibanaConfig(cfg *models.KibanaConfiguration) []interface{} {
	var m = make(map[string]interface{})
	if cfg == nil {
		return nil
	}

	if cfg.UserSettingsYaml != "" {
		m["user_settings_yaml"] = cfg.UserSettingsYaml
	}

	if cfg.UserSettingsOverrideYaml != "" {
		m["user_settings_override_yaml"] = cfg.UserSettingsOverrideYaml
	}

	if o := cfg.UserSettingsJSON; o != nil {
		if b, _ := json.Marshal(o); len(b) > 0 && !bytes.Equal([]byte("{}"), b) {
			m["user_settings_json"] = string(b)
		}
	}

	if o := cfg.UserSettingsOverrideJSON; o != nil {
		if b, _ := json.Marshal(o); len(b) > 0 && !bytes.Equal([]byte("{}"), b) {
			m["user_settings_override_json"] = string(b)
		}
	}

	if len(m) == 0 {
		return nil
	}

	return []interface{}{m}
}
