package argocd

import "regexp"

var repositoryRe = regexp.MustCompile(`(?P<protocol>\w+)@(?P<owner>[\w.]+)/(?P<name>\w+)`)

func (r repository) findItem(key string) string {
	match := repositoryRe.FindStringSubmatch(r.URL)

	for i, name := range match {
		if name == key {
			return match[i]
		}
	}

	return ""
}

func (r repository) Owner() string {
	return r.findItem("owner")
}

func (r repository) Name() string {
	return r.findItem("name")
}
