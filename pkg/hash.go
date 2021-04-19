package pkg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

func itemToBytes(item DBItem) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(item.Title),
			[]byte(item.Author),
			[]byte(item.PubDate),
			[]byte(item.Link),
			[]byte(item.Description),
		},
		[]byte{},
	)

	return data
}

func hashMD5(item DBItem) string {
	hash := md5.Sum(itemToBytes(item))
	return hex.EncodeToString(hash[:])
}
