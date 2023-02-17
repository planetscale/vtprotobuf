// Commercial secret, LLC "RevTech". Refer to CONFIDENTIAL file in the root for details

package common

import (
	"strings"
)

func ConvertIdent(ident string) string {
	return strings.ReplaceAll(ident, ".", "_")
}
