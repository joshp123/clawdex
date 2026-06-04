package contactexport

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecodeNormalizesContacts(t *testing.T) {
	got, err := Decode(strings.NewReader(`{"contacts":[{"display_name":" Ada Lovelace ","phone_numbers":[" +1 555 0100 ","","+1 555 0100"]},{"display_name":"","phone_numbers":[]}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Contacts) != 1 {
		t.Fatalf("contacts = %#v", got.Contacts)
	}
	if got.Contacts[0].DisplayName != "Ada Lovelace" {
		t.Fatalf("name = %q", got.Contacts[0].DisplayName)
	}
	if got.Contacts[0].PhoneNumbers[0] != "+1 555 0100" || len(got.Contacts[0].PhoneNumbers) != 1 {
		t.Fatalf("phones = %#v", got.Contacts[0].PhoneNumbers)
	}
}

func TestDecodeRejectsBadContacts(t *testing.T) {
	for _, input := range []string{
		`{`,
		`{"contacts":[{"display_name":"Ada","phone_numbers":[]}]}`,
		`{"contacts":[{"display_name":"","phone_numbers":["123"]}]}`,
		`{"contacts":[{"display_name":"Ada","phone_numbers":["123"],"extra":"x"}]}`,
		`{"contacts":[]}{"contacts":[]}`,
		`{"contacts":[]}
private junk`,
	} {
		t.Run(input, func(t *testing.T) {
			if _, err := Decode(strings.NewReader(input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestEncodeUsesContractFieldNames(t *testing.T) {
	var out bytes.Buffer
	err := Encode(&out, ContactExport{Contacts: []Contact{{DisplayName: "Ada", PhoneNumbers: []string{"123"}}}})
	if err != nil {
		t.Fatal(err)
	}
	text := out.String()
	if !strings.Contains(text, `"display_name": "Ada"`) || !strings.Contains(text, `"phone_numbers":`) {
		t.Fatalf("encoded = %s", text)
	}
}
