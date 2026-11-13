package base

import (
	
	/*
	M	"github.com/gerard94/util/misc"
	*/
	//"github.com/gerard94/util/alea"
	"testing"
	"fmt"
		"os"
	
)

var (
	
	locations = Location1ShortListT{
		
		{
			"France", "Seine",
		},
		{
			"Espagne", "Catalogne",
		},
		{
			"France", "Val de Marne",
		},
		{
			"France", "Seine et Marne",
		},
		{
			"", "Vaucluse",
		},
		{
			"Espagne", "",
		},
		
	} //locations
	
	services = Service1ListT {
		
		{
			"Fabrication", "Impression 3D",
		},
		{
			"Santé", "Sobriété",
		},
		{
			"Santé", "Plantes",
		},
		{
			"Sport", "Rugby",
		},
		{
			"", "Marché",
		},
		{
			"Électronique", "Régulation",
		},
		{
			"Rencontres", "",
		},
		{
			"Informatique", "Programmation",
		},
		{
			"Informatique", "Webmaster",
		},
		{
			"Alimentation", "Légumes",
		},
		{
			"Alimentation", "Pain",
		},
		{
			"Alimentation", "Pain",
		},
		{
			"Relations sociales", "Conférencier",
		},
		{
			"Hygiène", "Savon",
		},
	
	} //services
	
	identities = IdentityListT{
		
		{
			Nickname: "Anne-Marie",
			Password: "un",
			G1_nickname:"gerard94",
			Pubkey: "CRBxCJrTA6tmHsgt9cQh9SHcCc8w8q95YTp38CPHx2Uk",
			Email: "gemeu@free.fr",
			Locations: LocationListT{
				&LocationT{
					Country: "France",
					Regions: []*RegionT{
						&RegionT{
							Region: "Val de Marne",
							Cities: []string{"Champigny"},
						},
						&RegionT{
							Region: "Seine",
							Cities: []string{"Paris"},
						},
					},
				},
			},
			Transport: false,
			Pickup: true,
			Supplied_services: ServiceListT{
				&ServiceT{
					Category: "Informatique",
					Specialities: []string{"Programmation", "Webmaster"},
				},
			},
			Searched_services: ServiceListT{
				&ServiceT{
					Category: "Alimentation",
					Specialities: []string{"Légumes", "Pain"},
				},
				&ServiceT{
					Category: "Hygiène",
					Specialities: []string{"Savon"},
				},
			},
			Additional_informations: "Habite à Champigny",
		},
		
		{
			Nickname: "Geronimo",
			Password: "deux",
			G1_nickname:"gero",
			Pubkey: "jRRFieXn4F2i4gqBe8Tj111kBGqdZu2iBeBAXJgbj4n",
			Email: "gero@free.fr",
			Locations: LocationListT{
				&LocationT{
					Country: "France",
					Regions: []*RegionT{
						&RegionT{
							Region: "Seine et Marne",
							Cities: []string{"Lagny"},
						},
						&RegionT{
							Region: "Val de Marne",
							Cities: []string{"Champigny sur Marne"},
						},
						&RegionT{
							Region: "Seine",
							Cities: []string{"Paris"},
						},
					},
				},
			},
			Transport: true,
			Pickup: true,
			Supplied_services: ServiceListT{
				&ServiceT{
					Category: "Électronique",
					Specialities: []string{"Régulation"},
				},
				&ServiceT{
					Category: "Relations sociales",
					Specialities: []string{"Conférencier"},
				},
			},
			Searched_services: ServiceListT{
				&ServiceT{
					Category: "Fabrication",
					Specialities: []string{"Impression 3D"},
				},
				&ServiceT{
					Category: "Hygiène",
					Specialities: []string{"Savon"},
				},
			},
			Additional_informations: "Électronicien",
		},
	
	} //identities

)

func insertLocations (l Location1ShortListT) {
	for _, loc := range l {
		ll, err := InsertLocation(loc.Country, loc.Region)
		if err != nil {
			fmt.Println(err.Error(), err.ErrorNum())
		}
		if ll != nil {
			printLocation1Short(ll)
			fmt.Println()
		}
	}
}

func insertServices (l Service1ListT) {
	for _, s := range l {
		ss, err := InsertService(s.Category, s.Speciality)
		if err != nil {
			fmt.Println(err.Error(), err.ErrorNum())
		}
		if ss != nil {
			printService1(os.Stdout, ss)
			fmt.Println()
		}
	}
}

func insertIdentities (l IdentityListT) {
	for _, id := range l {
		id1, err := InsertIdentity(id)
		if err != nil {
			fmt.Println(err.Error(), err.ErrorNum())
		}
		if id1 != nil {
			printIdentity(os.Stdout, id1)
			fmt.Fprintln(os.Stdout)
		}
	}
}

func createBase () {
	fmt.Println()
	fmt.Println("locations")
	insertLocations(locations)
	
	fmt.Println()
	fmt.Println("services")
	insertServices(services)
	
	fmt.Println()
	fmt.Println("identities")
	insertIdentities(identities)
	}

func TestCreate (t *testing.T) {
	fmt.Println()
	fmt.Println("TestCreate")
	
	os.Remove(basePath)
	OpenBase()
	createBase()
	fmt.Println()
	fmt.Println("Identity list:")
	for _, id := range AllIdentities("") {
		fmt.Println(id.Nickname)
	}
	fmt.Println()
	fmt.Println("Identity list(a):")
	for _, id := range AllIdentities("a") {
		fmt.Println(id.Nickname)
	}
	fmt.Println()
	fmt.Println("Identity list(g):")
	for _, id := range AllIdentities("g") {
		fmt.Println(id.Nickname)
	}
	CloseBase()
}

func TestOpen (t *testing.T) {
	fmt.Println()
	fmt.Println("TestOpen")
	
	OpenBase()
}

func allLocations (co string) {
	l, err := AllLocations(co)
	if l != nil {
		printLocationShortList(os.Stdout, l)
		fmt.Println()
	}
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
}

func TestAllLocations (t *testing.T) {
	fmt.Println()
	fmt.Println("TestAllLocations")
	
	allLocations("France")
	allLocations("Espagne")
	allLocations("Belgique")
	allLocations("")
}

func allServices (ca string) {
	s, err := AllServices(ca)
	if s != nil {
		printServiceNoPList(os.Stdout, s)
		fmt.Println()
	}
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
}

func TestAllServices (t *testing.T) {
	fmt.Println()
	fmt.Println("TestAllServices")
	
	allServices("Santé")
	allServices("Sport")
	allServices("Loisirs")
	allServices("")
}

func serviceList (co, rg, cy, ca string, sc int8) {
	s, err := ServiceList(co, rg, cy, ca, sc)
	if s != nil {
		printService1List(os.Stdout, s)
		fmt.Println()
	}
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
}

func TestServiceList (t *testing.T) {
	fmt.Println()
	fmt.Println("TestServiceList")
	
	fmt.Println()
	serviceList("", "", "", "", SuppliedServices)
	fmt.Println()
	serviceList("", "", "", "", SearchedServices)
	fmt.Println()
	serviceList("", "", "", "", BothServices)
}

func locationList (co, rg string) {
	l, err := LocationList(co, rg, true)
	if l != nil {
		printLocation1List(os.Stdout, l)
		fmt.Println()
}
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
}

func TestLocationList (t *testing.T) {
	fmt.Println()
	fmt.Println("TestLocationList")
	
	locationList("", "")
}

func allIdentities (hint string) {
	ids := AllIdentities(hint)
	for _, id := range ids {
		printIdentity(os.Stdout, id)
		fmt.Println()
	}
}

func TestAllIdentities (t *testing.T) {
	fmt.Println()
	fmt.Println("TestAllIdentities")
	
	fmt.Println()
	allIdentities("")
}

func identityList (co, rg, cy string, ser ServiceListT, sc int8) {
	ids, err := IdentityList(co, rg, cy, ser, sc)
	if ids != nil {
		for _, id := range ids {
			printIdentity(os.Stdout, id)
			fmt.Println()
		}
	}
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
	fmt.Println()
}

func TestIdentityList (t *testing.T) {
	fmt.Println()
	fmt.Println("TestIdentityList")
	
	fmt.Println()
	identityList("", "", "", nil, SuppliedServices)
	fmt.Println()
	identityList("", "", "", nil, SearchedServices)
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	fmt.Println()
	identityList("France", "Seine et Marne", "", nil, BothServices)
	fmt.Println()
	identityList("France", "Seine", "", ServiceListT{{Category: "Informatique"}}, SuppliedServices)
	fmt.Println()
	identityList("France", "Seine et Marne", "", nil, SearchedServices)
}

func renameLocation (loc int, co, rg, cy, nN string) {
	fmt.Println(loc, co, rg, cy, "->", nN)
	err := RenameLocation(loc, co, rg, cy, nN)
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
}

func TestRenameLocation (t *testing.T) {
	fmt.Println()
	fmt.Println("TestRenameLocation")
	
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allLocations("")
	fmt.Println()
	renameLocation(City, "France", "Val de Marne", "Champigny", "Champigny sur Marne")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allLocations("")
	fmt.Println()
	renameLocation(Country, "France", "Seine et Marne", "", "France métropolitaine")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allLocations("")
}

func renameService (ser int, ca, sp, nN string) {
	fmt.Println(ser, ca, sp, "->", nN)
	err := RenameService(ser, ca, sp, nN)
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
}

func TestRenameService (t *testing.T) {
	fmt.Println()
	fmt.Println("TestRenameService")
	
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allServices("")
	fmt.Println()
	renameService(Category, "Hygiène", "", "Propreté")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allServices("")
	fmt.Println()
	renameService(Speciality, "Relations sociales", "Conférencier", "Administrateur")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allServices("")
	renameService(Category, "Alimentation", "Légumes", "Nourriture")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	allServices("")
}

func removeLocation1 (co, rg string) {
	err := RemoveLocation (co, rg)
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
	fmt.Println(co, rg)
}

func TestRemoveLocation (t *testing.T) {
	fmt.Println()
	fmt.Println("TestRemoveLocation")
	
	allLocations("")
	fmt.Println()
	removeLocation1("Brésil", "Mato Grosso")
	removeLocation1("France", "Landes")
	removeLocation1("Espagne", "Catalogne")
	fmt.Println()
	allLocations("")
}

func removeService1 (ca, sp string) {
	err := RemoveService(ca, sp)
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
	fmt.Println(ca, sp, "removed")
}

func TestRemoveService (t *testing.T) {
	fmt.Println()
	fmt.Println("TestRemoveService")
	
	allServices("")
	fmt.Println()
	removeService1("Transports", "Camion")
	removeService1("Santé", "Vaccination")
	removeService1("Sport", "Rugby")
	fmt.Println()
	allServices("")
	fmt.Println()
	removeService1("Santé", "Plantes")
	fmt.Println()
	allServices("")
	fmt.Println()
	removeService1("Précautions", "Sobriété")
	fmt.Println()
	allServices("")
}

/*
func changeIdentity (nN string, nI *IdentityT) {
	idO, err := ChangeIdentity(nN, nI)
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
	if idO != nil {
		fmt.Println(*idO)
		fmt.Println(nN, " changed")
	}
}

func TestChangeIdentity (t *testing.T) {
	fmt.Println()
	fmt.Println("TestChangeIdentity")
	
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	fmt.Println()
	changeIdentity("Anne-Marie", &IdentityT{Nickname: "Geronimo", Password: "fxfc"})
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	fmt.Println()
	changeIdentity("Anne-Marie", &IdentityT{Nickname: "Arthur", Password: "fxfc"})
	fmt.Println()
	identityList("", "", "", nil, BothServices)
}
*/

func removeIdentity1 (n string) {
	id, err := RemoveIdentity(n)
	if err != nil {
		fmt.Println(err, err.ErrorNum())
	}
	if id != nil {
		fmt.Println(*id)
		fmt.Println(n, " removed")
	}
}

func TestRemoveIdentity (t *testing.T) {
	fmt.Println()
	fmt.Println("TestRemoveIdentity")
	
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	fmt.Println()
	locationList("", "")
	fmt.Println()
	serviceList("", "", "", "", BothServices)
	fmt.Println()
	removeIdentity1("Arthur")
	fmt.Println()
	removeIdentity1("Anne-Marie")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	fmt.Println()
	locationList("", "")
	fmt.Println()
	serviceList("", "", "", "", BothServices)
	fmt.Println()
	removeIdentity1("Geronimo")
	fmt.Println()
	identityList("", "", "", nil, BothServices)
	fmt.Println()
	locationList("", "")
	fmt.Println()
	serviceList("", "", "", "", BothServices)
}

func TestClose (t *testing.T) {
	fmt.Println()
	fmt.Println("TestClose")
	
	CloseBase()
}

