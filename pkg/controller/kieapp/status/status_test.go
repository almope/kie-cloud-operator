package status

import (
	"fmt"
	"testing"

	api "github.com/kiegroup/kie-cloud-operator/pkg/apis/app/v2"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetDeployed(t *testing.T) {
	now := metav1.Now()
	cr := &api.KieApp{}

	assert.True(t, SetDeployed(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	assert.Equal(t, api.DeployedConditionType, cr.Status.Conditions[0].Type)
	assert.Equal(t, api.DeployedConditionType, cr.Status.Phase)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[0].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[0].LastTransitionTime))
}

func TestSetDeployedSkipUpdate(t *testing.T) {
	cr := &api.KieApp{}
	SetDeployed(cr)

	assert.NotEmpty(t, cr.Status.Conditions)
	condition := cr.Status.Conditions[0]

	assert.False(t, SetDeployed(cr))
	assert.Equal(t, 1, len(cr.Status.Conditions))
	assert.Equal(t, condition, cr.Status.Conditions[0])
	assert.Equal(t, condition.Type, cr.Status.Phase)
}

func TestSetProvisioning(t *testing.T) {
	now := metav1.Now()
	cr := &api.KieApp{}
	assert.True(t, SetProvisioning(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	assert.Equal(t, api.ProvisioningConditionType, cr.Status.Conditions[0].Type)
	assert.Equal(t, api.ProvisioningConditionType, cr.Status.Phase)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[0].Status)
	assert.True(t, now.Before(&cr.Status.Conditions[0].LastTransitionTime))
}

func TestSetProvisioningSkipUpdate(t *testing.T) {
	cr := &api.KieApp{}
	assert.True(t, SetProvisioning(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	condition := cr.Status.Conditions[0]

	assert.False(t, SetProvisioning(cr))
	assert.Equal(t, 1, len(cr.Status.Conditions))
	assert.Equal(t, condition, cr.Status.Conditions[0])
	assert.Equal(t, condition.Type, cr.Status.Phase)
}

func TestSetProvisioningAndThenDeployed(t *testing.T) {
	now := metav1.Now()
	cr := &api.KieApp{}

	assert.True(t, SetProvisioning(cr))
	assert.True(t, SetDeployed(cr))

	assert.NotEmpty(t, cr.Status.Conditions)
	condition := cr.Status.Conditions[0]
	assert.Equal(t, 2, len(cr.Status.Conditions))
	assert.Equal(t, api.ProvisioningConditionType, condition.Type)
	assert.Equal(t, corev1.ConditionTrue, condition.Status)
	assert.True(t, now.Before(&condition.LastTransitionTime))

	assert.Equal(t, api.DeployedConditionType, cr.Status.Conditions[1].Type)
	assert.Equal(t, corev1.ConditionTrue, cr.Status.Conditions[1].Status)
	assert.True(t, condition.LastTransitionTime.Before(&cr.Status.Conditions[1].LastTransitionTime))
	assert.Equal(t, api.DeployedConditionType, cr.Status.Phase)
}

func TestBuffer(t *testing.T) {
	cr := &api.KieApp{}
	for i := 0; i < maxBuffer+2; i++ {
		SetFailed(cr, api.UnknownReason, fmt.Errorf("Error %d", i))
	}
	size := len(cr.Status.Conditions)
	assert.Equal(t, maxBuffer, size)
	assert.Equal(t, "Error 31", cr.Status.Conditions[size-1].Message)
}
