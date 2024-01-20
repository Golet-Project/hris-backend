package primitive

import "regexp"

var TenantDomainPattern, _ = regexp.Compile(`^[a-zA-Z0-9]+$`)
