package passcracker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBruteForce_Crack(t *testing.T) {
	testCases := []struct {
		hash     string
		password string
	}{
		{
			hash:     "8254c329a92850f6d539dd376f4816ee2764517da5e0235514af433164480d7a",
			password: "k",
		},
		{
			hash:     "21e721c35a5823fdb452fa2f9f0a612c74fb952e06927489c6b27a43b817bed4",
			password: "cd",
		},
		{
			hash:     "9834876dcfb05cb167a5c24953eba58c4ac89b1adf57f28f2f9d09af107ee8f0",
			password: "aaa",
		},
		{
			hash:     "77af778b51abd4a3c51c5ddd97204a9c3ae614ebccb75a606c3b6865aed6744e",
			password: "cat",
		},
		{
			hash:     "17f165d5a5ba695f27c023a83aa2b3463e23810e360b7517127e90161eebabda",
			password: "zzz",
		},
		{
			hash:     "cc967443070ab409a57a455dc8c2405a10b6ba96ca75f87bcf18ca368bb090a8",
			password: "lena",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			bfCracker := NewBruteForce("abcdefghijklmnopqrstuvwxyz", 1, 1*time.Second)
			result := bfCracker.Crack(tc.hash)
			assert.Equal(t, tc.password, result.Password)
		})
	}
}

func TestBruteForce_updateBaseWord(t *testing.T) {
	testCases := []struct {
		baseWord string
		expected string
	}{
		{
			baseWord: "",
			expected: "a",
		},
		{
			baseWord: "a",
			expected: "b",
		},
		{
			baseWord: "z",
			expected: "aa",
		},
		{
			baseWord: "aa",
			expected: "ab",
		},
		{
			baseWord: "abc",
			expected: "abd",
		},
		{
			baseWord: "zzz",
			expected: "aaaa",
		},
		{
			baseWord: "azza",
			expected: "azzb",
		},
		{
			baseWord: "azzz",
			expected: "baaa",
		},
		{
			baseWord: "eezz",
			expected: "efaa",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.baseWord, " should return ", tc.expected), func(t *testing.T) {
			w := worker{searchSpace: "abcdefghijklmnopqrstuvwxyz", workerTotal: 1, searchDuration: time.Second}
			actual := w.updateBaseWord(tc.baseWord, 0, 26)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
