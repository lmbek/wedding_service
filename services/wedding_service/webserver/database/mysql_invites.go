package database

import (
	"errors"
	"strings"
	"time"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// normalizeSpace trims and collapses internal whitespace to a single space
// to make invitation code and names matching robust against stray spaces.
func normalizeSpace(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.Join(strings.Fields(s), " ")
}

// mysqlInvites is an unexported GORM-backed implementation of Invites.
type mysqlInvites struct {
	db *gorm.DB
}

// GORM models (unexported)
type invitedModel struct {
	ID         uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	Code       string `gorm:"column:code;size:64;not null"`
	InviteName string `gorm:"column:invite_name;size:255;not null"`
	MemberName string `gorm:"column:member_name;size:255;not null"`
	Accepted   *bool  `gorm:"column:accepted"`
}

func (invitedModel) TableName() string { return "invited" }

type invitedEventModel struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Code       string    `gorm:"column:code;size:64;not null"`
	MemberName string    `gorm:"column:member_name;size:255;not null"`
	Action     string    `gorm:"column:action;size:32;not null"`
	At         time.Time `gorm:"column:at;autoCreateTime"`
}

func (invitedEventModel) TableName() string { return "invited_events" }

type inviteVisitModel struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Code      string    `gorm:"column:code;size:64;not null"`
	IP        string    `gorm:"column:ip;size:64;not null"`
	UserAgent string    `gorm:"column:user_agent;size:512;not null"`
	Referer   string    `gorm:"column:referer;size:1024;not null"`
	Path      string    `gorm:"column:path;size:512;not null"`
	At        time.Time `gorm:"column:at;autoCreateTime"`
}

func (inviteVisitModel) TableName() string { return "invite_visits" }

// NewInvitesMySQL creates a GORM-backed Invites using the provided DSN.
// DSN example: user:pass@tcp(host:port)/dbname?parseTime=true&charset=utf8mb4,utf8
func NewInvitesMySQL(dsn string) (Invites, error) {
	if dsn == "" {
		return nil, errors.New("empty DSN")
	}
	// Ensure parseTime & utf8mb4 present
	hasParse := strings.Contains(dsn, "parseTime=")
	if !hasParse {
		if strings.Contains(dsn, "?") {
			dsn = dsn + "&parseTime=true"
		} else {
			dsn = dsn + "?parseTime=true"
		}
	}
	if !strings.Contains(dsn, "charset=") {
		if strings.Contains(dsn, "?") {
			dsn = dsn + "&charset=utf8mb4,utf8"
		} else {
			dsn = dsn + "?charset=utf8mb4,utf8"
		}
	}
	// Retry open + ping for robustness during MySQL startup
	var gdb *gorm.DB
	var err error
	start := time.Now()
	for {
		gdb, err = gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Verify ping on underlying sql.DB
			sqldb, e2 := gdb.DB()
			if e2 == nil {
				e3 := sqldb.Ping()
				if e3 == nil {
					break
				}
				err = e3
			} else {
				err = e2
			}
		}
		if time.Since(start) > 60*time.Second {
			return nil, err
		}
		time.Sleep(2 * time.Second)
	}
	mi := &mysqlInvites{db: gdb}
	err = mi.EnsureSchema()
	if err != nil {
		return nil, err
	}
	return mi, nil
}

func (m *mysqlInvites) EnsureSchema() error {
	err := m.db.AutoMigrate(&invitedModel{}, &invitedEventModel{}, &inviteVisitModel{})
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlInvites) FindByCode(code string) (Invite, bool) {
	code = normalizeSpace(code)
	if code == "" {
		return Invite{}, false
	}
	var rows []invitedModel
	err := m.db.Where("TRIM(code) = TRIM(?)", code).Order("member_name asc").Find(&rows).Error
	if err != nil {
		return Invite{}, false
	}
	if len(rows) == 0 {
		return Invite{}, false
	}
	inv := Invite{Code: code}
	inv.Name = rows[0].InviteName
	seen := make(map[string]struct{}, len(rows))
	members := make([]string, 0, len(rows))
	for _, r := range rows {
		if _, ok := seen[r.MemberName]; ok {
			continue
		}
		seen[r.MemberName] = struct{}{}
		members = append(members, r.MemberName)
	}
	inv.Members = members
	return inv, true
}

func (m *mysqlInvites) ListAccepted(code string) ([]string, error) {
	code = normalizeSpace(code)
	var rows []invitedModel
	err := m.db.Where("TRIM(code) = TRIM(?) AND accepted = ?", code, true).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(rows))
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if _, ok := seen[r.MemberName]; ok {
			continue
		}
		seen[r.MemberName] = struct{}{}
		out = append(out, r.MemberName)
	}
	return out, nil
}

func (m *mysqlInvites) ListAllAccepted() ([]string, error) {
	var rows []invitedModel
	err := m.db.Where("accepted = ?", true).Order("member_name asc").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(rows))
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if _, ok := seen[r.MemberName]; ok {
			continue
		}
		seen[r.MemberName] = struct{}{}
		out = append(out, r.MemberName)
	}
	return out, nil
}

func (m *mysqlInvites) Accept(code, name string) error {
	code = normalizeSpace(code)
	name = normalizeSpace(name)
	if code == "" || name == "" {
		return errors.New("code and name required")
	}
	// Find existing row (must exist); do NOT insert new rows
	var rec invitedModel
	err := m.db.Where("TRIM(code) = TRIM(?) AND TRIM(member_name) = TRIM(?)", code, name).Take(&rec).Error
	if err != nil {
		return err
	}
	acc := true
	err = m.db.Model(&invitedModel{}).Where("id = ?", rec.ID).Update("accepted", acc).Error
	if err != nil {
		return err
	}
	// Log event
	evt := invitedEventModel{Code: code, MemberName: name, Action: "accept", At: time.Now()}
	return m.db.Create(&evt).Error
}

func (m *mysqlInvites) Decline(code, name string) error {
	code = normalizeSpace(code)
	name = normalizeSpace(name)
	if code == "" || name == "" {
		return errors.New("code and name required")
	}
	// Only update existing row; do not insert placeholders
	var rec invitedModel
	err := m.db.Where("TRIM(code) = TRIM(?) AND TRIM(member_name) = TRIM(?)", code, name).Take(&rec).Error
	if err != nil {
		return err
	}
	err = m.db.Model(&invitedModel{}).Where("id = ?", rec.ID).Update("accepted", nil).Error
	if err != nil {
		return err
	}
	// Log event
	evt := invitedEventModel{Code: code, MemberName: name, Action: "decline", At: time.Now()}
	return m.db.Create(&evt).Error
}

func (m *mysqlInvites) TrackVisit(code, ip, userAgent, referer, path string) error {
	code = normalizeSpace(code)
	if code == "" {
		return errors.New("code required")
	}
	rec := inviteVisitModel{Code: code, IP: ip, UserAgent: userAgent, Referer: referer, Path: path, At: time.Now()}
	return m.db.Create(&rec).Error
}
