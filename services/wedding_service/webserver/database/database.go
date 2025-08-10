package database

import (
	"errors"
	"strings"
	"time"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBInit provides an interface to ensure the MySQL database exists using GORM.
// Interface-first: do not expose concrete structs from this package.
// It connects without selecting a database and issues a CREATE DATABASE IF NOT EXISTS.
// This is helpful when the MySQL data dir is empty or DB was not yet created.
// Note: proper privileges are required for the configured user.
// We avoid returning *gorm.DB to keep struct types unexposed.
//
// Usage example (call once on startup before opening schema-specific connections):
//   _ = database.NewDBInit().EnsureDatabase(host, port, user, pass, dbname)
//
// It is safe to ignore the error if init SQL or the container already created the DB.
// The init scripts under services/mysql_service/initdb handle schema and seed data.
// This initializer only ensures the database container has the database itself.
//
// We keep it minimal and dependency-injected via the returned interface.
// No short if form is used per project guidelines.
//
//go:generate echo "database.DBInit is interface-first; no codegen required"

type DBInit interface {
	EnsureDatabase(host, port, user, pass, dbname string) error
}

type dbinit struct{}

func NewDBInit() DBInit { return &dbinit{} }

type UserInit interface {
	EnsureUser(host, port, rootPass, dbname, user, pass string) error
}

type userinit struct{}

func NewUserInit() UserInit { return &userinit{} }

func (u *userinit) EnsureUser(host, port, rootPass, dbname, user, pass string) error {
	if host == "" || port == "" || user == "" || dbname == "" {
		return errors.New("missing user init parameters")
	}
	if rootPass == "" {
		return nil
	}
	// root DSN without db selected
	dsn := "root:" + rootPass + "@tcp(" + host + ":" + port + ")/"
	args := []string{"parseTime=true", "charset=utf8mb4,utf8"}
	suffix := strings.Join(args, "&")
	dsn = dsn + "?" + suffix
	// retry up to ~60s
	var lastErr error
	start := time.Now()
	for {
		gdb, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Verify server reachable
			res := gdb.Exec("SELECT 1")
			if res.Error == nil {
				// Ensure user at '%' and localhost and grant on db
				st1 := gdb.Exec("CREATE USER IF NOT EXISTS `"+user+"`@'%' IDENTIFIED BY ?", pass)
				if st1.Error != nil {
					lastErr = st1.Error
				} else {
					st1b := gdb.Exec("ALTER USER `"+user+"`@'%' IDENTIFIED BY ?", pass)
					if st1b.Error != nil {
						lastErr = st1b.Error
					}
				}
				st2 := gdb.Exec("GRANT ALL PRIVILEGES ON `" + dbname + "`.* TO `" + user + "`@'%'")
				if st2.Error != nil {
					lastErr = st2.Error
				}
				st3 := gdb.Exec("CREATE USER IF NOT EXISTS `"+user+"`@'localhost' IDENTIFIED BY ?", pass)
				if st3.Error != nil {
					lastErr = st3.Error
				} else {
					st3b := gdb.Exec("ALTER USER `"+user+"`@'localhost' IDENTIFIED BY ?", pass)
					if st3b.Error != nil {
						lastErr = st3b.Error
					}
				}
				st4 := gdb.Exec("GRANT ALL PRIVILEGES ON `" + dbname + "`.* TO `" + user + "`@'localhost'")
				if st4.Error != nil {
					lastErr = st4.Error
				}
				st5 := gdb.Exec("FLUSH PRIVILEGES")
				if st5.Error == nil {
					// close pool and return
					sqldb, e2 := gdb.DB()
					if e2 == nil {
						_ = sqldb.Close()
					}
					return nil
				}
				lastErr = st5.Error
			} else {
				lastErr = res.Error
			}
			// close before retry
			sqldb, e2 := gdb.DB()
			if e2 == nil {
				_ = sqldb.Close()
			}
		} else {
			lastErr = err
		}
		if time.Since(start) > 60*time.Second {
			return lastErr
		}
		time.Sleep(2 * time.Second)
	}
}

func (d *dbinit) EnsureDatabase(host, port, user, pass, dbname string) error {
	// Validate inputs
	if host == "" || port == "" || user == "" || dbname == "" {
		return errors.New("missing database connection parameters")
	}
	// Build DSN without database segment. MySQL allows a DSN ending with '/'.
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/"
	// Make sure parseTime and utf8mb4 are present for consistency
	args := []string{"parseTime=true", "charset=utf8mb4,utf8"}
	suffix := strings.Join(args, "&")
	dsn = dsn + "?" + suffix

	// Retry loop: wait up to ~60s for MySQL to accept connections
	var lastErr error
	start := time.Now()
	for {
		gdb, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Try a simple statement to ensure the server is reachable
			res := gdb.Exec("SELECT 1")
			if res.Error == nil {
				// Create DB if not exists
				res2 := gdb.Exec("CREATE DATABASE IF NOT EXISTS `" + dbname + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
				if res2.Error == nil {
					// Close pool
					sqldb, err2 := gdb.DB()
					if err2 == nil {
						_ = sqldb.Close()
					}
					return nil
				}
				lastErr = res2.Error
			} else {
				lastErr = res.Error
			}
			// Best-effort close before retry
			sqldb, err2 := gdb.DB()
			if err2 == nil {
				_ = sqldb.Close()
			}
		} else {
			lastErr = err
		}
		if time.Since(start) > 60*time.Second {
			return lastErr
		}
		time.Sleep(2 * time.Second)
	}
}
