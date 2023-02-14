package stringcases_test

import (
	"testing"

	"github.com/alextanhongpin/stringcases"
	"github.com/stretchr/testify/assert"
)

func TestStringCase(t *testing.T) {
	t.Run("from kebab", func(t *testing.T) {
		assert := assert.New(t)

		r := "user-api"

		k := stringcases.ToKebab(r)
		p := stringcases.ToPascal(k)
		s := stringcases.ToSnake(p)
		c := stringcases.ToCamel(s)

		assert.Equal("user-api", k)
		assert.Equal("UserAPI", p)
		assert.Equal("user_api", s)
		assert.Equal("userAPI", c)
	})

	t.Run("from pascal", func(t *testing.T) {
		assert := assert.New(t)

		r := "UserAPI"

		k := stringcases.ToKebab(r)
		p := stringcases.ToPascal(k)
		s := stringcases.ToSnake(p)
		c := stringcases.ToCamel(s)

		assert.Equal("user-api", k)
		assert.Equal("UserAPI", p)
		assert.Equal("user_api", s)
		assert.Equal("userAPI", c)
	})

	t.Run("from snake", func(t *testing.T) {
		assert := assert.New(t)

		r := "user_api"

		k := stringcases.ToKebab(r)
		p := stringcases.ToPascal(k)
		s := stringcases.ToSnake(p)
		c := stringcases.ToCamel(s)

		assert.Equal("user-api", k)
		assert.Equal("UserAPI", p)
		assert.Equal("user_api", s)
		assert.Equal("userAPI", c)
	})

	t.Run("from camel", func(t *testing.T) {
		assert := assert.New(t)

		r := "userAPI"

		k := stringcases.ToKebab(r)
		p := stringcases.ToPascal(k)
		s := stringcases.ToSnake(p)
		c := stringcases.ToCamel(s)

		assert.Equal("user-api", k)
		assert.Equal("UserAPI", p)
		assert.Equal("user_api", s)
		assert.Equal("userAPI", c)
	})
}

func TestStringCaseCommonInitialism(t *testing.T) {
	tests := []struct {
		scenario string
		text     string
		kebab    string
		snake    string
		camel    string
		pascal   string
	}{
		{"single word", "id", "id", "id", "id", "ID"},
		{"suffix", "userId", "user-id", "user_id", "userID", "UserID"},
		{"prefix", "jsonSerializer", "json-serializer", "json_serializer", "jsonSerializer", "JSONSerializer"},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			assert := assert.New(t)

			k := stringcases.ToKebab(test.text)
			p := stringcases.ToPascal(test.text)
			s := stringcases.ToSnake(test.text)
			c := stringcases.ToCamel(test.text)

			assert.Equal(test.kebab, k)
			assert.Equal(test.pascal, p)
			assert.Equal(test.snake, s)
			assert.Equal(test.camel, c)
		})
	}
}
