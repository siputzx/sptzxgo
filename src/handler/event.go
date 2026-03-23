package handler

import (
	"fmt"
	"sptzx/src/commands/games"
	"sptzx/src/core"
	"sptzx/src/serialize"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type EventHandler struct {
	bot         *core.Bot
	registry    *core.Registry
	groupEvents *GroupEventHandler
}

func NewEventHandler(bot *core.Bot, registry *core.Registry) *EventHandler {
	return &EventHandler{
		bot:         bot,
		registry:    registry,
		groupEvents: NewGroupEventHandler(bot),
	}
}

func (h *EventHandler) Handle(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.Message:
		go h.onMessage(evt)

	case *events.GroupInfo:
		go h.groupEvents.OnGroupInfo(evt)
	case *events.Picture:
		go h.groupEvents.OnPicture(evt)
	case *events.JoinedGroup:
		go h.groupEvents.OnJoinedGroup(evt)

	case *events.Connected:
		h.bot.Log.Infof("✅ Connected to WhatsApp")

	case *events.Disconnected:
		h.bot.Log.Warnf("⚠️ Disconnected from WhatsApp")

	case *events.LoggedOut:
		h.bot.Log.Errorf("❌ Logged out: %s", evt.Reason.String())

	case *events.PairSuccess:
		h.bot.Log.Infof("✅ Paired: %s", evt.ID.String())

	case *events.KeepAliveTimeout:
		h.bot.Log.Warnf("⚠️ Keepalive timeout — errors: %d, last ok: %s ago",
			evt.ErrorCount, evt.LastSuccess.String())

	case *events.KeepAliveRestored:
		h.bot.Log.Infof("✅ Keepalive restored")

	case *events.OfflineSyncCompleted:
		h.bot.Log.Infof("📥 Offline sync done — %d messages", evt.Count)

	case *events.TemporaryBan:
		h.bot.Log.Errorf("🚫 Temporary ban: %s (expires in %s)", evt.Code, evt.Expire)

	case *events.ClientOutdated:
		h.bot.Log.Errorf("❌ Client outdated — update whatsmeow")

	case *events.StreamReplaced:
		h.bot.Log.Warnf("⚠️ Stream replaced — another session opened")

	case *events.UndecryptableMessage:
		h.bot.Log.Warnf("⚠️ Undecryptable message from %s", evt.Info.Sender.String())

	case *events.IdentityChange:
		h.bot.Log.Warnf("🔑 Identity changed: %s (implicit: %v)", evt.JID.String(), evt.Implicit)

	case *events.PrivacySettings:
		h.bot.Log.Debugf("🔒 Privacy settings updated")

	case *events.PushName:
		h.bot.Log.Debugf("📛 Push name update: %s → %s", evt.JID.String(), evt.NewPushName)

	case *events.Receipt:
		h.bot.Log.Debugf("📨 Receipt [%s] from %s", evt.Type, evt.MessageSource.Sender.String())

	case *events.CallOffer:
		h.bot.Log.Infof("📞 Incoming call from %s", evt.From.String())

	case *events.CallTerminate:
		h.bot.Log.Infof("📵 Call ended from %s", evt.From.String())
	}
}

func resolveSender(info types.MessageInfo) (phoneJID types.JID, displayName string) {
	sender := info.Sender
	alt := info.SenderAlt

	if sender.Server == types.HiddenUserServer {
		if !alt.IsEmpty() && alt.Server == types.DefaultUserServer {
			phoneJID = types.NewJID(alt.User, types.DefaultUserServer)
		} else {
			phoneJID = types.NewJID(sender.User, types.DefaultUserServer)
		}
	} else {
		phoneJID = types.NewJID(sender.User, types.DefaultUserServer)
	}

	if info.PushName != "" && info.PushName != "-" {
		displayName = info.PushName
	} else {
		displayName = phoneJID.User
	}
	return
}

func isCommand(text string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(text, p) {
			return true
		}
	}
	return false
}

func getReplyInfo(msg *waE2E.Message) (quotedID, userAnswer string) {
	if msg == nil {
		return "", ""
	}
	switch {
	case msg.ExtendedTextMessage != nil:
		ctx := msg.ExtendedTextMessage.GetContextInfo()
		if ctx.GetStanzaID() == "" {
			return "", ""
		}
		return ctx.GetStanzaID(), strings.TrimSpace(msg.ExtendedTextMessage.GetText())

	case msg.ImageMessage != nil:
		ctx := msg.ImageMessage.GetContextInfo()
		return ctx.GetStanzaID(), strings.TrimSpace(msg.ImageMessage.GetCaption())

	case msg.VideoMessage != nil:
		ctx := msg.VideoMessage.GetContextInfo()
		return ctx.GetStanzaID(), strings.TrimSpace(msg.VideoMessage.GetCaption())

	case msg.AudioMessage != nil:
		ctx := msg.AudioMessage.GetContextInfo()
		return ctx.GetStanzaID(), ""

	case msg.DocumentMessage != nil:
		ctx := msg.DocumentMessage.GetContextInfo()
		return ctx.GetStanzaID(), strings.TrimSpace(msg.DocumentMessage.GetCaption())
	}
	return "", ""
}

func (h *EventHandler) onMessage(evt *events.Message) {
	if evt == nil || evt.Message == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			h.bot.Log.Errorf("Recovered from panic in onMessage: %v", r)
		}
	}()

	quotedID, userAnswer := getReplyInfo(evt.Message)
	if quotedID != "" && userAnswer != "" && !isCommand(userAnswer, h.bot.Config.Prefixes) {
		if res, sess, found := games.CheckAnswer(quotedID, userAnswer); found {
			phoneJID, _ := resolveSender(evt.Info)
			ptz := core.NewPtz(h.bot, evt)

			switch res {
			case games.AnswerCorrect:
				text := fmt.Sprintf(
					"✅ *Benar!* Selamat @%s!\n\n🎮 Jawaban: *%s*",
					phoneJID.User, sess.Answer,
				)
				serialize.SendTextMentionReplyToID(
					h.bot.Client,
					evt.Info.Chat,
					text,
					[]types.JID{phoneJID},
					sess.QuestionID,
					phoneJID.String(),
				)
				games.DeleteSession(sess)
			case games.AnswerVeryClose:
				ptz.ReplyText("⚠️ *Dikit lagi bener!* Jawaban kamu udah semakin mirip.")
			case games.AnswerGettingClose:
				ptz.ReplyText("🤔 *Jawaban mendekati.* Coba pikirin lagi kata-katanya.")
			case games.AnswerWrong:
				ptz.ReplyText("❌ *Salah!* Coba lagi atau tunggu soal expired.")
			}
			return
		}

		ptzForGame := core.NewPtz(h.bot, evt)
		if handled, err := games.ProcessCcsdAnswer(ptzForGame, quotedID, userAnswer); handled {
			if err != nil {
				h.bot.Log.Errorf("ccsd answer error: %v", err)
			}
			return
		}
		if handled, err := games.ProcessFamily100Answer(ptzForGame, quotedID, userAnswer); handled {
			if err != nil {
				h.bot.Log.Errorf("family100 answer error: %v", err)
			}
			return
		}
	}

	if evt.Info.IsFromMe {
		return
	}

	h.logMessage(evt)

	if !h.bot.Antispam.Check(evt.Info.Sender.User) {
		return
	}

	ptz := core.NewPtz(h.bot, evt)

	if h.bot.BotConfig.GetSelfMode() && !ptz.IsOwner() {
		return
	}

	if h.bot.BotConfig.GetPrivateOnly() && ptz.IsGroup {
		return
	}

	if h.bot.BotConfig.GetGroupOnly() && !ptz.IsGroup {
		return
	}

	if ptz.Command == "" {
		return
	}

	cmd, ok := h.registry.Get(ptz.Command)
	if !ok {
		return
	}

	if err := cmd.Execute(ptz); err != nil {
		h.bot.Log.Errorf("Command %s error: %v", cmd.Name, err)
	}
}

func (h *EventHandler) logMessage(evt *events.Message) {
	info := evt.Info
	msg := evt.Message

	msgType, content := serialize.ClassifyMessage(msg)

	location := "DM"
	sender := info.Sender.User
	if info.IsGroup {
		location = fmt.Sprintf("Group[%s]", info.Chat.User)
		if info.PushName != "" && info.PushName != "-" {
			sender = info.PushName
		}
	}

	h.bot.Log.Infof("[%s] %s → %s: %s", location, sender, msgType, content)
}
