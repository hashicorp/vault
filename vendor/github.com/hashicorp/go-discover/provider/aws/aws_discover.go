// Package aws provides node discovery for Amazon AWS.
package aws

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Amazon AWS:

    provider:          "aws"
    region:            The AWS region. Default to region of instance.
    tag_key:           The tag key to filter on
    tag_value:         The tag value to filter on
    addr_type:         "private_v4", "public_v4" or "public_v6". Defaults to "private_v4".
    access_key_id:     The AWS access key to use
    secret_access_key: The AWS secret access key to use

    The only required IAM permission is 'ec2:DescribeInstances'. If the Consul agent is
    running on AWS instance it is recommended you use an IAM role, otherwise it is
    recommended you make a dedicated IAM user and access key used only for auto-joining.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "aws" {
		return nil, fmt.Errorf("discover-aws: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	addrType := args["addr_type"]
	accessKey := args["access_key_id"]
	secretKey := args["secret_access_key"]

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-aws: Address type %s is not supported. Valid values are {private_v4,public_v4,public_v6}. Falling back to 'private_v4'", addrType)
		addrType = "private_v4"
	}

	if addrType == "" {
		l.Printf("[DEBUG] discover-aws: Address type not provided. Using 'private_v4'")
		addrType = "private_v4"
	}

	l.Printf("[DEBUG] discover-aws: Using region=%s tag_key=%s tag_value=%s addr_type=%s", region, tagKey, tagValue, addrType)
	if accessKey == "" && secretKey == "" {
		l.Printf("[DEBUG] discover-aws: No static credentials")
		l.Printf("[DEBUG] discover-aws: Using environment variables, shared credentials or instance role")
	} else {
		l.Printf("[DEBUG] discover-aws: Static credentials provided")
	}

	if region == "" {
		l.Printf("[INFO] discover-aws: Region not provided. Looking up region in metadata...")
		ec2meta := ec2metadata.New(session.New())
		identity, err := ec2meta.GetInstanceIdentityDocument()
		if err != nil {
			return nil, fmt.Errorf("discover-aws: GetInstanceIdentityDocument failed: %s", err)
		}
		region = identity.Region
	}
	l.Printf("[INFO] discover-aws: Region is %s", region)

	l.Printf("[DEBUG] discover-aws: Creating session...")
	svc := ec2.New(session.New(), &aws.Config{
		Region: &region,
		Credentials: credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     accessKey,
						SecretAccessKey: secretKey,
					},
				},
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{},
				defaults.RemoteCredProvider(*(defaults.Config()), defaults.Handlers()),
			}),
	})

	l.Printf("[INFO] discover-aws: Filter instances with %s=%s", tagKey, tagValue)
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:" + tagKey),
				Values: []*string{aws.String(tagValue)},
			},
			&ec2.Filter{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("discover-aws: DescribeInstancesInput failed: %s", err)
	}

	l.Printf("[DEBUG] discover-aws: Found %d reservations", len(resp.Reservations))
	var addrs []string
	for _, r := range resp.Reservations {
		l.Printf("[DEBUG] discover-aws: Reservation %s has %d instances", *r.ReservationId, len(r.Instances))
		for _, inst := range r.Instances {
			id := *inst.InstanceId
			l.Printf("[DEBUG] discover-aws: Found instance %s", id)

			switch addrType {
			case "public_v6":
				l.Printf("[DEBUG] discover-aws: Instance %s has %d network interfaces", id, len(inst.NetworkInterfaces))

				for _, networkinterface := range inst.NetworkInterfaces {
					l.Printf("[DEBUG] discover-aws: Checking NetworInterfaceId %s on Instance %s", *networkinterface.NetworkInterfaceId, id)
					// Check if instance got any ipv6
					if networkinterface.Ipv6Addresses == nil {
						l.Printf("[DEBUG] discover-aws: Instance %s has no IPv6 on NetworkInterfaceId %s", id, *networkinterface.NetworkInterfaceId)
						continue
					}
					for _, ipv6address := range networkinterface.Ipv6Addresses {
						l.Printf("[INFO] discover-aws: Instance %s has IPv6 %s on NetworkInterfaceId %s", id, *ipv6address.Ipv6Address, *networkinterface.NetworkInterfaceId)
						addrs = append(addrs, *ipv6address.Ipv6Address)
					}
				}

			case "public_v4":
				if inst.PublicIpAddress == nil {
					l.Printf("[DEBUG] discover-aws: Instance %s has no public IPv4", id)
					continue
				}

				l.Printf("[INFO] discover-aws: Instance %s has public ip %s", id, *inst.PublicIpAddress)
				addrs = append(addrs, *inst.PublicIpAddress)

			default:
				// EC2-Classic don't have the PrivateIpAddress field
				if inst.PrivateIpAddress == nil {
					l.Printf("[DEBUG] discover-aws: Instance %s has no private ip", id)
					continue
				}

				l.Printf("[INFO] discover-aws: Instance %s has private ip %s", id, *inst.PrivateIpAddress)
				addrs = append(addrs, *inst.PrivateIpAddress)
			}
		}
	}

	l.Printf("[DEBUG] discover-aws: Found ip addresses: %v", addrs)
	return addrs, nil
}
