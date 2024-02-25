package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNamespaces(
	ctx context.Context,
	clientset *kubernetes.Clientset,
) (
	[]string,
	error,
) {

	list, err := clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces: %w", err)
	}

	namespaces := []string{}

	for _, item := range list.Items {
		namespaces = append(namespaces, item.Name)
	}

	return namespaces, nil
}
