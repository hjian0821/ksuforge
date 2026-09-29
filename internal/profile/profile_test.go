package profile

import "testing"

func TestXiaomiBrands(t *testing.T) {
	for _, brand := range []string{"Xiaomi", "Redmi", "POCO"} {
		if _, ok := Match("", brand); !ok {
			t.Errorf("%s should match", brand)
		}
	}
	if _, ok := Match("Google", "Pixel"); ok {
		t.Error("Pixel should not match")
	}
}
