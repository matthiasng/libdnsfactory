package libdnsfactory

import "github.com/libdns/libdns"

// Provider summarizes all libdns interfaces
type Provider interface {
	libdns.RecordAppender
	libdns.RecordDeleter
	libdns.RecordGetter
	libdns.RecordSetter
}
