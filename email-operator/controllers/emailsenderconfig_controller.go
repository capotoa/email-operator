package controllers

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	examplev1 "email-operator/api/v1"
)

// EmailSenderConfigReconciler reconciles a EmailSenderConfig object
type EmailSenderConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=emailsenderconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=emailsenderconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.com,resources=emailsenderconfigs/finalizers,verbs=update

func (r *EmailSenderConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Reconciling EmailSenderConfig %s/%s\n", req.Namespace, req.Name)

	var emailSenderConfig examplev1.EmailSenderConfig
	if err := r.Get(ctx, req.NamespacedName, &emailSenderConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Printf("EmailSenderConfig %s/%s: created/updated\n", req.Namespace, req.Name)

	// Perform any additional logic here (e.g., validation, external API calls)

	return ctrl.Result{}, nil
}

func (r *EmailSenderConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.EmailSenderConfig{}).
		Complete(r)
}
