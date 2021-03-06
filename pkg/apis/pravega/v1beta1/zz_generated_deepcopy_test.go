/**
 * Copyright (c) 2018 Dell Inc., or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */
package v1beta1_test

import (
	"fmt"

	"github.com/pravega/pravega-operator/pkg/apis/pravega/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PravegaCluster DeepCopy", func() {
	Context("with defaults", func() {
		var (
			str1, str2 string
			str3, str4 v1.PullPolicy
			p1, p2     *v1beta1.PravegaCluster
		)
		BeforeEach(func() {
			p1 = &v1beta1.PravegaCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "default",
				},
			}
			p1.Spec = v1beta1.ClusterSpec{
				Pravega: &v1beta1.PravegaSpec{
					CacheVolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						VolumeName: "abc",
					},
					SegmentStoreSecret: &v1beta1.SegmentStoreSecret{
						Secret:    "seg-secret",
						MountPath: "",
					},
					InfluxDBSecret: &v1beta1.InfluxDBSecret{
						Secret:    "influx-secret",
						MountPath: "",
					},
					SegmentStoreInitContainers: []v1.Container{
						v1.Container{
							Name:    "testing",
							Image:   "dummy-image",
							Command: []string{"sh", "-c", "ls;pwd"},
						},
					},
					ControllerInitContainers: []v1.Container{
						v1.Container{
							Name:    "testing",
							Image:   "dummy-image",
							Command: []string{"sh", "-c", "ls;pwd"},
						},
					},
					AuthImplementations: &v1beta1.AuthImplementationSpec{
						MountPath: "/ifs/data",
						AuthHandlers: []v1beta1.AuthHandlerSpec{
							{
								Image:  "testimage1",
								Source: "source1",
							},
							{
								Image:  "testimage2",
								Source: "source2",
							},
						},
					},
				},
			}
			p1.WithDefaults()
			p1.Status.Init()
			no := int64(0)
			securitycontext := corev1.PodSecurityContext{
				RunAsUser: &no,
			}
			p1.Spec.Pravega.SegmentStoreSecurityContext = &securitycontext
			p1.Spec.Pravega.ControllerSecurityContext = &securitycontext
			temp := *p1.DeepCopy()
			p2 = &temp
			str1 = p1.Spec.Pravega.Image.Repository
			str2 = p2.Spec.Pravega.Image.Repository
			p1.Spec.Pravega.Image.Repository = "pravega/exmple"
			p1.Spec.Pravega.DeepCopyInto(p2.Spec.Pravega)
			str3 = p1.Spec.Pravega.Image.PullPolicy
			str4 = p2.Spec.Pravega.Image.PullPolicy
			p1.Spec.Pravega.Image.PullPolicy = "PullIfNotPresent"
			p1.Spec.Pravega.Image.DeepCopyInto(p2.Spec.Pravega.Image)
			p2.Spec.Authentication = p1.Spec.Authentication.DeepCopy()
			p1.Spec.Authentication.Enabled = true
			temp1 := *p1.Spec.Authentication.DeepCopy()
			p2.Spec.Authentication = &temp1
			p1.Spec.Version = "0.4.0"
			p2.Spec = *p1.Spec.DeepCopy()
			p1.Status.SetPodsReadyConditionTrue()
			p2.Status.Conditions[0] = *p1.Status.Conditions[0].DeepCopy()
			p1.Status.VersionHistory = []string{"0.6.0", "0.5.0"}
			p2.Status = *p1.Status.DeepCopy()
			p1.Status.Members.Ready = []string{"bookie-0", "bookie-1"}
			p1.Status.Members.Unready = []string{"bookie-3", "bookie-2"}
			p2.Status.Members = *p1.Status.Members.DeepCopy()
			p1.Spec.ExternalAccess.DomainName = "example.com"
			p2.Spec.ExternalAccess = p1.Spec.ExternalAccess.DeepCopy()
			p1.Spec.TLS.Static.ControllerSecret = "controller-secret"
			p2.Spec.TLS = p1.Spec.TLS.DeepCopy()
			p1.Spec.TLS.Static.SegmentStoreSecret = "segmentstore-secret"
			p2.Spec.TLS.Static = p1.Spec.TLS.Static.DeepCopy()
			p1.Spec.Pravega.Image.Repository = "pravega/exmple"
			p2.Spec.Pravega.Image = p1.Spec.Pravega.Image.DeepCopy()
			p1.Spec.Pravega.LongTermStorage.FileSystem.PersistentVolumeClaim.ClaimName = "fs"
			p2.Spec.Pravega.LongTermStorage.FileSystem = p1.Spec.Pravega.LongTermStorage.FileSystem.DeepCopy()
			p1.Spec.Pravega.Options["key"] = "value"
			p1.Spec.Pravega.SegmentStoreServiceAnnotations["user"] = "test"
			p1.Spec.Pravega.ControllerServiceAnnotations["user"] = "test1"
			p1.Spec.Pravega.ControllerPodLabels["user"] = "test2"
			p1.Spec.Pravega.SegmentStorePodLabels["user"] = "test2"
			p1.Spec.Pravega.ControllerPodAnnotations["user"] = "test2"
			p1.Spec.Pravega.SegmentStorePodAnnotations["user"] = "test2"

			p2.Spec.Pravega = p1.Spec.Pravega.DeepCopy()
			p2.Spec.Pravega.SegmentStoreSecret = p1.Spec.Pravega.SegmentStoreSecret.DeepCopy()
			p2.Spec.Pravega.InfluxDBSecret = p1.Spec.Pravega.InfluxDBSecret.DeepCopy()
			p2.Spec.Pravega.AuthImplementations = p1.Spec.Pravega.AuthImplementations.DeepCopy()
			p2.Spec.Pravega.AuthImplementations.AuthHandlers[0] = *p1.Spec.Pravega.AuthImplementations.AuthHandlers[0].DeepCopy()

			p2.Spec.Pravega.LongTermStorage = p1.Spec.Pravega.LongTermStorage.DeepCopy()
			p1.Spec.Pravega.LongTermStorage = &v1beta1.LongTermStorageSpec{
				Ecs: &v1beta1.ECSSpec{
					ConfigUri:   "configUri",
					Bucket:      "bucket",
					Prefix:      "prefix",
					Credentials: "credentials",
				},
			}
			p2.Spec.Pravega.LongTermStorage.Ecs = p1.Spec.Pravega.LongTermStorage.Ecs.DeepCopy()
			p1.Spec.Pravega.LongTermStorage = &v1beta1.LongTermStorageSpec{
				FileSystem: &v1beta1.FileSystemSpec{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "claim",
						ReadOnly:  true,
					},
				},
			}
			p2.Spec.Pravega.LongTermStorage.FileSystem = p1.Spec.Pravega.LongTermStorage.FileSystem.DeepCopy()

			p1.Spec.Pravega.LongTermStorage = &v1beta1.LongTermStorageSpec{
				Hdfs: &v1beta1.HDFSSpec{
					Uri:               "uri",
					Root:              "root",
					ReplicationFactor: 1,
				},
			}
			p2.Spec.Pravega.LongTermStorage.Hdfs = p1.Spec.Pravega.LongTermStorage.Hdfs.DeepCopy()
		})
		It("value of str1 and str2 should be equal", func() {
			??(str2).To(Equal(str1))
		})
		It("value of str3 and str4 should be equal", func() {
			??(str3).To(Equal(str4))
		})
		It("Authentication enabled should be true", func() {
			??(p2.Spec.Authentication.Enabled).To(Equal(true))
		})

		It("checking spec version", func() {
			??(p2.Spec.Version).To(Equal("0.4.0"))
		})
		It("checking status conditions", func() {
			??(p2.Status.Conditions[0].Reason).To(Equal(p1.Status.Conditions[0].Reason))
		})
		It("checking  version history", func() {
			??(p2.Status.VersionHistory[0]).To(Equal("0.6.0"))
		})
		It("checking ready members", func() {
			??(p2.Status.Members.Ready[0]).To(Equal("bookie-0"))
		})
		It("checking  unready members", func() {
			??(p2.Status.Members.Unready[0]).To(Equal("bookie-3"))
		})
		It("checking  external access domain name", func() {
			??(p2.Spec.ExternalAccess.DomainName).To(Equal("example.com"))
		})
		It("checking  controller secret", func() {
			??(p2.Spec.TLS.Static.ControllerSecret).To(Equal("controller-secret"))
		})
		It("checking  segmentstore secret", func() {
			??(p2.Spec.TLS.Static.SegmentStoreSecret).To(Equal("segmentstore-secret"))
		})
		It("checking  segmentstore secret inside pravega spec", func() {
			??(p2.Spec.Pravega.SegmentStoreSecret.Secret).To(Equal("seg-secret"))
		})
		It("checking  image  repository", func() {
			??(p2.Spec.Pravega.Image.Repository).To(Equal("pravega/exmple"))
		})
		It("checking SegmentStoreSecurityContext", func() {
			??(fmt.Sprintf("%v", *p2.Spec.Pravega.SegmentStoreSecurityContext.RunAsUser)).To(Equal("0"))
		})
		It("checking ControllerSecurityContext", func() {
			??(fmt.Sprintf("%v", *p2.Spec.Pravega.ControllerSecurityContext.RunAsUser)).To(Equal("0"))
		})
		It("checking ControllerInitContainer", func() {
			??(p2.Spec.Pravega.ControllerInitContainers[0].Name).To(Equal("testing"))
		})
		It("checking AuthHandlerDetails", func() {
			??(p2.Spec.Pravega.AuthImplementations.AuthHandlers[0].Image).To(Equal("testimage1"))
			??(p2.Spec.Pravega.AuthImplementations.MountPath).To(Equal("/ifs/data"))

		})
		It("checking SegementStoreInitContainer", func() {
			??(p2.Spec.Pravega.SegmentStoreInitContainers[0].Name).To(Equal("testing"))
		})
		It("checking  pravega options", func() {
			??(p2.Spec.Pravega.Options["key"]).To(Equal("value"))
		})
		It("checking LongTermStorage ECS", func() {
			??(p2.Spec.Pravega.LongTermStorage.Ecs.ConfigUri).To(Equal("configUri"))
		})
		It("checking LongTermStorage Hdfs", func() {
			??(p2.Spec.Pravega.LongTermStorage.Hdfs.Uri).To(Equal("uri"))
		})
		It("checking LongTermStorage Filesystem", func() {
			??(p2.Spec.Pravega.LongTermStorage.FileSystem.PersistentVolumeClaim.ClaimName).To(Equal("claim"))
		})
		It("checking for nil authentication", func() {
			p1.Spec.Authentication = nil
			??(p1.Spec.Authentication.DeepCopy()).Should(BeNil())
		})

		It("checking for nil TLS", func() {
			p1.Spec.TLS = nil
			??(p1.Spec.TLS.DeepCopy()).Should(BeNil())
		})
		It("checking for nil Pravega", func() {
			p1.Spec.Pravega = nil
			??(p1.Spec.Pravega.DeepCopy()).Should(BeNil())
		})

		It("checking for nil External access", func() {
			p1.Spec.ExternalAccess = nil
			??(p1.Spec.ExternalAccess.DeepCopy()).Should(BeNil())
		})
		It("checking for nil TLS Static", func() {
			p1.Spec.TLS.Static = nil
			??(p1.Spec.TLS.Static.DeepCopy()).Should(BeNil())
		})
		It("checking for nil InfluxDBsecret", func() {
			p1.Spec.Pravega.InfluxDBSecret = nil
			??(p1.Spec.Pravega.InfluxDBSecret.DeepCopy()).Should(BeNil())
		})
		It("checking for nil SegmentStore secret", func() {
			p1.Spec.Pravega.SegmentStoreSecret = nil
			??(p1.Spec.Pravega.SegmentStoreSecret.DeepCopy()).Should(BeNil())
		})
		It("checking for nil Pravega LongTermStorage", func() {
			p1.Spec.Pravega.LongTermStorage = nil
			??(p1.Spec.Pravega.LongTermStorage.DeepCopy()).Should(BeNil())
		})
		It("checking for nil Pravega Image", func() {
			p1.Spec.Pravega.Image = nil
			??(p1.Spec.Pravega.Image.DeepCopy()).Should(BeNil())
		})
		It("checking for nil ECS", func() {
			p1.Spec.Pravega.LongTermStorage.Ecs = nil
			??(p1.Spec.Pravega.LongTermStorage.Ecs.DeepCopy()).Should(BeNil())
		})
		It("checking for nil Hdfs", func() {
			p1.Spec.Pravega.LongTermStorage.Hdfs = nil
			??(p1.Spec.Pravega.LongTermStorage.Hdfs.DeepCopy()).Should(BeNil())
		})
		It("checking for nil filesystem", func() {
			p1.Spec.Pravega.LongTermStorage.FileSystem = nil
			??(p1.Spec.Pravega.LongTermStorage.FileSystem.DeepCopy()).Should(BeNil())
		})
		It("checking for nil LongTermStorage", func() {
			p1.Spec.Pravega.LongTermStorage = nil
			??(p1.Spec.Pravega.LongTermStorage.DeepCopy()).Should(BeNil())
		})
		It("checking for nil AuthImplementations", func() {
			p1.Spec.Pravega.AuthImplementations = nil
			??(p1.Spec.Pravega.AuthImplementations.DeepCopy()).Should(BeNil())
		})
		It("checking for nil member status", func() {
			var memberstatus *v1beta1.MembersStatus
			memberstatus2 := memberstatus.DeepCopy()
			??(memberstatus2).To(BeNil())
		})
		It("checking for nil cluster status", func() {
			var clusterstatus *v1beta1.ClusterStatus
			clusterstatus2 := clusterstatus.DeepCopy()
			??(clusterstatus2).To(BeNil())
		})
		It("checking for nil cluster spec", func() {
			var clusterspec *v1beta1.ClusterSpec
			clusterspec2 := clusterspec.DeepCopy()
			??(clusterspec2).To(BeNil())
		})
		It("checking for nil cluster condition", func() {
			var clustercond *v1beta1.ClusterCondition
			clustercond2 := clustercond.DeepCopy()
			??(clustercond2).To(BeNil())
		})
		It("checking for nil pravega cluster", func() {
			var cluster *v1beta1.PravegaCluster
			cluster2 := cluster.DeepCopy()
			??(cluster2).To(BeNil())
		})
		It("checking for nil imagespec", func() {
			var imagespec *v1beta1.ImageSpec
			imagespec2 := imagespec.DeepCopy()
			??(imagespec2).To(BeNil())
		})
		It("checking for nil clusterlist", func() {
			var clusterlist *v1beta1.PravegaClusterList
			clusterlist2 := clusterlist.DeepCopy()
			??(clusterlist2).To(BeNil())
		})
		It("checking for deepcopy for clusterlist", func() {
			var clusterlist v1beta1.PravegaClusterList
			clusterlist.ResourceVersion = "v1beta1"
			clusterlist2 := clusterlist.DeepCopy()
			??(clusterlist2.ResourceVersion).To(Equal("v1beta1"))
		})
		It("checking for Deepcopy object", func() {
			pravega := p2.DeepCopyObject()
			??(pravega.GetObjectKind().GroupVersionKind().Version).To(Equal(""))
		})
		It("checking for nil pravega cluster deepcopyobject", func() {
			var cluster *v1beta1.PravegaCluster
			cluster2 := cluster.DeepCopyObject()
			??(cluster2).To(BeNil())
		})
		It("checking for deepcopyobject for clusterlist", func() {
			var clusterlist v1beta1.PravegaClusterList
			clusterlist.ResourceVersion = "v1beta1"
			clusterlist2 := clusterlist.DeepCopyObject()
			??(clusterlist2).ShouldNot(BeNil())
		})
		It("checking for nil pravega clusterlist deepcopyobject", func() {
			var clusterlist *v1beta1.PravegaClusterList
			clusterlist2 := clusterlist.DeepCopyObject()
			??(clusterlist2).To(BeNil())
		})
		It("checking for deepcopyobject for clusterlist with items", func() {
			var clusterlist v1beta1.PravegaClusterList
			clusterlist.ResourceVersion = "v1beta1"
			clusterlist.Items = []v1beta1.PravegaCluster{
				{
					Spec: v1beta1.ClusterSpec{},
				},
			}
			clusterlist2 := clusterlist.DeepCopyObject()
			??(clusterlist2).ShouldNot(BeNil())
		})
	})
})
