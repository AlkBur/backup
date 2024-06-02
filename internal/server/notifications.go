package server

import (
	"backup/internal/config"
	"log"
	"net/http"
	"net/url"
	"path"
)

// SendNotifications will send a notification to all registered apps
func (s *Server) SendNotifications(title, message, priority string) {
	go NotifyGotify(s.cfgInfo, title, message, priority)
}

// NotifyGotify messages gotify
func NotifyGotify(config config.Notification, title, message, priority string) {
	if config.URL != "" && config.Token != "" {
		u, err := url.Parse(config.URL)
		if err != nil {
			log.Println("Gotify:", err)
			return
		}

		u.Path = path.Join(u.Path, "message")
		queryString := u.Query()
		queryString.Set("token", config.Token)
		u.RawQuery = queryString.Encode()
		s := u.String()

		log.Println("Gotify url:", s)

		_, err = http.PostForm(s,
			url.Values{"message": {message}, "title": {title}, "priority": {priority}})
		if err != nil {
			log.Println("Gotify:", err)
			return
		}
	}
}
