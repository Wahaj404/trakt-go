package core

import "strings"

// ExtendedLevel represents the ?extended= query param value.
type ExtendedLevel string

const (
	ExtendedFull       ExtendedLevel = "full"
	ExtendedMetadata   ExtendedLevel = "metadata"
	ExtendedEpisodes   ExtendedLevel = "episodes"
	ExtendedNoSeasons  ExtendedLevel = "noseasons"
	ExtendedGuestStars ExtendedLevel = "guest_stars"
	ExtendedVIP        ExtendedLevel = "vip"
)

// ListOptions captures shared pagination/extended params for list endpoints.
// All fields are optional; zero value means "don't send".
type ListOptions struct {
	Page     int             // 1-indexed; 0 = omit
	Limit    int             // 0 = omit (Trakt default is 10)
	Extended []ExtendedLevel // empty = omit
}

// toQuery converts ListOptions to a map[string]any suitable for Client.Get.
// Returns nil when the receiver is nil OR when no fields are set, so endpoint
// methods can do `client.GetInto(ctx, path, opts.toQuery(), &out)` without a
// nil check.
func (o *ListOptions) toQuery() map[string]any {
	if o == nil {
		return nil
	}
	q := map[string]any{}
	if o.Page != 0 {
		q["page"] = o.Page
	}
	if o.Limit != 0 {
		q["limit"] = o.Limit
	}
	if len(o.Extended) > 0 {
		parts := make([]string, len(o.Extended))
		for i, e := range o.Extended {
			parts[i] = string(e)
		}
		q["extended"] = strings.Join(parts, ",")
	}
	if len(q) == 0 {
		return nil
	}
	return q
}
