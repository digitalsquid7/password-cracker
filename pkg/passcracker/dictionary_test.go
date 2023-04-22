package passcracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDictionary_Crack(t *testing.T) {
	testCases := []struct {
		hash       string
		password   string
		dictionary []string
	}{
		{
			hash:       "708e4b2a324e91922e63e65f519c6d206b0e9053654c853c65b45d09faa88368",
			password:   "beavis",
			dictionary: []string{"beavis"},
		},
		{
			hash:       "d6e21286621a8586f7e54720ab5a39c93acc2f4b8fba7a16ec1c24d69a08c613",
			password:   "house",
			dictionary: []string{"beavis", "teller", "house"},
		},
		{
			hash:       "aee408847d35e44e99430f0979c3357b85fe8dbb4535a494301198adbee85f27",
			password:   "success",
			dictionary: []string{"beavis", "teller", "success", "house"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			bfCracker := NewDictionary(tc.dictionary, 1)
			result := bfCracker.Crack(tc.hash)
			assert.Equal(t, tc.password, result.Password)
		})
	}
}
