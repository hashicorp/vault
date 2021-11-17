package awsutil

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
)

// getOpts iterates the inbound Options and returns a struct
func getOpts(opt ...Option) (options, error) {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o == nil {
			continue
		}
		if err := o(&opts); err != nil {
			return options{}, err
		}
	}
	return opts, nil
}

// Option - how Options are passed as arguments
type Option func(*options) error

// options = how options are represented
type options struct {
	withEnvironmentCredentials bool
	withSharedCredentials      bool
	withAwsSession             *session.Session
	withClientType             string
	withUsername               string
	withAccessKey              string
	withSecretKey              string
	withLogger                 hclog.Logger
	withStsEndpoint            string
	withIamEndpoint            string
	withMaxRetries             *int
	withRegion                 string
	withHttpClient             *http.Client
	withValidityCheckTimeout   time.Duration
}

func getDefaultOptions() options {
	return options{
		withEnvironmentCredentials: true,
		withSharedCredentials:      true,
		withClientType:             "iam",
	}
}

// WithEnvironmentCredentials allows controlling whether environment credentials
// are used
func WithEnvironmentCredentials(with bool) Option {
	return func(o *options) error {
		o.withEnvironmentCredentials = with
		return nil
	}
}

// WithSharedCredentials allows controlling whether shared credentials are used
func WithSharedCredentials(with bool) Option {
	return func(o *options) error {
		o.withSharedCredentials = with
		return nil
	}
}

// WithAwsSession allows controlling the session passed into the client
func WithAwsSession(with *session.Session) Option {
	return func(o *options) error {
		o.withAwsSession = with
		return nil
	}
}

// WithClientType allows choosing the client type to use
func WithClientType(with string) Option {
	return func(o *options) error {
		switch with {
		case "iam", "sts":
		default:
			return fmt.Errorf("unsupported client type %q", with)
		}
		o.withClientType = with
		return nil
	}
}

// WithUsername allows passing the user name to use for an operation
func WithUsername(with string) Option {
	return func(o *options) error {
		o.withUsername = with
		return nil
	}
}

// WithAccessKey allows passing an access key to use for operations
func WithAccessKey(with string) Option {
	return func(o *options) error {
		o.withAccessKey = with
		return nil
	}
}

// WithSecretKey allows passing a secret key to use for operations
func WithSecretKey(with string) Option {
	return func(o *options) error {
		o.withSecretKey = with
		return nil
	}
}

// WithStsEndpoint allows passing a custom STS endpoint
func WithStsEndpoint(with string) Option {
	return func(o *options) error {
		o.withStsEndpoint = with
		return nil
	}
}

// WithIamEndpoint allows passing a custom IAM endpoint
func WithIamEndpoint(with string) Option {
	return func(o *options) error {
		o.withIamEndpoint = with
		return nil
	}
}

// WithRegion allows passing a custom region
func WithRegion(with string) Option {
	return func(o *options) error {
		o.withRegion = with
		return nil
	}
}

// WithLogger allows passing a logger to use
func WithLogger(with hclog.Logger) Option {
	return func(o *options) error {
		o.withLogger = with
		return nil
	}
}

// WithMaxRetries allows passing custom max retries to set
func WithMaxRetries(with *int) Option {
	return func(o *options) error {
		o.withMaxRetries = with
		return nil
	}
}

// WithHttpClient allows passing a custom client to use
func WithHttpClient(with *http.Client) Option {
	return func(o *options) error {
		o.withHttpClient = with
		return nil
	}
}

// WithValidityCheckTimeout allows passing a timeout for operations that can wait
// on success.
func WithValidityCheckTimeout(with time.Duration) Option {
	return func(o *options) error {
		o.withValidityCheckTimeout = with
		return nil
	}
}
