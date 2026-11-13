/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package gqFieldsRes

import (
	
	A	"github.com/gerard94/util/avl"
	B	"github.com/gerard94/skills/base"
	G	"github.com/gerard94/util/graphQL"
	GQ	"github.com/gerard94/skills/gqReceiver"
	M	"github.com/gerard94/util/misc"
	
	/*
	"fmt"
	*/

)

var (
	
	lg = M.GetLog()

)

func error (ts G.TypeSystem, err M.Nerror) {
	if err != nil {
		ts.Error(err.Error(), nil, nil)
		lg.Printf("Error number %v\n", err.ErrorNum())
	}
} //error

func getStringParam (argumentValues *A.Tree, name string) string {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v)
	if ok {
		switch s := v.(type) {
		case *G.StringValue:
			return s.String.S
		default:
			M.Halt(s, 100)
		}
	}
	return ""
} //getStringParam


func getBooleanParam (argumentValues *A.Tree, name string) bool {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v)
	if ok {
		switch b := v.(type) {
		case *G.BooleanValue:
			return b.Boolean
		default:
			M.Halt(b, 100)
		}
	}
	return false
} //getBooleanParam

func getLocation1ShortListParam (argumentValues *A.Tree, name string) B.Location1ShortListT {
	var v G.Value
	if !G.GetValue(argumentValues, name, &v) {
		return nil
	}
	switch l := v.(type) {
	case *G.ListValue:
		locs := make(B.Location1ShortListT, l.Len())
		i := 0
		for e := l.First(); e != nil; e = l.Next(e) {
			o := e.Value.(*G.InputObjectValue)
			var val G.Value
			b := G.GetObjectValueInputField(o, "country", &val); M.Assert(b, 100)
			co := val.(*G.StringValue).String
			if co.S == "" {
				return nil
			}
			b = G.GetObjectValueInputField(o, "region", &val)
			rg := val.(*G.StringValue).String
			if rg.S == "" {
				return nil
			}
			locs[i] = &B.Location1ShortT{Country: co.S, Region: rg.S}
			i++
		}
		return locs
	default:
		M.Halt(l, 100)
		return nil
	}
} //getLocation1ShortListParam

func getService1ListParam (argumentValues *A.Tree, name string) B.Service1ListT {
	var v G.Value
	if !G.GetValue(argumentValues, name, &v) {
		return nil
	}
	switch l := v.(type) {
	case *G.ListValue:
		sers := make(B.Service1ListT, l.Len())
		i := 0
		for e := l.First(); e != nil; e = l.Next(e) {
			o := e.Value.(*G.InputObjectValue)
			var val G.Value
			b := G.GetObjectValueInputField(o, "category", &val); M.Assert(b, 100)
			ca := val.(*G.StringValue).String
			if ca.S == "" {
				return nil
			}
			b = G.GetObjectValueInputField(o, "speciality", &val)
			sp := val.(*G.StringValue).String
			if sp.S == "" {
				return nil
			}
			sers[i] = &B.Service1T{Category: ca.S, Speciality: sp.S}
			i++
		}
		return sers
	default:
		M.Halt(l, 100)
		return nil
	}
} //getService1ListParam

func getServiceQParam (argumentValues *A.Tree, name string) (services B.ServiceListT, caVoid bool) {
	services, caVoid = nil, false
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v)
	if !ok {
		return
	}
	switch l := v.(type) {
	case *G.ListValue:
		s := make(B.ServiceListT, l.Len())
		i := 0
		for e := l.First(); e != nil; e = l.Next(e) {
			o := e.Value.(*G.InputObjectValue)
			var (val G.Value; sps []string)
			b := G.GetObjectValueInputField(o, "category", &val); M.Assert(b, 100)
			ca := val.(*G.StringValue).String
			if ca.S == "" {
				caVoid = true
				return
			}
			b = G.GetObjectValueInputField(o, "specialities", &val)
			if !b {
				sps = nil
			} else {
				lsp := val.(*G.ListValue)
				sps = make([]string, lsp.Len())
				j := 0
				for esp := lsp.First(); esp != nil; esp = lsp.Next(esp) {
					sps[j] = esp.Value.(*G.StringValue).String.S
					j++
				}
			}
			s[i] = &B.ServiceT{Category: ca.S, Specialities: sps}
			i++
		}
		services = s
		return
	default:
		M.Halt(l, 100)
		return
	}
} //getServiceQParam

func getSOrSParam (argumentValues *A.Tree, name string) int8 {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v)
	if !ok {
		return B.BothServices
	}
	switch e := v.(type) {
	case *G.EnumValue:
		switch e.Enum.S {
		case "SUPPLIED_SERVICES":
			return B.SuppliedServices
		case "SEARCHED_SERVICES":
			return B.SearchedServices
		case "BOTH_SERVICES":
			return B.BothServices
		default:
			M.Halt(101)
			return 0
		}
	default:
		M.Halt(e, 102)
		return 0
	}
} //getSOrSParam

func getLocParam (argumentValues *A.Tree, name string) int {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v); M.Assert(ok, 100)
	switch e := v.(type) {
	case *G.EnumValue:
		switch e.Enum.S {
		case "COUNTRY":
			return B.Country
		case "REGION":
			return B.Region
		case "CITY":
			return B.City
		default:
			M.Halt(e.Enum.S, 101)
			return 0
		}
	default:
		M.Halt(e, 102)
		return 0
	}
} //getLocParam

func getSerParam (argumentValues *A.Tree, name string) int {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v); M.Assert(ok, 100)
	switch e := v.(type) {
	case *G.EnumValue:
		switch e.Enum.S {
		case "CATEGORY":
			return B.Category
		case "SPECIALITY":
			return B.Speciality
		default:
			M.Halt(e.Enum.S, 101)
			return 0
		}
	default:
		M.Halt(e, 102)
		return 0
	}
} //getSerParam

func idGQLToGo (ts G.TypeSystem, idg *G.InputObjectValue) *B.IdentityT {
	
	getString := func (o *G.InputObjectValue, name string, needed bool) string {
		var v G.Value
		b := G.GetObjectValueInputField(o, name, &v); M.Assert(!needed || b, 100)
		if !b {
			return ""
		}
		switch sg := v.(type) {
		case *G.StringValue:
			return sg.String.S
		default:
			M.Halt(sg, 100)
			return ""
		}
	} //getString
	
	getBool := func (o *G.InputObjectValue, name string, needed bool) bool {
		var v G.Value
		b := G.GetObjectValueInputField(o, name, &v); M.Assert(!needed || b, 100)
		if !b {
			return false
		}
		switch bg:= v.(type) {
		case *G.BooleanValue:
			return bg.Boolean
		default:
			M.Halt(bg, 100)
			return false
		}
	} //getBool
	
	getList := func (o *G.InputObjectValue, name string, needed bool) *G.ListValue {
		var v G.Value
		b := G.GetObjectValueInputField(o, name, &v); M.Assert(!needed || b, 100)
		if !b {
			return nil
		}
		switch lg:= v.(type) {
		case *G.ListValue:
			return lg
		default:
			M.Halt(lg, 100)
			return nil
		}
	} //getList
	
	getLocations := func (o *G.InputObjectValue, name string) B.LocationListT{
		lg := getList(o, name, true)
		l := make(B.LocationListT, lg.Len())
		i := 0
		for e := lg.First(); e != nil; e = lg.Next(e) {
			co := e.Value.(*G.InputObjectValue)
			l[i] = new(B.LocationT)
			l[i].Country = getString(co, "country", true)
			if l[i].Country == "" {
				error(ts, B.GetError(4, "insertIdentity"))
				return nil
			}
			rgs := getList(co, "regions", true)
			l[i].Regions = make([]*B.RegionT, rgs.Len())
			j := 0
			for e := rgs.First(); e != nil; e = rgs.Next(e) {
				rg := e.Value.(*G.InputObjectValue)
				l[i].Regions[j] = new(B.RegionT)
				l[i].Regions[j].Region = getString(rg, "region", true)
				if l[i].Regions[j].Region == "" {
					error(ts, B.GetError(5, "insertIdentity"))
					return nil
				}
				cys := getList(rg, "cities", true)
				l[i].Regions[j].Cities = make([]string, cys.Len())
				k := 0
				for e := cys.First(); e != nil; e = cys.Next(e) {
					cy := e.Value.(*G.StringValue)
					l[i].Regions[j].Cities[k] = cy.String.S
					if l[i].Regions[j].Cities[k] == "" {
						error(ts, B.GetError(18, "insertIdentity"))
						return nil
					}
					k++
				}
				j++
			}
			i++
		}
		return l
	} //getLocations
	
	getServices := func (o *G.InputObjectValue, name string) B.ServiceListT {
		lg := getList(o, name, true)
		l := make(B.ServiceListT, lg.Len())
		i := 0
		for e := lg.First(); e != nil; e = lg.Next(e) {
			ca := e.Value.(*G.InputObjectValue)
			l[i] = new(B.ServiceT)
			l[i].Category = getString(ca, "category", true)
			if l[i].Category == "" {
				error(ts, B.GetError(6, "insertIdentity"))
				return nil
			}
			sps := getList(ca, "specialities", true)
			l[i].Specialities = make([]string, sps.Len())
			j := 0
			for e := sps.First(); e != nil; e = sps.Next(e) {
				sp := e.Value.(*G.StringValue)
				l[i].Specialities[j] = sp.String.S
				if l[i].Specialities[j] == "" {
					error(ts, B.GetError(7, "insertIdentity"))
					return nil
				}
				j++
			}
			i++
		}
		return l
	} //getServices
	
	//idGQLToGo
	id := new(B.IdentityT)
	
	id.Nickname = getString(idg, "nickname", true)
	if id.Nickname == "" {
		error(ts, B.GetError(19, "insertIdentity"))
		return nil
	}

	id.Password = getString(idg, "password", true)
	if id.Password == "" {
		error(ts, B.GetError(20, "insertIdentity"))
		return nil
	}
	
	id.G1_nickname = getString(idg, "g1_nickname", false)
	
	p := getString(idg, "pubkey", false)
	id.Pubkey = B.PubkeyT(p)
	
	id.Email = getString(idg, "email", false)
	
	id.Locations = getLocations(idg, "locations")
	if id.Locations == nil {
		return nil
	}
	
	id.Transport = getBool(idg, "transport", true)
	
	id.Pickup = getBool(idg, "pickup", true)
	
	id.Supplied_services = getServices(idg, "supplied_services")
	if id.Supplied_services == nil {
		return nil
	}
	
	id.Searched_services = getServices(idg, "searched_services")
	if id.Searched_services == nil {
		return nil
	}
	
	id.Additional_informations = getString(idg, "additional_informations", false)
	
	return id
} //idGQLToGo

func getIdParam (ts G.TypeSystem, argumentValues *A.Tree, name string) *B.IdentityT {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v); M.Assert(ok, 100)
	switch idg := v.(type) {
	case *G.InputObjectValue:
		return idGQLToGo(ts, idg)
	default:
		M.Halt(idg, 100)
		return nil
	}
} //getIdParam

func getIdsParam (ts G.TypeSystem, argumentValues *A.Tree, name string) B.IdentityListT {
	var v G.Value
	ok := G.GetValue(argumentValues, name, &v)
	if !ok {
		return nil
	}
	switch l := v.(type) {
	case *G.ListValue:
		ids := make(B.IdentityListT, l.Len())
		i := 0
		for e := l.First(); e != nil; e = l.Next(e) {
			switch idg := e.Value.(type) {
			case *G.InputObjectValue:
				ids[i] = idGQLToGo(ts, idg)
			default:
				M.Halt(idg, 100)
			}
			i++
		}
		return ids
	default:
		M.Halt(l, 101)
		return nil
	}
} //getIdsParam

func allLocationsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	country := getStringParam(argumentValues, "country")
	l , err := B.AllLocations(country)
	M.Assert(err == nil || err.ErrorNum() == 0, err, 100)
	error(ts, err)
	if l == nil {
		return G.MakeNullValue()
	}
	gco := G.NewListValue()
	for _, co := range l {
		gco.Append(G.Wrap(co.Country, co.Regions))
	}
	return gco
} //allLocationsR

func lsCoR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch co := G.Unwrap(rootValue, 0).(type) {
	case string:
		return G.MakeStringValue(co)
	default:
		M.Halt(co, 100)
		return nil
	}
} //lsCoR

func lsRgR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch rgs := G.Unwrap(rootValue, 1).(type) {
	case []string:
		grg := G.NewListValue()
		for _, rg := range rgs {
			grg.Append(G.MakeStringValue(rg))
		}
		return grg
	default:
		M.Halt(rgs, 100)
		return nil
	}
} //lsRgR

func locListR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	country := getStringParam(argumentValues, "country")
	region := getStringParam(argumentValues, "region")
	short := getBooleanParam(argumentValues, "short")
	l , err := B.LocationList(country, region, short)
	M.Assert(err == nil || err.ErrorNum() == 8, err, 100)
	error(ts, err)
	if l == nil {
		return G.MakeNullValue()
	}
	gloc := G.NewListValue()
	for _, loc := range l {
		gloc.Append(G.Wrap(loc.Country, loc.Region, loc.City))
	}
	return gloc
} //locListR

func ls1CoR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch co := G.Unwrap(rootValue, 0).(type) {
	case string:
		return G.MakeStringValue(co)
	default:
		M.Halt(co, 100)
		return nil
	}
} //ls1CoR

func ls1RgR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch rg := G.Unwrap(rootValue, 1).(type) {
	case string:
		return G.MakeStringValue(rg)
	default:
		M.Halt(rg, 100)
		return nil
	}
} //ls1RgR

func ls1CyR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch cy := G.Unwrap(rootValue, 2).(type) {
	case string:
		return G.MakeStringValue(cy)
	default:
		M.Halt(cy, 100)
		return nil
	}
} //ls1CyR

func allServicesR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	category := getStringParam(argumentValues, "category")
	l , err := B.AllServices(category)
	M.Assert(err == nil || err.ErrorNum() == 0, err, 100)
	error(ts, err)
	if l == nil {
		return G.MakeNullValue()
	}
	gca := G.NewListValue()
	for _, ca := range l {
		gca.Append(G.Wrap(ca.Category, ca.Specialities))
	}
	return gca
} //allServicesR

func ssCaR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch ca := G.Unwrap(rootValue, 0).(type) {
	case string:
		return G.MakeStringValue(ca)
	default:
		M.Halt(ca, 100)
		return nil
	}
} //ssCaR

func ssSpR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch sps := G.Unwrap(rootValue, 1).(type) {
	case []string:
		gsp := G.NewListValue()
		for _, sp := range sps {
			gsp.Append(G.MakeStringValue(sp))
		}
		return gsp
	default:
		M.Halt(sps, 100)
		return nil
	}
} //ssSpR

func serListR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	country := getStringParam(argumentValues, "country")
	region := getStringParam(argumentValues, "region")
	city := getStringParam(argumentValues, "city")
	category := getStringParam(argumentValues, "category")
	searchedOrSupplied := getSOrSParam(argumentValues, "searchedOrSupplied")
	l , err := B.ServiceList(country, region, city, category, searchedOrSupplied)
	error(ts, err)
	if l == nil {
		return G.MakeNullValue()
	}
	gser := G.NewListValue()
	for _, s := range l {
		gser.Append(G.Wrap(s.Category, s.Speciality))
	}
	return gser
} //serListR

func ss1CaR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch ca := G.Unwrap(rootValue, 0).(type) {
	case string:
		return G.MakeStringValue(ca)
	default:
		M.Halt(ca, 100)
		return nil
	}
} //ss1CaR

func ss1SpR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch sp := G.Unwrap(rootValue, 1).(type) {
	case string:
		return G.MakeStringValue(sp)
	default:
		M.Halt(sp, 100)
		return nil
	}
} //ss1SpR

func idOfR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	nickname := getStringParam(argumentValues, "nickname")
	id, err := B.IdentityOf(nickname)
	error(ts, err)
	if id == nil {
		return G.MakeNullValue()
	}
	return G.Wrap(id)
} //idOfR

func allIdsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	hint := getStringParam(argumentValues, "hint")
	l := B.AllIdentities(hint)
	gl := G.NewListValue()
	for _, id := range l {
		gl.Append(G.Wrap(id))
	}
	return gl
} //allIdsR

func idNR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeStringValue(id.Nickname)
	default:
		M.Halt(id, 101)
		return nil
	}
} //idNR

func idPassR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeStringValue(id.Password)
	default:
		M.Halt(id, 101)
		return nil
	}
} //idPassR

func idGNR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeStringValue(id.G1_nickname)
	default:
		M.Halt(id, 101)
		return nil
	}
} //idGNR

func idPubR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeStringValue(string(id.Pubkey))
	default:
		M.Halt(id, 101)
		return nil
	}
} //idPubR

func idMailR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeStringValue(id.Email)
	default:
		M.Halt(id, 101)
		return nil
	}
} //idMailR

func idLocsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		l := G.NewListValue()
		for _, ca := range id.Locations {
			l.Append(G.Wrap(ca.Country, ca.Regions))
		}
		return l
	default:
		M.Halt(id, 101)
		return nil
	}
} //idLocsR

func locCoR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch co := G.Unwrap(rootValue, 0).(type) {
	case string:
		return G.MakeStringValue(co)
	default:
		M.Halt(co, 100)
		return nil
	}
} //locCoR

func locRgsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch rgs := G.Unwrap(rootValue, 1).(type) {
	case []*B.RegionT:
		l := G.NewListValue()
		for _, rg := range rgs {
			l.Append(G.Wrap(rg.Region, rg.Cities))
		}
		return l
	default:
		M.Halt(rgs, 101)
		return nil
	}
} //locRgsR

func rgRgR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch rg := G.Unwrap(rootValue, 0).(type) {
	case string:
		return G.MakeStringValue(rg)
	default:
		M.Halt(rg, 100)
		return nil
	}
} //rgRgR

func rgCyR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	switch cys := G.Unwrap(rootValue, 1).(type) {
	case []string:
		l := G.NewListValue()
		for _, cy := range cys {
			l.Append(G.MakeStringValue(cy))
		}
		return l
	default:
		M.Halt(cys, 101)
		return nil
	}
} //rgCyR

func idTrprtR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeBooleanValue(id.Transport)
	default:
		M.Halt(id, 100)
		return nil
	}
} //idTrprtR

func idPickR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeBooleanValue(id.Pickup)
	default:
		M.Halt(id, 100)
		return nil
	}
} //idPickR

func idSupSerR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		s := id.Supplied_services
		gca := G.NewListValue()
		if s == nil {
			return gca
		}
		for _, ca := range s {
			gca.Append(G.Wrap(ca.Category, ca.Specialities))
		}
		return gca
	default:
		M.Halt(id, 100)
		return nil
	}
} //idSupSerR

func idSearchSerR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		s := id.Searched_services
		gca := G.NewListValue()
		if s == nil {
			return gca
		}
		for _, ca := range s {
			gca.Append(G.Wrap(ca.Category, ca.Specialities))
		}
		return gca
	default:
		M.Halt(id, 100)
		return nil
	}
} //idSearchSerR

func idAddR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := G.Unwrap(rootValue, 0)
	M.Assert(id != nil, 100)
	switch id := id.(type) {
	case *B.IdentityT:
		return G.MakeStringValue(id.Additional_informations)
	default:
		M.Halt(id, 101)
		return nil
	}
} //idAddR

func idListR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	country := getStringParam(argumentValues, "country")
	region := getStringParam(argumentValues, "region")
	city := getStringParam(argumentValues, "city")
	services, caPos := getServiceQParam(argumentValues, "services")
	if caPos {
		error(ts, M.NewError("|voidCa|", 6))
		return G.MakeNullValue()
	}
	searchedOrSupplied := getSOrSParam(argumentValues, "searchedOrSupplied")
	ids, err := B.IdentityList(country, region, city, services, searchedOrSupplied)
	error(ts, err)
	if ids == nil {
		return G.MakeNullValue()
	}
	l := G.NewListValue()
	for _, id := range ids {
		l.Append(G.Wrap(id))
	}
	return l
} //idListR

func locInsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	co := getStringParam(argumentValues, "country")
	rg := getStringParam(argumentValues, "region")
	loc, err := B.InsertLocation(co, rg)
	error(ts, err)
	M.Assert(loc != nil, 100)
	return G.Wrap(loc.Country, loc.Region)
} //locInsR

func locsInsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	locs := getLocation1ShortListParam(argumentValues, "locations")
	err := B.InsertLocations(locs)
	res := err == nil
	error(ts, err)
	return G.MakeBooleanValue(res)
} //locsInsR

func serInsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	ca := getStringParam(argumentValues, "category")
	sp := getStringParam(argumentValues, "speciality")
	ser, err := B.InsertService(ca, sp)
	error(ts, err)
	M.Assert(ser != nil, 100)
	return G.Wrap(ser.Category, ser.Speciality)
} //serInsR

func sersInsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	sers := getService1ListParam(argumentValues, "services")
	err := B.InsertServices(sers)
	res := err == nil
	error(ts, err)
	return G.MakeBooleanValue(res)
} //sersInsR

func IdInsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	id := getIdParam(ts, argumentValues, "id")
	if id == nil {
		return G.MakeNullValue()
	}
	idd, err := B.InsertIdentity(id)
	error(ts, err)
	if idd == nil {
		return G.MakeNullValue()
	}
	return G.Wrap(idd)
} //IdInsR

func IdsInsR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	ids := getIdsParam(ts, argumentValues, "ids")
	if ids == nil {
		return G.MakeNullValue()
	}
	err := B.InsertIdentities(ids)
	error(ts, err)
	if err != nil {
		return G.MakeNullValue()
	}
	l := G.NewListValue()
	for _, id := range ids {
		l.Append(G.Wrap(id))
	}
	return l
} //IdsInsR

func locRemR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	co := getStringParam(argumentValues, "country")
	rg := getStringParam(argumentValues, "region")
	err := B.RemoveLocation(co, rg)
	res := err == nil
	error(ts, err)
	return G.MakeBooleanValue(res)
} //locRemR

func serRemR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	ca := getStringParam(argumentValues, "category")
	sp := getStringParam(argumentValues, "speciality")
	err := B.RemoveService(ca, sp)
	res := err == nil
	error(ts, err)
	return G.MakeBooleanValue(res)
} //serRemR

func locRenR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	loc := getLocParam(argumentValues, "location")
	co := getStringParam(argumentValues, "country")
	rg := getStringParam(argumentValues, "region")
	cy := getStringParam(argumentValues, "city")
	nn := getStringParam(argumentValues, "newName")
	err := B.RenameLocation(loc, co, rg, cy, nn)
	res := err == nil
	error(ts, err)
	return G.MakeBooleanValue(res)
} //locRenR

func serRenR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	ser := getSerParam(argumentValues, "service")
	ca := getStringParam(argumentValues, "category")
	sp := getStringParam(argumentValues, "speciality")
	nn := getStringParam(argumentValues, "newName")
	err := B.RenameService(ser, ca, sp, nn)
	res := err == nil
	error(ts, err)
	return G.MakeBooleanValue(res)
} //serRenR

func idRemR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	nN := getStringParam(argumentValues, "nickname")
	id, err := B.RemoveIdentity(nN)
	error(ts, err)
	if id == nil {
		return G.MakeNullValue()
	}
	return G.Wrap(id)
} //idRemR

func idChgR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	nN := getStringParam(argumentValues, "nickname")
	nId := getIdParam(ts, argumentValues, "newId")
	id, err := B.ChangeIdentity(nN, nId)
	error(ts, err)
	if id == nil {
		return G.MakeNullValue()
	}
	return G.Wrap(id)
} //idChgR

func idsChgR (ts G.TypeSystem, rootValue *G.OutputObjectValue, argumentValues *A.Tree) G.Value {
	ids := getIdsParam(ts, argumentValues, "ids")
	if ids == nil {
		return G.MakeNullValue()
	}
	err := B.ChangeIdentities(ids)
	error(ts, err)
	if err != nil {
		return G.MakeNullValue()
	}
	l := G.NewListValue()
	for _, id := range ids {
		l.Append(G.Wrap(id))
	}
	return l
} //idsChgR

func fixFieldResolvers (ts G.TypeSystem) {
	
	ts.FixFieldResolver("Query", "allLocations", allLocationsR)
	ts.FixFieldResolver("Query", "locationList", locListR)
	ts.FixFieldResolver("Query", "allServices", allServicesR)
	ts.FixFieldResolver("Query", "serviceList", serListR)
	ts.FixFieldResolver("Query", "identityOf", idOfR)
	ts.FixFieldResolver("Query", "allIdentities", allIdsR)
	ts.FixFieldResolver("Query", "identityList", idListR)
	
	ts.FixFieldResolver("Mutation", "insertLocation", locInsR)
	ts.FixFieldResolver("Mutation", "insertLocations", locsInsR)
	ts.FixFieldResolver("Mutation", "insertService", serInsR)
	ts.FixFieldResolver("Mutation", "insertServices", sersInsR)
	ts.FixFieldResolver("Mutation", "insertIdentity", IdInsR)
	ts.FixFieldResolver("Mutation", "insertIdentities", IdsInsR)
	ts.FixFieldResolver("Mutation", "removeLocation", locRemR)
	ts.FixFieldResolver("Mutation", "removeService", serRemR)
	ts.FixFieldResolver("Mutation", "renameLocation", locRenR)
	ts.FixFieldResolver("Mutation", "renameService", serRenR)
	ts.FixFieldResolver("Mutation", "removeIdentity", idRemR)
	ts.FixFieldResolver("Mutation", "changeIdentity", idChgR)
	ts.FixFieldResolver("Mutation", "changeIdentities", idsChgR)
	
	ts.FixFieldResolver("Location", "country", locCoR)
	ts.FixFieldResolver("Location", "regions", locRgsR)
	
	ts.FixFieldResolver("Region", "region", rgRgR)
	ts.FixFieldResolver("Region", "cities", rgCyR)
	
	ts.FixFieldResolver("LocationShort", "country", lsCoR)
	ts.FixFieldResolver("LocationShort", "regions", lsRgR)
	
	ts.FixFieldResolver("Location1", "country", ls1CoR)
	ts.FixFieldResolver("Location1", "region", ls1RgR)
	ts.FixFieldResolver("Location1", "city", ls1CyR)
	
	ts.FixFieldResolver("Location1Short", "country", ls1CoR)
	ts.FixFieldResolver("Location1Short", "region", ls1RgR)
	
	ts.FixFieldResolver("Service", "category", ssCaR)
	ts.FixFieldResolver("Service", "specialities", ssSpR)
	
	ts.FixFieldResolver("Service1", "category", ss1CaR)
	ts.FixFieldResolver("Service1", "speciality", ss1SpR)
	
	ts.FixFieldResolver("Identity", "nickname", idNR)
	ts.FixFieldResolver("Identity", "password", idPassR)
	ts.FixFieldResolver("Identity", "g1_nickname", idGNR)
	ts.FixFieldResolver("Identity", "pubkey", idPubR)
	ts.FixFieldResolver("Identity", "email", idMailR)
	ts.FixFieldResolver("Identity", "locations", idLocsR)
	ts.FixFieldResolver("Identity", "transport", idTrprtR)
	ts.FixFieldResolver("Identity", "pickup", idPickR)
	ts.FixFieldResolver("Identity", "supplied_services", idSupSerR)
	ts.FixFieldResolver("Identity", "searched_services", idSearchSerR)
	ts.FixFieldResolver("Identity", "additional_informations", idAddR)

} //fixFieldResolvers

func init () {
	fixFieldResolvers(GQ.TS())
} //init
