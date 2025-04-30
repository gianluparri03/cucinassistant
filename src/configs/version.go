package configs

import "fmt"

const VersionCode = 7
const VersionName = "Ciliegia"

var Version string = fmt.Sprintf("%d (%s)", VersionCode, VersionName)
