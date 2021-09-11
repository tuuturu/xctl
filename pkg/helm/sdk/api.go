package sdk

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/yaml"
)

const defaultHelmDriver = "secrets"

func (h helmSDKClient) Install(plugin v1alpha1.Plugin) error {
	chart, err := loader.Load(plugin.Spec.HelmChart)
	if err != nil {
		return fmt.Errorf("loading Helm chart: %w", err)
	}

	actionConfig, err := h.generateActionConfig(plugin.Metadata.Namespace)
	if err != nil {
		return fmt.Errorf("generating action config: %w", err)
	}

	installAction := action.NewInstall(actionConfig)
	installAction.Namespace = plugin.Metadata.Namespace
	installAction.ReleaseName = plugin.Metadata.Name

	valuesMap, err := valuesStringToMap(plugin.Spec.Values)
	if err != nil {
		return fmt.Errorf("converting values to map: %w", err)
	}

	_, err = installAction.Run(chart, valuesMap)
	if err != nil {
		return fmt.Errorf("installing chart: %w", err)
	}

	return nil
}

func (h helmSDKClient) Delete(plugin v1alpha1.Plugin) error {
	actionConfig, err := h.generateActionConfig(plugin.Metadata.Namespace)
	if err != nil {
		return fmt.Errorf("generating action config: %w", err)
	}

	uninstallAction := action.NewUninstall(actionConfig)

	uninstallAction.KeepHistory = false
	uninstallAction.Timeout = config.DefaultHelmActionTimeout

	_, err = uninstallAction.Run(plugin.Metadata.Name)
	if err != nil {
		return fmt.Errorf("uninstalling chart: %w", err)
	}

	return nil
}

func (h helmSDKClient) Exists(plugin v1alpha1.Plugin) (bool, error) {
	_, err := h.findRelease(plugin.Metadata.Name, plugin.Metadata.Namespace)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("finding release: %w", err)
	}

	return true, nil
}

func (h helmSDKClient) generateActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	kubeConfigPath := config.GetAbsoluteKubeconfigPath()

	restClient := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfigPath,
		Namespace:  &namespace,
	}

	err := actionConfig.Init(restClient, namespace, defaultHelmDriver, h.debugPrinterf)
	if err != nil {
		return nil, fmt.Errorf("initializing action config: %w", err)
	}

	return actionConfig, nil
}

func (h helmSDKClient) findRelease(releaseName, namespace string) (*release.Release, error) {
	actionConfig, err := h.generateActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("generating action config: %w", err)
	}

	listAction := action.NewList(actionConfig)
	listAction.All = true

	listAction.SetStateMask()

	releases, err := listAction.Run()
	if err != nil {
		return nil, fmt.Errorf("running list action: %w", err)
	}

	for _, r := range releases {
		if r.Name == releaseName {
			return r, nil
		}
	}

	return nil, config.ErrNotFound
}

func (h helmSDKClient) debugPrinterf(format string, v ...interface{}) {
	if !h.debug {
		return
	}

	fmt.Fprintf(h.debugOut, format, v...)
}

func valuesStringToMap(rawValues string) (map[string]interface{}, error) {
	var values map[string]interface{}

	err := yaml.Unmarshal([]byte(rawValues), &values)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling values: %w", err)
	}

	return values, nil
}

func NewHelmClient(opts NewHelmClientOpts) helm.Client {
	return &helmSDKClient{
		debugOut: opts.DebugOut,
		debug:    opts.Debug,
	}
}
