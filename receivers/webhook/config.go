package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/grafana/alerting/receivers"
	"github.com/grafana/alerting/templates"
)

type Config struct {
	URL        string
	HTTPMethod string
	MaxAlerts  int
	// Authorization Header.
	AuthorizationScheme      string
	AuthorizationCredentials string
	// HTTP Basic Authentication.
	User     string
	Password string

	Title   string
	Message string
}

func ValidateConfig(factoryConfig receivers.FactoryConfig) (Config, error) {
	settings := Config{}
	rawSettings := struct {
		URL                      string           `json:"url,omitempty" yaml:"url,omitempty"`
		HTTPMethod               string           `json:"httpMethod,omitempty" yaml:"httpMethod,omitempty"`
		MaxAlerts                json.Number      `json:"maxAlerts,omitempty" yaml:"maxAlerts,omitempty"`
		AuthorizationScheme      string           `json:"authorization_scheme,omitempty" yaml:"authorization_scheme,omitempty"`
		AuthorizationCredentials receivers.Secret `json:"authorization_credentials,omitempty" yaml:"authorization_credentials,omitempty"`
		User                     receivers.Secret `json:"username,omitempty" yaml:"username,omitempty"`
		Password                 receivers.Secret `json:"password,omitempty" yaml:"password,omitempty"`
		Title                    string           `json:"title,omitempty" yaml:"title,omitempty"`
		Message                  string           `json:"message,omitempty" yaml:"message,omitempty"`
	}{}

	err := factoryConfig.Marshaller.Unmarshal(factoryConfig.Config.Settings, &rawSettings)
	if err != nil {
		return settings, fmt.Errorf("failed to unmarshal settings: %w", err)
	}
	if rawSettings.URL == "" {
		return settings, errors.New("required field 'url' is not specified")
	}
	settings.URL = rawSettings.URL

	if rawSettings.HTTPMethod == "" {
		rawSettings.HTTPMethod = http.MethodPost
	}
	settings.HTTPMethod = rawSettings.HTTPMethod

	if rawSettings.MaxAlerts != "" {
		settings.MaxAlerts, _ = strconv.Atoi(rawSettings.MaxAlerts.String())
	}

	if settings.AuthorizationCredentials != "" && settings.AuthorizationScheme == "" {
		settings.AuthorizationScheme = "Bearer"
	}
	if settings.User != "" && settings.Password != "" && settings.AuthorizationScheme != "" && settings.AuthorizationCredentials != "" {
		return settings, errors.New("both HTTP Basic Authentication and Authorization Header are set, only 1 is permitted")
	}
	settings.Title = rawSettings.Title
	if settings.Title == "" {
		settings.Title = templates.DefaultMessageTitleEmbed
	}
	settings.Message = rawSettings.Message
	if settings.Message == "" {
		settings.Message = templates.DefaultMessageEmbed
	}
	return settings, err
}
