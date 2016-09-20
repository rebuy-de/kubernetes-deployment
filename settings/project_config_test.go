package settings

import (
	"testing"
	"github.com/rebuy-de/kubernetes-deployment/util"
	"time"
)

func TestMergeConfig_Kubeconfig(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	*pc_local.Settings.Kubeconfig = "Kubeconfig_DEFAULT"
	pc_default.MergeConfig(pc_local)
	util.AssertStringEquals(t, "Kubeconfig_DEFAULT", *pc_default.Settings.Kubeconfig, "Kubeconfig")
}

func TestMergeConfig_Kubeconfig_nill(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local.Settings.Kubeconfig = nil
	pc_default.MergeConfig(pc_local)
	util.AssertStringEquals(t, "test-fixtures/kubeconfig.yml", *pc_default.Settings.Kubeconfig, "Kubeconfig")
}

func TestMergeConfig_Output(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	*pc_local.Settings.Output = "Output_DEFAULT"
	pc_default.MergeConfig(pc_local)
	util.AssertStringEquals(t, "Output_DEFAULT", *pc_default.Settings.Output, "Output")
}

func TestMergeConfig_Output_nill(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local.Settings.Output = nil
	pc_default.MergeConfig(pc_local)
	util.AssertStringEquals(t, "target/test-output", *pc_default.Settings.Output, "Output")
}

func TestMergeConfig_Sleep(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	*pc_local.Settings.Sleep = 1000 * time.Second
	pc_default.MergeConfig(pc_local)
	util.AssertDurationEquals(t, 1000 * time.Second, *pc_default.Settings.Sleep, "Sleep")
}

func TestMergeConfig_Sleep_nill(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local.Settings.Sleep = nil
	pc_default.MergeConfig(pc_local)
	util.AssertDurationEquals(t, 1 * time.Second, *pc_default.Settings.Sleep, "Sleep")
}

func TestMergeConfig_RetrySleep(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	*pc_local.Settings.RetrySleep = 2000 * time.Second
	pc_default.MergeConfig(pc_local)
	util.AssertDurationEquals(t, 2000 * time.Second, *pc_default.Settings.RetrySleep, "RetrySleep")
}

func TestMergeConfig_RetrySleep_nill(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local.Settings.RetrySleep = nil
	pc_default.MergeConfig(pc_local)
	util.AssertDurationEquals(t, 250 * time.Millisecond, *pc_default.Settings.RetrySleep, "RetrySleep")
}

func TestMergeConfig_RetryCount(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	*pc_local.Settings.RetryCount = 10001
	pc_default.MergeConfig(pc_local)
	util.AssertIntEquals(t, 10001, *pc_default.Settings.RetryCount, "RetryCount")
}

func TestMergeConfig_RetryCount_nill(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_local.Settings.RetryCount = nil
	pc_default.MergeConfig(pc_local)
	util.AssertIntEquals(t, 3, *pc_default.Settings.RetryCount, "RetryCount")
}

func TestMergeConfig_templateValues(t *testing.T) {
	pc_default, err := ReadProjectConfigFrom("../config/services.yaml")
	util.AssertNoError(t, err)
	pc_local, err := ReadProjectConfigFrom("../config/services_test.yaml")
	util.AssertNoError(t, err)
	pc_default.MergeConfig(pc_local)
	util.AssertStringEquals(t, "unit-test.rebuy.de", pc_default.Settings.TemplateValuesMap["clusterDomain"], "clusterDomainValue")
}


