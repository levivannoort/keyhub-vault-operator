// Copyright 2021 Topicus Security BV
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"

	keyhubv1alpha1 "github.com/topicuskeyhub/keyhub-vault-operator/api/v1alpha1"
	controllerMetrics "github.com/topicuskeyhub/keyhub-vault-operator/controllers/metrics"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	controllers_test "github.com/topicuskeyhub/keyhub-vault-operator/controllers/test"
)

var _ = Describe("KeyHubSecret Controller", func() {

	const timeout = time.Second * 10
	const interval = time.Second * 1

	var manifestToLog interface{}

	BeforeEach(func() {
		manifestToLog = nil

		// Flush caches
		policyEngine.Flush()
		vaultIndexCache.Flush()

		// Failed test runs that don't clean up leave resources behind.
		cfg := controllers_test.BeforeEachInputs{Client: k8sClient}
		controllers_test.CleanUp(&cfg)

		// Reset Prometheus collectors
		controllerMetrics.Reset()
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
		controllers_test.LogManifest(manifestToLog)
	})

	Context("htpasswd secret", func() {
		It("Should handle htpasswd credentials correctly", func() {
			spec := keyhubv1alpha1.KeyHubSecretSpec{
				Template: keyhubv1alpha1.SecretTemplate{
					Type: keyhubv1alpha1.SecretTypeApachePasswordFile,
				},
				Data: []keyhubv1alpha1.SecretKeyReference{
					{Name: "auth", Record: "00000000-0000-0000-1001-000000000002"},
				},
			}

			key := types.NamespacedName{
				Name:      "sample-ks",
				Namespace: "default",
			}

			toCreate := &keyhubv1alpha1.KeyHubSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sample-ks",
					Namespace: "default",
				},
				Spec: spec,
			}

			By("By creating a new KeyHubSecret")
			Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())

			By("By checking the Secret is created correctly")
			fetched := &corev1.Secret{}
			Eventually(func() bool {
				k8sClient.Get(context.Background(), key, fetched)
				manifestToLog = fetched

				if fetched.Type != keyhubv1alpha1.SecretTypeApachePasswordFile ||
					len(fetched.Data) != 1 {
					return false
				}

				auth := strings.Split(strings.TrimSuffix(string(fetched.Data["users"]), "\n"), ":")
				if len(auth) != 2 {
					return false
				}

				return auth[0] == "admin" &&
					bcrypt.CompareHashAndPassword([]byte(auth[1]), []byte("test1234")) == nil
			}, timeout, interval).Should(BeTrue())
			manifestToLog = nil

			By("Deleting the KeyHubSecret and Secret")
			Eventually(func() error {
				f := &keyhubv1alpha1.KeyHubSecret{}
				k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())
			Eventually(func() error {
				f := &corev1.Secret{}
				k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())
		})

		// 	It("Should handle KeyHubSecret updates correctly", func() {
		// 		spec := keyhubv1alpha1.KeyHubSecretSpec{
		// 			Template: keyhubv1alpha1.SecretTemplate{
		// 				Type: corev1.SecretTypeBasicAuth,
		// 			},
		// 			Data: []keyhubv1alpha1.SecretKeyReference{
		// 				{Name: "auth", Record: "1001-0002"},
		// 			},
		// 		}

		// 		key := types.NamespacedName{
		// 			Name:      "sample-ks",
		// 			Namespace: "default",
		// 		}

		// 		toCreate := &keyhubv1alpha1.KeyHubSecret{
		// 			ObjectMeta: metav1.ObjectMeta{
		// 				Name:      "sample-ks",
		// 				Namespace: "default",
		// 			},
		// 			Spec: spec,
		// 		}

		// 		By("By creating a new KeyHubSecret")
		// 		Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())

		// 		By("By checking the Secret is created correctly")
		// 		fetched := &corev1.Secret{}
		// 		Eventually(func() bool {
		// 			k8sClient.Get(context.Background(), key, fetched)
		// 			manifestToLog = fetched

		// 			return fetched.Type == corev1.SecretTypeBasicAuth &&
		// 				len(fetched.Data) == 2 &&
		// 				string(fetched.Data["username"]) == "admin" &&
		// 				string(fetched.Data["password"]) == "test1234"
		// 		}, timeout, interval).Should(BeTrue())
		// 		manifestToLog = nil

		// 		By("By checking the KeyHubSecret status")
		// 		fetchedKeyHubSecret := &keyhubv1alpha1.KeyHubSecret{}
		// 		Eventually(func() bool {
		// 			k8sClient.Get(context.Background(), key, fetchedKeyHubSecret)
		// 			manifestToLog = fetchedKeyHubSecret

		// 			records := fetchedKeyHubSecret.Status.VaultRecordStatuses
		// 			keys := fetchedKeyHubSecret.Status.SecretKeyStatuses

		// 			encKey0 := make([]byte, base64.StdEncoding.EncodedLen(len("admin")))
		// 			base64.StdEncoding.Encode(encKey0, []byte("admin"))

		// 			encKey1 := make([]byte, base64.StdEncoding.EncodedLen(len("test1234")))
		// 			base64.StdEncoding.Encode(encKey1, []byte("test1234"))

		// 			return len(records) == 1 &&
		// 				len(keys) == 2 &&
		// 				records[0].RecordID == "1001-0002" &&
		// 				records[0].Name == "Username + password" &&
		// 				keys[0].Key == "username" &&
		// 				bcrypt.CompareHashAndPassword(keys[0].Hash, encKey0) == nil &&
		// 				keys[1].Key == "password" &&
		// 				bcrypt.CompareHashAndPassword(keys[1].Hash, encKey1) == nil
		// 		}, timeout, interval).Should(BeTrue())
		// 		manifestToLog = nil

		// 		By("By updating the KeyHubSecret")
		// 		fetchedKeyHubSecret.Spec.Data = []keyhubv1alpha1.SecretKeyReference{
		// 			{Name: "auth", Record: "1001-0003"},
		// 		}
		// 		manifestToLog = fetchedKeyHubSecret
		// 		k8sClient.Update(context.Background(), fetchedKeyHubSecret)

		// 		By("By checking the Secret is updated correctly")
		// 		fetched = &corev1.Secret{}
		// 		Eventually(func() bool {
		// 			k8sClient.Get(context.Background(), key, fetched)
		// 			manifestToLog = fetched

		// 			return fetched.Type == corev1.SecretTypeBasicAuth &&
		// 				len(fetched.Data) == 2 &&
		// 				string(fetched.Data["username"]) == "example.io" &&
		// 				string(fetched.Data["password"]) == "test5678"
		// 		}, timeout, interval).Should(BeTrue())
		// 		manifestToLog = nil

		// 		By("By checking the KeyHubSecret status is updated correctly")
		// 		fetchedKeyHubSecret = &keyhubv1alpha1.KeyHubSecret{}
		// 		Eventually(func() bool {
		// 			k8sClient.Get(context.Background(), key, fetchedKeyHubSecret)
		// 			manifestToLog = fetchedKeyHubSecret

		// 			records := fetchedKeyHubSecret.Status.VaultRecordStatuses
		// 			keys := fetchedKeyHubSecret.Status.SecretKeyStatuses

		// 			encKey0 := make([]byte, base64.StdEncoding.EncodedLen(len("example.io")))
		// 			base64.StdEncoding.Encode(encKey0, []byte("example.io"))

		// 			encKey1 := make([]byte, base64.StdEncoding.EncodedLen(len("test5678")))
		// 			base64.StdEncoding.Encode(encKey1, []byte("test5678"))

		// 			return len(records) == 1 &&
		// 				len(keys) == 2 &&
		// 				records[0].RecordID == "1001-0003" &&
		// 				records[0].Name == "Certificate" &&
		// 				keys[0].Key == "username" &&
		// 				bcrypt.CompareHashAndPassword(keys[0].Hash, encKey0) == nil &&
		// 				keys[1].Key == "password" &&
		// 				bcrypt.CompareHashAndPassword(keys[1].Hash, encKey1) == nil
		// 		}, timeout, interval).Should(BeTrue())
		// 		manifestToLog = nil

		// 		By("Deleting the KeyHubSecret and Secret")
		// 		Eventually(func() error {
		// 			f := &keyhubv1alpha1.KeyHubSecret{}
		// 			k8sClient.Get(context.Background(), key, f)
		// 			return k8sClient.Delete(context.Background(), f)
		// 		}, timeout, interval).Should(Succeed())
		// 		Eventually(func() error {
		// 			f := &corev1.Secret{}
		// 			k8sClient.Get(context.Background(), key, f)
		// 			return k8sClient.Delete(context.Background(), f)
		// 		}, timeout, interval).Should(Succeed())
		// 	})

		// 	It("Should not make excessive api calls to KeyHub", func() {
		// 		spec := keyhubv1alpha1.KeyHubSecretSpec{
		// 			Template: keyhubv1alpha1.SecretTemplate{
		// 				Type: corev1.SecretTypeBasicAuth,
		// 				Metadata: keyhubv1alpha1.SecretTemplateMetadata{
		// 					Labels: map[string]string{
		// 						"iteration": "0",
		// 					},
		// 				},
		// 			},
		// 			Data: []keyhubv1alpha1.SecretKeyReference{
		// 				{Name: "auth", Record: "1001-0002"},
		// 			},
		// 		}

		// 		key := types.NamespacedName{
		// 			Name:      "sample-ks",
		// 			Namespace: "default",
		// 		}

		// 		toCreate := &keyhubv1alpha1.KeyHubSecret{
		// 			ObjectMeta: metav1.ObjectMeta{
		// 				Name:      "sample-ks",
		// 				Namespace: "default",
		// 			},
		// 			Spec: spec,
		// 		}

		// 		By("By checking the metrics are reset")
		// 		var keyhubApiRequests float64
		// 		collectedMetrics, _ := metrics.Registry.Gather()
		// 		for _, metricFamily := range collectedMetrics {
		// 			if "keyhub_api_request_total" == *metricFamily.Name {
		// 				for _, metric := range metricFamily.GetMetric() {
		// 					keyhubApiRequests += metric.GetCounter().GetValue()
		// 				}
		// 			}
		// 		}
		// 		Expect(keyhubApiRequests).ToNot(BeNil())
		// 		Expect(keyhubApiRequests).To(Equal(0.0))

		// 		By("By creating a new KeyHubSecret")
		// 		Expect(k8sClient.Create(context.Background(), toCreate)).Should(Succeed())

		// 		By("By checking the Secret is created correctly")
		// 		fetched := &corev1.Secret{}
		// 		Eventually(func() bool {
		// 			k8sClient.Get(context.Background(), key, fetched)
		// 			manifestToLog = fetched

		// 			return fetched.Type == corev1.SecretTypeBasicAuth &&
		// 				len(fetched.Data) == 2 &&
		// 				string(fetched.Data["username"]) == "admin"
		// 		}, timeout, interval).Should(BeTrue())
		// 		manifestToLog = nil

		// 		var actualKeyhubApiCalls float64 = 0.0
		// 		collectedMetrics, _ = metrics.Registry.Gather()
		// 		for _, metricFamily := range collectedMetrics {
		// 			if "keyhub_api_request_total" == *metricFamily.Name {
		// 				for _, metric := range metricFamily.GetMetric() {
		// 					actualKeyhubApiCalls += metric.GetCounter().GetValue()
		// 				}
		// 			}
		// 		}

		// 		fmt.Fprintf(GinkgoWriter, "Initial number of KeyHub API requests: %.0f\n", actualKeyhubApiCalls)
		// 		Expect(actualKeyhubApiCalls).ToNot(Equal(0.0))
		// 		allowedKeyhubApiCalls := actualKeyhubApiCalls

		// 		By("By repeatedly triggering a reconcile")
		// 		for i := 1; i <= 3; i++ {
		// 			By("By setting the 'iteration' label")
		// 			fetchedKeyHubSecret := &keyhubv1alpha1.KeyHubSecret{}
		// 			k8sClient.Get(context.Background(), key, fetchedKeyHubSecret)
		// 			fetchedKeyHubSecret.Spec.Template.Metadata.Labels["iteration"] = strconv.Itoa(i)
		// 			k8sClient.Update(context.Background(), fetchedKeyHubSecret)

		// 			By("By checking the Secret is updated correctly")
		// 			fetched := &corev1.Secret{}
		// 			Eventually(func() bool {
		// 				k8sClient.Get(context.Background(), key, fetched)
		// 				manifestToLog = fetched

		// 				return fetched.Type == corev1.SecretTypeBasicAuth &&
		// 					len(fetched.Data) == 2 &&
		// 					string(fetched.Data["username"]) == "admin" &&
		// 					len(fetched.Labels) == 1 &&
		// 					fetched.Labels["iteration"] == strconv.Itoa(i)
		// 			}, timeout, interval).Should(BeTrue())
		// 			manifestToLog = nil

		// 			By(fmt.Sprintf("Checking the number of KeyHub API calls (iteration %d)", i))
		// 			var actualKeyhubApiCalls float64 = 0.0
		// 			collectedMetrics, _ := metrics.Registry.Gather()
		// 			for _, metricFamily := range collectedMetrics {
		// 				if "keyhub_api_request_total" == *metricFamily.Name {
		// 					for _, metric := range metricFamily.GetMetric() {
		// 						actualKeyhubApiCalls += metric.GetCounter().GetValue()
		// 					}
		// 				}
		// 			}
		// 			fmt.Fprintf(GinkgoWriter, "Total number of KeyHub API requests (iteration %d): expected=%.0f, actual=%.0f\n", i, allowedKeyhubApiCalls, actualKeyhubApiCalls)
		// 			Expect(actualKeyhubApiCalls).To(Equal(allowedKeyhubApiCalls))
		// 		}

		// 		By("Deleting the KeyHubSecret and Secret")
		// 		Eventually(func() error {
		// 			f := &keyhubv1alpha1.KeyHubSecret{}
		// 			k8sClient.Get(context.Background(), key, f)
		// 			return k8sClient.Delete(context.Background(), f)
		// 		}, timeout, interval).Should(Succeed())
		// 		Eventually(func() error {
		// 			f := &corev1.Secret{}
		// 			k8sClient.Get(context.Background(), key, f)
		// 			return k8sClient.Delete(context.Background(), f)
		// 		}, timeout, interval).Should(Succeed())
		// 	})

	})
})
