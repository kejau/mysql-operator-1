// Copyright 2018 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
)

// UpdateStatefulSet performs a direct update for the specified StatefulSet.
func UpdateStatefulSet(kubeClient kubernetes.Interface, newData *v1beta1.StatefulSet) (*v1beta1.StatefulSet, error) {
	result, err := kubeClient.AppsV1beta1().StatefulSets(newData.Namespace).Update(newData)
	if err != nil {
		glog.Errorf("Failed to update StatefulSet: %v", err)
		return nil, err
	}

	return result, nil
}

// PatchStatefulSet performs a direct patch update for the specified StatefulSet.
func PatchStatefulSet(kubeClient kubernetes.Interface, oldData *v1beta1.StatefulSet, newData *v1beta1.StatefulSet) (*v1beta1.StatefulSet, error) {
	originalJSON, err := json.Marshal(oldData)
	if err != nil {
		return nil, err
	}

	updatedJSON, err := json.Marshal(newData)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(
		originalJSON, updatedJSON, v1beta1.StatefulSet{})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Patching StatefulSet %s/%s: %s", oldData.Name, oldData.ObjectMeta.Namespace, string(patchBytes))

	result, err := kubeClient.AppsV1beta1().StatefulSets(oldData.Namespace).Patch(oldData.Name, types.StrategicMergePatchType, patchBytes)
	if err != nil {
		glog.Errorf("Failed to patch StatefulSet: %v", err)
		return nil, err
	}

	return result, nil
}

// UpdatePod performs a direct update for the specified Pod.
func UpdatePod(kubeClient kubernetes.Interface, newData *v1.Pod) (*v1.Pod, error) {
	result, err := kubeClient.CoreV1().Pods(newData.Namespace).Update(newData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update pod")
	}

	return result, nil
}

// PatchPod perform a direct patch update for the specified Pod.
func PatchPod(kubeClient kubernetes.Interface, oldData *v1.Pod, newData *v1.Pod) (*v1.Pod, error) {
	currentPodJSON, err := json.Marshal(oldData)
	if err != nil {
		return nil, err
	}

	updatedPodJSON, err := json.Marshal(newData)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(currentPodJSON, updatedPodJSON, v1.Pod{})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Patching Pod %s/%s: %s", oldData.Name, oldData.Namespace, string(patchBytes))

	result, err := kubeClient.CoreV1().Pods(oldData.Namespace).Patch(oldData.Name, types.StrategicMergePatchType, patchBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to patch pod")
	}

	return result, nil
}
