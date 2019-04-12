package pgx

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func parsepgpass(line, cfgHost, cfgPort, cfgDatabase, cfgUsername string) *string {
	const (
		backslash = "\r"
		colon     = "\n"
	)
	const (
		host int = iota
		port
		database
		username
		pw
	)
	if strings.HasPrefix(line, "#") {
		return nil
	}
	line = strings.Replace(line, `\:`, colon, -1)
	line = strings.Replace(line, `\\`, backslash, -1)
	parts := strings.Split(line, `:`)
	if len(parts) != 5 {
		return nil
	}
	for i := range parts {
		if parts[i] == `*` {
			continue
		}
		parts[i] = strings.Replace(strings.Replace(parts[i], backslash, `\`, -1), colon, `:`, -1)
		switch i {
		case host:
			if parts[i] != cfgHost {
				return nil
			}
		case port:
			if parts[i] != cfgPort {
				return nil
			}
		case database:
			if parts[i] != cfgDatabase {
				return nil
			}
		case username:
			if parts[i] != cfgUsername {
				return nil
			}
		}
	}
	return &parts[4]
}

func pgpass(cfg *ConnConfig) (found bool) {
	passfile := os.Getenv("PGPASSFILE")
	if passfile == "" {
		u, err := user.Current()
		if err != nil {
			return
		}
		passfile = filepath.Join(u.HomeDir, ".pgpass")
	}
	f, err := os.Open(passfile)
	if err != nil {
		return
	}
	defer f.Close()

	host := cfg.Host
	if _, err := os.Stat(host); err == nil {
		host = "localhost"
	}
	port := fmt.Sprintf(`%v`, cfg.Port)
	if port == "0" {
		port = "5432"
	}
	username := cfg.User
	if username == "" {
		user, err := user.Current()
		if err != nil {
			return
		}
		username = user.Username
	}
	database := cfg.Database
	if database == "" {
		database = username
	}

	scanner := bufio.NewScanner(f)
	var pw *string
	for scanner.Scan() {
		pw = parsepgpass(scanner.Text(), host, port, database, username)
		if pw != nil {
			cfg.Password = *pw
			return true
		}
	}
	return false
}
