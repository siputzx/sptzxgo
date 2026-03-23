package games

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
)

var gameTypeNames = map[string]string{
	"caklontong":      "Cak Lontong",
	"tebakbendera":    "Tebak Bendera",
	"tebakkata":       "Tebak Kata",
	"tebaklagu":       "Tebak Lagu",
	"tekateki":        "Teka-Teki",
	"asahotak":        "Asah Otak",
	"siapakahaku":     "Siapakah Aku",
	"susunkata":       "Susun Kata",
	"tebakgambar":     "Tebak Gambar",
	"tebakkimia":      "Tebak Kimia",
	"tebakkalimat":    "Tebak Kalimat",
	"lengkapikalimat": "Lengkapi Kalimat",
	"tebaktebakan":    "Tebak Tebakan",
	"tebaklogo":       "Tebak Logo",
}

func buildClue(answer string, clueCount int) string {
	runes := []rune(answer)
	n := len(runes)

	revealed := make([]bool, n)

	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			revealed[i] = true
		}
	}

	totalLetters := 0
	for _, r := range revealed {
		if !r {
			totalLetters++
		}
	}

	revealCount := 0
	switch {
	case clueCount == 1:
		revealCount = max(1, totalLetters/4)
	case clueCount == 2:
		revealCount = max(1, totalLetters/2)
	default:
		revealCount = max(1, (totalLetters*3)/4)
	}

	hiddenIdx := make([]int, 0, totalLetters)
	for i, r := range revealed {
		if !r {
			hiddenIdx = append(hiddenIdx, i)
		}
	}

	rand.Shuffle(len(hiddenIdx), func(i, j int) {
		hiddenIdx[i], hiddenIdx[j] = hiddenIdx[j], hiddenIdx[i]
	})
	for i := 0; i < revealCount && i < len(hiddenIdx); i++ {
		revealed[hiddenIdx[i]] = true
	}

	var sb strings.Builder
	for i, r := range runes {
		if revealed[i] {
			if unicode.IsSpace(r) {
				sb.WriteRune(' ')
			} else {
				sb.WriteString(strings.ToUpper(string(r)))
			}
		} else {
			sb.WriteString("_")
		}
		if i < n-1 && unicode.IsLetter(r) && unicode.IsLetter(runes[i+1]) {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

func clueLabel(clueCount int) string {
	switch clueCount {
	case 1:
		return "💡 *Clue 1* (~25% terbuka)"
	case 2:
		return "💡 *Clue 2* (~50% terbuka)"
	default:
		return "💡 *Clue 3* (~75% terbuka)"
	}
}

func formatClueMessage(sess *Session) string {
	name := gameTypeNames[sess.GameType]
	if name == "" {
		name = sess.GameType
	}
	clue := buildClue(sess.Answer, sess.ClueCount)
	remaining := getRemainingTime(sess.ChatJID, sess.StarterJID, sess.GameType)
	return fmt.Sprintf(
		"%s\n\n`%s`\n\n_Sisa waktu: %s_",
		clueLabel(sess.ClueCount),
		clue,
		remaining,
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
