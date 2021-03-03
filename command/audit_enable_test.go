package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/Shopify/sarama"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/mitchellh/cli"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func testAuditEnableCommand(tb testing.TB) (*cli.MockUi, *AuditEnableCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuditEnableCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuditEnableCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"empty",
			nil,
			"Missing TYPE!",
			1,
		},
		{
			"not_a_valid_type",
			[]string{"nope_definitely_not_a_valid_type_like_ever"},
			"",
			2,
		},
		{
			"enable",
			[]string{"file", "file_path=discard"},
			"Success! Enabled the file audit device at: file/",
			0,
		},
		{
			"enable_path",
			[]string{
				"-path", "audit_path",
				"file",
				"file_path=discard",
			},
			"Success! Enabled the file audit device at: audit_path/",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testAuditEnableCommand(t)
			cmd.client = client

			code := cmd.Run(tc.args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q to contain %q", combined, tc.out)
			}
		})
	}

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testAuditEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-path", "audit_enable_integration/",
			"-description", "The best kind of test",
			"file",
			"file_path=discard",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Enabled the file audit device at: audit_enable_integration/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		audits, err := client.Sys().ListAudit()
		if err != nil {
			t.Fatal(err)
		}

		auditInfo, ok := audits["audit_enable_integration/"]
		if !ok {
			t.Fatalf("expected audit to exist")
		}
		if exp := "file"; auditInfo.Type != exp {
			t.Errorf("expected %q to be %q", auditInfo.Type, exp)
		}
		if exp := "The best kind of test"; auditInfo.Description != exp {
			t.Errorf("expected %q to be %q", auditInfo.Description, exp)
		}

		filePath, ok := auditInfo.Options["file_path"]
		if !ok || filePath != "discard" {
			t.Errorf("missing some options: %#v", auditInfo)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuditEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error enabling audit device: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuditEnableCommand(t)
		assertNoTabs(t, cmd)
	})

	t.Run("mount_all", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerAllBackends(t)
		defer closer()
		cleanup, kafkaAddress := prepareKafkaTestContainer(t)
		defer cleanup()

		files, err := ioutil.ReadDir("../builtin/audit")
		if err != nil {
			t.Fatal(err)
		}

		var backends []string
		for _, f := range files {
			if f.IsDir() {
				backends = append(backends, f.Name())
			}
		}

		for _, b := range backends {
			ui, cmd := testAuditEnableCommand(t)
			cmd.client = client

			args := []string{
				b,
			}
			switch b {
			case "file":
				args = append(args, "file_path=discard")
			case "socket":
				args = append(args, "address=127.0.0.1:8888",
					"skip_test=true")
			case "kafka":
				args = append(args, fmt.Sprintf("address=%s", kafkaAddress), "topic=vault", "tls_disabled=true")
			case "syslog":
				if _, exists := os.LookupEnv("WSLENV"); exists {
					t.Log("skipping syslog test on WSL")
					continue
				}
				if os.Getenv("CIRCLECI") == "true" {
					// TODO install syslog in docker image we run our tests in
					t.Log("skipping syslog test on CircleCI")
					continue
				}
			}
			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Errorf("type %s, expected %d to be %d - %s", b, code, exp, ui.OutputWriter.String()+ui.ErrorWriter.String())
			}
		}
	})
}

func prepareKafkaTestContainer(t *testing.T) (cleanup func(), retURL string) {
	if os.Getenv("KAFKA_SERVER") != "" {
		return func() {}, os.Getenv("KAFKA_SERVER")
	}

	randName, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("Unable to create UUID: %s", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}
	client := pool.Client
	net, err := client.CreateNetwork(docker.CreateNetworkOptions{
		Name: fmt.Sprintf("vault-kafka-%s", randName),
	})
	if err != nil {
		t.Fatalf("Failed to create docker network: %s", err)
	}

	zookeeperResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       fmt.Sprintf("zookeeper-%s", randName),
		Repository: "confluentinc/cp-zookeeper",
		Tag:        "5.5.1",
		NetworkID:  net.ID,
		Env: []string{
			"ZOOKEEPER_CLIENT_PORT=2181",
		},
	})
	if err != nil {
		t.Fatalf("Could not start local zookeeper docker container: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       fmt.Sprintf("kafka-%s", randName),
		Repository: "confluentinc/cp-kafka",
		Tag:        "5.5.1",
		NetworkID:  net.ID,
		Env: []string{
			"KAFKA_BROKER_ID=1",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
			"KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092",
			fmt.Sprintf("KAFKA_ZOOKEEPER_CONNECT=zookeeper-%s:2181", randName),
		},
	})

	if err != nil {
		t.Fatalf("Could not start local kafka docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(zookeeperResource)
		if err != nil {
			t.Fatalf("Failed to cleanup local containers: %s", err)
		}

		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local containers: %s", err)
		}

		client.RemoveNetwork(net.ID)
		if err != nil {
			t.Fatalf("Failed to cleanup network: %s", err)
		}
	}

	kafkaServerAddress := fmt.Sprintf("localhost:%s", resource.GetPort("9092/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		c, err := sarama.NewClient([]string{kafkaServerAddress}, nil)
		if err != nil {
			return err
		}

		return c.Close()
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to kakfa docker container: %s", err)
	}

	return cleanup, kafkaServerAddress
}
