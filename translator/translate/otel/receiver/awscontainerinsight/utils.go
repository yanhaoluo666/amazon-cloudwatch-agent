// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package awscontainerinsight

import (
	"log"
	"time"

	"go.opentelemetry.io/collector/confmap"

	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

const (
	BaseContainerInsights            = iota + 1
	DefaultMetricsCollectionInterval = time.Minute
)

func EnhancedContainerInsightsEnabled(conf *confmap.Conf) bool {
	isSet := common.GetOrDefaultBool(conf, common.ConfigKey(common.LogsKey, common.MetricsCollectedKey, common.KubernetesKey, common.EnhancedContainerInsights), false)
	if !isSet {
		levelFloat := common.GetOrDefaultNumber(conf, common.ConfigKey(common.LogsKey, common.MetricsCollectedKey, common.KubernetesKey, common.ContainerInsightsMetricGranularity), 1)
		if levelFloat > BaseContainerInsights {
			isSet = true
		}
	}
	return isSet
}

func AcceleratedComputeMetricsEnabled(conf *confmap.Conf) bool {
	return common.GetOrDefaultBool(conf, common.ConfigKey(common.LogsKey, common.MetricsCollectedKey, common.KubernetesKey, common.EnableAcceleratedComputeMetric), true)
}

func GetAcceleratedComputeGPUMetricsCollectionInterval(conf *confmap.Conf) time.Duration {
	return common.GetOrDefaultDuration(conf, []string{
		common.ConfigKey(common.LogsKey, common.MetricsCollectedKey, common.KubernetesKey, common.AcceleratedComputeGPUMetricsCollectionInterval),
	}, DefaultMetricsCollectionInterval)
}

func IsHighFrequencyGPUMetricsEnabled(conf *confmap.Conf) bool {
	enhancedEnabled := EnhancedContainerInsightsEnabled(conf)
	acceleratedEnabled := AcceleratedComputeMetricsEnabled(conf)

	// Check if accelerated_compute_gpu_metrics_collection_interval exists in config
	gpuMetricsCollectionIntervalKey := common.ConfigKey(common.LogsKey, common.MetricsCollectedKey, common.KubernetesKey, common.AcceleratedComputeGPUMetricsCollectionInterval)
	gpuMetricsCollectionIntervalExists := conf.IsSet(gpuMetricsCollectionIntervalKey)

	// Get the collection interval
	gpuMetricsCollectionInterval := GetAcceleratedComputeGPUMetricsCollectionInterval(conf)
	isHighFrequency := gpuMetricsCollectionInterval < DefaultMetricsCollectionInterval

	// Log the configuration details
	if gpuMetricsCollectionIntervalExists {
		// Log that the config exists and its value
		log.Printf("[DEBUG] accelerated_compute_gpu_metrics_collection_interval exists with value: %s", gpuMetricsCollectionInterval.String())
	} else {
		log.Printf("[DEBUG] accelerated_compute_gpu_metrics_collection_interval does not exist, using default: %s", DefaultMetricsCollectionInterval.String())
	}

	log.Printf("[DEBUG] IsHighFrequencyGPUMetricsEnabled: enhancedEnabled=%v acceleratedEnabled=%v isHighFrequency=%v",
		enhancedEnabled, acceleratedEnabled, isHighFrequency)

	return enhancedEnabled && acceleratedEnabled && isHighFrequency
}
