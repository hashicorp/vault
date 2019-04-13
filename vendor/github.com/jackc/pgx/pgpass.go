package pgx

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func parsepgpass(cfg *ConnConfig, line string) *string {
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
			if parts[i] != cfg.Host {
				return nil
			}
		case port:
			portstr := fmt.Sprintf(`%v`, cfg.Port)
			if portstr == "0" {
				portstr = "5432"
			}
			if parts[i] != portstr {
				return nil
			}
		case database:
			if parts[i] != cfg.Database {
				return nil
			}
		case username:
			if parts[i] != cfg.User {
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
	scanner := bufio.NewScanner(f)
	var pw *string
	for scanner.Scan() {
		pw = parsepgpass(cfg, scanner.Text())
		if pw != nil {
			cfg.Password = *pw
			return true
		}
	}
	return false
}
