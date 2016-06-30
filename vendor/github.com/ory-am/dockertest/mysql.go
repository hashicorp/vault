package dockertest

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

const (
	defaultMySQLDBName = "mysql"
)

// SetupMySQLContainer sets up a real MySQL instance for testing purposes,
// using a Docker container. It returns the container ID and its IP address,
// or makes the test fail on error.
func SetupMySQLContainer() (c ContainerID, ip string, port int, err error) {
	port = RandomPort()
	forward := fmt.Sprintf("%d:%d", port, 3306)
	if BindDockerToLocalhost != "" {
		forward = "127.0.0.1:" + forward
	}
	c, ip, err = SetupContainer(MySQLImageName, port, 10*time.Second, func() (string, error) {
		return run("--name", GenerateContainerID(), "-d", "-p", forward, "-e", fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", MySQLPassword), MySQLImageName)
	})
	return
}

// ConnectToMySQL starts a MySQL image and passes the database url to the connector callback function.
// The url will match the username:password@tcp(ip:port) pattern (e.g. `root:root@tcp(123.123.123.123:3131)`)
func ConnectToMySQL(tries int, delay time.Duration, connector func(url string) bool) (c ContainerID, err error) {
	c, ip, port, err := SetupMySQLContainer()
	if err != nil {
		return c, fmt.Errorf("Could not set up MySQL container: %v", err)
	}

	for try := 0; try <= tries; try++ {
		time.Sleep(delay)
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", MySQLUsername, MySQLPassword, ip, port)
		if connector(url) {
			return c, nil
		}
		log.Printf("Try %d failed. Retrying.", try)
	}
	return c, errors.New("Could not set up MySQL container.")
}

// SetUpMySQLDatabase connects mysql container with given $connectURL and also creates a new database named $databaseName
// A modified url used to connect the created database will be returned
func SetUpMySQLDatabase(databaseName, connectURL string) (url string, err error) {
	if databaseName == defaultMySQLDBName {
		return connectURL, nil
	}

	db, err := sql.Open("mysql", connectURL)
	if err != nil {
		return "", err
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", databaseName))
	if err != nil {
		return "", err
	}

	// parse dsn
	config, err := mysql.ParseDSN(connectURL)
	if err != nil {
		return "", err
	}
	config.DBName = databaseName // overwrite database name
	return config.FormatDSN(), nil
}
