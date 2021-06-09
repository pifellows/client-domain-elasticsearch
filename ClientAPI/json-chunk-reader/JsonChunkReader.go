package json_chunk_reader

import (
	"encoding/json"
	"io"
)

func NewJsonChunkReader(reader io.Reader) (*JsonChunkReader, error) {
	r := &JsonChunkReader{decoder: json.NewDecoder(reader)}

	// prepare by reading the first token
	_, err := r.decoder.Token()
	if err != nil {
		return nil, err
	}
	return r, nil
}

type JsonChunkReader struct {
	decoder *json.Decoder
}

func (jcr *JsonChunkReader) ReadItem(v interface{}) (key interface{}, err error) {

	// when reading the map, we first read the key, and then the body
	key, err = jcr.decoder.Token()

	// if we encounter an error, return it to the caller
	// the caller should check for io.EOF
	if err != nil {
		return nil, err
	}

	err = jcr.decoder.Decode(v)
	if err != nil {
		return nil, err
	}
	return key, nil
}
