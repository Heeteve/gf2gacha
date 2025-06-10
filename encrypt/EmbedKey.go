package encrypt

import _ "embed"

//go:embed capture.log.key
var AesKey []byte
