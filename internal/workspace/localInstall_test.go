package workspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotMountedMessages(t *testing.T) {
	assert.True(t, umountNotMounted("umount: /myworkspace: not currently mounted", "/myworkspace"))
	assert.True(t, umountNotMounted("umount: /myworkspace: not mounted", "/myworkspace"))
	assert.True(t, fusermountNotMounted("fusermount: entry for /myworkspace not found", "/myworkspace"))

	assert.False(t, umountNotMounted("umount: /myworkspace: currently mounted", "/myworkspace"))
	assert.False(t, umountNotMounted("umount: /myworkspace: mounted", "/myworkspace"))
	assert.False(t, fusermountNotMounted("fusermount: entry for /myworkspace found", "/myworkspace"))

}
