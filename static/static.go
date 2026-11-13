/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package static

import (
	
	F	"path/filepath"
	G	"github.com/gerard94/util/graphQL"
		"io"
		"io/ioutil"
		"strings"

)

var (
	
	oldLinkDocs G.LinkGQL

)

func linkStringDocs (path string) io.ReadCloser {
	switch F.Base(path) {
	case "typeSystem.txt":
		return ioutil.NopCloser(strings.NewReader(typeSystem))
	default:
		return oldLinkDocs(path)
	}
}

func init () {
	oldLinkDocs = G.SetLGQL(linkStringDocs)
}
