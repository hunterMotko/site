package server

import "testing"

func TestResolvePort(t *testing.T) {
	// The bug this guards: an unset or non-numeric variable made Atoi return 0,
	// which binds a random free port. The site came up, reported no error, and
	// was unreachable at the address anything else expected.
	tests := []struct {
		name string
		dev  bool
		dp   string
		pp   string
		want int
	}{
		{name: "development reads DP", dev: true, dp: "3000", pp: "80", want: 3000},
		{name: "production reads PP", dev: false, dp: "3000", pp: "8081", want: 8081},
		{name: "unset falls back", dev: false, want: defaultPort},
		{name: "empty falls back", dev: true, dp: "", want: defaultPort},
		{name: "non-numeric falls back", dev: true, dp: "eighty-eighty", want: defaultPort},
		{name: "zero falls back", dev: true, dp: "0", want: defaultPort},
		{name: "negative falls back", dev: true, dp: "-1", want: defaultPort},
		{name: "above the port range falls back", dev: true, dp: "70000", want: defaultPort},
		{name: "whitespace falls back", dev: true, dp: " 8080 ", want: defaultPort},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("DP", tc.dp)
			t.Setenv("PP", tc.pp)
			if got := resolvePort(tc.dev); got != tc.want {
				t.Errorf("resolvePort(dev=%v) with DP=%q PP=%q = %d, want %d",
					tc.dev, tc.dp, tc.pp, got, tc.want)
			}
		})
	}
}

func TestMetaCanonicalURL(t *testing.T) {
	// An absent canonical is correct when SITE_URL is unset; a wrong one
	// actively hands ranking signal to another URL.
	tests := []struct {
		name    string
		siteURL string
		path    string
		want    string
	}{
		{name: "joined when set", siteURL: "https://huntermotko.dev", path: "/work", want: "https://huntermotko.dev/work"},
		{name: "root path", siteURL: "https://huntermotko.dev", path: "/", want: "https://huntermotko.dev/"},
		{name: "empty when unset", siteURL: "", path: "/work", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := Meta{SiteURL: tc.siteURL, Path: tc.path}
			if got := m.CanonicalURL(); got != tc.want {
				t.Errorf("CanonicalURL() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMetaIsCurrent(t *testing.T) {
	// Previously aria-current="page" was hardcoded on the Home link, so every
	// other page announced itself as the home page.
	m := Meta{Path: "/about"}
	if !m.IsCurrent("/about") {
		t.Error("IsCurrent(/about) = false on the /about page, want true")
	}
	if m.IsCurrent("/") {
		t.Error("IsCurrent(/) = true on the /about page, want false")
	}
	if m.IsCurrent("/work") {
		t.Error("IsCurrent(/work) = true on the /about page, want false")
	}
}

func TestSecureEqual(t *testing.T) {
	if !secureEqual("hunter2", "hunter2") {
		t.Error("secureEqual with matching tokens = false, want true")
	}
	if secureEqual("hunter2", "hunter3") {
		t.Error("secureEqual with differing tokens = true, want false")
	}
	if secureEqual("hunter", "hunter2") {
		t.Error("secureEqual with a prefix = true, want false")
	}
	// An unset STATS_TOKEN must never authorise, or an empty query parameter
	// would be a valid credential on any host that forgot to configure one.
	if secureEqual("", "") {
		t.Error("secureEqual with an empty expected token = true, want false")
	}
}
