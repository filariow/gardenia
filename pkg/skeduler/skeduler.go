package skeduler

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Skeduler interface {
	AddJob(context.Context, string) (string, error)
	RemoveJob(context.Context, string) error
	GetJob(context.Context, string) (*batchv1.CronJob, error)
	ListJobs(context.Context) (*batchv1.CronJobList, error)
}

func New(application string) (Skeduler, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &skeduler{
		clientset: clientset,
	}, nil
}

type skeduler struct {
	application string
	clientset   *kubernetes.Clientset
}

func (s *skeduler) AddJob(ctx context.Context, schedule string) (string, error) {
	startingDeadlineSeconds := int64(60)
	successfulJobsHistoryLimit := int32(10)
	failedJobsHistoryLimit := int32(10)
	suspend := false
	parallelism := int32(1)

	j := batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.application,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:                   schedule,
			StartingDeadlineSeconds:    &startingDeadlineSeconds,
			ConcurrencyPolicy:          batchv1.ForbidConcurrent,
			Suspend:                    &suspend,
			SuccessfulJobsHistoryLimit: &successfulJobsHistoryLimit,
			FailedJobsHistoryLimit:     &failedJobsHistoryLimit,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: s.application,
				},
				Spec: batchv1.JobSpec{
					Parallelism: &parallelism,
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: s.application,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  s.application,
									Image: s.application,
								},
							},
						},
					},
				},
			},
		},
	}

	cj, err := s.clientset.BatchV1().CronJobs(s.application).Create(ctx, &j, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return cj.GetName(), nil
}

func (s *skeduler) RemoveJob(ctx context.Context, name string) error {
	return s.clientset.BatchV1().CronJobs(s.application).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *skeduler) GetJob(ctx context.Context, name string) (*batchv1.CronJob, error) {
	return s.clientset.BatchV1().CronJobs(s.application).Get(ctx, name, metav1.GetOptions{})
}

func (s *skeduler) ListJobs(ctx context.Context) (*batchv1.CronJobList, error) {
	return s.clientset.BatchV1().CronJobs(s.application).List(ctx, metav1.ListOptions{})
}
