# Email Operator

## Overview

The Email Operator is a Kubernetes operator written in Go that manages custom resources for configuring email sending and sending emails via transactional email providers like MailerSend and Mailgun. The operator supports cross-namespace operations and demonstrates sending emails from multiple providers.

## Features

- Manages email sending configurations using custom resources.
- Sends emails via MailerSend and Mailgun.
- Securely handles sensitive information using Kubernetes secrets.
- Logs actions for created and updated configurations.
- Updates the status of email resources to reflect delivery status and message ID.

## Prerequisites

- Kubernetes cluster (e.g., Minikube, Kind, or a cloud provider's Kubernetes service)
- kubectl command-line tool
- Docker
- Go programming language
- Kubebuilder

## Setup

### 1. Install Dependencies

Ensure you have the following installed:

- [Go](https://golang.org/doc/install)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html)

### 2. Clone the Repository

git clone https://github.com/capotoa/email-operator.git
cd email-operator

### 3. Build and Push the Docker Image
Build and push the Docker image for the operator:

```
 - make docker-build docker-push IMG=email-operator:latest
```

### Deployment
#### 1. Apply CRDs and RBAC Configurations
Apply the necessary CustomResourceDefinitions (CRDs) and Role-Based Access Control (RBAC) configurations:

```
kubectl apply -k config/default
```

#### 2. Deploy the Operator
Modify the config/manager/manager.yaml file to use your Docker image:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: email-operator
  namespace: system
spec:
  replicas: 1
  selector:
    matchLabels:
      control-plane: controller-manager
  template:
    metadata:
      labels:
        control-plane: controller-manager
    spec:
      containers:
      - name: manager
        image: email-operator:latest
        command:
        - /manager
        args:
        - --enable-leader-election
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
        env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: SECRET_NAME
            value: controller-leader-election-helper
        volumeMounts:
        - mountPath: /tmp/k8s-webhook-server/serving-certs
          name: cert
          readOnly: true
      volumes:
        - name: cert
          secret:
            secretName: webhook-server-cert
```
Apply the manager configuration:


```
kubectl apply -k config/default
```

### Usage

#### 1. Create a Secret for the API Token
First, encode your MailerSend or Mailgun API token in base64:

```
echo -n 'your-api-token' | base64
```

Create a secret YAML file named mailersend-secret.yaml (or mailgun-secret.yaml):
```
apiVersion: v1
kind: Secret
metadata:
  name: mailersend-api-token
  namespace: default
type: Opaque
data:
  apiToken: <base64-encoded-token>
```
Apply the secret:
```
kubectl apply -f mailersend-secret.yaml
```
#### 2. Create the EmailSenderConfig Resource
Create a file named emailsenderconfig.yaml with the following content:

```
apiVersion: example.com/v1
kind: EmailSenderConfig
metadata:
  name: mailersend-config
  namespace: default
spec:
  apiTokenSecretRef: mailersend-api-token
  senderEmail: "devinveslaw@gmail.com"
```
Apply the resource:

```
kubectl apply -f emailsenderconfig.yaml
```
#### 3. Create the Email Resource
Create a file named email.yaml with the following content:

```
apiVersion: example.com/v1
kind: Email
metadata:
  name: test-email
  namespace: default
spec:
  senderConfigRef: mailersend-config  # or mailgun-config
  recipientEmail: "recipient@example.com"
  subject: "Test Email"
  body: "This is a test email."
```
Apply the resource:

```
kubectl apply -f email.yaml
```
#### 4. Verify Email Delivery
Check the status of the Email resource to verify that the email was sent successfully:

```
kubectl get email test-email -o yaml
```
The output should include fields like deliveryStatus, messageId, and error in the status section, indicating the result of the email sending operation.
