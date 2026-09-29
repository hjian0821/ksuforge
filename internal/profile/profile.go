package profile

import "strings"

type Profile interface {
	Name() string
	Matches(manufacturer, brand string) bool
	ValidateCodename(codename string) error
}

type Xiaomi struct{}

func (Xiaomi) Name() string { return "xiaomi" }
func (Xiaomi) Matches(manufacturer, brand string) bool {
	v := strings.ToLower(manufacturer + " " + brand)
	return strings.Contains(v, "xiaomi") || strings.Contains(v, "redmi") || strings.Contains(v, "poco")
}
func (Xiaomi) ValidateCodename(string) error { return nil }

func Match(manufacturer, brand string) (Profile, bool) {
	profiles := []Profile{Xiaomi{}}
	for _, p := range profiles {
		if p.Matches(manufacturer, brand) {
			return p, true
		}
	}
	return nil, false
}
