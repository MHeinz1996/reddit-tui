package config

const defaultConfiguration = `
#
# Default configuration for reddittui.
# Uncomment to configure
#

#[core]
#bypassCache = false
#logLevel = "Warn"

#[filter]
#keywords = ["drama"]
#subreddits = ["news", "politics"]

#[client]
#timeoutSeconds = 10
#cacheTtlSeconds = 3600

#[server]
#domain = "old.reddit.com"
#type = "old"

#
# old.reddit.com redirects logged-out requests to a login page, so a session
# is required. Copy the "reddit_session" cookie value from your browser:
# Chrome DevTools -> Application -> Cookies -> https://old.reddit.com
#
# This value grants full access to your account. Keep this file private;
# reddittui sets it to mode 0600 when a cookie is configured. Replace the
# value and restart reddittui when the session expires. The
# REDDITTUI_SESSION_COOKIE environment variable overrides this setting.
#
#[auth]
#sessionCookie = ""
`
