// Copyright 2020 Topicus Security BV
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	keyhub "github.com/topicuskeyhub/go-keyhub"
	keyhubmodel "github.com/topicuskeyhub/go-keyhub/model"
	"github.com/topicuskeyhub/keyhub-vault-operator/controllers/metrics"
	"github.com/topicuskeyhub/keyhub-vault-operator/controllers/settings"
	"gopkg.in/yaml.v2"
)

type comment struct {
	Policies []policy `yaml:"policies"`
}

type PolicyLoader interface {
	Load() (*[]Policy, error)
}

type policyLoader struct {
	log             logr.Logger
	settingsManager settings.SettingsManager
	client          *keyhub.Client
}

func NewPolicyLoader(log logr.Logger, settingsMgr settings.SettingsManager) PolicyLoader {
	return &policyLoader{
		log:             log,
		settingsManager: settingsMgr,
	}
}

func (pl *policyLoader) Load() (*[]Policy, error) {
	if pl.client == nil {
		err := pl.init()
		if err != nil {
			return nil, err
		}
	}

	metrics.KeyHubApiRequests.WithLabelValues("group", "list").Inc()
	groups, err := pl.client.Groups.List()
	if err != nil {
		return nil, err
	}

	var policies []Policy
	policies = make([]Policy, 0)
	for _, group := range groups {
		err := pl.loadPolicies(&policies, group)
		if err != nil {
			return nil, err
		}
	}
	return &policies, nil
}

func (pl *policyLoader) loadPolicies(policies *[]Policy, group keyhubmodel.Group) error {
	pl.log.Info("loading group policies", "uuid", group.UUID, "name", group.Name)
	metrics.KeyHubApiRequests.WithLabelValues("vault", "list").Inc()
	records, err := pl.client.Vaults.GetRecords(&group)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Color == keyhubmodel.VaultRecordColorRed {
			continue
		}

		recordUUID, err := uuid.Parse(record.UUID)
		if err != nil {
			return err
		}
		metrics.KeyHubApiRequests.WithLabelValues("vault", "get").Inc()
		rec, err := pl.client.Vaults.GetByUUID(&group, recordUUID, &keyhubmodel.VaultRecordAdditionalQueryParams{Audit: true, Secret: true})
		if err != nil {
			return err
		}

		if rec.Username == "" || rec.Password() == nil || rec.Comment() == nil || !strings.HasPrefix(*rec.Comment(), "policies:") {
			pl.log.Info("record missing username, password or policy comment", "uuid", rec.UUID)
			continue
		}

		comment := comment{}
		err = yaml.Unmarshal([]byte(*rec.Comment()), &comment)
		if err != nil {
			pl.log.Info("unmarshal error", "uuid", rec.UUID, "err", err)
			continue
		}

		for _, policy := range comment.Policies {
			*policies = append(*policies, Policy{policy: policy, Credentials: ClientCredentials{ClientID: rec.Username, ClientSecret: *rec.Password()}})
		}
	}

	pl.log.Info("group policies loaded", "count", len(*policies))

	return nil
}

func (pl *policyLoader) init() error {
	settings, err := pl.settingsManager.GetSettings()
	if err != nil {
		return err
	}

	pl.log.Info("creating KeyHub client", "URI", settings.URI, "ClientID", settings.ClientID)
	client, err := keyhub.NewClient(http.DefaultClient, settings.URI, settings.ClientID, settings.ClientSecret)
	if err != nil {
		return fmt.Errorf("Failed to create KeyHub client: %w", err)
	}

	keyhubVersionInfo, err := client.Version.Get()
	if err != nil {
		return fmt.Errorf("Failed to connect to KeyHub: %w", err)
	}
	pl.log.Info(fmt.Sprintf("KeyHub Version: %v", keyhubVersionInfo.KeyhubVersion))

	pl.client = client

	return nil
}
