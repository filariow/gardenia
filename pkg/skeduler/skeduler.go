package skeduler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	EnvJobDuration    = "DURATION_IN_SEC"
	EnvValvedAddress  = "VALVED_ADDRESS"
	EnvValvedUnixAddr = "VALVED_ADDRESS_UNIX"
)

type Job struct {
	JobName  string
	Schedule string
	Duration uint64
}

type Skeduler interface {
	AddJob(context.Context, string, uint64) (string, error)
	RemoveJob(context.Context, string) error
	GetJob(context.Context, string) (*Job, error)
	ListJobs(context.Context) ([]Job, error)
}

func New(application string, jobImage string, valvedAddress *string, localAddress *string) (Skeduler, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if valvedAddress == nil && localAddress == nil {
		return nil, fmt.Errorf("one Environment variable among '%s' and '%s' must be defined", EnvValvedAddress, EnvValvedUnixAddr)
	}

	return &skeduler{
		application:   application,
		clientset:     clientset,
		jobImage:      jobImage,
		valvedAddress: valvedAddress,
		localAddress:  localAddress,
	}, nil
}

type skeduler struct {
	application   string
	jobImage      string
	valvedAddress *string
	localAddress  *string
	clientset     *kubernetes.Clientset
}

func (s *skeduler) AddJob(ctx context.Context, schedule string, durationSec uint64) (string, error) {
	startingDeadlineSeconds := int64(60)
	successfulJobsHistoryLimit := int32(10)
	failedJobsHistoryLimit := int32(10)
	suspend := false
	parallelism := int32(1)
	privileged := true

	j := batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.application + "-" + strconv.FormatInt(time.Now().Unix(), 10),
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
							RestartPolicy: "Never",
							Containers: []corev1.Container{
								{
									Name:            s.application,
									Image:           s.jobImage,
									ImagePullPolicy: corev1.PullIfNotPresent,
									SecurityContext: &corev1.SecurityContext{
										Privileged: &privileged,
									},
									Env: []corev1.EnvVar{
										{
											Name:  EnvJobDuration,
											Value: strconv.FormatUint(durationSec, 10),
										},
										{
											Name:  EnvValvedAddress,
											Value: *s.valvedAddress,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if s.localAddress != nil && *s.localAddress != "" {
		la := *s.localAddress
		volumeName := "valvedsock-privileged"
		j.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env = append(
			j.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env,
			corev1.EnvVar{
				Name:  EnvValvedUnixAddr,
				Value: "unix:/var/valved.sock",
			})

		j.Spec.JobTemplate.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: la,
					},
				},
			},
		}

		j.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				MountPath: "/var/valved.sock",
				Name:      volumeName,
			},
		}
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

func (s *skeduler) GetJob(ctx context.Context, name string) (*Job, error) {
	j, err := s.clientset.BatchV1().CronJobs(s.application).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return s.mapCronJobToJob(j)
}

func (s *skeduler) ListJobs(ctx context.Context) ([]Job, error) {
	jj, err := s.clientset.BatchV1().CronJobs(s.application).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	mjj := make([]Job, len(jj.Items))
	for i, j := range jj.Items {
		mj, err := s.mapCronJobToJob(&j)
		if err != nil {
			return nil, err
		}
		mjj[i] = *mj
	}

	return mjj, nil
}

func (s *skeduler) mapCronJobToJob(job *batchv1.CronJob) (*Job, error) {
	ee := job.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env
	v := func() string {
		for _, e := range ee {
			if e.Name == EnvJobDuration {
				return e.Value
			}
		}
		return ""
	}()

	if v == "" {
		return nil, fmt.Errorf("job %s: can not find container with Env var with Name %s", job.GetName(), EnvJobDuration)
	}

	d, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("job %s: can not parse duration from container's Env var %s: %w", job.GetName(), EnvJobDuration, err)
	}

	return &Job{
		JobName:  job.GetName(),
		Schedule: job.Spec.Schedule,
		Duration: d,
	}, nil
}
