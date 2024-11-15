package gcputil

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const (
	resourceIdRegex   = "^[^\t\n\f\r]+$"
	collectionIdRegex = "^[a-z][a-zA-Z]*$"

	fullResourceNameRegex = `^//([a-z]+)\.googleapis\.com/(.+)$`
	selfLinkMarker        = "projects/"
)

var singleCollectionIds = map[string]struct{}{
	"global": {},
}

type RelativeResourceName struct {
	Name                 string
	TypeKey              string
	IdTuples             map[string]string
	OrderedCollectionIds []string
}

func ParseRelativeName(resource string) (*RelativeResourceName, error) {
	resourceRe := regexp.MustCompile(resourceIdRegex)
	collectionRe := regexp.MustCompile(collectionIdRegex)

	tokens := strings.Split(resource, "/")
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid relative resource name %s (too few tokens)", resource)
	}

	ids := map[string]string{}
	typeKey := ""
	currColId := ""
	for idx, v := range tokens {
		if len(currColId) == 0 {
			if _, ok := singleCollectionIds[v]; ok {
				// Ignore 'single' collectionIds like Global, but error if they are the last ID
				if idx == len(tokens)-1 {
					return nil, fmt.Errorf("invalid relative resource name %s (last collection '%s' has no ID)", resource, currColId)
				}
				continue
			}
			if len(collectionRe.FindAllString(v, 1)) == 0 {
				return nil, fmt.Errorf("invalid relative resource name %s (invalid collection ID %s)", resource, v)
			}
			currColId = v
			typeKey += currColId + "/"
		} else {
			if len(resourceRe.FindAllString(v, 1)) == 0 {
				return nil, fmt.Errorf("invalid relative resource name %s (invalid resource sub-ID %s)", resource, v)
			}
			ids[currColId] = v
			currColId = ""
		}
	}

	typeKey = typeKey[:len(typeKey)-1]
	collectionIds := strings.Split(typeKey, "/")
	resourceName := tokens[len(tokens)-2]
	return &RelativeResourceName{
		Name:                 resourceName,
		TypeKey:              typeKey,
		OrderedCollectionIds: collectionIds,
		IdTuples:             ids,
	}, nil
}

type FullResourceName struct {
	Service string
	*RelativeResourceName
}

func ParseFullResourceName(name string) (*FullResourceName, error) {
	fullRe := regexp.MustCompile(fullResourceNameRegex)
	matches := fullRe.FindAllStringSubmatch(name, 1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid full name '%s'", name)
	}

	if len(matches[0]) != 3 {
		return nil, fmt.Errorf("invalid full name '%s'", name)
	}

	serviceName := matches[0][1]
	relName, err := ParseRelativeName(strings.Trim(matches[0][2], "/"))
	if err != nil {
		return nil, fmt.Errorf("error parsing relative resource path in full resource name '%s': %v", name, err)
	}

	return &FullResourceName{
		Service:              serviceName,
		RelativeResourceName: relName,
	}, nil
}

type SelfLink struct {
	Prefix string
	*RelativeResourceName
}

func ParseProjectResourceSelfLink(link string) (*SelfLink, error) {
	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid self link '%s' must have scheme/host", link)
	}

	split := strings.SplitAfterN(link, selfLinkMarker, 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("self link '%s' is not for project-level resource, must contain '%s')", link, selfLinkMarker)
	}

	relName, err := ParseRelativeName(selfLinkMarker + split[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing relative resource path in self-link '%s': %v", link, err)
	}

	return &SelfLink{
		Prefix:               strings.TrimSuffix(split[0], selfLinkMarker),
		RelativeResourceName: relName,
	}, nil
}
