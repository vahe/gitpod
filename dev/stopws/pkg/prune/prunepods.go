// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package prune

import (
	"context"
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/gitpod-io/gitpod/stopws/pkg/stop"
)

var (
	namespace                 = "default"
	ContainerIsGoneAnnotation = "gitpod.io/containerIsGone"
	DisposalStatusAnnotation  = "gitpod.io/disposalStatus"
)

func GetPods(ctx context.Context, client kubernetes.Interface) []corev1.Pod {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"component": "workspace"}}

	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
		//FieldSelector: "status.phase=Running",
		//Limit:         100, // TODO
	})
	if err != nil {
		log.Fatal(err)
	}

	var r []corev1.Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		if pod.DeletionTimestamp == nil || time.Now().Sub(pod.DeletionTimestamp.Time) < 1*time.Hour {
			continue
		}
		if _, ok := pod.Annotations[ContainerIsGoneAnnotation]; !ok {
			continue
		}
		if _, ok := pod.Annotations[DisposalStatusAnnotation]; !ok {
			continue
		}
		isReady, err := stop.IsWorkspacePodReady(pod)
		if err != nil {
			fmt.Println(pod.Name, "error isReady", err)
			continue
		}
		if isReady {
			continue
		}
		r = append(r, pod)
	}
	return r
}

func ListPods() {
	ctx := context.Background()

	client, err := stop.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range GetPods(ctx, client) {
		fmt.Println(pod.Name, pod.DeletionTimestamp)
		//fmt.Println(pod.Name, pod.Annotations[ContainerIsGoneAnnotation])
		//fmt.Println(pod.Name, pod.Annotations[DisposalStatusAnnotation])
	}
}

func PrunePods() {
	ctx := context.Background()

	client, err := stop.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	grace := int64(0)
	for _, pod := range GetPods(ctx, client) {
		fmt.Println(pod.Name, pod.DeletionTimestamp)

		err = client.CoreV1().Pods(namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{GracePeriodSeconds: &grace})
		if err != nil {
			fmt.Println("error force deleting pod")
		}
	}
}