package metaads

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// CredentialProvider abstracts how Meta API credentials are obtained.
type CredentialProvider interface {
	GetCredentials(ctx context.Context, accountName string) (*Credentials, error)
	ListAccounts(ctx context.Context) ([]AccountInfo, error)
}

// Credentials holds all Meta API credentials for one account.
type Credentials struct {
	AccessToken string
	AppSecret   string
	AppID       string
	AdAccountID string
	BusinessID  string
}

// AccountInfo holds non-secret account metadata for display.
type AccountInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// EnvCredentialProvider reads Meta API credentials from environment variables.
type EnvCredentialProvider struct{}

func (p *EnvCredentialProvider) GetCredentials(_ context.Context, _ string) (*Credentials, error) {
	token := os.Getenv("META_ADS_ACCESS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("META_ADS_ACCESS_TOKEN environment variable is not set")
	}
	return &Credentials{
		AccessToken: token,
		AppSecret:   os.Getenv("META_APP_SECRET"),
		AppID:       os.Getenv("META_APP_ID"),
		AdAccountID: os.Getenv("META_AD_ACCOUNT_ID"),
		BusinessID:  os.Getenv("META_BUSINESS_ID"),
	}, nil
}

func (p *EnvCredentialProvider) ListAccounts(_ context.Context) ([]AccountInfo, error) {
	if os.Getenv("META_ADS_ACCESS_TOKEN") == "" {
		return nil, nil
	}
	return []AccountInfo{{Name: "default", DisplayName: "Default"}}, nil
}

// Config manages Meta Ads accounts and their API clients.
type Config struct {
	provider CredentialProvider
	clients  map[string]*Client
	mu       sync.Mutex
}

// NewConfig creates a Config with the given credential provider.
func NewConfig(provider CredentialProvider) *Config {
	return &Config{
		provider: provider,
		clients:  make(map[string]*Client),
	}
}

// NewEnvConfig creates a Config that reads credentials from environment variables.
func NewEnvConfig() *Config {
	return NewConfig(&EnvCredentialProvider{})
}

// GetClient returns the API client for the named account, creating it on first access.
func (c *Config) GetClient(ctx context.Context, name string) (*Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if name == "" {
		name = "default"
	}

	if cl, ok := c.clients[name]; ok {
		return cl, nil
	}

	creds, err := c.provider.GetCredentials(ctx, name)
	if err != nil {
		names, _ := c.AccountNames(ctx)
		return nil, fmt.Errorf("account %q: %w (available: %s)", name, err, strings.Join(names, ", "))
	}

	client := NewClient(creds.AccessToken, creds.AppSecret, creds.AdAccountID, creds.BusinessID)
	c.clients[name] = client
	return client, nil
}

// Accounts returns non-secret metadata for all configured accounts.
func (c *Config) Accounts(ctx context.Context) []AccountInfo {
	accounts, _ := c.provider.ListAccounts(ctx)
	return accounts
}

// AccountNames returns all account names.
func (c *Config) AccountNames(ctx context.Context) ([]string, error) {
	accounts, err := c.provider.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(accounts))
	for i, a := range accounts {
		names[i] = a.Name
	}
	return names, nil
}
