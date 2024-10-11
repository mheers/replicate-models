package replicate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	ig := getFakeReplicate(t)

	err := ig.DownloadFromID("vnrptc56j5rg80chjjkv7g1d7c", "/tmp/moamklzbhgbjqarcfpsmys3ys4.png", 0)
	require.NoError(t, err)
}

func getFakeReplicate(t *testing.T) *Replicate {
	t.Helper()
	ig, err := NewReplicate(getReplicateToken(t), "")
	require.NoError(t, err)
	require.NotNil(t, ig)
	return ig
}

func getReplicateToken(t *testing.T) string {
	t.Helper()
	token := os.Getenv("AIS_REPLICATE_API_KEY")
	if token == "" {
		t.Skip("AIS_REPLICATE_API_KEY not set")
	}
	return token
}
