/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	echov1alpha1 "github.com/martineg/echo-operator/api/v1alpha1"
)

// EchoReconciler reconciles a Echo object
type EchoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=echo.martineg.net,resources=echoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=echo.martineg.net,resources=echoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=echo.martineg.net,resources=echoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// The reconcile function follows this pattern:
// 1. Fetch the Echo resource
// 2. Handle deletion cases (resource not found)
// 3. Implement the state machine based on current phase
// 4. Update status to reflect current state
// 5. Return appropriate result (success, error, or requeue)
func (r *EchoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Log the reconcile request for debugging
	log.Info("Starting reconcile", "echo", req.NamespacedName)

	// Step 1: Fetch the Echo resource
	echo := &echov1alpha1.Echo{}
	err := r.Get(ctx, req.NamespacedName, echo)
	if err != nil {
		if errors.IsNotFound(err) {
			// Resource was deleted, nothing to do
			log.Info("Echo resource not found. Assuming it was deleted", "echo", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		log.Error(err, "Failed to get Echo resource", "echo", req.NamespacedName)
		return ctrl.Result{}, err
	}

	// Log successful fetch
	log.Info("Successfully fetched Echo resource",
		"echo", req.NamespacedName,
		"message", echo.Spec.Message,
		"currentPhase", echo.Status.Phase)

	// Step 2: Implement state machine logic based on current phase
	switch echo.Status.Phase {
	case "": // Empty phase means new resource
		log.Info("New Echo resource detected, setting phase to Pending")
		return r.handleNewEcho(ctx, echo)

	case echov1alpha1.EchoPhasePending:
		log.Info("Echo is pending, creating Job")
		return r.handlePendingEcho(ctx, echo)

	case echov1alpha1.EchoPhaseRunning:
		log.Info("Echo is running, checking Job status")
		return r.handleRunningEcho(ctx, echo)

	case echov1alpha1.EchoPhaseCompleted:
		log.Info("Echo completed, no action needed")
		return ctrl.Result{}, nil

	case echov1alpha1.EchoPhaseFailed:
		log.Info("Echo failed, no action needed")
		return ctrl.Result{}, nil

	default:
		// Unknown phase - this shouldn't happen with our enum validation
		log.Info("Echo has unknown phase, resetting to Pending", "phase", echo.Status.Phase)
		return r.handleNewEcho(ctx, echo)
	}
}

// handleNewEcho sets the initial status for a new Echo resource
func (r *EchoReconciler) handleNewEcho(ctx context.Context, echo *echov1alpha1.Echo) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Update status to Pending
	echo.Status.Phase = echov1alpha1.EchoPhasePending
	echo.Status.Message = "Echo resource created, preparing to execute"

	// Update the status subresource
	err := r.Status().Update(ctx, echo)
	if err != nil {
		log.Error(err, "Failed to update Echo status to Pending")
		return ctrl.Result{}, err
	}

	log.Info("Echo status updated to Pending")
	// Requeue immediately to handle the Pending phase
	return ctrl.Result{Requeue: true}, nil
}

// handlePendingEcho creates a Job for the Echo and updates status to Running
func (r *EchoReconciler) handlePendingEcho(ctx context.Context, echo *echov1alpha1.Echo) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO: Create Kubernetes Job here
	// For now, we'll just simulate job creation and move to Running

	jobName := fmt.Sprintf("echo-job-%s", echo.Name)
	log.Info("Would create Job", "jobName", jobName, "message", echo.Spec.Message)

	// Update status to Running
	echo.Status.Phase = echov1alpha1.EchoPhaseRunning
	echo.Status.JobName = jobName
	echo.Status.Message = fmt.Sprintf("Job %s created for message: %s", jobName, echo.Spec.Message)

	err := r.Status().Update(ctx, echo)
	if err != nil {
		log.Error(err, "Failed to update Echo status to Running")
		return ctrl.Result{}, err
	}

	log.Info("Echo status updated to Running", "jobName", jobName)
	// TODO: Requeue to check job status later
	// For now, we'll simulate immediate completion
	return ctrl.Result{Requeue: true}, nil
}

// handleRunningEcho checks the Job status and updates Echo status accordingly
func (r *EchoReconciler) handleRunningEcho(ctx context.Context, echo *echov1alpha1.Echo) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO: Check actual Job status here
	// For now, we'll simulate successful completion

	log.Info("Checking Job status", "jobName", echo.Status.JobName)

	// Simulate job completion
	echo.Status.Phase = echov1alpha1.EchoPhaseCompleted
	echo.Status.Message = fmt.Sprintf("Job %s completed successfully. Message '%s' was echoed.",
		echo.Status.JobName, echo.Spec.Message)

	err := r.Status().Update(ctx, echo)
	if err != nil {
		log.Error(err, "Failed to update Echo status to Completed")
		return ctrl.Result{}, err
	}

	log.Info("Echo status updated to Completed")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EchoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&echov1alpha1.Echo{}).
		Named("echo").
		Complete(r)
}
