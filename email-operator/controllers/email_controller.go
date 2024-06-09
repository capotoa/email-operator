package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	examplev1 "email-operator/api/v1"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=emails,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=emails/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.com,resources=emails/finalizers,verbs=update

func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Reconciling Email %s/%s\n", req.Namespace, req.Name)

	var email examplev1.Email
	if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var emailSenderConfig examplev1.EmailSenderConfig
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: email.Spec.SenderConfigRef}, &emailSenderConfig); err != nil {
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = fmt.Sprintf("Could not find EmailSenderConfig: %v", err)
		r.Status().Update(ctx, &email)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var apiTokenSecret corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: emailSenderConfig.Spec.APITokenSecretRef}, &apiTokenSecret); err != nil {
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = fmt.Sprintf("Could not find API token secret: %v", err)
		r.Status().Update(ctx, &email)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	apiToken := string(apiTokenSecret.Data["apiToken"])

	err := sendEmail(apiToken, emailSenderConfig.Spec.SenderEmail, email.Spec.RecipientEmail, email.Spec.Subject, email.Spec.Body)
	if err != nil {
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = fmt.Sprintf("Failed to send email: %v", err)
	} else {
		email.Status.DeliveryStatus = "Sent"
		email.Status.MessageID = "some-message-id" // Get the actual message ID from the email provider
	}

	r.Status().Update(ctx, &email)
	return ctrl.Result{}, nil
}

func sendEmail(apiToken, senderEmail, recipientEmail, subject, body string) error {
	mailerSendAPI := "https://api.mailersend.com/v1/email"
	requestBody := fmt.Sprintf(`{
        "from": {
            "email": "%s"
        },
        "to": [{
            "email": "%s"
        }],
        "subject": "%s",
        "html": "%s"
    }`, senderEmail, recipientEmail, subject, body)

	req, err := http.NewRequest("POST", mailerSendAPI, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Failed to send email: %s", string(bodyBytes))
	}

	return nil
}

func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.Email{}).
		Complete(r)
}
