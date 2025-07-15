package datastore

import (
	"os"

	portainer "github.com/portainer/portainer/api"
)

// Init creates the default data set.
func (store *Store) Init() error {
	err := store.checkOrCreateDefaultSettings()
	if err != nil {
		return err
	}

	err = store.checkOrCreateDefaultSSLSettings()
	if err != nil {
		return err
	}

	err = store.checkOrCreateDefaultRoles()
	if err != nil {
		return err
	}

	return store.checkOrCreateDefaultData()
}

func (store *Store) checkOrCreateDefaultRoles() error {
	roles, err := store.RoleService.ReadAll()
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		defaultRoles := []portainer.Role{
			{
				ID:          1,
				Name:        "Environment administrator",
				Description: "Full control of all resources in an environment",
				// Authorizations: заполните по необходимости
				Priority: 1,
			},
			{
				ID:          2,
				Name:        "Helpdesk",
				Description: "Read-only access of all resources in an environment",
				Priority:    2,
			},
			{
				ID:          3,
				Name:        "Standard user",
				Description: "Full control of assigned resources in an environment",
				Priority:    3,
			},
			{
				ID:          4,
				Name:        "Read-only user",
				Description: "Read-only access of assigned resources in an environment",
				Priority:    4,
			},
			{
				ID:          5,
				Name:        "Operator",
				Description: "Operational Control of all existing resources in an environment",
				Priority:    5,
			},
		}
		for _, role := range defaultRoles {
			err := store.RoleService.Create(&role)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (store *Store) checkOrCreateDefaultSettings() error {
	isDDExtention := false
	if _, ok := os.LookupEnv("DOCKER_EXTENSION"); ok {
		isDDExtention = true
	}

	// TODO: these need to also be applied when importing
	settings, err := store.SettingsService.Settings()
	if store.IsErrObjectNotFound(err) {
		defaultSettings := &portainer.Settings{
			EnableTelemetry:      false,
			AuthenticationMethod: portainer.AuthenticationInternal,
			BlackListedLabels:    make([]portainer.Pair, 0),
			InternalAuthSettings: portainer.InternalAuthSettings{
				RequiredPasswordLength: 12,
			},
			LDAPSettings: portainer.LDAPSettings{
				AnonymousMode:   true,
				AutoCreateUsers: true,
				TLSConfig:       portainer.TLSConfiguration{},
				SearchSettings: []portainer.LDAPSearchSettings{
					{},
				},
				GroupSearchSettings: []portainer.LDAPGroupSearchSettings{
					{},
				},
			},
			OAuthSettings: portainer.OAuthSettings{
				SSO: true,
			},
			SnapshotInterval:         portainer.DefaultSnapshotInterval,
			EdgeAgentCheckinInterval: portainer.DefaultEdgeAgentCheckinIntervalInSeconds,
			TemplatesURL:             "",
			HelmRepositoryURL:        portainer.DefaultHelmRepositoryURL,
			UserSessionTimeout:       portainer.DefaultUserSessionTimeout,
			KubeconfigExpiry:         portainer.DefaultKubeconfigExpiry,
			KubectlShellImage:        *store.flags.KubectlShellImage,

			IsDockerDesktopExtension: isDDExtention,
		}

		return store.SettingsService.UpdateSettings(defaultSettings)
	}
	if err != nil {
		return err
	}

	if settings.UserSessionTimeout == "" {
		settings.UserSessionTimeout = portainer.DefaultUserSessionTimeout
		return store.Settings().UpdateSettings(settings)
	}

	return nil
}

func (store *Store) checkOrCreateDefaultSSLSettings() error {
	_, err := store.SSLSettings().Settings()
	if store.IsErrObjectNotFound(err) {
		defaultSSLSettings := &portainer.SSLSettings{
			HTTPEnabled: true,
		}

		return store.SSLSettings().UpdateSettings(defaultSSLSettings)
	}

	return err
}

func (store *Store) checkOrCreateDefaultData() error {
	groups, err := store.EndpointGroupService.ReadAll()
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		unassignedGroup := &portainer.EndpointGroup{
			Name:               "Unassigned",
			Description:        "Unassigned environments",
			Labels:             []portainer.Pair{},
			UserAccessPolicies: portainer.UserAccessPolicies{},
			TeamAccessPolicies: portainer.TeamAccessPolicies{},
			TagIDs:             []portainer.TagID{},
		}

		err = store.EndpointGroupService.Create(unassignedGroup)
		if err != nil {
			return err
		}
	}

	return nil
}
