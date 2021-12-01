package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_ValidForm(t *testing.T) {
	r := httptest.NewRequest("POST", "/any-url", nil)
	form := NewForm(r.PostForm)

	isValid := form.ValidForm()
	if !isValid {
		t.Error("returned invalid when should have been valid")
	}
}

func TestForm_RequiredFields(t *testing.T) {
	r := httptest.NewRequest("POST", "/any-url", nil)
	form := NewForm(r.PostForm)

	form.RequiredFields("a", "b", "c")
	if form.ValidForm() {
		t.Error("form shows as valid when required fiels are missing !")
	}

	// try with some posted data
	postedData := url.Values{}
	postedData.Add("a", "1")
	postedData.Add("b", "2")
	postedData.Add("c", "3")

	r = httptest.NewRequest("POST", "/any-url", nil)
	r.PostForm = postedData
	form = NewForm(r.PostForm)
	form.RequiredFields("a", "b", "c")
	if !form.ValidForm() {
		t.Error("shows that it does not have required fields !")
	}
}

func TestForm_HasField(t *testing.T) {
	r := httptest.NewRequest("POST", "/any-url", nil)
	form := NewForm(r.PostForm)

	// test with field known not to be present
	has := form.HasField("anything")
	if has {
		t.Error("form shows it has field when it does not")
	}

	// try with some posted data
	postedData := url.Values{}
	postedData.Add("x", "1")
	postedData.Add("y", "2")

	form = NewForm(postedData)

	has = form.HasField("y")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/any-url", nil)
	form := NewForm(r.PostForm)

	// initially test for length of a non-existent field
	form.MinLength("z", 8)
	if form.ValidForm() {
		t.Error("form shows minimum length for a non-existent field")
	}

	// test GetErrMsg() from errors.go when there IS an error
	isErr := form.Errors.GetErrMsg("z")
	if isErr == "" {
		t.Error("this should have an error but did not get one")
	}

	// try with some posted data
	postedData := url.Values{}
	postedData.Add("x", "field1")
	postedData.Add("y", "field999")

	form = NewForm(postedData)

	form.MinLength("y", 100)
	if form.ValidForm() {
		t.Error("shows minlength of 100 is met when actual length is 8")
	}

	// try with some posted data
	postedData = url.Values{}
	postedData.Add("x", "field1")
	postedData.Add("y", "field999")

	form = NewForm(postedData)

	form.MinLength("x", 1)
	if !form.ValidForm() {
		t.Error("shows minlength of 1 is NOT met when actual length is 6")
	}

	// test GetErrMsg() from errors.go when there is NO error
	isErr = form.Errors.GetErrMsg("x")
	if isErr != "" {
		t.Error("this should NOT have an error but did get one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := NewForm(postedData)

	form.IsEmail("xyz@123.com")
	if form.ValidForm() {
		t.Error("form shows valid email for non-existent field")
	}

	// try with some posted data
	postedData = url.Values{}
	postedData.Add("x", "field1")
	postedData.Add("email", "hello@world.com")

	form = NewForm(postedData)

	form.IsEmail("email")
	if !form.ValidForm() {
		t.Error("gave invalid email when it is actually valid")
	}

	// try again with some posted data & an invalid email address
	postedData = url.Values{}
	postedData.Add("x", "field1")
	postedData.Add("email", "hello.world.com")

	form = NewForm(postedData)

	form.IsEmail("email")
	if form.ValidForm() {
		t.Error("gave valid when emial address is actually invalid")
	}

}
