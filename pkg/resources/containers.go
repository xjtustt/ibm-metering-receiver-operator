//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Linter doesn't like "Secret" in string var names that are assigned a value,
// so use concatenation to create the value.
// Example:  const MySecretName = "metering-secret" + ""

package resources

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// SecretCheckData contains info about additional secrets for the secret-check container.
// Names will be added to the SECRET_LIST env var.
// Dirs will be added to the SECRET_DIR_LIST env var.
// VolumeMounts contains the volume mounts associated with the secrets.
type SecretCheckData struct {
	Names        string
	Dirs         string
	VolumeMounts []corev1.VolumeMount
}

const DefaultImageRegistry = "quay.io/opencloudio"
const DefaultReceiverImageName = "metering-data-manager"

const DefaultClusterIssuer = "cs-ca-clusterissuer"

// starting with Common Services 3.4, images can be pulled by SHA or tag.
// run scripts/get-image-sha.sh to update operator.yaml with the SHA values.
// a SHA value looks like this: "sha256:nnnnnnnn"
// a tag value looks like this: "3.5.0".
const DefaultReceiverImageTag = "3.6.0"

// define the env vars that contain either the SHA or the tag
const VarImageSHAforReceiver = "IMAGE_SHA_OR_TAG_RECEIVER"

// use concatenation so linter won't complain about "Secret" vars
const DefaultAPIKeySecretName = "icp-serviceid-apikey-secret" + ""
const DefaultPlatformOidcSecretName = "platform-oidc-credentials" + ""

var TrueVar = true
var FalseVar = false
var Replica1 int32 = 1
var Seconds60 int64 = 60

var cpu100 = resource.NewMilliQuantity(100, resource.DecimalSI)          // 100m
var cpu1000 = resource.NewMilliQuantity(1000, resource.DecimalSI)        // 1000m
var memory100 = resource.NewQuantity(100*1024*1024, resource.BinarySI)   // 100Mi
var memory256 = resource.NewQuantity(256*1024*1024, resource.BinarySI)   // 256Mi
var memory2560 = resource.NewQuantity(2560*1024*1024, resource.BinarySI) // 2560Mi

const DefaultClusterName = "mycluster"

var ArchitectureList = []string{
	"amd64",
	"ppc64le",
	"s390x",
}

var CommonEnvVars = []corev1.EnvVar{
	{
		Name:  "NODE_TLS_REJECT_UNAUTHORIZED",
		Value: "0",
	},
}

var SecretCheckCmd = `set -- $SECRET_LIST; ` +
	`for secretDirName in $SECRET_DIR_LIST; do` +
	`  while true; do` +
	`    echo ` + "`date`" + `: Checking for secret $1;` +
	`    ls /sec/$secretDirName/* && break;` +
	`    echo ` + "`date`" + `: Required secret $1 not found ... try again in 30s;` +
	`    sleep 30;` +
	`  done;` +
	`  echo ` + "`date`" + `: Secret $1 found;` +
	`  shift; ` +
	`done; ` +
	`echo ` + "`date`" + `: All required secrets exist`

// API certificate definition
const APICertName = "icp-metering-api-ca-cert"
const APICertCommonName = "metering-server"
const APICertDirName = "metering-api"

// use concatenation so linter won't complain about "Secret" vars
const APICertSecretName = "icp-metering-api-secret" + ""
const APICertVolumeName = "icp-metering-api-certs"

var APICertVolumeMount = corev1.VolumeMount{
	Name:      APICertVolumeName,
	MountPath: "/sec/" + APICertDirName,
}

var APICertVolume = corev1.Volume{
	Name: APICertVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName:  APICertSecretName,
			DefaultMode: &DefaultMode,
			Optional:    &TrueVar,
		},
	},
}

// Receiver certificate definition
const ReceiverCertName = "icp-metering-receiver-ca-cert"
const ReceiverCertCommonName = "metering-receiver"
const ReceiverCertDirName = "metering-receiver"

// use concatenation so linter won't complain about "Secret" vars
const ReceiverCertSecretName = "icp-metering-receiver-secret" + ""
const ReceiverCertVolumeName = "icp-metering-receiver-certs"

var ReceiverCertVolumeMountForSecretCheck = corev1.VolumeMount{
	Name:      ReceiverCertVolumeName,
	MountPath: "/sec/" + ReceiverCertDirName,
}
var ReceiverCertVolumeMountForMain = corev1.VolumeMount{
	Name:      ReceiverCertVolumeName,
	MountPath: "/certs/" + ReceiverCertDirName,
}
var ReceiverCertVolume = corev1.Volume{
	Name: ReceiverCertVolumeName,
	VolumeSource: corev1.VolumeSource{
		Secret: &corev1.SecretVolumeSource{
			SecretName:  ReceiverCertSecretName,
			DefaultMode: &DefaultMode,
			Optional:    &TrueVar,
		},
	},
}

var ReceiverSslEnvVars = []corev1.EnvVar{
	{
		Name:  "HC_RECEIVER_SSL_CA",
		Value: "/certs/" + ReceiverCertDirName + "/ca.crt",
	},
	{
		Name:  "HC_RECEIVER_SSL_CERT",
		Value: "/certs/" + ReceiverCertDirName + "/tls.crt",
	},
	{
		Name:  "HC_RECEIVER_SSL_KEY",
		Value: "/certs/" + ReceiverCertDirName + "/tls.key",
	},
}

// Common definitions
var commonInitVolumeMounts = []corev1.VolumeMount{
	{
		Name:      "mongodb-ca-cert",
		MountPath: "/certs/mongodb-ca",
	},
	{
		Name:      "mongodb-client-cert",
		MountPath: "/certs/mongodb-client",
	},
}

var commonInitResources = corev1.ResourceRequirements{
	Limits: map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    *cpu100,
		corev1.ResourceMemory: *memory100},
	Requests: map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    *cpu100,
		corev1.ResourceMemory: *memory100},
}

var commonSecurityContext = corev1.SecurityContext{
	AllowPrivilegeEscalation: &FalseVar,
	Privileged:               &FalseVar,
	ReadOnlyRootFilesystem:   &TrueVar,
	RunAsNonRoot:             &TrueVar,
	Capabilities: &corev1.Capabilities{
		Drop: []corev1.Capability{
			"ALL",
		},
	},
}

var CommonMainVolumeMounts = []corev1.VolumeMount{
	{
		Name:      "mongodb-ca-cert",
		MountPath: "/certs/mongodb-ca",
	},
	{
		Name:      "mongodb-client-cert",
		MountPath: "/certs/mongodb-client",
	},
}

var LoglevelVolumeMount = corev1.VolumeMount{
	Name:      "loglevel",
	MountPath: "/etc/config",
}

var Log4jsVolumeMount = corev1.VolumeMount{
	Name:      "log4js",
	MountPath: "/etc/config",
}

//***************************************************************
// Container definitions
//***************************************************************

var ReceiverMainContainer = corev1.Container{
	Image:           "metering-data-manager",
	Name:            "metering-receiver",
	ImagePullPolicy: corev1.PullAlways,
	VolumeMounts: []corev1.VolumeMount{
		LoglevelVolumeMount,
	},
	// CommonEnvVars, IAMEnvVars and mongoDBEnvVars will be added by the controller.
	// HC_DM_MCM_RECEIVER_ENABLED will be set by BuildReceiverEnvVars().
	// Removed ICP_API_KEY.
	Env: []corev1.EnvVar{
		{
			Name:  "METERING_API_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_USE_HTTPS",
			Value: "false",
		},
		{
			Name:  "HC_DM_MCM_SENDER_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_STORAGEREADER_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_REPORTER2_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_PURGER2_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_PREAGGREGATOR_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_METRICS_ENABLED",
			Value: "false",
		},
		{
			Name:  "HC_DM_SELFMETER_PURGER_ENABLED",
			Value: "false",
		},
	},
	Ports: []corev1.ContainerPort{
		{ContainerPort: 3000},
		{ContainerPort: 5000},
	},
	LivenessProbe: &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/livenessProbe",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 3000,
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 305,
		TimeoutSeconds:      5,
		PeriodSeconds:       300,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	},
	ReadinessProbe: &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/readinessProbe",
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 3000,
				},
				Scheme: corev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: 15,
		TimeoutSeconds:      15,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	},
	Resources: corev1.ResourceRequirements{
		Limits: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    *cpu1000,
			corev1.ResourceMemory: *memory2560},
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    *cpu100,
			corev1.ResourceMemory: *memory256},
	},
	SecurityContext: &commonSecurityContext,
}
