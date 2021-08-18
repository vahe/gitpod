// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package stop

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

var (
	namespace                 = "default"
	ContainerIsGoneAnnotation = "gitpod.io/containerIsGone"
	DisposalStatusAnnotation  = "gitpod.io/disposalStatus"
)

func GetPod(ctx context.Context, client kubernetes.Interface, name string) *corev1.Pod {
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return pod
}

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
		if pod.DeletionTimestamp == nil || time.Now().Sub(pod.DeletionTimestamp.Time) < 1*time.Hour {
			continue
		}
		if _, ok := pod.Annotations[ContainerIsGoneAnnotation]; ok {
			continue
		}
		if _, ok := pod.Annotations[DisposalStatusAnnotation]; ok {
			continue
		}
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		isReady, err := IsWorkspacePodReady(pod)
		if err != nil {
			fmt.Println(pod.Name, "error isReady", err, pod.DeletionTimestamp, pod.Spec.NodeName, pod.Status.Phase)
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

	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range GetPods(ctx, client) {
		fmt.Println(pod.Name, pod.DeletionTimestamp)
	}
}

func Single(name string) {
	ctx := context.Background()

	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	pod := GetPod(ctx, client, name)
	fmt.Println(pod.Name)

	ensurePodGetsDeleted(ctx, client, pod)
}

func All() {
	ctx := context.Background()

	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	pods := GetPods(ctx, client)
	// nodeName := ""
	// if len(pods) > 0 {
	// 	nodeName = pods[0].Spec.NodeName
	// }

	// fmt.Println("nodeName", nodeName)

	i := 0
	for _, pod := range pods {
		// if pod.Spec.NodeName == nodeName {
		if pod.Name != "ws-ddb6e58f-9056-4070-812d-099bc6a3ff46" {
			fmt.Println(pod.Name, pod.DeletionTimestamp)
			ensurePodGetsDeleted(ctx, client, &pod)

			i++
			if i > 50 {
				break
			}
		// }
		}
	}
}

func ensurePodGetsDeleted(ctx context.Context, client kubernetes.Interface, pod *corev1.Pod) {
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		pod.Annotations[ContainerIsGoneAnnotation] = "true"
		_, err := client.CoreV1().Pods(namespace).Update(ctx, pod, metav1.UpdateOptions{})
		return err
	})
	if err != nil {
		fmt.Println("setting annotation failed")
		return
	}
}

func IsWorkspacePodReady(pod corev1.Pod) (bool, error) {
	for _, s := range pod.Status.ContainerStatuses {
		if s.Name == "workspace" {
			return s.Ready, nil
		}
	}
	return false, errors.New("not found")
}

func NewClient() (kubernetes.Interface, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}
