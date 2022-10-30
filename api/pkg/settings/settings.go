package settings

import "os"

var (
	Development = os.Getenv("MODE") == "development"
)
