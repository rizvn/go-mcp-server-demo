package util_test

import (
	"strings"
	"testing"

	"github.com/rizvn/go-mcp/util"
)

func TestUtil(t *testing.T) {

	t.Run("Test extract tag", func(t *testing.T) {
		text := "Here is some text in side a tag <taga>bla bla bla</taga>"
		extracted := util.ExtractTag("taga", text)

		if extracted != "bla bla bla" {
			t.Errorf("Expected 'bla bla bla', but got '%s'", extracted)
		}
	})

	t.Run("Remove tag", func(t *testing.T) {
		text := "Here is some text in side a tag <taga>bla bla bla</taga>"
		tagRemoved := util.RemoveTag("taga", text)

		if strings.Contains("taga", tagRemoved) {
			t.Errorf("Expected taga removed', but got '%s'", tagRemoved)
		}
	})
}
