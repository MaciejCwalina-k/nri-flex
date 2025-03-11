package runtime

import (
	"os"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/load"
)

type AzureFunction struct {
	configDir string
}

func (azf *AzureFunction) isAvailable() bool {
	if os.Getenv("FUNCTIONS_WORKER_RUNTIME") == "" {
		return false
	}

	load.ServerlessName = "test"          //TODO: probabbly use the WEBSITE_CONTENTSHARE env var
	load.ServerlessExecutionEnv = "azure" //TODO: Get the runtime env var and add Azure Function to it
	azf.SetConfigDir("./flexConfigs")     //TODO: Maybe make this an env var ? dont think it matters
	err := azf.init()
	if err != nil {
		log.Error("AzureFunction.IsAvailable: azf.init() threw some errors ", err.Error())
		return false
	}

	return true
}

func (azf *AzureFunction) loadConfigs(configs *[]load.Config) error {
	load.Logrus.Info("AzureFunction.loadConfigs: running as Azure Function")
	errors := addConfigsFromPath(azf.configDir, configs)
	if len(errors) > 0 {
		log.Error("AzureFunction.loadConfigs: failed to read some configuration files, please review them")
	}

	isSyncGitConfigured, err := config.SyncGitConfigs("/tmp/")
	if err != nil {
		log.WithError(err).Warn("AzureFunction.loadConfigs: failed to sync git configs")
	} else if isSyncGitConfigured {
		errors = addConfigsFromPath("/tmp/", configs)
		if len(errors) > 0 {
			log.Error("AzureFunction.loadConfigs: failed to load git sync configuration files, ignoring and continuing")
		}
	}

	return nil
}

func (azf *AzureFunction) SetConfigDir(config string) {
	azf.configDir = config
}

func (azf *AzureFunction) init() error {
	return nil
}
