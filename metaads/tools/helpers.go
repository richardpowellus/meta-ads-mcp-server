package tools

import "net/url"

// setOptional sets a url.Values key only if value is non-empty.
func setOptional(v url.Values, key, value string) {
	if value != "" {
		v.Set(key, value)
	}
}

// insightParams builds url.Values for Meta Ads insights endpoints.
func insightParams(fields, datePreset, timeRange, breakdowns, level, filtering, sort string) url.Values {
	vals := url.Values{}
	if fields != "" {
		vals.Set("fields", fields)
	}
	if datePreset != "" {
		vals.Set("date_preset", datePreset)
	}
	if timeRange != "" {
		vals.Set("time_range", timeRange)
	}
	if breakdowns != "" {
		vals.Set("breakdowns", breakdowns)
	}
	if level != "" {
		vals.Set("level", level)
	}
	if filtering != "" {
		vals.Set("filtering", filtering)
	}
	if sort != "" {
		vals.Set("sort", sort)
	}
	return vals
}
