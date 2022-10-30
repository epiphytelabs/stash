package store

import "strings"

func (s *Store) TokenList(hash string) (map[string]int, error) {
	if err := hashValidate(hash); err != nil {
		return nil, err
	}

	if err := s.BlobExists(hash); err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT token, count FROM tokens WHERE hash = ?", hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := map[string]int{}

	for rows.Next() {
		var token string
		var count int

		if err := rows.Scan(&token, &count); err != nil {
			return nil, err
		}

		tokens[token] = count
	}

	return tokens, nil
}

func tokenize(s string) map[string]int {
	ts := map[string]int{}

	for _, f := range strings.Fields(s) {
		ts[f]++
	}

	return ts
}
