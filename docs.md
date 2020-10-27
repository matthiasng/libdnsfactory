libdns providers
=======================

# Index
1. [cloudflare](#cloudflare)
1. [digitalocean](#digitalocean)
1. [dnspod](#dnspod)
1. [gandi](#gandi)
1. [hetzner](#hetzner)
1. [route53](#route53)

# Providers

## cloudflare

This package supports API **token** authentication.

You will need to create a token with the following permissions:

- Zone / Zone / Read
- Zone / DNS / Edit

The first permission is needed to get the zone ID, and the second permission is obviously necessary to edit the DNS records. If you're only using the `GetRecords()` method, you can change the second permission to Read to guarantee no changes will be made.

To clarify, do NOT use API keys, which are globally-scoped:

![Don't use API keys](https://user-images.githubusercontent.com/1128849/81196485-556aca00-8f7c-11ea-9e13-c6a8a966f689.png)

DO use scoped API tokens:

![Don't use API keys](https://user-images.githubusercontent.com/1128849/81196503-5c91d800-8f7c-11ea-93cc-ad7d73420fab.png)

**Variables:**
| Name | Description | Type | Required |
|------|-------------|------|----------|
| APIToken | API token is used for authentication. Make sure to use a<br>scoped API **token**, NOT a global API **key**. It will<br>need two permissions: Zone-Zone-Read and Zone-DNS-Edit,<br>unless you are only using `GetRecords()`, in which case<br>the second can be changed to Read. | string | false |

**Example:**
```go
provider, err := libdnsfactory.NewProvider("cloudflare", map[string]string{
    "APIToken": "...",
})
```

**Repository**: [https://github.com/libdns/cloudflare](https://github.com/libdns/cloudflare)

## digitalocean

To authenticate you need to supply a DigitalOcean API token.

**Variables:**
| Name | Description | Type | Required |
|------|-------------|------|----------|
| APIToken | APIToken is the DigitalOcean API token - see https://www.digitalocean.com/docs/apis-clis/api/create-personal-access-token/ | string | true |

**Example:**
```go
provider, err := libdnsfactory.NewProvider("digitalocean", map[string]string{
    "APIToken": "...",
})
```

**Repository**: [https://github.com/libdns/digitalocean](https://github.com/libdns/digitalocean)

## dnspod

To authenticate you need to supply a [DNSPOD API token](https://support.dnspod.cn/Kb/showarticle/tsid/227/).

**Variables:**
| Name | Description | Type | Required |
|------|-------------|------|----------|
| APIToken | APIToken is the DNSPOD API token - see https://www.dnspod.cn/docs/info.html#common-parameters | string | true |

**Example:**
```go
provider, err := libdnsfactory.NewProvider("dnspod", map[string]string{
    "APIToken": "...",
})
```

**Repository**: [https://github.com/libdns/dnspod](https://github.com/libdns/dnspod)

## gandi

This package supports **API Key authentication** but does not yet support **Sharing ID authentication**. Refer to the [LiveDNS documentation](https://doc.livedns.gandi.net/) for more information.

Start by [retrieving your API key](https://account.gandi.net/) from the _Security_ section in Gandi account admin panel to be able to make authenticated requests to the API.

**Variables:**
| Name | Description | Type | Required |
|------|-------------|------|----------|
| APIToken | - | string | false |

**Example:**
```go
provider, err := libdnsfactory.NewProvider("gandi", map[string]string{
    "APIToken": "...",
})
```

**Repository**: [https://github.com/libdns/gandi](https://github.com/libdns/gandi)

## hetzner

To authenticate you need to supply a Hetzner [Auth-API-Token](https://dns.hetzner.com/api-docs#section/Authentication/Auth-API-Token).

**Variables:**
| Name | Description | Type | Required |
|------|-------------|------|----------|
| AuthAPIToken | AuthAPIToken is the Hetzner Auth API token - see https://dns.hetzner.com/api-docs#section/Authentication/Auth-API-Token | string | true |

**Example:**
```go
provider, err := libdnsfactory.NewProvider("hetzner", map[string]string{
    "AuthAPIToken": "...",
})
```

**Repository**: [https://github.com/libdns/hetzner](https://github.com/libdns/hetzner)

## route53

This package supports all the credential configuration methods described in the [AWS Developer Guide](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html), such as `Environment Variables`, `EC2 Instance Profile` and the `AWS Credentials file` located in `.aws/credentials`

The following IAM policy is a minimal working example to give `libdns` permissions to manage DNS records:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:ListResourceRecordSets",
                "route53:GetChange",
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/ZABCD1EFGHIL",
                "arn:aws:route53:::change/*"
            ]
        },
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZonesByName",
                "route53:ListHostedZones"
            ],
            "Resource": "*"
        }
    ]
}
```

**Variables:**
| Name | Description | Type | Required |
|------|-------------|------|----------|
| MaxRetries | - | int | false |

**Example:**
```go
provider, err := libdnsfactory.NewProvider("route53", map[string]string{
    "MaxRetries": "...",
})
```

**Repository**: [https://github.com/libdns/route53](https://github.com/libdns/route53)
