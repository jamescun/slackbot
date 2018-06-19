package slackbot

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Request is a the parsed result of a received webhook slash command.
type Request struct {
	TriggerID string

	Token string

	Team
	Enterprise
	Channel
	User

	Command string
	Text    string

	ResponseURL string
}

// ParseRequest extracts known values from a slack command webhook into
// a new Request object.
func ParseRequest(v url.Values) *Request {
	return &Request{
		TriggerID: v.Get("trigger_id"),
		Token:     v.Get("token"),
		Team: Team{
			ID:     v.Get("team_id"),
			Domain: v.Get("team_domain"),
		},
		Enterprise: Enterprise{
			ID:   v.Get("enterprise_id"),
			Name: v.Get("enterprise_name"),
		},
		Channel: Channel{
			ID:   v.Get("channel_id"),
			Name: v.Get("channel_name"),
		},
		User: User{
			ID:   v.Get("user"),
			Name: v.Get("user"),
		},
		Command:     v.Get("command"),
		Text:        v.Get("text"),
		ResponseURL: v.Get("response_url"),
	}
}

// Author will display a small section at the top of a message attachment.
type Author struct {
	// Small text used to display the author's name.
	Name string `json:"author_name,omitempty"`

	// A valid URL that will hyperlink the author_name text mentioned above.
	// Will only work if author_name is present.
	Link string `json:"author_link,omitempty"`

	// A valid URL that displays a small 16x16px image to the left of the
	// author_name text. Will only work if author_name is present.
	Icon string `json:"author_icon,omitempty"`
}

// Title is displayed as larger, bold text near the top of a message attachment
type Title struct {
	Title string `json:"title,omitempty"`
	Link  string `json:"title_link,omitempty"`
}

// Field is defined as a dictionary with key-value pairs. Fields get displayed
// in a table-like way.
type Field struct {
	// Shown as a bold heading above the value text. It cannot contain markup.
	Title string `json:"title,omitempty"`

	// The text value of the field. It may contain standard message markup and
	// must be escaped as normal. May be multi-line.
	Value string `json:"value,omitempty"`

	// An optional flag indicating whether the value is short enough to be
	// displayed side-by-side with other values.
	Short bool `json:"short"`
}

// Attachment is added context to the response given to a Command.
type Attachment struct {
	// A plain-text summary of the attachment. This text will be used in
	// clients that don't show formatted text (eg. IRC, mobile notifications)
	// and should not contain any markup.
	Fallback string `json:"fallback"`

	// An optional value that can either be one of good, warning, danger, or
	// any hex color code (eg. #439FE0). This value is used to color the border
	// along the left side of the message attachment.
	Color string `json:"color,omitempty"`

	// This is optional text that appears above the message attachment block.
	Pretext string `json:"pretext,omitempty"`

	// Optional text that appears within the attachment.
	Text string `json:"text,omitempty"`

	*Author
	*Title

	// Fields are defined as an array and get displayed in a table-like way.
	Fields []*Field `json:"fields,omitempty"`

	// A valid URL to an image file that will be displayed inside a message
	// attachment.
	ImageURL string `json:"image_url,omitempty"`

	// A valid URL to an image file that will be displayed as a thumbnail on
	// the right side of a message attachment.
	ThumbURL string `json:"thumb_url,omitempty"`

	// Brief text to help contextualize and identify an attachment. Limited to
	// 300 characters, and may be truncated further when displayed to users in
	// environments with limited screen real estate.
	Footer string `json:"footer,omitempty"`

	// A small icon beside your footer text, provide a publicly accessible URL
	// string in the footer_icon field. You must also provide a footer for the
	// field to be recognized.
	FooterIcon string `json:"footer_icon,omitempty"`

	// attachment will display an additional timestamp value as part of the
	// attachment's footer.
	Timestamp int64 `json:"ts,omitempty"`
}

// Response is the optional response to a Request.
type Response struct {
	Text string `json:"text,omitempty"`

	Attachments []*Attachment `json:"attachments,omitempty"`
}

// HandlerFunc is invoked by Handle on each slash command webhook.
type HandlerFunc func(*Request) *Response

// Handle invokes the given HandlerFunc for each slack command webhook
// received. It does not attempt to verify the security token.
func Handle(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()

		req := ParseRequest(r.PostForm)

		res := fn(req)
		if res != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			json.NewEncoder(w).Encode(res)
		}
	}
}
