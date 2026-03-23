package stalk

import (
	"fmt"
	"sptzx/src/serialize"
	"strings"
	"time"

	"sptzx/src/api"
	"sptzx/src/core"
)

func init() {
	core.Use(&core.Command{
		Name:        "stalktwitter",
		Aliases:     []string{"twitterstalk", "stalkx"},
		Description: "Stalk profil Twitter/X",
		Usage:       "stalktwitter <username>",
		Category:    "stalk",
		Handler: func(ptz *core.Ptz) error {
			if len(ptz.Args) == 0 {
				return ptz.ReplyText("*stalktwitter* — Stalk profil Twitter/X\n\nUsage: .stalktwitter <username>")
			}

			ptz.React("⏳")
			defer ptz.Unreact()

			username := strings.Join(ptz.Args, "")
			username = strings.TrimPrefix(username, "@")

			ctx, cancel := ptz.ContextWithTimeout(30 * time.Second)
			defer cancel()

			type TwitterData struct {
				ID          string `json:"id"`
				Username    string `json:"username"`
				Name        string `json:"name"`
				Verified    bool   `json:"verified"`
				Description string `json:"description"`
				Location    string `json:"location"`
				CreatedAt   string `json:"created_at"`
				Stats       struct {
					Tweets    float64 `json:"tweets"`
					Following float64 `json:"following"`
					Followers float64 `json:"followers"`
					Likes     float64 `json:"likes"`
					Media     float64 `json:"media"`
				} `json:"stats"`
				Profile struct {
					Image  string  `json:"image"`
					Banner *string `json:"banner"`
				} `json:"profile"`
			}

			data, err := api.Request[TwitterData](ctx, ptz.Bot.API, "/api/stalk/twitter", map[string]string{"user": username})
			if err != nil {
				ptz.Bot.Log.Errorf("Twitter stalk error: %v", err)
				return ptz.ReplyText("❌ Terjadi kesalahan pada server, silakan coba lagi nanti.")
			}

			verified := ""
			if data.Verified {
				verified = " ✅"
			}

			caption := fmt.Sprintf("🐦 *Twitter / X*\n\n"+
				"👤 %s (@%s)%s\n",
				data.Name, data.Username, verified)

			if data.Description != "" {
				caption += "📝 " + data.Description + "\n"
			}
			if data.Location != "" {
				caption += "📍 " + data.Location + "\n"
			}

			caption += fmt.Sprintf("\n"+
				"🐦 *%s* Tweets\n"+
				"👥 *%s* Followers\n"+
				"➕ *%s* Following\n"+
				"❤️ *%s* Likes\n"+
				"🖼 *%s* Media",
				serialize.NumFmt(data.Stats.Tweets),
				serialize.NumFmt(data.Stats.Followers),
				serialize.NumFmt(data.Stats.Following),
				serialize.NumFmt(data.Stats.Likes),
				serialize.NumFmt(data.Stats.Media),
			)

			if len(data.CreatedAt) >= 10 {
				caption += "\n📅 Joined: " + data.CreatedAt[:10]
			}

			imgData, err := serialize.Fetch(data.Profile.Image)
			if err != nil {
				return ptz.ReplyText(caption)
			}
			return ptz.ReplyImage(imgData, "image/jpeg", caption)
		},
	})
}
