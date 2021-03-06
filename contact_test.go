package main

import (
	"testing"

	"github.com/joshsoftware/curem/config"
	"labix.org/v2/mgo/bson"
)

func TestNewContact(t *testing.T) {
	fakeContact, err := NewContact(
		"Encom Inc.",
		"Flynn",
		"flynn@encom.com",
		"",
		"",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	var fetchedContact contact
	err = config.ContactsCollection.Find(bson.M{}).One(&fetchedContact)
	if err != nil {
		t.Errorf("%s", err)
	}

	// fakeContact is a pointer, because NewContact returns a pointer to a struct of contact type.
	// That's why we check fetchedContact with *fakeContact.

	if fetchedContact != *fakeContact {
		t.Errorf("inserted contact is not the fetched contact")
	}
	dropContactsCollection(t)
}

func TestValidateNewContact(t *testing.T) {
	_, err := NewContact(
		"Encom Ic.",
		"",
		"flynn@encom.com",
		"",
		"",
		"USA",
	)
	if err == nil {
		t.Errorf("%s", "error shouldn't be nil when person is empty")
	}
	_, err = NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"",
		"",
		"",
		"USA",
	)
	if err == nil {
		t.Errorf("%s", "error shouldn't be nil when email is empty")
	}
	_, err = NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"x@.xyzc.com",
		"",
		"",
		"USA",
	)
	if err == nil {
		t.Errorf("%s", "error shouldn't be nil when email is invalid")
	}
}

func TestGetContactByID(t *testing.T) {
	fakeContact, err := NewContact(
		"Encom Inc.",
		"Flynn",
		"flynn@encom.com",
		"",
		"",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	id := fakeContact.ID
	fetchedContact, err := GetContactByID(id)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *fakeContact != *fetchedContact {
		t.Errorf("Expected %+v, but got %+v\n", *fakeContact, *fetchedContact)
	}
	dropContactsCollection(t)
}

func TestGetNonExistingContactByID(t *testing.T) {
	_, err := GetContactByID(bson.ObjectIdHex("53b112bde3bdea2642000002"))
	if err == nil {
		t.Errorf("%s", "error shouldn't be nil when we try a fetch a non existent contact")
	}
}

func TestGetContactBySlug(t *testing.T) {
	c, err := NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"samflynn@encom.com",
		"103-345-456",
		"sam_flynn",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	f, err := GetContactBySlug(c.Slug)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *c != *f {
		t.Errorf("expected %+v, but got %+v", *c, *f)
	}
	dropContactsCollection(t)
}

func TestGetNonExistingContactBySlug(t *testing.T) {
	_, err := GetContactBySlug("nlvnjrelvenliqas")
	if err == nil {
		t.Errorf("%s", "error shouldn't be nil when we try a fetch a non existent contact")
	}
}

func TestGetAllContacts(t *testing.T) {
	_, err := NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"samflynn@encom.com",
		"103-345-456",
		"sam_flynn",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = NewContact(
		"Encom Inc.",
		"Kevin Flynn",
		"kevinflynn@encom.com",
		"234-877-988",
		"kevin_flynn",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}

	fetchedContacts, err := GetAllContacts()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(fetchedContacts) != 2 {
		t.Errorf("expected 2 contacts, but got %d", len(fetchedContacts))
	}
	dropContactsCollection(t)
}

func TestUpdateContact(t *testing.T) {
	fakeContact, err := NewContact(
		"Encom Inc.",
		"Flynn",
		"flynn@encom.com",
		"",
		"",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	fakeContact.Country = "India"
	fakeContact.Update()
	fetchedContact, err := GetContactByID(fakeContact.ID)
	if err != nil {
		t.Errorf("%s", err)
	}
	if fetchedContact.Country != "India" {
		t.Errorf("%s", "contact not updated")
	}
	dropContactsCollection(t)
}

func TestUpdateContactValidationCheck(t *testing.T) {
	fakeContact, err := NewContact(
		"Encom Inc.",
		"Flynn",
		"flynn@encom.com",
		"",
		"",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	fakeContact.Email = "India"
	err = fakeContact.Update()
	if err == nil {
		t.Errorf("error shouldn't be nil when updated with an invalid email address")
	}
	dropContactsCollection(t)
}

func TestDelete(t *testing.T) {
	fakeContact, err := NewContact(
		"Encom Inc.",
		"Flynn",
		"flynn@encom.com",
		"",
		"",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = fakeContact.Delete()
	if err != nil {
		t.Errorf("%s", err)
	}
	n, err := config.ContactsCollection.Count()
	if err != nil {
		t.Errorf("%s", err)
	}
	if n != 0 {
		t.Errorf("expected 0 documents in the collection, but found %d", n)
	}
	dropContactsCollection(t)
}

func TestSlugifyContact(t *testing.T) {
	c, err := NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"samflynn@encom.com",
		"103-345-456",
		"sam_flynn",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	if c.Slug != "sam-flynn" {
		t.Errorf("expected slug to be %s, but got %s", "sam-flynn", c.Slug)
	}
	d := &contact{
		Person: "Sam Flynn",
		Email:  "sam@example.com",
	}
	slugifyContact(d)
	if d.Slug == "sam-flynn" {
		t.Errorf("expected something other than %s as slug", "sam-flynn")
	}
	dropContactsCollection(t)
}

func TestContactSlugExists(t *testing.T) {
	c, err := NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"samflynn@encom.com",
		"103-345-456",
		"sam_flynn",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	if !contactSlugExists(c.Slug) {
		t.Errorf("%s", "expected contactSlugExists to return true but returns false")
	}
	dropContactsCollection(t)
}

func TestGetLeadsOfContact(t *testing.T) {
	c, err := NewContact(
		"Encom Inc.",
		"Sam Flynn",
		"samflynn@encom.com",
		"103-345-456",
		"sam_flynn",
		"USA",
	)
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = NewLead(c.Slug, "Web", "Gautam", "Won", 3, 5, 2, "3rd July, 2014", nil)
	_, err = NewLead(c.Slug, "Referral", "Sethu", "Warming Up", 3, 5, 2, "3rd July, 2014", nil)
	x, err := c.Leads()
	if len(x) != 2 {
		t.Errorf("expecting 2 leads but got %d", len(x))
	}
	dropCollections(t)
}
