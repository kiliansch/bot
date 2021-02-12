package topflop

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/neovg/kmptnzbot/internal/db"
	"github.com/neovg/kmptnzbot/internal/matcher/abstract"
	"github.com/neovg/kmptnzbot/internal/telegram"
)

// Each matcher extends the abstract matcher
type Matcher struct {
	abstract.Matcher
}

// Return the identifier of this matcher for use in logging
func (m Matcher) Identifier() string {
	return "topflop"
}

// Process a message received from Telegram
func (m Matcher) ProcessRequestMessage(requestMessage telegram.RequestMessage) error {
	// Check if text starts with /top or /flop and if not, ignore it
	match := m.getMatch(requestMessage.Text)
	if len(match) == 0 {
		return nil
	}

	var records []db.Plusplus
	switch match {
	case "top":
		records = db.FindTops()
	case "flop":
		records = db.FindFlops()
	}

	// Choose one option and send the result
	return m.sendResultResponse(requestMessage, records)
}

// Check if a text starts with /top or /flop
func (m Matcher) getMatch(text string) string {
	// Initialize the regular expression
	r := regexp.MustCompile(`^(?:/)(top|flop)(\s|$)`)

	// Find either "top" or "flop"
	match := r.FindString(text)
	match = strings.TrimLeft(match, "/")

	return match
}

// Send the result to the user who sent the request message
func (m Matcher) sendResultResponse(requestMessage telegram.RequestMessage, records []db.Plusplus) error {
	responseText := "```"

	// Add one record per record
	for _, record := range records {
		responseText = responseText + fmt.Sprintf("\n%5d | %s", record.Value, record.Name)
	}

	responseText = responseText + "```"

	responseMessage := telegram.Message{
		Text:             responseText,
		ReplyToMessageID: requestMessage.ID,
		ParseMode:        "Markdown",
	}

	return telegram.SendMessage(requestMessage, responseMessage)
}