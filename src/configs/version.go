package configs

import "fmt"

const VersionCode = 8
const VersionName = "Banana"

var Version string = fmt.Sprintf("%d (%s)", VersionCode, VersionName)
