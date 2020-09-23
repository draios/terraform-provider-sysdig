package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestAccMacro(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	fixedRandomText := rText()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: macroWithName(rText()),
			},
			{
				Config: macroWithName(fixedRandomText),
			},
			{
				Config: macroUpdatedWithName(fixedRandomText),
			},
			{
				Config: macroAppendToDefault(),
			},
			{
				Config: macroWithMacro(rText(), rText()),
			},
			{
				Config: macroWithMacroAndList(rText(), rText(), rText()),
			},
		},
	})
}

func macroWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample" {
  name = "terraform_test_%s"
  condition = "always_true"
}
`, name)
}

func macroUpdatedWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample" {
  name = "terraform_test_%s"
  condition = "never_true"
}
`, name)
}

func macroAppendToDefault() string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample2" {
  name = "container"
  condition = "and always_true"
  append = true
}
`)
}

func macroWithMacro(name1, name2 string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_macro" "sample3" {
  name = "terraform_test_%s"
  condition = "always_true"
}

resource "sysdig_secure_macro" "sample4" {
  name = "terraform_test_%s"
  condition = "never_true and ${sysdig_secure_macro.sample3.name}"
}
`, name1, name2)
}

func macroWithMacroAndList(name1, name2, name3 string) string {
	return fmt.Sprintf(`
%s

resource "sysdig_secure_macro" "sample5" {
  name = "terraform_test_%s"
  condition = "fd.name in (${sysdig_secure_list.sample.name})"
}

resource "sysdig_secure_macro" "sample6" {
  name = "terraform_test_%s"
  condition = "never_true and ${sysdig_secure_macro.sample5.name}"
}
`, listWithName(name3), name1, name2)
}
