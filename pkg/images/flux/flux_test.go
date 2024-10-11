package flux

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFlux(t *testing.T) {
	f, err := NewFlux(getReplicateToken(t))
	require.NoError(t, err)
	require.NotNil(t, f)

	tmpDir := t.TempDir()
	fileName := "test.webp"
	dstPath := path.Join(tmpDir, fileName)

	err = f.Create("Text 'Villa JoMaNeKy' written in the sand on a beach with palms and a sunset in the background", dstPath, nil)
	require.NoError(t, err)
}

func getReplicateToken(t *testing.T) string {
	t.Helper()
	token := os.Getenv("AIS_REPLICATE_API_KEY")
	if token == "" {
		t.Skip("AIS_REPLICATE_API_KEY not set")
	}
	return token
}
