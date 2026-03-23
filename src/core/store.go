package core

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type SettingsStore struct {
	mu       sync.RWMutex
	db       *sql.DB
	log      waLog.Logger
	settings map[types.JID]*GroupSettings
}

type GroupSettings struct {
	WelcomeEnabled  bool
	WelcomeMessage  string
	GoodbyeEnabled  bool
	GoodbyeMessage  string
	AntispamEnabled bool
}

type BotSettings struct {
	mu          sync.RWMutex
	SelfMode    bool
	PrivateOnly bool
	GroupOnly   bool
}

func NewSettingsStore(db *sql.DB, log waLog.Logger) *SettingsStore {
	s := &SettingsStore{
		db:       db,
		log:      log,
		settings: make(map[types.JID]*GroupSettings),
	}
	s.migrate()
	s.loadAll()
	return s
}

func (s *SettingsStore) migrate() {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS group_settings (
		chat_jid         TEXT PRIMARY KEY,
		welcome_enabled  INTEGER NOT NULL DEFAULT 0,
		welcome_message  TEXT NOT NULL DEFAULT '',
		goodbye_enabled  INTEGER NOT NULL DEFAULT 0,
		goodbye_message  TEXT NOT NULL DEFAULT '',
		antispam_enabled INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil && s.log != nil {
		s.log.Errorf("DB migrate error: %v", err)
	}
}

func (s *SettingsStore) loadAll() {
	rows, err := s.db.Query(`SELECT chat_jid, welcome_enabled, welcome_message,
		goodbye_enabled, goodbye_message, antispam_enabled FROM group_settings`)
	if err != nil {
		if s.log != nil {
			s.log.Errorf("DB loadAll error: %v", err)
		}
		return
	}
	defer rows.Close()
	for rows.Next() {
		var jidStr string
		var we, ge, ae int
		gs := &GroupSettings{}
		if err := rows.Scan(&jidStr, &we, &gs.WelcomeMessage, &ge, &gs.GoodbyeMessage, &ae); err != nil {
			if s.log != nil {
				s.log.Errorf("DB scan error: %v", err)
			}
			continue
		}
		gs.WelcomeEnabled = we == 1
		gs.GoodbyeEnabled = ge == 1
		gs.AntispamEnabled = ae == 1
		jid, err := types.ParseJID(jidStr)
		if err != nil {
			continue
		}
		s.settings[jid] = gs
	}
	if err := rows.Err(); err != nil && s.log != nil {
		s.log.Errorf("DB rows error: %v", err)
	}
}

func (s *SettingsStore) GetGroupSettings(jid types.JID) *GroupSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if gs, ok := s.settings[jid]; ok {
		return gs
	}
	return &GroupSettings{
		WelcomeEnabled:  false,
		WelcomeMessage:  "Selamat datang @user! Semoga betah di grup ini.",
		GoodbyeEnabled:  false,
		GoodbyeMessage:  "Selamat tinggal @user!",
		AntispamEnabled: false,
	}
}

func (s *SettingsStore) SetGroupSettings(jid types.JID, gs *GroupSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings[jid] = gs
	we, ge, ae := 0, 0, 0
	if gs.WelcomeEnabled {
		we = 1
	}
	if gs.GoodbyeEnabled {
		ge = 1
	}
	if gs.AntispamEnabled {
		ae = 1
	}
	_, err := s.db.Exec(`INSERT OR REPLACE INTO group_settings
		(chat_jid, welcome_enabled, welcome_message, goodbye_enabled, goodbye_message, antispam_enabled)
		VALUES (?, ?, ?, ?, ?, ?)`,
		jid.String(), we, gs.WelcomeMessage, ge, gs.GoodbyeMessage, ae,
	)
	if err != nil && s.log != nil {
		s.log.Errorf("DB save error: %v", err)
	}
}

func NewBotSettings() *BotSettings {
	return &BotSettings{}
}

func (bs *BotSettings) GetSelfMode() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.SelfMode
}

func (bs *BotSettings) SetSelfMode(v bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.SelfMode = v
}

func (bs *BotSettings) GetPrivateOnly() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.PrivateOnly
}

func (bs *BotSettings) SetPrivateOnly(v bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.PrivateOnly = v
}

func (bs *BotSettings) GetGroupOnly() bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.GroupOnly
}

func (bs *BotSettings) SetGroupOnly(v bool) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.GroupOnly = v
}
