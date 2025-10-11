package configs

import "fmt"

const VersionCode = 9
const VersionName = "Maracuja"

var Version string = fmt.Sprintf("%d (%s)", VersionCode, VersionName)
