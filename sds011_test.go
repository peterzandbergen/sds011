package sds011

import "testing"

func TestOpenPort(t *testing.T) {
	s := &Client{}
	s, err := Open(WithPort("usb"))
	// s, err := Open(WithPort("usb"), WithSerialSettings(nil))
	if err != nil {
		t.Errorf("error opening Client on port: %s", err.Error())
	}
	defer s.Close()
}

func TestWorkingModes(t *testing.T) {
	cases := []struct {
		name  string
		value WorkingMode
		real  WorkingMode
	}{
		{
			name:  "Work",
			value: Work,
			real:  0,
		},
		{
			name:  "Sleep",
			value: Sleep,
			real:  1,
		},
	}
	for _, c := range cases {

		if c.value != c.real {
			t.Errorf("%s: expected %d, got %d", c.name, c.value, c.real)
		}
	}
}

func TestReportingModes(t *testing.T) {
	cases := []struct {
		name  string
		value ReportingMode
		real  ReportingMode
	}{
		{
			name:  "ActiveMode",
			value: ActiveMode,
			real:  0,
		},
		{
			name:  "QueryMode",
			value: QueryMode,
			real:  1,
		},
	}
	for _, c := range cases {

		if c.value != c.real {
			t.Errorf("%s: expected %d, got %d", c.name, c.value, c.real)
		}
	}
}
