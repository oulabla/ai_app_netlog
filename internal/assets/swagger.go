package assets

import "embed"

//go:embed openapi/all-apis.swagger.json
var SwaggerFS embed.FS
