package core

import (
	"database/sql"
	"sptzx/src/api"
	"sptzx/src/config"
	"sptzx/src/middleware"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	Client    *whatsmeow.Client
	Config    *config.Config
	Antispam  *middleware.Antispam
	Settings  *SettingsStore
	BotConfig *BotSettings
	Registry  *Registry
	Container *sqlstore.Container
	Log       waLog.Logger
	API       *api.Client
}

func NewBot(cfg *config.Config, container *sqlstore.Container, client *whatsmeow.Client, log waLog.Logger, db *sql.DB) *Bot {
	apiClient := api.NewClient(cfg.SiputzX.BaseURL)
	apiClient.SetLogger(log)

	return &Bot{
		Client:    client,
		Config:    cfg,
		Antispam:  middleware.NewAntispam(cfg.Antispam.MaxMsgPerSecond, cfg.Antispam.MaxMsgPerMinute, cfg.Antispam.BanDurationSecs),
		Settings:  NewSettingsStore(db, log),
		BotConfig: NewBotSettings(),
		Container: container,
		Log:       log,
		API:       apiClient,
	}
}

func (b *Bot) GetPrefix() string {
	if len(b.Config.Prefixes) > 0 {
		return b.Config.Prefixes[0]
	}
	return "!"
}
