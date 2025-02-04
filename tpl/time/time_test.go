// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package time

import (
	"testing"
	"time"

	translators "github.com/bep/gotranslators"
)

func TestTimeLocation(t *testing.T) {
	t.Parallel()

	loc, _ := time.LoadLocation("America/Antigua")
	ns := New(translators.GetTranslator("en"), loc)

	for i, test := range []struct {
		value    string
		location interface{}
		expect   interface{}
	}{
		{"2020-10-20", "", "2020-10-20 00:00:00 +0000 UTC"},
		{"2020-10-20", nil, "2020-10-20 00:00:00 -0400 AST"},
		{"2020-10-20", "America/New_York", "2020-10-20 00:00:00 -0400 EDT"},
		{"2020-01-20", "America/New_York", "2020-01-20 00:00:00 -0500 EST"},
		{"2020-10-20 20:33:59", "", "2020-10-20 20:33:59 +0000 UTC"},
		{"2020-10-20 20:33:59", "America/New_York", "2020-10-20 20:33:59 -0400 EDT"},
		// The following have an explicit offset specified. In this case, it overrides timezone
		{"2020-09-23T20:33:44-0700", "", "2020-09-23 20:33:44 -0700 -0700"},
		{"2020-09-23T20:33:44-0700", "America/New_York", "2020-09-23 20:33:44 -0700 -0700"},
		{"2020-01-20", "invalid-timezone", false}, // unknown time zone invalid-timezone
		{"invalid-value", "", false},
	} {
		var args []interface{}
		if test.location != nil {
			args = append(args, test.location)
		}
		result, err := ns.AsTime(test.value, args...)
		if b, ok := test.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] AsTime didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] AsTime failed: %s", i, err)
				continue
			}
			if result.(time.Time).String() != test.expect {
				t.Errorf("[%d] AsTime got %v but expected %v", i, result, test.expect)
			}
		}
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	ns := New(translators.GetTranslator("en"), time.UTC)

	for i, test := range []struct {
		layout string
		value  interface{}
		expect interface{}
	}{
		{"Monday, Jan 2, 2006", "2015-01-21", "Wednesday, Jan 21, 2015"},
		{"Monday, Jan 2, 2006", time.Date(2015, time.January, 21, 0, 0, 0, 0, time.UTC), "Wednesday, Jan 21, 2015"},
		{"This isn't a date layout string", "2015-01-21", "This isn't a date layout string"},
		// The following test case gives either "Tuesday, Jan 20, 2015" or "Monday, Jan 19, 2015" depending on the local time zone
		{"Monday, Jan 2, 2006", 1421733600, time.Unix(1421733600, 0).Format("Monday, Jan 2, 2006")},
		{"Monday, Jan 2, 2006", 1421733600.123, false},
		{time.RFC3339, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "2016-03-03T04:05:00Z"},
		{time.RFC1123, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "Thu, 03 Mar 2016 04:05:00 UTC"},
		{time.RFC3339, "Thu, 03 Mar 2016 04:05:00 UTC", "2016-03-03T04:05:00Z"},
		{time.RFC1123, "2016-03-03T04:05:00Z", "Thu, 03 Mar 2016 04:05:00 UTC"},
		// Custom layouts, as introduced in Hugo 0.87.
		{":date_medium", "2015-01-21", "Jan 21, 2015"},
	} {
		result, err := ns.Format(test.layout, test.value)
		if b, ok := test.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] DateFormat didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] DateFormat failed: %s", i, err)
				continue
			}
			if result != test.expect {
				t.Errorf("[%d] DateFormat got %v but expected %v", i, result, test.expect)
			}
		}
	}
}

func TestDuration(t *testing.T) {
	t.Parallel()

	ns := New(translators.GetTranslator("en"), time.UTC)

	for i, test := range []struct {
		unit   interface{}
		num    interface{}
		expect interface{}
	}{
		{"nanosecond", 10, 10 * time.Nanosecond},
		{"ns", 10, 10 * time.Nanosecond},
		{"microsecond", 20, 20 * time.Microsecond},
		{"us", 20, 20 * time.Microsecond},
		{"µs", 20, 20 * time.Microsecond},
		{"millisecond", 20, 20 * time.Millisecond},
		{"ms", 20, 20 * time.Millisecond},
		{"second", 30, 30 * time.Second},
		{"s", 30, 30 * time.Second},
		{"minute", 20, 20 * time.Minute},
		{"m", 20, 20 * time.Minute},
		{"hour", 20, 20 * time.Hour},
		{"h", 20, 20 * time.Hour},
		{"hours", 20, false},
		{"hour", "30", 30 * time.Hour},
	} {
		result, err := ns.Duration(test.unit, test.num)
		if b, ok := test.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Duration didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] Duration failed: %s", i, err)
				continue
			}
			if result != test.expect {
				t.Errorf("[%d] Duration got %v but expected %v", i, result, test.expect)
			}
		}
	}
}
