package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	endpointURISchemeDTN uint64 = 1
	endpointURISchemeIPN uint64 = 2
)

type EndpointID struct {
	SchemeName         uint64
	SchemeSpecificPart interface{}
}

func newEndpointIDDTN(ssp string) EndpointID {
	var sspRaw interface{}
	if ssp == "none" {
		sspRaw = uint64(0)
	} else {
		sspRaw = string(ssp)
	}

	return EndpointID{
		SchemeName:         endpointURISchemeDTN,
		SchemeSpecificPart: sspRaw,
	}
}

func newEndpointIDIPN(ssp string) EndpointID {
	re := regexp.MustCompile(`^(\d+)\.(\d+)$`)
	matches := re.FindStringSubmatch(ssp)

	nodeNo, _ := strconv.ParseUint(matches[1], 10, 64)
	serviceNo, _ := strconv.ParseUint(matches[2], 10, 64)

	return EndpointID{
		SchemeName:         endpointURISchemeIPN,
		SchemeSpecificPart: [2]uint64{nodeNo, serviceNo},
	}
}

func NewEndpointID(eid string) EndpointID {
	re := regexp.MustCompile(`^([[:alnum:]]+):(.+)$`)
	matches := re.FindStringSubmatch(eid)

	name := matches[1]
	ssp := matches[2]

	switch name {
	case "dtn":
		return newEndpointIDDTN(ssp)
	case "ipn":
		return newEndpointIDIPN(ssp)
	default:
		return EndpointID{}
	}
}

func (eid EndpointID) String() string {
	var b strings.Builder

	switch eid.SchemeName {
	case endpointURISchemeDTN:
		b.WriteString("dtn")
	case endpointURISchemeIPN:
		b.WriteString("ipn")
	default:
		fmt.Fprintf(&b, "unknown_%d", eid.SchemeName)
	}
	b.WriteRune(':')

	switch t := eid.SchemeSpecificPart.(type) {
	case uint64:
		if eid.SchemeName == endpointURISchemeDTN && eid.SchemeSpecificPart.(uint64) == 0 {
			b.WriteString("none")
		} else {
			fmt.Fprintf(&b, "%d", eid.SchemeSpecificPart.(uint64))
		}

	case string:
		b.WriteString(eid.SchemeSpecificPart.(string))

	case [2]uint64:
		var ssp [2]uint64 = eid.SchemeSpecificPart.([2]uint64)
		if eid.SchemeName == endpointURISchemeIPN {
			fmt.Fprintf(&b, "%d.%d", ssp[0], ssp[1])
		} else {
			fmt.Fprintf(&b, "%v", ssp)
		}

	default:
		fmt.Fprintf(&b, "unknown %T: %v", t, eid.SchemeSpecificPart)
	}

	return b.String()
}
