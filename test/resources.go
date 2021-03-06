// Copyright (c) 2020 Red Hat, Inc.
//
// The Universal Permissive License (UPL), Version 1.0
//
// Subject to the condition set forth below, permission is hereby granted to any
// person obtaining a copy of this software, associated documentation and/or data
// (collectively the "Software"), free of charge and under any and all copyright
// rights in the Software, and any and all patent rights owned or freely
// licensable by each licensor hereunder covering either (i) the unmodified
// Software as contributed to or provided by such licensor, or (ii) the Larger
// Works (as defined below), to deal in both
//
// (a) the Software, and
// (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
// one is included with the Software (each a "Larger Work" to which the Software
// is contributed by such licensors),
//
// without restriction, including without limitation the rights to copy, create
// derivative works of, display, perform, and distribute the Software and make,
// use, sell, offer for sale, import, export, have made, and have sold the
// Software and the Larger Work(s), and to sublicense the foregoing rights on
// either these or other terms.
//
// This license is subject to the following condition:
// The above copyright notice and either this complete permission notice or at
// a minimum a reference to the UPL must be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package test

import (
	"time"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"github.com/onsi/gomega"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/rh-jmc-team/container-jfr-operator/pkg/apis"
	rhjmcv1beta1 "github.com/rh-jmc-team/container-jfr-operator/pkg/apis/rhjmc/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func NewTestScheme() *runtime.Scheme {
	s := scheme.Scheme

	// Add all APIs used by the operator to the scheme
	err := apis.AddToScheme(s)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	err = certv1.AddToScheme(s)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	err = routev1.AddToScheme(s)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	return s
}

func NewContainerJFR() *rhjmcv1beta1.ContainerJFR {
	return &rhjmcv1beta1.ContainerJFR{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr",
			Namespace: "default",
		},
		Spec: rhjmcv1beta1.ContainerJFRSpec{
			Minimal: false,
		},
	}
}

func NewMinimalContainerJFR() *rhjmcv1beta1.ContainerJFR {
	return &rhjmcv1beta1.ContainerJFR{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr",
			Namespace: "default",
		},
		Spec: rhjmcv1beta1.ContainerJFRSpec{
			Minimal: true,
		},
	}
}

func NewFlightRecorder() *rhjmcv1beta1.FlightRecorder {
	return newFlightRecorder(&rhjmcv1beta1.JMXAuthSecret{
		SecretName: "test-jmx-auth",
	})
}

func NewFlightRecorderNoJMXAuth() *rhjmcv1beta1.FlightRecorder {
	return newFlightRecorder(nil)
}

func NewFlightRecorderBadJMXUserKey() *rhjmcv1beta1.FlightRecorder {
	key := "not-username"
	return newFlightRecorder(&rhjmcv1beta1.JMXAuthSecret{
		SecretName:  "test-jmx-auth",
		UsernameKey: &key,
	})
}

func NewFlightRecorderBadJMXPassKey() *rhjmcv1beta1.FlightRecorder {
	key := "not-password"
	return newFlightRecorder(&rhjmcv1beta1.JMXAuthSecret{
		SecretName:  "test-jmx-auth",
		PasswordKey: &key,
	})
}

func NewFlightRecorderForCJFR() *rhjmcv1beta1.FlightRecorder {
	userKey := "CONTAINER_JFR_RJMX_USER"
	passKey := "CONTAINER_JFR_RJMX_PASS"
	recorder := newFlightRecorder(&rhjmcv1beta1.JMXAuthSecret{
		SecretName:  "containerjfr-jmx-auth",
		UsernameKey: &userKey,
		PasswordKey: &passKey,
	})
	recorder.Name = "containerjfr-pod"
	recorder.Labels = map[string]string{"app": "containerjfr-pod"}
	recorder.OwnerReferences[0].Name = "containerjfr-pod"
	recorder.Spec.RecordingSelector.MatchLabels = map[string]string{"rhjmc.redhat.com/flightrecorder": "containerjfr-pod"}
	return recorder
}

func newFlightRecorder(jmxAuth *rhjmcv1beta1.JMXAuthSecret) *rhjmcv1beta1.FlightRecorder {
	return &rhjmcv1beta1.FlightRecorder{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FlightRecorder",
			APIVersion: "rhjmc.redhat.com/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test-pod",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Pod",
					Name:       "test-pod",
					UID:        "",
				},
			},
		},
		Spec: rhjmcv1beta1.FlightRecorderSpec{
			JMXCredentials:    jmxAuth,
			RecordingSelector: metav1.AddLabelToSelector(&metav1.LabelSelector{}, rhjmcv1beta1.RecordingLabel, "test-pod"),
		},
		Status: rhjmcv1beta1.FlightRecorderStatus{
			Target: &corev1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Pod",
				Name:       "test-pod",
				Namespace:  "default",
			},
			Port: 8001,
		},
	}
}

func NewRecording() *rhjmcv1beta1.Recording {
	return newRecording(getDuration(false), nil, nil, false)
}

func NewContinuousRecording() *rhjmcv1beta1.Recording {
	return newRecording(getDuration(true), nil, nil, false)
}

func NewRunningRecording() *rhjmcv1beta1.Recording {
	running := rhjmcv1beta1.RecordingStateRunning
	return newRecording(getDuration(false), &running, nil, false)
}

func NewRunningContinuousRecording() *rhjmcv1beta1.Recording {
	running := rhjmcv1beta1.RecordingStateRunning
	return newRecording(getDuration(true), &running, nil, false)
}

func NewRecordingToStop() *rhjmcv1beta1.Recording {
	running := rhjmcv1beta1.RecordingStateRunning
	stopped := rhjmcv1beta1.RecordingStateStopped
	return newRecording(getDuration(true), &running, &stopped, false)
}

func NewStoppedRecordingToArchive() *rhjmcv1beta1.Recording {
	stopped := rhjmcv1beta1.RecordingStateStopped
	return newRecording(getDuration(false), &stopped, nil, true)
}

func NewRecordingToStopAndArchive() *rhjmcv1beta1.Recording {
	running := rhjmcv1beta1.RecordingStateRunning
	stopped := rhjmcv1beta1.RecordingStateStopped
	return newRecording(getDuration(true), &running, &stopped, true)
}

func NewArchivedRecording() *rhjmcv1beta1.Recording {
	stopped := rhjmcv1beta1.RecordingStateStopped
	rec := newRecording(getDuration(false), &stopped, nil, true)
	savedDownloadURL := "http://path/to/saved-test-recording.jfr"
	savedReportURL := "http://path/to/saved-test-recording.html"
	rec.Status.DownloadURL = &savedDownloadURL
	rec.Status.ReportURL = &savedReportURL
	return rec
}

func NewDeletedArchivedRecording() *rhjmcv1beta1.Recording {
	rec := NewArchivedRecording()
	delTime := metav1.Unix(0, 1598045501618*int64(time.Millisecond))
	rec.DeletionTimestamp = &delTime
	return rec
}

func newRecording(duration time.Duration, currentState *rhjmcv1beta1.RecordingState,
	requestedState *rhjmcv1beta1.RecordingState, archive bool) *rhjmcv1beta1.Recording {
	finalizers := []string{}
	status := rhjmcv1beta1.RecordingStatus{}
	if currentState != nil {
		downloadUrl := "http://path/to/test-recording.jfr"
		reportUrl := "http://path/to/test-recording.html"
		finalizers = append(finalizers, "recording.finalizer.rhjmc.redhat.com")
		status = rhjmcv1beta1.RecordingStatus{
			State:       currentState,
			StartTime:   metav1.Unix(0, 1597090030341*int64(time.Millisecond)),
			Duration:    metav1.Duration{Duration: duration},
			DownloadURL: &downloadUrl,
			ReportURL:   &reportUrl,
		}
	}
	return &rhjmcv1beta1.Recording{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "my-recording",
			Namespace:  "default",
			Finalizers: finalizers,
		},
		Spec: rhjmcv1beta1.RecordingSpec{
			Name: "test-recording",
			EventOptions: []string{
				"jdk.socketRead:enabled=true",
				"jdk.socketWrite:enabled=true",
			},
			Duration: metav1.Duration{Duration: duration},
			State:    requestedState,
			Archive:  archive,
			FlightRecorder: &corev1.LocalObjectReference{
				Name: "test-pod",
			},
		},
		Status: status,
	}
}

func getDuration(continuous bool) time.Duration {
	seconds := 0
	if !continuous {
		seconds = 30
	}
	return time.Duration(seconds) * time.Second
}

func NewTargetPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			PodIP: "1.2.3.4",
		},
	}
}

func NewContainerJFRPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			PodIP: "1.2.3.4",
		},
	}
}

func NewTestEndpoints() *corev1.Endpoints {
	target := &corev1.ObjectReference{
		Kind:      "Pod",
		Name:      "test-pod",
		Namespace: "default",
	}
	ports := []corev1.EndpointPort{
		{
			Name: "jfr-jmx",
			Port: 1234,
		},
		{
			Name: "other-port",
			Port: 9091,
		},
	}
	return newTestEndpoints(target, ports)
}

func NewTestEndpointsNoTargetRef() *corev1.Endpoints {
	ports := []corev1.EndpointPort{
		{
			Name: "jfr-jmx",
			Port: 1234,
		},
		{
			Name: "other-port",
			Port: 9091,
		},
	}
	return newTestEndpoints(nil, ports)
}

func NewTestEndpointsNoPorts() *corev1.Endpoints {
	target := &corev1.ObjectReference{
		Kind:      "Pod",
		Name:      "test-pod",
		Namespace: "default",
	}
	return newTestEndpoints(target, nil)
}

func NewTestEndpointsNoJMXPort() *corev1.Endpoints {
	target := &corev1.ObjectReference{
		Kind:      "Pod",
		Name:      "test-pod",
		Namespace: "default",
	}
	ports := []corev1.EndpointPort{
		{
			Name: "other-port",
			Port: 9091,
		},
	}
	return newTestEndpoints(target, ports)
}

func newTestEndpoints(targetRef *corev1.ObjectReference, ports []corev1.EndpointPort) *corev1.Endpoints {
	return &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-svc",
			Namespace: "default",
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP:        "1.2.3.4",
						Hostname:  "test-pod",
						TargetRef: targetRef,
					},
				},
				Ports: ports,
			},
		},
	}
}

func NewContainerJFREndpoints() *corev1.Endpoints {
	return &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr",
			Namespace: "default",
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP:       "1.2.3.4",
						Hostname: "containerjfr-pod",
						TargetRef: &corev1.ObjectReference{
							Kind:      "Pod",
							Name:      "containerjfr-pod",
							Namespace: "default",
						},
					},
				},
				Ports: []corev1.EndpointPort{
					{
						Name: "jfr-jmx",
						Port: 1234,
					},
				},
			},
		},
	}
}

func NewContainerJFRService() *corev1.Service {
	c := true
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: rhjmcv1beta1.SchemeGroupVersion.String(),
					Kind:       "ContainerJFR",
					Name:       "containerjfr",
					UID:        "",
					Controller: &c,
				},
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "1.2.3.4",
			Ports: []corev1.ServicePort{
				{
					Name: "export",
					Port: 8181,
				},
			},
		},
	}
}

func NewTestService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-svc",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "1.2.3.4",
			Ports: []corev1.ServicePort{
				{
					Name: "test",
					Port: 8181,
				},
			},
		},
	}
}

func NewCACert() *certv1.Certificate {
	return &certv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr-ca",
			Namespace: "default",
		},
		Spec: certv1.CertificateSpec{
			SecretName: "containerjfr-ca",
		},
	}
}

func newCASecret(certData []byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr-ca",
			Namespace: "default",
		},
		Data: map[string][]byte{
			corev1.TLSCertKey: certData,
		},
	}
}

func NewJMXAuthSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-jmx-auth",
			Namespace: "default",
		},
		Data: map[string][]byte{
			rhjmcv1beta1.DefaultUsernameKey: []byte("hello"),
			rhjmcv1beta1.DefaultPasswordKey: []byte("world"),
		},
	}
}

func NewJMXAuthSecretForCJFR() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "containerjfr-jmx-auth",
			Namespace: "default",
		},
		Data: map[string][]byte{
			rhjmcv1beta1.DefaultUsernameKey: []byte("hello"),
			rhjmcv1beta1.DefaultPasswordKey: []byte("world"),
		},
	}
}
