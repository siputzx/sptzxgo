package games

import (
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"sptzx/src/core"
	"sptzx/src/serialize"
)

var allGameTypes = []string{
	"caklontong", "tebakbendera", "tebakkata", "tebaklagu", "tekateki",
}

func findMySession(ptz *core.Ptz) (*Session, bool) {
	if ptz.IsGroup {
		return GetActiveChatSessionAny(ptz.Chat.String())
	}
	return GetActiveUserSessionAny(ptz.Chat.String(), ptz.Sender.String())
}

func init() {
	core.Use(&core.Command{
		Name:        "clue",
		Aliases:     []string{"hint"},
		Description: "Minta clue untuk soal aktif (harus reply ke pesan soal)",
		Usage:       "clue",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			quotedID := getQuotedID(ptz)
			if quotedID == "" {
				return ptz.ReplyText("💡 Untuk minta clue, *reply pesan soal* lalu ketik !clue")
			}

			sess, ok := MatchByQuestionID(quotedID)
			if !ok {
				return ptz.ReplyText("❓ Tidak ada soal aktif pada pesan itu.")
			}

			if ptz.IsGroup {
				if sess.ChatJID != ptz.Chat.String() {
					return ptz.ReplyText("❓ Soal ini bukan untuk grup ini.")
				}
			} else {
				if sess.StarterJID != ptz.Sender.String() {
					return ptz.ReplyText("❓ Soal ini bukan milikmu.")
				}
			}

			if sess.ClueCount >= 3 {
				return ptz.ReplyText("🚫 Clue sudah habis! Maksimal 3 clue per soal.\nKetik *!nyerah* (reply soal) untuk menyerah.")
			}

			sess.ClueCount++
			UpdateClueCount(sess)

			return ptz.ReplyText(formatClueMessage(sess))
		},
	})

	core.Use(&core.Command{
		Name:        "nyerah",
		Aliases:     []string{"giveup", "skip"},
		Description: "Menyerah dan lihat jawaban (harus reply ke pesan soal)",
		Usage:       "nyerah",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			quotedID := getQuotedID(ptz)
			if quotedID == "" {
				return ptz.ReplyText("🏳️ Untuk menyerah, *reply pesan soal* lalu ketik !nyerah")
			}

			sess, ok := MatchByQuestionID(quotedID)
			if !ok {
				return ptz.ReplyText("❓ Tidak ada soal aktif pada pesan itu.")
			}

			if ptz.IsGroup {
				if sess.ChatJID != ptz.Chat.String() {
					return ptz.ReplyText("❓ Soal ini bukan untuk grup ini.")
				}
			} else {
				if sess.StarterJID != ptz.Sender.String() {
					return ptz.ReplyText("❓ Soal ini bukan milikmu.")
				}
			}

			name := gameTypeNames[sess.GameType]
			if name == "" {
				name = sess.GameType
			}

			DeleteSession(sess)

			return ptz.ReplyText(fmt.Sprintf(
				"🏳️ *Nyerah!*\n\nGame: %s\nJawaban: *%s*\n\n_Semangat untuk soal berikutnya!_",
				name, sess.Answer,
			))
		},
	})

	core.Use(&core.Command{
		Name:        "soalku",
		Aliases:     []string{"myquestion", "aktif"},
		Description: "Lihat soal aktif saat ini",
		Usage:       "soalku",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			sess, ok := findMySession(ptz)
			if !ok {
				return ptz.ReplyText("❓ Tidak ada soal aktif saat ini.")
			}

			name := gameTypeNames[sess.GameType]
			if name == "" {
				name = sess.GameType
			}

			remaining := getRemainingTimeFromSess(sess)

			var clueInfo string
			if sess.ClueCount > 0 {
				clueInfo = fmt.Sprintf("\nClue dipakai: %d/3", sess.ClueCount)
			}

			var starterInfo string
			if ptz.IsGroup {
				starterInfo = fmt.Sprintf("\nDimulai oleh: @%s", sess.StarterJID)
			}

			return ptz.ReplyText(fmt.Sprintf(
				"🎮 *Soal Aktif*\n\nGame: %s\nSisa waktu: %s%s%s\n\n_Reply pesan soal untuk menjawab, minta clue, atau menyerah_",
				name, remaining, clueInfo, starterInfo,
			))
		},
	})

	core.Use(&core.Command{
		Name:        "daftargame",
		Aliases:     []string{"listgame", "games"},
		Description: "Lihat semua game yang tersedia",
		Usage:       "daftargame",
		Category:    "games",
		Handler: func(ptz *core.Ptz) error {
			p := ptz.Bot.GetPrefix()
			var sb strings.Builder
			sb.WriteString("🎮 *Daftar Game*\n\n")
			for _, gt := range allGameTypes {
				name := gameTypeNames[gt]
				sb.WriteString(fmt.Sprintf("• *%s%s* — %s\n", p, gt, name))
			}
			sb.WriteString(fmt.Sprintf("\n📌 *Command Pendukung*\n"))
			sb.WriteString(fmt.Sprintf("• *%sclue* — minta petunjuk, maks 3x (reply soal)\n", p))
			sb.WriteString(fmt.Sprintf("• *%snyerah* — lihat jawaban & akhiri (reply soal)\n", p))
			sb.WriteString(fmt.Sprintf("• *%ssoalku* — cek soal aktif\n", p))
			return ptz.ReplyText(sb.String())
		},
	})
}

func getQuotedID(ptz *core.Ptz) string {
	if ptz.Message == nil {
		return ""
	}
	if ptz.Message.ExtendedTextMessage != nil {
		return ptz.Message.ExtendedTextMessage.GetContextInfo().GetStanzaID()
	}
	return ""
}

func mentionJID(phone string) types.JID {
	return types.NewJID(phone, types.DefaultUserServer)
}

func replyWithMention(ptz *core.Ptz, text string, phones ...string) error {
	jids := make([]types.JID, 0, len(phones))
	for _, p := range phones {
		jids = append(jids, mentionJID(p))
	}
	return serialize.SendTextReplyMention(
		ptz.Bot.Client, ptz.Chat, text, jids, ptz.Message, ptz.Info,
	)
}
