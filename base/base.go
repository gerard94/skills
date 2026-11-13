/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package base

// Package base implements the database of the server

import (
	
	A	"github.com/gerard94/util/avlT"
	B	"github.com/gerard94/util/gbTree"
	BA	"github.com/gerard94/skills/basic"
	F	"path/filepath"
	M	"github.com/gerard94/util/misc"
	S	"github.com/gerard94/util/sets3"
	ST	"github.com/gerard94/util/strings"
		"fmt"
		"io"
		"os"
		"github.com/gerard94/util/sort2"
		"strings"

)

const (
	
	// Name of the database and of its two copies in the system directory
	baseName = "base.data"
	copyName = "base.data.bak"
	copy1Name = "base.data.bak1"

)

const (
	
	// Numbers of the places of main indexes (see doc/base/indexes.txt) in base
	location1Place = iota // location1Index
	service1Place // service1Index
	identityPlace // identityIndex
	identitiesPlace // identitiesIndex
	
	placeNb // Number of places

)

const (
	
	// SearchedOrSupplied
	SuppliedServices = iota
	SearchedServices
	BothServices

)

const (
	
	// Locations
	Country = iota
	Region
	City
	
	locationsNb

)

const (
	
	// Services
	Category = iota
	Speciality
	
	servicesNb

)

type (
	
	PubkeyT string
	
	Location1ShortT struct {
		Country,
		Region string
	}
	
	Location1ShortListT []*Location1ShortT
	
	Location1T struct {
		Country,
		Region,
		City string
	}
	
	Location1ListT []Location1T
	
	LocationShortT struct {
		Country string
		Regions []string
	}
	
	LocationShortListT []LocationShortT
	
	Service1T struct {
		Category,
		Speciality string
	}
	
	Service1ListT []*Service1T
	
	ServiceT struct {
		Category string
		Specialities []string
	}
	
	ServiceNoPListT []ServiceT
	
	RegionT struct {
		Region string
		Cities []string
	}
	
	LocationT struct {
		Country string
		Regions []*RegionT
	}
	
	LocationListT []*LocationT
	
	ServiceListT []*ServiceT
	
	serviceDispatchT struct {
		supplied,
		searched B.FilePos
	}
	
	serviceDispatchFacT struct {
	}
	
	IdentityT struct {
		Nickname,
		Password,
		G1_nickname string
		Pubkey PubkeyT
		Email string
		Locations LocationListT
		Transport,
		Pickup bool
		Supplied_services,
		Searched_services ServiceListT
		Additional_informations string
	}
	
	identityFacT struct {
	}
	
	IdentityListT []*IdentityT
	
	serviceSetT struct{
		*A.TreeS[*Service1T]
	}
	
	stringKeyManager struct {
	}

)

var (
	
	systemPath = BA.GetSystemPath()
	
	basePath = F.Join(systemPath, baseName)
	copyPath = F.Join(systemPath, copyName)
	copy1Path = F.Join(systemPath, copy1Name)
	
	lg = M.GetLog()

)

var (
	
	pageNb = BA.PageNb()
	
	base *B.Database
	
	serviceDispatchFac serviceDispatchFacT
	identityFac identityFacT
	
	serviceDispatchMan,
	identityMan *B.DataMan
	
	stringKeyMan = B.MakeKM(new(stringKeyManager))
	stringFac B.StringFac
	
	location1Index,
	service1Index,
	identityIndex,
	identitiesIndex *B.Index

)

var errorStrings = [...]string {
	"dontExistCo", //0	Country ^0 doesn't exist in the base
	"existsId", //1	nickname ^0 already exists in the base
	"dontExist2Loc", //2	Location ^0/^1 doesn't exist in the base
	"exists2Loc", //3	Location ^0 / ^1 already exists in the base
	"voidCo", //4	Void country
	"voidRg", //5	Void region
	"voidCa", //6	Void category
	"voidSp", //7	Void speciality
	"inVoidCo", //8	Region ^0 in void country
	"inVoidRg", //9	City ^0 in void region
	"idInLoc", //10	There are identities in location ^0 / ^1
	"inVoidCa", //11	Speciality ^0 in void category
	"nnNull", //12	New name is void
	"dontExistCa", //13	Category ^0 doesn't exist in the base
	"dontExistId", //14	Nickname ^0 doesn't exist in the base
	"dontExist2Ser", //15	Service ^0/^1 doesn't exist in the base
	"exists2Ser", //16	Service ^0 / ^1 already exists in the base
	"idInSer", //17	There are identities in service ^0 / ^1
	"voidCy", //18	Void city
	"voidNN", //19	Void nickname
	"voidPwd", //20	Void password
	"absentLoc", //21	Location ^0/^1 in identities doesn't exist in new locations
	"absentSer", //22	Service ^0/^1 in identities doesn't exist in new services
}

func GetError (n int, f string, p ...string) M.Nerror {
	M.Assert(n >= 0 && n < len(errorStrings), 100)
	var b strings.Builder
	b.WriteRune('|')
	b.WriteString(errorStrings[n])
	for _, pi := range p {
		b.WriteRune('|')
		b.WriteString(pi)
	}
	b.WriteRune('|')
	return M.NewError(b.String(), n)
} //GetError

// Utilities for debugging

func printLocation1Short (l *Location1ShortT) {
	fmt.Printf("{Country: %v, Region: %v}", l.Country, l.Region)
} //printLocation1Short

func printLocation1ShortList (l Location1ShortListT) {
	for i, li := range l {
		if i > 0 {
			fmt.Println()
		}
		printLocation1Short(li)
	}
} //printLocation1ShortList

func printLocation1 (out io.Writer, l *Location1T) {
	fmt.Fprintf(out, "{Country: %v, Region: %v, City: %v}", l.Country, l.Region, l.City)
} //printLocation1

func printLocation1List (out io.Writer, l Location1ListT) {
	for i, li := range l {
		if i > 0 {
			fmt.Fprintln(out)
		}
		printLocation1(out, &li)
	}
} //printLocation1List

func printLocationShort (out io.Writer, l *LocationShortT) {
	fmt.Fprintf(out, "{Country: %v", l.Country)
	for _, li := range l.Regions {
		fmt.Fprintf(out, "\n\tRegion: %v", li)
	}
	fmt.Fprintln(out, "\n}")
} //printLocationShort

func printLocationShortList (out io.Writer, l LocationShortListT) {
	for i, li := range l {
		if i > 0 {
			fmt.Fprintln(out)
		}
		printLocationShort(out, &li)
	}
} //printLocationShortList

func printRegion (out io.Writer, r *RegionT) {
	fmt.Fprintf(out, "\t{Region: %v", r.Region)
	for _, ci := range r.Cities {
		fmt.Fprintf(out, "\n\t\tCity: %v", ci)
	}
	fmt.Fprint(out, "\n\t}")
} //printRegion

func printLocation (out io.Writer, l *LocationT) {
	fmt.Fprintf(out, "{Country: %v", l.Country)
	for _, r := range l.Regions {
		fmt.Fprintln(out)
		printRegion(out, r)
	}
	fmt.Fprint(out, "\n}")
} //printLocation

func printLocationList (out io.Writer, l LocationListT) {
	for i, li := range l {
		if i > 0 {
			fmt.Fprintln(out)
		}
		printLocation(out, li)
	}
} //printLocationList

func printService1 (out io.Writer, s *Service1T) {
	fmt.Fprintf(out, "{Category: %v, Speciality: %v}", s.Category, s.Speciality)
} //printService1

func printService1List (out io.Writer, l Service1ListT) {
	for i, li := range l {
		if i > 0 {
			fmt.Fprintln(out)
		}
		printService1(out, li)
	}
} //printService1List

func printService (out io.Writer, s *ServiceT) {
	fmt.Fprintf(out, "{Category: %v", s.Category)
	for _, sp := range s.Specialities {
		fmt.Fprintf(out, "\n\tSpeciality: %v", sp)
	}
	fmt.Fprint(out, "\n}")
} //printService

func printServiceList (out io.Writer, l ServiceListT) {
	for i, s := range l {
		if i > 0 {
			fmt.Fprintln(out)
		}
		printService(out, s)
	}
} //printServiceList

func printServiceNoPList (out io.Writer, l ServiceNoPListT) {
	for i, s := range l {
		if i > 0 {
			fmt.Fprintln(out)
		}
		printService(out, &s)
	}
} //printServiceNoPList

func printIdentity (out io.Writer, id *IdentityT) {
	fmt.Fprintf(out, "nickname: %v", id.Nickname)
	fmt.Fprintf(out, "\npassword: %v", id.Password)
	fmt.Fprintf(out, "\ng1_nickname: %v", id.G1_nickname)
	fmt.Fprintf(out, "\npubkey: %v", id.Pubkey)
	fmt.Fprintf(out, "\nemail: %v", id.Email)
	fmt.Fprint(out, "\nlocations:\n")
	printLocationList(out, id.Locations)
	fmt.Fprintf(out, "\ntransport: %v", id.Transport)
	fmt.Fprintf(out, "\npickup: %v", id.Pickup)
	fmt.Fprint(out, "\nSupplied_services:\n")
	printServiceList(out, id.Supplied_services)
	fmt.Fprint(out, "\nSearched_services:\n")
	printServiceList(out, id.Searched_services)
	fmt.Fprintf(out, "\nadditional_informations: %v", id.Additional_informations)
} //printIdentity

// End of utilities for debugging

func (sd *serviceDispatchT) Read (r *B.Reader) {
	sd.supplied = r.InFilePos()
	sd.searched = r.InFilePos()
} //Read

func (rg *RegionT) Read (r *B.Reader) {
	rg.Region = r.InString()
	rg.Cities = make([]string, r.InInt16())
	for i := range rg.Cities {
		rg.Cities[i] = r.InString()
	}
} //Read

func (l *LocationT) Read (r *B.Reader) {
	l.Country = r.InString()
	l.Regions = make([]*RegionT, r.InInt16())
	for i := range l.Regions {
		l.Regions[i] = new(RegionT)
		(l.Regions[i]).Read(r)
	}
} //Read

func (s *ServiceT) Read (r *B.Reader) {
	s.Category = r.InString()
	s.Specialities = make([]string, r.InInt16())
	for i := range s.Specialities {
		s.Specialities[i] = r.InString()
	}
} //Read

func (id *IdentityT) Read (r *B.Reader) {
	id.Nickname = r.InString()
	id.Password = r.InString()
	id.G1_nickname = r.InString()
	id.Pubkey = PubkeyT(r.InString())
	id.Email = r.InString()
	id.Locations = make([]*LocationT, r.InInt16())
	for i := range id.Locations {
		id.Locations[i] = new(LocationT)
		(id.Locations[i]).Read(r)
	}
	id.Transport = r.InBool()
	id.Pickup  = r.InBool()
	id.Supplied_services = make([]*ServiceT, r.InInt16())
	for i := range id.Supplied_services {
		id.Supplied_services[i] = new(ServiceT)
		(id.Supplied_services[i]).Read(r)
	}
	id.Searched_services = make([]*ServiceT, r.InInt16())
	for i := range id.Searched_services {
		id.Searched_services[i] = new(ServiceT)
		(id.Searched_services[i]).Read(r)
	}
	id.Additional_informations = r.InString()
} //Read

func (sd *serviceDispatchT) Write (w *B.Writer) {
	w.OutFilePos(sd.supplied)
	w.OutFilePos(sd.searched)
} //Write

func (rg *RegionT) Write (w *B.Writer) {
	w.OutString(rg.Region)
	w.OutInt16(int16(len(rg.Cities)))
	for i := range rg.Cities {
		w.OutString(rg.Cities[i])
	}
} //Write

func (l *LocationT) Write (w *B.Writer) {
	w.OutString(l.Country)
	w.OutInt16(int16(len(l.Regions)))
	for i := range l.Regions {
		(l.Regions[i]).Write(w)
	}
} //Write

func (s *ServiceT) Write (w *B.Writer) {
	w.OutString(s.Category)
	w.OutInt16(int16(len(s.Specialities)))
	for i := range s.Specialities {
		w.OutString(s.Specialities[i])
	}
} //Write

func (id *IdentityT) Write (w *B.Writer) {
	w.OutString(id.Nickname)
	w.OutString(id.Password)
	w.OutString(id.G1_nickname)
	w.OutString(string(id.Pubkey))
	w.OutString(id.Email)
	w.OutInt16(int16(len(id.Locations)))
	for i := range id.Locations {
		(id.Locations[i]).Write(w)
	}
	w.OutBool(id.Transport)
	w.OutBool(id.Pickup)
	w.OutInt16(int16(len(id.Supplied_services)))
	for i := range id.Supplied_services {
		(id.Supplied_services[i]).Write(w)
	}
	w.OutInt16(int16(len(id.Searched_services)))
	for i := range id.Searched_services {
		(id.Searched_services[i]).Write(w)
	}
	w.OutString(id.Additional_informations)
} //Write

func (serviceDispatchFacT) New (int) B.Data {
	return new(serviceDispatchT)
} //New

func (identityFacT) New (int) B.Data {
	return new(IdentityT)
} //New

func (*stringKeyManager) CompP (key1, key2 B.Data) B.Comp {
	M.Assert(key1 != nil && key2 != nil, 20)
	k1 := key1.(*B.String); k2 := key2.(*B.String)
	return ST.CompP(k1.C, k2.C)
} //CompP

func (man *stringKeyManager) PrefP (key1 B.Data, key2 *B.Data) {
	*key2 = B.StringPrefP(key1.(*B.String), (*key2).(*B.String), man.CompP)
} //PrefP

func saveBase () {
	os.Remove(copy1Path)
	os.Rename(copyPath, copy1Path)
	const bufferSize = 0X800
	f, err := os.Open(basePath)
	if err == nil {
		defer f.Close()
		lg.Println("Making a copy of \"" + baseName + "\"")
		fC, err := os.Create(copyPath); M.Assert(err == nil, err, 100)
		defer fC.Close()
		buf := make([]byte, bufferSize)
		for n, err := f.Read(buf); err == nil; n, err = f.Read(buf) {
			_, err1 := fC.Write(buf[:n]); M.Assert(err1 == nil, err, 101)
		}
		lg.Println("Copy made")
	}
} //saveBase

func OpenBase () {
	M.Assert(base == nil, 100)
	B.Fac.CloseBase(basePath)
	base = B.Fac.OpenBase(basePath, pageNb)
	if base == nil {
		b := B.Fac.CreateBase(basePath, placeNb); M.Assert(b, 101)
		base = B.Fac.OpenBase(basePath, pageNb); M.Assert(base != nil, 102)
		base.WritePlace(location1Place, int64(base.CreateIndex(0)))
		base.WritePlace(service1Place, int64(base.CreateIndex(0)))
		base.WritePlace(identityPlace, int64(base.CreateIndex(0)))
		base.WritePlace(identitiesPlace, int64(base.CreateIndex(0)))
	}
	serviceDispatchMan = base.CreateDataMan(serviceDispatchFac)
	identityMan = base.CreateDataMan(identityFac)
	location1Index = base.OpenIndex(B.FilePos(base.ReadPlace(location1Place)), stringKeyMan, stringFac)
	service1Index = base.OpenIndex(B.FilePos(base.ReadPlace(service1Place)), stringKeyMan, stringFac)
	identityIndex = base.OpenIndex(B.FilePos(base.ReadPlace(identityPlace)), stringKeyMan, stringFac)
	identitiesIndex = base.OpenIndex(B.FilePos(base.ReadPlace(identitiesPlace)), stringKeyMan, stringFac)
} //OpenBase

func CloseBase () {
	M.Assert(base != nil, 100)
	base.CloseBase()
	base = nil
} //CloseBase

/*******************************************************************/

func (s1 *Service1T) Compare (s2 *Service1T) A.Comp {
	if s1.Category < s2.Category {
		return A.Lt
	}
	if s1.Category > s2.Category {
		return A.Gt
	}
	if s1.Speciality < s2.Speciality {
		return A.Lt
	}
	if s1.Speciality > s2.Speciality {
		return A.Gt
	}
	return A.Eq
} //Compare

func (s *serviceSetT) incl (e *Service1T) {
	s.TreeS.SearchIns(e)
} //incl

func (s *serviceSetT) list () Service1ListT {
	l := make(Service1ListT, s.NumberOfElems())
	i := 0
	for e := s.Next(nil); e != nil; e = s.Next(e) {
		l[i] = e.Val()
		i++
	}
	return l
} //list

// 'locationExists' returns true if 'country' is null, or exists and 'region' is null or 'exists in 'country'
func locationExists (country, region string) bool {
	M.Assert(base != nil, 100)
	if country == "" {
		return true
	}
	rCo := location1Index.NewReader()
	sCo := &B.String{country}
	if !rCo.Search(sCo) {
		return false
	}
	if region == "" {
		return true
	}
	iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
	rRg := iRg.NewReader()
	sRg := &B.String{region}
	return rRg.Search(sRg)
} //locationExists

// 'serviceExists' returns true if 'category' is null, or exists and 'speciality' is null or exists in 'category'
func serviceExists (category, speciality string) bool {
	M.Assert(base != nil, 100)
	if category == "" {
		return true
	}
	rCa := service1Index.NewReader()
	sCa := &B.String{category}
	if !rCa.Search(sCa) {
		return false
	}
	if speciality == "" {
		return true
	}
	iSp := base.OpenIndex(rCa.ReadValue(), stringKeyMan, stringFac)
	rSp := iSp.NewReader()
	sSp := &B.String{speciality}
	return rSp.Search(sSp)
} //serviceExists

func AllLocations (country string) (LocationShortListT, M.Nerror) {
	
	InsertRegions := func (ls *LocationShortT, co string, pos B.FilePos) {
		iRg := base.OpenIndex(pos, stringKeyMan, stringFac)
		(*ls).Country = co
		l := make([]string, iRg.NumberOfKeys())
		(*ls).Regions = l
		rRg := iRg.NewReader()
		i := 0
		for rRg.Next(); rRg.PosSet(); rRg.Next() {
			l[i] = rRg.CurrentKey().(*B.String).C
			i++
		}
	} //InsertRegions
	
	//AllLocations
	M.Assert(base != nil, 100)
	if country != "" &&  !locationExists(country, "") {
		return nil, GetError(0, "allLocations", country)
	}
	var l LocationShortListT
	rCo := location1Index.NewReader()
	if country == "" {
		l = make(LocationShortListT, location1Index.NumberOfKeys())
		i := 0
		for rCo.Next(); rCo.PosSet(); rCo.Next() {
			InsertRegions(&l[i], rCo.CurrentKey().(*B.String).C, rCo.ReadValue())
			i++
		}
	} else {
		 l = make(LocationShortListT, 1)
		sCo := &B.String{country}
		b := rCo.Search(sCo); M.Assert(b, 101)
		InsertRegions(&l[0], country, rCo.ReadValue())
	}
	return l, nil
} //AllLocations

func LocationList (country, region string, short bool) (Location1ListT, M.Nerror) {
	
	insertLoc := func (l *Location1ListT, co, rg, cy string) {
		*l = append(*l, Location1T{co, rg, cy})
	} //insertLoc
	
	//LocationList
	M.Assert(base != nil, 100)
	if country == "" && region != "" {
		return nil, GetError(8, "locationList", region)
	}
	
	l := make(Location1ListT, 0)
	
	rCo := identitiesIndex.NewReader()
	if country == "" {
		for rCo.Next(); rCo.PosSet(); rCo.Next() {
			co := rCo.CurrentKey().(*B.String).C
			if co != "" {
				iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
				rRg := iRg.NewReader()
				for rRg.Next(); rRg.PosSet(); rRg.Next() {
					rg := rRg.CurrentKey().(*B.String).C
					if short {
						insertLoc(&l, co, rg, "")
					} else {
						iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
						rCy := iCy.NewReader()
						for rCy.Next(); rCy.PosSet(); rCy.Next() {
							cy := rCy.CurrentKey().(*B.String).C
							insertLoc(&l, co, rg, cy)
						}
					}
				}
			}
		}
	} else {
		sCo := &B.String{country}
		if rCo.Search(sCo) {
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rRg := iRg.NewReader()
			if region == "" {
				for rRg.Next(); rRg.PosSet(); rRg.Next() {
					rg := rRg.CurrentKey().(*B.String).C
					if short {
						insertLoc(&l, country, rg, "")
					} else {
						iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
						rCy := iCy.NewReader()
						for rCy.Next(); rCy.PosSet(); rCy.Next() {
							cy := rCy.CurrentKey().(*B.String).C
							insertLoc(&l, country, rg, cy)
						}
					}
				}
			} else {
				sRg := &B.String{region}
				if rRg.Search(sRg) {
					if short {
						insertLoc(&l, country, region, "")
					} else {
						iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
						rCy := iCy.NewReader()
						for rCy.Next(); rCy.PosSet(); rCy.Next() {
							cy := rCy.CurrentKey().(*B.String).C
							insertLoc(&l, country, region, cy)
						}
					}
				}
			}
		}
	}
	return l, nil
} //LocationList

func AllServices (category string) (ServiceNoPListT, M.Nerror) {
	
	InsertSpecialities := func (ls *ServiceT, ca string, pos B.FilePos) {
		iSp := base.OpenIndex(pos, stringKeyMan, stringFac)
		(*ls).Category = ca
		l := make([]string, iSp.NumberOfKeys())
		(*ls).Specialities = l
		rSp := iSp.NewReader()
		i := 0
		for rSp.Next(); rSp.PosSet(); rSp.Next() {
			l[i] = rSp.CurrentKey().(*B.String).C
			i++
		}
	} //InsertSpecialities
	
	//AllServices
	M.Assert(base != nil, 100)
	if category != "" &&  !serviceExists(category, "") {
		return nil, GetError(13, "allServices", category)
	}
	var l ServiceNoPListT
	rCa := service1Index.NewReader()
	if category == "" {
		l = make(ServiceNoPListT, service1Index.NumberOfKeys())
		i := 0
		for rCa.Next(); rCa.PosSet(); rCa.Next() {
			InsertSpecialities(&l[i], rCa.CurrentKey().(*B.String).C, rCa.ReadValue())
			i++
		}
	} else {
		sCa := &B.String{category}
		b := rCa.Search(sCa); M.Assert(b, 101)
		 l = make(ServiceNoPListT, 1)
		InsertSpecialities(&l[0], category, rCa.ReadValue())
	}
	return l, nil
} //AllServices

func ServiceList (country, region, city, category string, searchedOrSupplied int8) (Service1ListT, M.Nerror) {
	
	insertSpecialities := func (s *serviceSetT, ca string, pos B.FilePos) {
		iSp := base.OpenIndex(pos, stringKeyMan, stringFac)
		rSp := iSp.NewReader()
		for rSp.Next(); rSp.PosSet(); rSp.Next() {
			sp := rSp.CurrentKey().(*B.String).C
			s.incl(&Service1T{ca, sp})
		}
	} //insertSpecialities
	
	insertCategories := func (s *serviceSetT, pos B.FilePos) {
		
		insertCat := func (pos B.FilePos) {
			iCa := base.OpenIndex(pos, stringKeyMan, stringFac)
			rCa := iCa.NewReader()
			if category == "" {
				for rCa.Next(); rCa.PosSet(); rCa.Next() {
					ca := rCa.CurrentKey().(*B.String).C
					if ca != "" {
						insertSpecialities(s, ca, rCa.ReadValue())
					}
				}
			} else {
				sCa := &B.String{category}
				if rCa.Search(sCa) {
					ca := rCa.CurrentKey().(*B.String).C
					insertSpecialities(s, ca, rCa.ReadValue())
				}
			}
		} //insertCat
		
		//insertCategories
		sd := serviceDispatchMan.ReadData(pos).(*serviceDispatchT)
		switch searchedOrSupplied {
		case SuppliedServices:
			insertCat(sd.supplied)
		case SearchedServices:
			insertCat(sd.searched)
		case BothServices:
			insertCat(sd.supplied)
			insertCat(sd.searched)
		default:
			M.Halt("searchedOrSupplied = ", searchedOrSupplied, 20)
		}
	} //insertCategories
	
	//ServiceList
	M.Assert(base != nil, 100)
	if country == "" && region != "" {
		return nil, GetError(8, "serviceList", region)
	}
	if region == "" && city != "" {
		return nil, GetError(9, "serviceList", city)
	}
	
	var s = &serviceSetT{&A.TreeS[*Service1T]{A.New[*Service1T]()}}
	
	rCo := identitiesIndex.NewReader()
	if country == "" {
		for rCo.Next(); rCo.PosSet(); rCo.Next() {
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rRg := iRg.NewReader()
			for rRg.Next(); rRg.PosSet(); rRg.Next() {
				iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
				rCy := iCy.NewReader()
				for rCy.Next(); rCy.PosSet(); rCy.Next() {
					insertCategories(s, rCy.ReadValue())
				}
			}
		}
	} else {
		sCo := &B.String{country}
		if rCo.Search(sCo) {
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rRg := iRg.NewReader()
			if region == "" {
				for rRg.Next(); rRg.PosSet(); rRg.Next() {
					iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
					rCy := iCy.NewReader()
					for rCy.Next(); rCy.PosSet(); rCy.Next() {
						insertCategories(s, rCy.ReadValue())
					}
				}
			} else {
				sRg := &B.String{region}
				if rRg.Search(sRg) {
					iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
					rCy := iCy.NewReader()
					if city == "" {
						for rCy.Next(); rCy.PosSet(); rCy.Next() {
							insertCategories(s, rCy.ReadValue())
						}
					} else {
						sCy := &B.String{""}
						if rCy.Search(sCy) {
							insertCategories(s, rCy.ReadValue())
						}
						sCy = &B.String{city}
						if rCy.Search(sCy) {
							insertCategories(s, rCy.ReadValue())
						}
					}
				}
			}
		}
	}
	
	return s.list(), nil
} //ServiceList

func IdentityOf (nickname string) (*IdentityT, M.Nerror) {
	rId := identityIndex.NewReader()
	sId := &B.String{nickname}
	if !rId.Search(sId) {
		return nil, GetError(14, "identityOf", nickname)
	}
	return identityMan.ReadData(rId.ReadValue()).(*IdentityT), nil
} //IdentityOf

func AllIdentities (hint string) IdentityListT {
	M.Assert(base != nil, 100)
	lId := make(IdentityListT, 0)
	rId := identityIndex.NewReader()
	hint = ST.ToDown(hint)
	sId := &B.String{hint}
	for rId.Search(sId); rId.PosSet(); rId.Next() {
		idk := rId.CurrentKey().(*B.String).C
		if !ST.Prefix(hint, ST.ToDown(idk)) {break}
		id := identityMan.ReadData(rId.ReadValue()).(*IdentityT)
		lId = append(lId, id)
	}
	return lId
} //AllIdentities

func IdentityList (country, region, city string, services ServiceListT, searchedOrSupplied int8) (IdentityListT, M.Nerror) {
	
	unionSet := func (s *S.Set, pos B.FilePos) {
		iIds := base.OpenIndex(pos, stringKeyMan, stringFac)
		rIds := iIds.NewReader()
		for rIds.Next(); rIds.PosSet(); rIds.Next() {
			id := identityMan.ReadData(rIds.ReadValue()).(*IdentityT)
			s.Incl(id.Nickname)
		}
	} //unionSet
	
	insertServices := func (s *S.Set, pos B.FilePos) {
		
		insert1 := func (pos B.FilePos) {
			iCa := base.OpenIndex(pos, stringKeyMan, stringFac)
			rCa := iCa.NewReader()
			if services == nil || len(services) == 0 {
				for rCa.Next(); rCa.PosSet(); rCa.Next() {
					iSp := base.OpenIndex(rCa.ReadValue(), stringKeyMan, stringFac)
					rSp := iSp.NewReader()
					for rSp.Next(); rSp.PosSet(); rSp.Next() {
						unionSet(s, rSp.ReadValue())
					}
				}
			} else {
				for _, ca := range services {
					M.Assert(ca.Category != "", 100)
					sCa := &B.String{ca.Category}
					if rCa.Search(sCa) {
						iSp := base.OpenIndex(rCa.ReadValue(), stringKeyMan, stringFac)
						rSp := iSp.NewReader()
						if ca.Specialities == nil {
							for rSp.Next(); rSp.PosSet(); rSp.Next() {
								unionSet(s, rSp.ReadValue())
							}
						} else {
							for _, sp := range ca.Specialities {
								sSp := &B.String{sp}
								if rSp.Search(sSp) {
									unionSet(s, rSp.ReadValue())
								}
							}
						}
					}
				}
			}
		} //insert1
		
		//insertServices
		sd := serviceDispatchMan.ReadData(pos).(*serviceDispatchT)
		switch searchedOrSupplied {
		case SuppliedServices:
			insert1(sd.supplied)
		case SearchedServices:
			insert1(sd.searched)
		case BothServices:
			insert1(sd.supplied)
			insert1(sd.searched)
		default:
			M.Halt("searchedOrSupplied = ", searchedOrSupplied, 20)
		}
	} //insertServices
	
	//IdentityList
	M.Assert(base != nil, 100)
	if country == "" && region != "" {
		return nil, GetError(8, "identityList", region)
	}
	if region == "" && city != "" {
		return nil, GetError(9, "identityList", city)
	}
	if services != nil {
		for _, ca := range services {
			M.Assert(ca.Category != "", 101)
		}
	}
	
	var s = S.NewSet()
	
	rCo := identitiesIndex.NewReader()
	if country == "" {
		for rCo.Next(); rCo.PosSet(); rCo.Next() {
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rRg := iRg.NewReader()
			for rRg.Next(); rRg.PosSet(); rRg.Next() {
				iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
				rCy := iCy.NewReader()
				for rCy.Next(); rCy.PosSet(); rCy.Next() {
					insertServices(s, rCy.ReadValue())
				}
			}
		}
	} else {
		sCo := &B.String{country}
		if rCo.Search(sCo) {
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rRg := iRg.NewReader()
			if region == "" {
				for rRg.Next(); rRg.PosSet(); rRg.Next() {
					iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
					rCy := iCy.NewReader()
					for rCy.Next(); rCy.PosSet(); rCy.Next() {
						insertServices(s, rCy.ReadValue())
					}
				}
			} else {
				sRg := &B.String{region}
				if rRg.Search(sRg) {
					iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
					rCy := iCy.NewReader()
					if city == "" {
						for rCy.Next(); rCy.PosSet(); rCy.Next() {
							insertServices(s, rCy.ReadValue())
						}
					} else {
						sCy := &B.String{city}
						if rCy.Search(sCy) {
							insertServices(s, rCy.ReadValue())
							sCy = &B.String{""} // the whole region
							if rCy.Search(sCy) {
								insertServices(s, rCy.ReadValue())
							}
						}
					}
				}
			}
		}
	}
	
	var l = make(IdentityListT, s.NbElems())
	it := s.Attach()
	i := 0
	for name, ok := it.FirstE(); ok; name, ok = it.NextE() {
		l[i], _ = IdentityOf(name)
		i++
	}
	return l, nil
} //IdentityList

func insertLocation (country, region string) (*Location1ShortT, M.Nerror) {
	M.Assert(base != nil, 100)
	loc := &Location1ShortT{Country: country, Region: region}
	if country == "" {
		return loc, GetError(4, "insertLocation")
	}
	if region == "" {
		return loc, GetError(5, "insertLocation")
	}
	wCo := location1Index.Writer()
	sCo := &B.String{country}
	if !wCo.SearchIns(sCo) {
		wCo.WriteValue(base.CreateIndex(0))
	}
	iRg := base.OpenIndex(wCo.ReadValue(), stringKeyMan, stringFac)
	wRg := iRg.Writer()
	sRg := &B.String{region}
	if wRg.SearchIns(sRg) {
		return loc, GetError(3, "insertLocation", country, region)
	}
	return loc, nil
} //insertLocation

func InsertLocation (country, region string) (*Location1ShortT, M.Nerror) {
	defer base.UpdateBase()
	return insertLocation(country, region)
} //InsertLocation

// Location1ShortListT is a 'sort.Finder'
func (l Location1ShortListT) Less (i, j int) bool {
	switch {
	case l[i].Country < l[j].Country:
		return true
	case l[i].Country > l[j].Country:
		return false
	default:
		return l[i].Region < l[j].Region
	}
} //Less

// Location1ShortListT is a 'sort.Sorter'
func (l Location1ShortListT) Swap (i, j int) {
	l[i], l[j] = l[j], l[i]
} //Swap

func InsertLocations (locations Location1ShortListT) M.Nerror {
	M.Assert(base != nil, 20)
	M.Assert(locations != nil, 21)
	locats := append(locations, nil) // add a void element without changing 'services'
	tf := sort.TF{Finder: locats}
	ts := sort.TS{Sorter: locats}
	n := len(locations)
	ts.QuickSort(0, n - 1) // sort locats
	locs, err := LocationList("", "", true); M.Assert(err == nil, err, 100)
	for _, l := range locs {
		locats[n] = &Location1ShortT{Country: l.Country, Region: l.Region}
		m := n
		tf.BinSearch(0, n - 1, &m)
		if m == n { // error if a used location is absent in the new locations
			return GetError(21, "InsertLocations", l.Country, l.Region)
		}
	}
	
	defer base.UpdateBase()
	base.DeleteIndex(B.FilePos(base.ReadPlace(location1Place)))
	pos := base.CreateIndex(0)
	base.WritePlace(location1Place, int64(pos))
	location1Index = base.OpenIndex(pos, stringKeyMan, stringFac)
	for _, l := range locations {
		_, err := insertLocation(l.Country, l.Region)
		if err != nil {
			return err
		}
	}
	return nil
} //InsertLocations

func insertService (category, speciality string) (*Service1T, M.Nerror) {
	M.Assert(base != nil, 20)
	s := &Service1T{Category: category, Speciality: speciality}
	if category == "" {
		return s, GetError(6, "insertService")
	}
	if speciality == "" {
		return s, GetError(7, "insertService")
	}
	wCa := service1Index.Writer()
	sCa := &B.String{category}
	if !wCa.SearchIns(sCa) {
		wCa.WriteValue(base.CreateIndex(0))
	}
	iSp := base.OpenIndex(wCa.ReadValue(), stringKeyMan, stringFac)
	wSp := iSp.Writer()
	sSp := &B.String{speciality}
	if wSp.SearchIns(sSp) {
		return s, GetError(16, "insertService",category , speciality)
	}
	return s, nil
} //insertService

func InsertService (category, speciality string) (*Service1T, M.Nerror) {
	defer base.UpdateBase()
	return insertService(category, speciality)
} //InsertService

// Service1ListT is a 'sort.Finder'
func (s Service1ListT) Less (i, j int) bool {
	switch {
	case s[i].Category < s[j].Category:
		return true
	case s[i].Category > s[j].Category:
		return false
	default:
		return s[i].Speciality < s[j].Speciality
	}
} //Less

// Service1ListT is a 'sort.Sorter'
func (s Service1ListT) Swap (i, j int) {
	s[i], s[j] = s[j], s[i]
} //Swap

func InsertServices (services Service1ListT) M.Nerror {
	M.Assert(base != nil, 20)
	M.Assert(services != nil, 21)
	servs := append(services, nil) // add a void element without changing 'services'
	tf := sort.TF{Finder: servs}
	ts := sort.TS{Sorter: servs}
	n := len(services)
	ts.QuickSort(0, n - 1) // sort servs
	sers, err := ServiceList("", "", "", "", BothServices); M.Assert(err == nil, err, 100)
	for _, s := range sers {
		servs[n] = s
		m := n
		tf.BinSearch(0, n - 1, &m)
		if m == n { // error if a used service is absent in the new services
			return GetError(22, "InsertServices", s.Category, s.Speciality)
		}
	}
	
	defer base.UpdateBase()
	base.DeleteIndex(B.FilePos(base.ReadPlace(service1Place)))
	pos := base.CreateIndex(0)
	base.WritePlace(service1Place, int64(pos))
	service1Index = base.OpenIndex(pos, stringKeyMan, stringFac)
	for _, s := range services {
		_, err := insertService(s.Category, s.Speciality)
		if err != nil {
			return err
		}
	}
	return nil
} //InsertServices

func verifyIdBeforeInsertion (id *IdentityT, nickHere bool) (*IdentityT, M.Nerror) {
	M.Assert(base != nil, 20)
	M.Assert(id != nil, 21)
	M.Assert(id.Nickname != "", 22)
	M.Assert(id.Password != "", 23)
	M.Assert(id.Locations != nil, 24)
	for _, co := range id.Locations {
		M.Assert(co.Country != "", 25)
		M.Assert(co.Regions != nil, 26)
		for _, rg := range co.Regions {
			M.Assert(rg.Region != "", 27)
			if !locationExists(co.Country, rg.Region) {
				return id, GetError(2, "insertIdentity", co.Country, rg.Region)
			}
		}
	}
	M.Assert(id.Supplied_services != nil, 29)
	for _, ca := range id.Supplied_services {
		M.Assert(ca.Category != "", 30)
		M.Assert(ca.Specialities != nil, 31)
		for _, sp := range ca.Specialities {
			M.Assert(sp != "", 32)
			if !serviceExists(ca.Category, sp) {
				return id, GetError(15, "insertIdentity", ca.Category, sp)
			}
		}
	}
	M.Assert(id.Searched_services != nil, 33)
	for _, ca := range id.Searched_services {
		M.Assert(ca.Category != "", 34)
		M.Assert(ca.Specialities != nil, 35)
		for _, sp := range ca.Specialities {
			M.Assert(sp != "", 36)
			if !serviceExists(ca.Category, sp) {
				return id, GetError(15, "insertIdentity", ca.Category, sp)
			}
		}
	}
	_, err := IdentityOf(id.Nickname)
	if nickHere {
		if err != nil {
			return nil, GetError(14, "changeIdentity", id.Nickname)
		}
	} else {
		if err == nil {
			return nil, GetError(1, "insertIdentity", id.Nickname)
		}
	}
	return id, nil
} //verifyIdBeforeInsertion

// 'insertId' inserts 'id' into the base and returns it. The nickname of 'id' is supposed NOT to be already in the base. if 'cities' is empty in 'id.Locations', an empty key is inserted in the corresponding index.
func insertId (id *IdentityT) *IdentityT {
	
	insertService := func (pos, idRef B.FilePos, sId B.Data, s ServiceListT) {
		iCa := base.OpenIndex(pos, stringKeyMan, stringFac)
		wCa := iCa.Writer()
		for _, ca := range s {
			sCa := &B.String{ca.Category}
			if !wCa.SearchIns(sCa) {
				wCa.WriteValue(base.CreateIndex(0))
			}
			iSp := base.OpenIndex(wCa.ReadValue(), stringKeyMan, stringFac)
			wSp := iSp.Writer()
			for _, sp := range ca.Specialities {
				sSp := &B.String{sp}
				if !wSp.SearchIns(sSp) {
					wSp.WriteValue(base.CreateIndex(0))
				}
				iIds := base.OpenIndex(wSp.ReadValue(), stringKeyMan, stringFac)
				wIds := iIds.Writer()
				if !wIds.SearchIns(sId) {
					wIds.WriteValue(idRef)
				}
			}
		}
	} //insertService
	
	sortLocations := func (locs LocationListT) {
		for ci := 1; ci < len(locs); ci++ {
			for cj := ci; cj > 0 && ST.CompP(locs[cj].Country, locs[cj - 1].Country) == ST.Lt; cj--{
				locs[cj], locs[cj - 1] = locs[cj - 1], locs[cj]
			}
		}
		for ci := 0; ci < len(locs); ci++ {
			rg := locs[ci].Regions
			for ri := 1; ri < len(rg); ri++ {
				for rj := ri; rj > 0 && ST.CompP(rg[rj].Region, rg[rj - 1].Region) == ST.Lt; rj-- {
					rg[rj], rg[rj - 1] = rg[rj - 1], rg[rj]
				}
			}
			for ri := 0; ri < len(rg); ri++ {
				cy := rg[ri].Cities
				for ci := 1; ci < len(cy); ci++ {
					for cj := ci; cj > 0 && ST.CompP(cy[cj], cy[cj - 1]) == ST.Lt; cj-- {
						cy[cj], cy[cj - 1] = cy[cj - 1], cy[cj]
					}
				}
			}
		}
	} //sortLocations
	
	sortServices := func (sers ServiceListT) {
		for ci := 1; ci < len(sers); ci++ {
			for cj := ci; cj > 0 && ST.CompP(sers[cj].Category, sers[cj - 1].Category) == ST.Lt; cj--{
				sers[cj], sers[cj - 1] = sers[cj - 1], sers[cj]
			}
		}
		for ci := 0; ci < len(sers); ci++ {
			sp := sers[ci].Specialities
			for si := 1; si < len(sp); si++ {
				for sj := si; sj > 0 && ST.CompP(sp[sj], sp[sj - 1]) == ST.Lt; sj-- {
					sp[sj], sp[sj - 1] = sp[sj - 1], sp[sj]
				}
			}
		}
	} //sortServices
	
	//insertId
	sortLocations(id.Locations)
	sortServices(id.Supplied_services)
	sortServices(id.Searched_services)
	
	wId := identityIndex.Writer()
	sId := &B.String{id.Nickname}
	b := wId.SearchIns(sId); M.Assert(!b, id.Nickname, 21)
	idRef := identityMan.WriteAllocateData(id)
	wId.WriteValue(idRef)
	
	wCo := identitiesIndex.Writer()
	for _, co := range id.Locations {
		sCo := &B.String{co.Country}
		if !wCo.SearchIns(sCo) {
			wCo.WriteValue(base.CreateIndex(0))
		}
		iRg := base.OpenIndex(wCo.ReadValue(), stringKeyMan, stringFac)
		wRg := iRg.Writer()
		for _, rg := range co.Regions {
			sRg := &B.String{rg.Region}
			if !wRg.SearchIns(sRg) {
				wRg.WriteValue(base.CreateIndex(0))
			}
			iCy := base.OpenIndex(wRg.ReadValue(), stringKeyMan, stringFac)
			wCy := iCy.Writer()
			var lCy []string
			if len(rg.Cities) == 0 {
				lCy = []string{""}
			} else {
				lCy = rg.Cities
			}
			for _, cy := range lCy {
				sCy := &B.String{cy}
				if !wCy.SearchIns(sCy) {
					sd := new(serviceDispatchT)
					sd.supplied = base.CreateIndex(0)
					sd.searched = base.CreateIndex(0)
					wCy.WriteValue(serviceDispatchMan.WriteAllocateData(sd))
				}
				sd := serviceDispatchMan.ReadData(wCy.ReadValue()).(*serviceDispatchT)
				insertService(sd.supplied, idRef, sId, id.Supplied_services)
				insertService(sd.searched, idRef, sId, id.Searched_services)
			}
		}
	}
	
	return id
} //insertId

func insertIdentity (id *IdentityT) (*IdentityT, M.Nerror) {
	if id, err := verifyIdBeforeInsertion(id, false); err != nil {
		return id, err
	}
	return insertId(id), nil
} //insertIdentity

func InsertIdentity (id *IdentityT) (*IdentityT, M.Nerror) {
	defer base.UpdateBase()
	return insertIdentity(id)
} //InsertIdentity

func InsertIdentities (l IdentityListT) M.Nerror {
	M.Assert(l != nil, 20)
	for _, id := range l {
		_, err := verifyIdBeforeInsertion(id, false)
		if err != nil {
			return err
		}
	}
	defer base.UpdateBase()
	for _, id := range l {
		insertId(id)
	}
	return nil
} //InsertIdentities

func removeLocation (country, region string) M.Nerror {
	
	testEmpty := func () bool {
		if country == "" {
			return identitiesIndex.IsEmpty()
		}
		rCo := identitiesIndex.NewReader()
		sCo := &B.String{country}
		if !rCo.Search(sCo) {
			return true
		}
		iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
		if region == "" {
			return iRg.IsEmpty()
		}
		rRg := iRg.NewReader()
		sRg := &B.String{region}
		return !rRg.Search(sRg)
	} //testEmpty
	
	removeLoc := func (sCo, sRg *B.String) {
		wCo := location1Index.Writer()
		b := wCo.Search(sCo); M.Assert(b, sCo.C, 100)
		iRgRef := wCo.ReadValue()
		iRg := base.OpenIndex(iRgRef, stringKeyMan, stringFac)
		wRg := iRg.Writer()
		b = wRg.Erase(sRg); M.Assert(b, sRg.C, 101)
		if iRg.IsEmpty() {
			b := wCo.Erase(sCo); M.Assert(b, sCo.C, 102)
			base.DeleteIndex(iRgRef)
		}
	} //removeLoc
	
	//removeLocation
	if country == "" && region != "" {
		return GetError(8, "removeLocation", region)
	}
	if !locationExists(country, region) {
		return GetError(2, "removeLocation", country, region)
	}
	if !testEmpty() {
		return GetError(10, "removeLocation", country, region)
	}
	rCo := location1Index.NewReader()
	if country == "" {
		for rCo.Next(); rCo.PosSet(); {
			sCo := rCo.CurrentKey().(*B.String)
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rCo.Next()
			rRg := iRg.NewReader()
			for rRg.Next(); rRg.PosSet(); {
				sRg := rRg.CurrentKey().(*B.String)
				rRg.Next()
				removeLoc(sCo, sRg)
			}
		}
	}  else {
		sCo := &B.String{country}
		b := rCo.Search(sCo); M.Assert(b, country, 100)
		iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
		rRg := iRg.NewReader()
		if region == "" {
			for rRg.Next(); rRg.PosSet(); {
				sRg := rRg.CurrentKey().(*B.String)
				rRg.Next()
				removeLoc(sCo, sRg)
			}
		} else {
			sRg := &B.String{region}
			removeLoc(sCo, sRg)
		}
	}
	return nil
} //removeLocation

func RemoveLocation (country, region string) M.Nerror {
	defer base.UpdateBase()
	return removeLocation (country, region)
} //RemoveLocation

func removeService (category, speciality string) M.Nerror {
	
	testEmpty := func () bool {
		
		isEmpty := func (pos B.FilePos) bool {
			iCa := base.OpenIndex(pos, stringKeyMan, stringFac)
			if category == "" {
				return iCa.IsEmpty()
			}
			rCa := iCa.NewReader()
			sCa := &B.String{category}
			if !rCa.Search(sCa) {
				return true
			}
			iSp := base.OpenIndex(rCa.ReadValue(), stringKeyMan, stringFac)
			if speciality == "" {
				return iSp.IsEmpty()
			}
			rSp := iSp.NewReader()
			sSp := &B.String{speciality}
			return !rSp.Search(sSp)
		} //isEmpty
		
		//testEmpty
		rCo := identitiesIndex.NewReader()
		for rCo.Next(); rCo.PosSet(); rCo.Next() {
			iRg := base.OpenIndex(rCo.ReadValue(), stringKeyMan, stringFac)
			rRg := iRg.NewReader()
			for rRg.Next(); rRg.PosSet(); rRg.Next() {
				iCy := base.OpenIndex(rRg.ReadValue(), stringKeyMan, stringFac)
				rCy := iCy.NewReader()
				for rCy.Next(); rCy.PosSet(); rCy.Next() {
					sd := serviceDispatchMan.ReadData(rCy.ReadValue()).(*serviceDispatchT)
					if !(isEmpty(sd.supplied) && isEmpty(sd.searched)) {
						return false
					}
				}
			}
		}
		return true
	} //testEmpty
	
	removeSpec := func (sCa, sSp *B.String) {
		wCa := service1Index.Writer()
		b := wCa.Search(sCa); M.Assert(b, sCa.C, 100)
		iSpRef := wCa.ReadValue()
		iSp := base.OpenIndex(iSpRef, stringKeyMan, stringFac)
		wSp := iSp.Writer()
		b = wSp.Erase(sSp); M.Assert(b, sSp.C, 101)
		if iSp.IsEmpty() {
			b := wCa.Erase(sCa); M.Assert(b, sCa.C, 102)
			base.DeleteIndex(iSpRef)
		}
	} //removeSpec
	
	//removeService
	if category == "" && speciality != "" {
		return GetError(11, "removeService", speciality)
	}
	if !serviceExists(category, speciality) {
		return GetError(15, "removeService", category, speciality)
	}
	if !testEmpty() {
		return GetError(17, "removeService", category, speciality)
	}
	rCa := service1Index.NewReader()
	if category == "" {
		for rCa.Next(); rCa.PosSet(); {
			sCa := rCa.CurrentKey().(*B.String)
			iSp := base.OpenIndex(rCa.ReadValue(), stringKeyMan, stringFac)
			rCa.Next()
			rSp := iSp.NewReader()
			for rSp.Next(); rSp.PosSet(); {
				sSp := rSp.CurrentKey().(*B.String)
				rSp.Next()
				removeSpec(sCa, sSp)
			}
		}
	}  else {
		sCa := &B.String{category}
		b := rCa.Search(sCa); M.Assert(b, category, 100)
		iSp := base.OpenIndex(rCa.ReadValue(), stringKeyMan, stringFac)
		rSp := iSp.NewReader()
		if speciality == "" {
			for rSp.Next(); rSp.PosSet(); {
				sSp := rSp.CurrentKey().(*B.String)
				rSp.Next()
				removeSpec(sCa, sSp)
			}
		} else {
			sSp := &B.String{speciality}
			removeSpec(sCa, sSp)
		}
	}
	return nil
} //removeService

func RemoveService (category, speciality string) M.Nerror {
	defer base.UpdateBase()
	return removeService(category, speciality)
} //RemoveService

func removeIdentity (nickname string) (*IdentityT, M.Nerror) {
	
	M.Assert(nickname != "", 100)
	sId := &B.String{nickname}
	wId := identityIndex.Writer()
	if !wId.Search(sId) {
		return nil, GetError(14, "removeIdentity", nickname)
	}
	idRef := wId.ReadValue()
	id := identityMan.ReadData(idRef).(*IdentityT)
	wId.Erase(sId)
	identityMan.EraseData(idRef)
	
	wCo := identitiesIndex.Writer()
	for _, co := range id.Locations {
		sCo := &B.String{co.Country}
		b := wCo.Search(sCo); M.Assert(b, "country =", co.Country, 102)
		iRgRef := wCo.ReadValue()
		iRg := base.OpenIndex(iRgRef, stringKeyMan, stringFac)
		wRg := iRg.Writer()
		for _, rg := range co.Regions {
			sRg := &B.String{rg.Region}
			b := wRg.Search(sRg); M.Assert(b, "region =", rg.Region, 103)
			iCyRef := wRg.ReadValue()
			iCy := base.OpenIndex(iCyRef, stringKeyMan, stringFac)
			wCy := iCy.Writer()
			var lCy []string
			if len(rg.Cities) == 0 {
				lCy = []string{""}
			} else {
				lCy = rg.Cities
			}
			for _, cy := range lCy {
				sCy := &B.String{cy}
				b := wCy.Search(sCy); M.Assert(b, "city =", cy, 104)
				sdRef := wCy.ReadValue()
				sd := serviceDispatchMan.ReadData(sdRef).(*serviceDispatchT)
				var (iCaRef B.FilePos; se ServiceListT; empties [2]bool)
				for i := 0; i < 2; i++ {
					if i == 1 {
						se = id.Searched_services
						iCaRef = sd.searched
					} else {
						se = id.Supplied_services
						iCaRef = sd.supplied
					}
					iCa := base.OpenIndex(iCaRef, stringKeyMan, stringFac)
					wCa := iCa.Writer()
					for _, ca := range se {
						sCa := &B.String{ca.Category}
						b := wCa.Search(sCa); M.Assert(b, "category =", ca.Category, 105)
						iSpRef := wCa.ReadValue()
						iSp := base.OpenIndex(iSpRef, stringKeyMan, stringFac)
						wSp := iSp.Writer()
						for _, sp := range ca.Specialities {
							sSp := &B.String{sp}
							b := wSp.Search(sSp); M.Assert(b, "speciality =", sp, 106)
							iIdsRef := wSp.ReadValue()
							iIds := base.OpenIndex(iIdsRef, stringKeyMan, stringFac)
							wIds := iIds.Writer()
							b = wIds.Erase(sId); M.Assert(b, 107)
							if iIds.IsEmpty() {
								base.DeleteIndex(iIdsRef)
								b := wSp.Erase(sSp); M.Assert(b, 108)
							}
						}
						if iSp.IsEmpty() {
							base.DeleteIndex(iSpRef)
							b := wCa.Erase(sCa); M.Assert(b, 109)
						}
					}
					if iCa.IsEmpty() {
						empties[i] = true
					}
				}
				if empties[0] && empties[1] {
					base.DeleteIndex(sd.supplied)
					base.DeleteIndex(sd.searched)
					serviceDispatchMan.EraseData(sdRef)
					b := wCy.Erase(sCy); M.Assert(b, 110)
				}
			}
			if iCy.IsEmpty() {
				base.DeleteIndex(iCyRef)
				b := wRg.Erase(sRg); M.Assert(b, 111)
			}
		}
		if iRg.IsEmpty() {
			base.DeleteIndex(iRgRef)
			b := wCo.Erase(sCo); M.Assert(b, 112)
		}
	}
	return id, nil
} //removeIdentity

func RemoveIdentity (nickname string) (*IdentityT, M.Nerror) {
	defer base.UpdateBase()
	return removeIdentity(nickname)
} //RemoveIdentity

func changeId (nickname string, id, newId *IdentityT) (*IdentityT, M.Nerror) {
	_, err := removeIdentity(nickname); M.Assert(err == nil, err, 100)
	if _, err := insertIdentity(newId); err != nil {
		_, err1 := insertIdentity(id); M.Assert(err1 == nil, err1, 101)
		return nil, err
	}
	return id, nil
} //changeId

func ChangeIdentity (nickname string, newId *IdentityT) (*IdentityT, M.Nerror) {
	M.Assert(newId != nil, 20)
	M.Assert(nickname != "", 21)
	id, err := IdentityOf(nickname)
	if err != nil {
		return nil, GetError(14, "changeIdentity", nickname)
	}
	b := newId.Nickname != nickname
	if b {
		_, err := IdentityOf(newId.Nickname)
		b = err == nil
	}
	if b {
		return nil, GetError(1, "changeIdentity", newId.Nickname)
	}
	
	defer base.UpdateBase()
	return changeId(nickname, id, newId)
} //ChangeIdentity

func ChangeIdentities (l IdentityListT) M.Nerror {
	M.Assert(l != nil, 20)
	for _, id := range l {
		_, err := verifyIdBeforeInsertion(id, true)
		if err != nil {
			return err
		}
	}
	defer base.UpdateBase()
	for _, id := range l {
		_, err := removeIdentity(id.Nickname); M.Assert(err == nil, err, 100)
		_, err = insertIdentity(id); M.Assert(err == nil, err, 101)
	}
	return nil
} //InsertIdentities

func (l1 *Location1T) Compare (l2 *Location1T) A.Comp {
	if l1.Country < l2.Country {
		return A.Lt
	}
	if l1.Country > l2.Country {
		return A.Gt
	}
	if l1.Region < l2.Region {
		return A.Lt
	}
	if l1.Region > l2.Region {
		return A.Gt
	}
	if l1.City < l2.City {
		return A.Lt
	}
	if l1.City > l2.City {
		return A.Gt
	}
	return A.Eq
} //Compare

func locationListToLocation1Tree (l LocationListT) *A.TreeS[*Location1T] {
	t := A.NewS[*Location1T]()
	for _, co := range l {
		for _, rg := range co.Regions {
			if len(rg.Cities) == 0 {
				t.SearchIns(&Location1T{Country: co.Country, Region: rg.Region})
			} else {
				for _, cy := range rg.Cities {
					t.SearchIns(&Location1T{Country: co.Country, Region: rg.Region, City: cy})
				}
			}
		}
	}
	return t
} //locationListToLocation1Tree

func location1TreeToLocationList (t *A.TreeS[*Location1T]) LocationListT {
	var (oldCo, oldRg string; iRg int)
	cos := make(LocationListT, 0)
	iCo := -1
	for e := t.Next(nil); e != nil; e = t.Next(e) {
		loc := e.Val()
		if loc.Country != oldCo {
			oldCo = loc.Country
			oldRg = ""
			iRg = -1
			rg := make([]*RegionT, 0)
			cos = append(cos, &LocationT{Country: loc.Country, Regions: rg})
			iCo++
		}
		if loc.Region != oldRg {
			oldRg = loc.Region
			cy := make([]string, 0)
			cos[iCo].Regions = append(cos[iCo].Regions, &RegionT{Region: loc.Region, Cities: cy})
			iRg++
		}
		if loc.City != "" {
			cos[iCo].Regions[iRg].Cities = append(cos[iCo].Regions[iRg].Cities, loc.City)
		}
	}
	return cos
} //location1TreeToLocationList

func RenameLocation (location int, country, region, city, newName string) M.Nerror {
	M.Assert(location >= 0 && location < locationsNb, "location =", location, 20)
	if country == "" {
		return GetError(4, "renameLocation")
	}
	switch location {
	case Country:
	case Region:
		if region == "" {
			return GetError(5, "renameLocation")
		}
	case City:
		if city == "" {
			return GetError(18, "renameLocation")
		}
	}
	if region == "" && city != "" {
		return GetError(9, "renameLocation", city)
	}
	if !locationExists(country, region) {
		return GetError(2, "renameLocation", country, region)
	}
	if newName == "" {
		return GetError(12, "renameLocation")
	}
	
	defer base.UpdateBase()
	if location == Country && newName != country || location == Region && newName != region || location == City && newName != city {
		l, err := IdentityList(country, region, city, nil, BothServices); M.Assert(err == nil, 100)
		for _, id := range l {
			t := locationListToLocation1Tree(id.Locations)
			var (e *A.Elem[*Location1T]; rank int; erased = make([]*Location1T, 0))
			if region == "" {
				e, _, rank = t.SearchNext(&Location1T{Country: country})
				for e != nil {
					loc := e.Val()
					if loc.Country != country {break}
					erased = append(erased, loc)
					ee := t.Next(e)
					t.Erase(rank)
					e = ee
				}
			} else if city == "" {
				e, _, rank = t.SearchNext(&Location1T{Country: country, Region: region})
				for e != nil {
					loc := e.Val()
					if loc.Country != country || loc.Region != region {break}
					erased = append(erased, loc)
					ee := t.Next(e)
					t.Erase(rank)
					e = ee
				}
			} else {
				loc := &Location1T{Country: country, Region: region, City: city}
				if t.Delete(loc) {
					erased = append(erased, loc)
				}
			}
			for _, loc := range erased {
				switch location {
				case Country:
					loc.Country = newName
				case Region:
					loc.Region = newName
				case City:
					loc.City = newName
				}
				t.SearchIns(loc)
			}
			id.Locations = location1TreeToLocationList(t)
			
			_, err := removeIdentity(id.Nickname); M.Assert(err == nil, err, 101)
		}
		if location == Country && newName != country || location == Region && newName != region {
			var list LocationShortListT
			if region == "" {
				list, err = AllLocations(country); M.Assert(err == nil, err, 102)
			}
			err := removeLocation(country, region); M.Assert(err == nil, err, 103)
			if region == "" {
				M.Assert(len(list) == 1, list, 104)
				for _, rg := range list[0].Regions {
					_, err := insertLocation(newName, rg); M.Assert(err == nil, err, 105)
				}
			} else {
				if location == Country {
					_, err := insertLocation(newName, region); M.Assert(err == nil, err, 106)
				} else {
					_, err := insertLocation(country, newName); M.Assert(err == nil || err.ErrorNum() == 3, err, 107)
				}
			}
		}
		for _, id := range l {
			_, err = insertIdentity(id); M.Assert(err == nil, err, 108)
		}
	}
	return nil
} //RenameLocation

func serviceListToService1Tree (l ServiceListT) *A.TreeS[*Service1T] {
	t := A.NewS[*Service1T]()
	for _, ca := range l {
		for _, sp := range ca.Specialities {
			t.SearchIns(&Service1T{Category: ca.Category, Speciality: sp})
		}
	}
	return t
} //serviceListToService1Tree

func service1TreeToServiceList (t *A.TreeS[*Service1T]) ServiceListT {
	var (oldCa string)
	cas := make(ServiceListT, 0)
	iCa := -1
	for e := t.Next(nil); e != nil; e = t.Next(e) {
		ser := e.Val()
		if ser.Category != oldCa {
			oldCa = ser.Category
			sp := make([]string, 0)
			cas = append(cas, &ServiceT{Category: ser.Category, Specialities: sp})
			iCa++
		}
		cas[iCa].Specialities = append(cas[iCa].Specialities, ser.Speciality)
	}
	return cas
} //service1TreeToServiceList

func RenameService (service int, category, speciality, newName string) M.Nerror {
	
	changeServices := func (s *ServiceListT) {
		t := serviceListToService1Tree(*s)
			var (e *A.Elem[*Service1T]; rank int; erased = make([]*Service1T, 0))
			if speciality == "" {
				e, _, rank = t.SearchNext(&Service1T{Category: category})
				for e != nil {
					ser := e.Val()
					if ser.Category != category {break}
					erased = append(erased, ser)
					ee := t.Next(e)
					t.Erase(rank)
					e = ee
				}
			} else {
				ser := &Service1T{Category: category, Speciality: speciality}
				if t.Delete(ser) {
					erased = append(erased, ser)
				}
			}
			for _, ser := range erased {
				if service == Category {
					ser.Category = newName
				} else {
					ser.Speciality = newName
				}
				t.SearchIns(ser)
			}
		*s = service1TreeToServiceList(t)
	} //changeServices
	
	//RenameService
	M.Assert(service >= 0 && service < servicesNb, "service =", service, 20)
	if category == "" {
		return GetError(6, "renameService")
	}
	if service == Speciality && speciality == "" {
		return GetError(7, "renameService")
	}
	if !serviceExists(category, speciality) {
		return GetError(15, "renameService", category, speciality)
	}
	if newName == "" {
		return GetError(12, "renameService")
	}
	
	defer base.UpdateBase()
	if service == Category && newName != category || service == Speciality && newName != speciality {
		var specs ServiceListT
		if speciality == "" {
			specs = ServiceListT{&ServiceT{Category: category}}
		} else {
			specs = ServiceListT{&ServiceT{Category: category, Specialities: []string{speciality}}}
		}
		l, err := IdentityList("", "", "", specs, BothServices); M.Assert(err == nil, 100)
		for _, id := range l {
			changeServices(&id.Supplied_services)
			changeServices(&id.Searched_services)
			_, err := removeIdentity(id.Nickname); M.Assert(err == nil, id.Nickname, 101)
		}
		var list ServiceNoPListT
		if speciality == "" {
			list, err = AllServices(category); M.Assert(err == nil, err, 102)
		}
		err = removeService(category, speciality); M.Assert(err == nil, err, 103)
		if speciality == "" {
			M.Assert(len(list) == 1, list, 104)
			for _, sp := range list[0].Specialities {
				_, err := insertService(newName, sp); M.Assert(err == nil, err, 105)
			}
		} else {
			if service == Category {
				_, err := insertService(newName, speciality); M.Assert(err == nil, 106)
			} else {
				_, err := insertService(category, newName); M.Assert(err == nil || err.ErrorNum() == 16, 107)
			}
		}
		for _, id := range l {
			_, err = insertIdentity(id); M.Assert(err == nil, id.Nickname, 108)
		}
	}
	return nil
} //RenameService

func init () {
	saveBase()
}
