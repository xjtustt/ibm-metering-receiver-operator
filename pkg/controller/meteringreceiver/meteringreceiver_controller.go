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

package meteringreceiver

import (
	"context"
	"reflect"
	"time"

	operatorv1alpha1 "github.com/ibm/ibm-metering-receiver-operator/pkg/apis/operator/v1alpha1"
	res "github.com/ibm/ibm-metering-receiver-operator/pkg/resources"
	certmgr "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const meteringReceiverCrType = "meteringreceiver_cr"

var commonVolumes = []corev1.Volume{}

var mongoDBEnvVars = []corev1.EnvVar{}

var log = logf.Log.WithName("controller_meteringreceiver")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MeteringReceiver Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMeteringReceiver{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	reqLogger := log.WithValues("func", "add")

	// Create a new controller
	c, err := controller.New("meteringreceiver-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MeteringReceiver
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.MeteringReceiver{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Deployment and requeue the owner MeteringReceiver
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.MeteringReceiver{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Service" and requeue the owner MeteringReceiver
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.MeteringReceiver{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource "Certificate" and requeue the owner MeteringReceiver
	err = c.Watch(&source.Kind{Type: &certmgr.Certificate{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &operatorv1alpha1.MeteringReceiver{},
	})
	if err != nil {
		reqLogger.Error(err, "Failed to watch Certificate")
		// CertManager might not be installed, so don't fail
		//CS??? return err
	}

	return nil
}

// blank assignment to verify that ReconcileMeteringReceiver implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMeteringReceiver{}

// ReconcileMeteringReceiver reconciles a MeteringReceiver object
type ReconcileMeteringReceiver struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MeteringReceiver object and makes changes based on the state read
// and what is in the MeteringReceiver.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMeteringReceiver) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MeteringReceiver")

	// if we need to create several resources, set a flag so we just requeue one time instead of after each create.
	needToRequeue := false

	// Fetch the MeteringReceiver CR instance
	instance := &operatorv1alpha1.MeteringReceiver{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("MeteringReceiver resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get MeteringReceiver CR")
		return reconcile.Result{}, err
	}

	version := instance.Spec.Version
	reqLogger.Info("got Metering instance, version=" + version)

	// set a default Status value
	if len(instance.Status.PodNames) == 0 {
		instance.Status.PodNames = res.DefaultStatusForCR
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to set MeteringReceiver default status")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Checking Services")
	// Check if the Receiver Services already exist. If not, create new ones.
	err = r.reconcileAllServices(instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Checking Receiver Deployment", "Deployment.Name", res.ReceiverDeploymentName)

	// set common MongoDB env vars based on the instance
	mongoDBEnvVars = res.BuildMongoDBEnvVars(instance.Spec.MongoDB)

	// set common Volumes based on the instance
	commonVolumes = res.BuildCommonVolumes(instance.Spec.MongoDB, res.ReceiverDeploymentName, "loglevel")

	// Check if the Receiver Deployment already exists, if not create a new one
	newReceiverDeployment, err := r.deploymentForReceiver(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	err = res.ReconcileDeployment(r.client, instance.Namespace, res.ReceiverDeploymentName, "Receiver", newReceiverDeployment, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Checking Certificates")
	// Check if the Certificates already exist, if not create new ones
	err = r.reconcileAllCertificates(instance, &needToRequeue)
	if err != nil {
		return reconcile.Result{}, err
	}

	if needToRequeue {
		// one or more resources was created, so requeue the request after 5 seconds
		reqLogger.Info("Requeue the request")
		// tried RequeueAfter but it is ignored because we're watching secondary resources.
		// so sleep instead to allow resources to be created by k8s.
		time.Sleep(5 * time.Second)
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Updating MeteringReceiver status")
	// Update the MeteringReceiver status with the pod names.
	// List the pods for this instance's Deployments

	podNames, err := r.getAllPodNames(instance)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods")
		return reconcile.Result{}, err
	}
	// if no pods were found set the default status
	if len(podNames) == 0 {
		podNames = res.DefaultStatusForCR
	}

	// Update status.PodNames if needed
	if !reflect.DeepEqual(podNames, instance.Status.PodNames) {
		instance.Status.PodNames = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update MeteringReceiver status")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Reconciliation completed")
	// since we updated the status in the MeteringReceiver CR, sleep 5 seconds to allow the CR to be refreshed.
	time.Sleep(5 * time.Second)
	return reconcile.Result{}, nil
}

// Check if the Services already exist. If not, create new ones.
// This function was created to reduce the cyclomatic complexity :)
func (r *ReconcileMeteringReceiver) reconcileAllServices(instance *operatorv1alpha1.MeteringReceiver, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileAllServices")

	reqLogger.Info("Checking Receiver Service", "Service.Name", res.ReceiverServiceName)
	// Check if the Receiver Service already exists, if not create a new one
	newReceiverService, err := r.serviceForReceiver(instance)
	if err != nil {
		return err
	}
	err = res.ReconcileService(r.client, instance.Namespace, res.ReceiverServiceName, "Receiver", newReceiverService, needToRequeue)
	if err != nil {
		return err
	}

	return nil
}

// Check if the Certificates already exist, if not create new ones.
// This function was created to reduce the cyclomatic complexity :)
func (r *ReconcileMeteringReceiver) reconcileAllCertificates(instance *operatorv1alpha1.MeteringReceiver, needToRequeue *bool) error {
	reqLogger := log.WithValues("func", "reconcileAllCertificates")

	certificateList := []res.CertificateData{}
	// need to create the receiver certificate
	certificateList = append(certificateList, res.ReceiverCertificateData)
	for _, certData := range certificateList {
		reqLogger.Info("Checking Certificate", "Certificate.Name", certData.Name)
		newCertificate := res.BuildCertificate(instance.Namespace, instance.Spec.ClusterIssuer, certData)
		// Set Metering instance as the owner and controller of the Certificate
		err := controllerutil.SetControllerReference(instance, newCertificate, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Failed to set owner for Certificate", "Certificate.Namespace", newCertificate.Namespace,
				"Certificate.Name", newCertificate.Name)
			return err
		}
		err = res.ReconcileCertificate(r.client, instance.Namespace, certData.Name, newCertificate, needToRequeue)
		if err != nil {
			return err
		}
	}
	return nil
}

// deploymentForReceiver returns a Receiver Deployment object
func (r *ReconcileMeteringReceiver) deploymentForReceiver(instance *operatorv1alpha1.MeteringReceiver) (*appsv1.Deployment, error) {
	reqLogger := log.WithValues("func", "deploymentForReceiver", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.ReceiverDeploymentName)
	selectorLabels := res.LabelsForSelector(res.ReceiverDeploymentName, meteringReceiverCrType, instance.Name)
	podLabels := res.LabelsForPodMetadata(res.ReceiverDeploymentName, meteringReceiverCrType, instance.Name)

	receiverImage := res.GetImageID(instance.Spec.ImageRegistry, instance.Spec.ImageTagPostfix,
		res.DefaultImageRegistry, res.DefaultReceiverImageName, res.VarImageSHAforReceiver, res.DefaultReceiverImageTag)
	reqLogger.Info("receiverImage=" + receiverImage)

	var additionalInfo res.SecretCheckData
	var additionalInfoPtr *res.SecretCheckData
	// add to the SECRET_LIST env var
	additionalInfo.Names = res.ReceiverCertSecretName
	// add to the SECRET_DIR_LIST env var
	additionalInfo.Dirs = res.ReceiverCertDirName
	// add the volume mount for the receiver cert
	additionalInfo.VolumeMounts = []corev1.VolumeMount{res.ReceiverCertVolumeMountForSecretCheck}
	additionalInfoPtr = &additionalInfo

	receiverSecretCheckContainer := res.BuildSecretCheckContainer(res.ReceiverDeploymentName, receiverImage,
		res.SecretCheckCmd, instance.Spec.MongoDB, additionalInfoPtr)

	initEnvVars := []corev1.EnvVar{
		{
			Name:  "MCM_VERBOSE",
			Value: "true",
		},
	}
	initEnvVars = append(initEnvVars, res.CommonEnvVars...)
	initEnvVars = append(initEnvVars, mongoDBEnvVars...)
	receiverInitContainer := res.BuildInitContainer(res.ReceiverDeploymentName, receiverImage, initEnvVars)

	receiverEnvVars := []corev1.EnvVar{
		{
			Name:  "HC_DM_MCM_RECEIVER_ENABLED",
			Value: "true",
		},
	}

	receiverEnvVars = append(receiverEnvVars, res.ReceiverSslEnvVars...)
	receiverMainContainer := res.ReceiverMainContainer
	receiverMainContainer.Image = receiverImage
	receiverMainContainer.Name = res.ReceiverDeploymentName

	receiverMainContainer.Env = append(receiverMainContainer.Env, receiverEnvVars...)
	receiverMainContainer.Env = append(receiverMainContainer.Env, res.CommonEnvVars...)
	receiverMainContainer.Env = append(receiverMainContainer.Env, mongoDBEnvVars...)

	receiverVolumes := commonVolumes
	receiverMainContainer.VolumeMounts = append(receiverMainContainer.VolumeMounts, res.ReceiverCertVolumeMountForMain)
	receiverVolumes = append(receiverVolumes, res.ReceiverCertVolume)
	receiverMainContainer.VolumeMounts = append(receiverMainContainer.VolumeMounts, res.CommonMainVolumeMounts...)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.ReceiverDeploymentName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &res.Replica1,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      podLabels,
					Annotations: res.AnnotationsForPod(),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:            res.GetServiceAccountName(),
					HostNetwork:                   false,
					HostPID:                       false,
					HostIPC:                       false,
					TerminationGracePeriodSeconds: &res.Seconds60,
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "beta.kubernetes.io/arch",
												Operator: corev1.NodeSelectorOpIn,
												Values:   res.ArchitectureList,
											},
										},
									},
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "dedicated",
							Operator: corev1.TolerationOpExists,
							Effect:   corev1.TaintEffectNoSchedule,
						},
						{
							Key:      "CriticalAddonsOnly",
							Operator: corev1.TolerationOpExists,
						},
					},
					Volumes: receiverVolumes,
					InitContainers: []corev1.Container{
						receiverSecretCheckContainer,
						receiverInitContainer,
					},
					Containers: []corev1.Container{
						receiverMainContainer,
					},
				},
			},
		},
	}
	// Set Metering instance as the owner and controller of the Deployment
	err := controllerutil.SetControllerReference(instance, deployment, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for Receiver Deployment")
		return nil, err
	}
	return deployment, nil
}

// serviceForReceiver returns a Receiver Service object
func (r *ReconcileMeteringReceiver) serviceForReceiver(instance *operatorv1alpha1.MeteringReceiver) (*corev1.Service, error) {
	reqLogger := log.WithValues("func", "serviceForReceiver", "instance.Name", instance.Name)
	metaLabels := res.LabelsForMetadata(res.ReceiverDeploymentName)
	selectorLabels := res.LabelsForSelector(res.ReceiverDeploymentName, meteringReceiverCrType, instance.Name)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      res.ReceiverServiceName,
			Namespace: instance.Namespace,
			Labels:    metaLabels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     "metering-receiver",
					Protocol: corev1.ProtocolTCP,
					Port:     5000,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 5000,
					},
				},
			},
			Selector: selectorLabels,
		},
	}

	// Set Metering instance as the owner and controller of the Service
	err := controllerutil.SetControllerReference(instance, service, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to set owner for Receiver Service")
		return nil, err
	}
	return service, nil
}

// getAllPodNames returns the list of pod names for the associated deployments
func (r *ReconcileMeteringReceiver) getAllPodNames(instance *operatorv1alpha1.MeteringReceiver) ([]string, error) {
	reqLogger := log.WithValues("func", "getAllPodNames")
	// List the pods for this instance's Receiver Deployment
	receiverPodList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(res.LabelsForSelector(res.ReceiverDeploymentName, meteringReceiverCrType, instance.Name)),
	}
	if err := r.client.List(context.TODO(), receiverPodList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "MeteringReceiver.Namespace", instance.Namespace, "Deployment.Name", res.ReceiverDeploymentName)
		return nil, err
	}

	podNames := res.GetPodNames(receiverPodList.Items)

	return podNames, nil
}
