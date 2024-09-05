package resources

type OpenSource struct {
  Alt string
  Open string
  Href string
}

type Resources struct {
  Alts []OpenSource
}

func GetResources() Resources {
  return Resources{
    Alts: []OpenSource{
      {Alt: "Microsoft Office", Open: "LibreOffice", Href: "https://github.com/LibreOffice"},
      {Alt: "Slack", Open: "Mattermost", Href: "https://github.com/mattermost/mattermost"},
      {Alt: "Airtable", Open: "Nocodb", Href: "https://github.com/nocodb/nocodb"},
      {Alt: "Jira", Open: "Plane", Href: "https://github.com/makeplane/plane"},
      {Alt: "Notion", Open: "Appflowy", Href: "https://github.com/AppFlowy-IO/AppFlowy"},
      {Alt: "Zoom", Open: "Jitsi", Href: "https://github.com/jitsi"},
      {Alt: "Salesforce", Open: "ERPNext", Href: "https://github.com/frappe/erpnext"},
      {Alt: "Vercel", Open: "Coolify", Href: "https://github.com/coollabsio/coolify"},
      {Alt: "Heroku", Open: "Dokku", Href: "https://github.com/dokku/dokku"},
      {Alt: "Firebase", Open: "Instant", Href: "https://github.com/instantdb/instant"},
      {Alt: "Datadog", Open: "Prometheus", Href: "https://github.com/prometheus/prometheus"},
      {Alt: "Dropbox", Open: "NextCloud", Href: "https://github.com/nextcloud "},
      {Alt: "Mailchimp", Open: "Mautic", Href: "https://github.com/mautic/mautic"},
      {Alt: "Trello", Open: "Wekan", Href: "https://github.com/wekan/wekan "},
      {Alt: "Adobe Primiere", Open: "Davinci", Href: "Resolve https://www.blackmagicdesign.com/products/davinciresolve"},
      {Alt: "Adobe Illistrator", Open: "Krita", Href: "https://krita.org/en/"},
      {Alt: "Adobe Affter Effects", Open: "Blender", Href: "https://www.blender.org/"},
    },
  }
}
