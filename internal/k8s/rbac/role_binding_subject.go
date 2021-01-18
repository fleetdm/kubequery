/**
 * Copyright (c) 2020-present, The kubequery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package rbac

import (
	"context"

	"github.com/Uptycs/kubequery/internal/k8s"
	"github.com/kolide/osquery-go/plugin/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type roleBindingSubject struct {
	k8s.CommonNamespacedFields
	RoleName         string
	RoleKind         string
	SubjectName      string
	SubjectKind      string
	SubjectNamespace string
}

// RoleBindingSubjectColumns TODO
func RoleBindingSubjectColumns() []table.ColumnDefinition {
	return k8s.GetSchema(&roleBindingSubject{})
}

// RoleBindingSubjectsGenerate TODO
func RoleBindingSubjectsGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	options := metav1.ListOptions{}
	results := make([]map[string]string, 0)

	for {
		rbs, err := k8s.GetClient().RbacV1().RoleBindings(metav1.NamespaceAll).List(context.TODO(), options)
		if err != nil {
			return nil, err
		}

		for _, rb := range rbs.Items {
			for _, s := range rb.Subjects {
				item := &roleBindingSubject{
					CommonNamespacedFields: k8s.GetCommonNamespacedFields(rb.ObjectMeta),
					RoleName:               rb.RoleRef.Name,
					RoleKind:               rb.RoleRef.Kind,
					SubjectName:            s.Name,
					SubjectKind:            s.Kind,
					SubjectNamespace:       s.Namespace,
				}
				results = append(results, k8s.ToMap(item))
			}
		}

		if rbs.Continue == "" {
			break
		}
		options.Continue = rbs.Continue
	}

	return results, nil
}