package util_test

import (
	"strings"
	"testing"

	"github.com/rizvn/dbchat/dbchat_graph"
	"github.com/rizvn/dbchat/util"
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

	t.Run("Test GetMethodName", func(t *testing.T) {
		g := &dbchat_graph.DbChatGraph{}
		methodName := util.GetMethodName(g.StepAddRelevantExamples)
		if !strings.HasSuffix(methodName, "StepAddRelevantExamples-fm") {
			t.Errorf("Expected method name to end with 'StepAddRelevantExamples-fm', but got '%s'", methodName)
		}
	})

}
