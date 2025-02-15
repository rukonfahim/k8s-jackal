package main

import (
	"context"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr" 
	//"k8s.io/client-go/kubernetes" 
	//"k8s.io/client-go/rest"	
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CTFOperator struct {
	K8sClient client.Client
}

func (o *CTFOperator) DeployCTFContainer(name, image string, sshPort int) error {
	ctx := context.Background()

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{ContainerPort: int32(sshPort)},
							},
						},
					},
				},
			},
		},
	}

	if err := o.K8sClient.Create(ctx, deployment); err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": name},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(sshPort),
					TargetPort: intstrFromInt(sshPort),
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	if err := o.K8sClient.Create(ctx, service); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	fmt.Printf("CTF container deployed. Access it via SSH: ssh user@<node-ip> -p %d\n", sshPort)
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
func intstrFromInt(i int) intstr.IntOrString {
	return intstr.IntOrString{Type: intstr.Int, IntVal: int32(i)}
}

func main() {
	operator := &CTFOperator{}
	if err := operator.DeployCTFContainer("ctf-challenge", "your-ctf-image", 2222); err != nil {
		fmt.Println("Error deploying CTF container:", err)
		os.Exit(1)
	}
}
