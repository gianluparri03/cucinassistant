package configs

import "fmt"

const VersionCode = 7
const VersionName = "?"

var Version string = fmt.Sprintf("%d (%s)", VersionCode, VersionName)
