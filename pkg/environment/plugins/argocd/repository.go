package argocd

import "regexp"

var repositoryRe = regexp.MustCompile(`(?P<protocol>\w+)@(?P<host>[\w.]+):(?P<owner>\w+)/(?P<name>\w+)`)

func (r repository) findItem(key string) string {
	groupNames := repositoryRe.SubexpNames()
	matches := repositoryRe.FindAllStringSubmatch(r.URL, -1)

	for _, match := range matches {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]

			if name == key {
				return group
			}
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
