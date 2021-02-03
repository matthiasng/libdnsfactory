// Code generated by go run ./generate/main.go; DO NOT EDIT.

package libdnsfactory

import (
	"fmt"
    "github.com/libdns/alidns"
    "github.com/libdns/azure"
    "github.com/libdns/cloudflare"
    "github.com/libdns/digitalocean"
    "github.com/libdns/dnspod"
    "github.com/libdns/gandi"
    "github.com/libdns/hetzner"
    "github.com/libdns/openstack-designate"
    "github.com/libdns/route53"
    "github.com/libdns/transip"
    "github.com/libdns/vultr"
)

// NewProvider creates a new Provider. See https://github.com/matthiasng/libdnsfactory/blob/master/docs.md for a complete list of supported providers.
func NewProvider(name string, config map[string]string) (Provider, error) {
	switch name {
        case "alidns":
            var err error
            p := &alidns.Provider{}
            p.AccKeyID, err = getValueString("AccKeyID", true, config)
            if err != nil {
                return nil, fmt.Errorf("alidns [AccKeyID]: %w", err)
            }
            p.AccKeySecret, err = getValueString("AccKeySecret", true, config)
            if err != nil {
                return nil, fmt.Errorf("alidns [AccKeySecret]: %w", err)
            }
            p.RegionID, err = getValueString("RegionID", false, config)
            if err != nil {
                return nil, fmt.Errorf("alidns [RegionID]: %w", err)
            }

            return p, nil
        case "azure":
            var err error
            p := &azure.Provider{}
            p.TenantId, err = getValueString("TenantId", false, config)
            if err != nil {
                return nil, fmt.Errorf("azure [TenantId]: %w", err)
            }
            p.ClientId, err = getValueString("ClientId", false, config)
            if err != nil {
                return nil, fmt.Errorf("azure [ClientId]: %w", err)
            }
            p.ClientSecret, err = getValueString("ClientSecret", false, config)
            if err != nil {
                return nil, fmt.Errorf("azure [ClientSecret]: %w", err)
            }
            p.SubscriptionId, err = getValueString("SubscriptionId", false, config)
            if err != nil {
                return nil, fmt.Errorf("azure [SubscriptionId]: %w", err)
            }
            p.ResourceGroupName, err = getValueString("ResourceGroupName", false, config)
            if err != nil {
                return nil, fmt.Errorf("azure [ResourceGroupName]: %w", err)
            }

            return p, nil
        case "cloudflare":
            var err error
            p := &cloudflare.Provider{}
            p.APIToken, err = getValueString("APIToken", false, config)
            if err != nil {
                return nil, fmt.Errorf("cloudflare [APIToken]: %w", err)
            }

            return p, nil
        case "digitalocean":
            var err error
            p := &digitalocean.Provider{}
            p.APIToken, err = getValueString("APIToken", true, config)
            if err != nil {
                return nil, fmt.Errorf("digitalocean [APIToken]: %w", err)
            }

            return p, nil
        case "dnspod":
            var err error
            p := &dnspod.Provider{}
            p.APIToken, err = getValueString("APIToken", true, config)
            if err != nil {
                return nil, fmt.Errorf("dnspod [APIToken]: %w", err)
            }

            return p, nil
        case "gandi":
            var err error
            p := &gandi.Provider{}
            p.APIToken, err = getValueString("APIToken", false, config)
            if err != nil {
                return nil, fmt.Errorf("gandi [APIToken]: %w", err)
            }

            return p, nil
        case "hetzner":
            var err error
            p := &hetzner.Provider{}
            p.AuthAPIToken, err = getValueString("AuthAPIToken", true, config)
            if err != nil {
                return nil, fmt.Errorf("hetzner [AuthAPIToken]: %w", err)
            }

            return p, nil
        case "openstack-designate":
            var err error
            p := &openstack-designate.Provider{}
            p.AuthOpenStack, err = getValueAuthOpenStack("AuthOpenStack", true, config)
            if err != nil {
                return nil, fmt.Errorf("openstack-designate [AuthOpenStack]: %w", err)
            }

            return p, nil
        case "route53":
            var err error
            p := &route53.Provider{}
            p.MaxRetries, err = getValueInt("MaxRetries", false, config)
            if err != nil {
                return nil, fmt.Errorf("route53 [MaxRetries]: %w", err)
            }
            p.AWSProfile, err = getValueString("AWSProfile", false, config)
            if err != nil {
                return nil, fmt.Errorf("route53 [AWSProfile]: %w", err)
            }
            p.AccessKeyId, err = getValueString("AccessKeyId", false, config)
            if err != nil {
                return nil, fmt.Errorf("route53 [AccessKeyId]: %w", err)
            }
            p.SecretAccessKey, err = getValueString("SecretAccessKey", false, config)
            if err != nil {
                return nil, fmt.Errorf("route53 [SecretAccessKey]: %w", err)
            }

            return p, nil
        case "transip":
            var err error
            p := &transip.Provider{}
            p.AccountName, err = getValueString("AccountName", true, config)
            if err != nil {
                return nil, fmt.Errorf("transip [AccountName]: %w", err)
            }
            p.PrivateKeyPath, err = getValueString("PrivateKeyPath", true, config)
            if err != nil {
                return nil, fmt.Errorf("transip [PrivateKeyPath]: %w", err)
            }

            return p, nil
        case "vultr":
            var err error
            p := &vultr.Provider{}
            p.APIToken, err = getValueString("APIToken", true, config)
            if err != nil {
                return nil, fmt.Errorf("vultr [APIToken]: %w", err)
            }

            return p, nil
	default:
		return nil, fmt.Errorf("Unknown provider: %s", name)
	}
}
