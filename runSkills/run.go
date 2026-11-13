/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package main

import (
	
	_	"github.com/gerard94/util/babel/static"
	_	"github.com/gerard94/util/json/static"
	_	"github.com/gerard94/util/graphQL/static"
	_	"github.com/gerard94/skills/static"
	
	A	"github.com/gerard94/util/avl"
	G	"github.com/gerard94/util/graphQL"
	GQ	"github.com/gerard94/skills/gqReceiver"
	M	"github.com/gerard94/util/misc"
	
	_	"github.com/gerard94/skills/gqFieldsRes"
	
		"fmt"

)

const (
	
	version = "1.3.0"

)

func versionR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	return G.MakeStringValue(version)
}

func main () {
	GQ.Start()
} //main

func init () {
	fmt.Println("skills Server version", version, "Tools version", M.Version(), "\n")
	GQ.TS().FixFieldResolver("Query", "version", versionR)
}
