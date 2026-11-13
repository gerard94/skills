/* 
Skills

Copyright (C) 2023 GérardMeunier

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License  for more details.

You should have received a copy of the GNU General Public License along with this program; if not, write to the Free Software Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
*/

package static

const (
	
	typeSystem = `
# Quelle sont les compétences (membres) de ma région ?
# Qui a telles compétences dans ma région ou au national ou à l'international ?
# Quelles compétences a ce membre ?

type Identity {
	nickname: String!
	password: String! #crypted
	g1_nickname: String
	pubkey: Pubkey!
	email: String
	locations: [Location!]!
	transport: Boolean!
	pickup: Boolean!
	supplied_services: [Service!]!
	searched_services: [Service!]!
	additional_informations: String
}

input Identity_input {
	nickname: String!
	password: String! #crypted
	g1_nickname: String
	pubkey: Pubkey! = ""
	email: String
	locations: [Location_input!]! = []
	transport: Boolean! = false
	pickup: Boolean! = false
	supplied_services: [Service_input!]! = []
	searched_services: [Service_input!]! = []
	additional_informations: String
}

enum SearchedOrSupplied {
	SUPPLIED_SERVICES
	SEARCHED_SERVICES
	BOTH_SERVICES
}

enum Locations {
	COUNTRY
	REGION
	CITY
}

enum Services {
	CATEGORY
	SPECIALITY
}

type Location {
	country: String!
	regions: [Region!]!
}

type Region {
	region: String!
	cities: [String!]!
}

type Location1Short {
	country: String!
	region: String!
}

input Location1Short_input {
	country: String!
	region: String!
}

type Location1 {
	country: String!
	region: String!
	city: String
}

type LocationShort {
	country: String!
	regions: [String!]!
}

input Location_input {
	country: String!
	regions: [Region_input!]!
}

input Region_input {
	region: String!
	cities: [String!]!
}

type Service {
	category: String!
	specialities: [String!]!
}

type Service1 {
	category: String!
	speciality: String!
}

input Service1_input {
	category: String!
	speciality: String!
}

input Service_input {
	category: String!
	specialities: [String!]!
}

input ServiceQ_input {
	category: String!
	specialities: [String!]
}

# If a String is said null, it means that the String is null or empty, and vice versa; a string field marked as non-null must also be non-empty.

type Query {
	
	"'allLocations' returns the list of all available locations in 'country'. It raises an error and returns null if 'country' doesn't exist in the base. If 'country' is null, the list returned is for all countries."
	allLocations (country: String): [LocationShort!]
	
	"'locationList' returns the list of used locations. If 'region' is null and 'country' is not null, the list returned is for all regions in country. If 'region' is null and 'country' is null, the list returned is for all countries. If 'region' is not null and 'country' is null, an error is raised and 'locationList' returns null. If 'short' is not null and true, only the fields 'country' and 'region' are returned in the result and 'city' is always empty"
	locationList (country: String, region: String, short: Boolean): [Location1!]
	
	"'allServices' returns the list of all available services in 'category'. It raises an error and returns null if 'category' doesn't exist in the base. If 'category' is null, the list returned is for all categories"
	allServices (category: String): [Service!]
	
	"'serviceList' returns the list of used services in the 'city' of the  'region' of 'country', and in the 'category'. If 'region' is null and 'country' is not null, the list returned is for all regions in country. If 'region' is null and 'country' is null, the list returned is for all countries. If 'region' is not null and 'country' is null, an error is raised. The same for 'region' and 'city'. If 'category' is null, the list is returned for all categories. If 'searchedOrSupplied' is SUPPLIED_SERVICES, the returned value is for supplied services, or for searched services if 'searchedOrSupplied' is SEARCHED_SERVICES, or for both services if 'searchedOrSupplied' is BOTH_SERVICES. The returned value is null if an error has occured"
	serviceList (country: String, region: String, city: String, category: String, searchedOrSupplied: SearchedOrSupplied! = SUPPLIED_SERVICES): [Service1!]
	
	"'identityOf' returns the identity whose nickname is 'nickname'. If the identity doesn't exist, returns null"
	identityOf (nickname: String!): Identity
	
	"'allIdentities' returns all identities, in the base, whose 'nickname' is a suffix of 'hint'"
	allIdentities(hint: String): [Identity!]!
	
	"'identityList' returns all identities in the 'city' of the 'region' of 'country' who supply and/or search at least one service in 'services'. The returned value is for supplied services if 'searchedOrSupplied' is SUPPLIED_SERVICES; otherwise, if 'searchedOrSupplied' is SEARCHED_SERVICES, it's for searched services, or both if 'searchedOrSupplied' is BOTH_SERVICES. If 'region' is null and 'country' is not null, the list returned is for all regions in country. If 'region' is null and 'country' is null, the list returned is for all countries. If 'region' is not null and 'country' is null, an error is raised. The same for 'region' and 'city'. If 'services' is null or empty, all identities in the location, with some service supplied, searched, or both, according to 'searchedOrSupplied', are returned. If 'services' is not null nor empty, no 'category' in it should be null, otherwise an error is raised; if a 'specialities' is null, all identities in the corresponding 'category' are returned. The value returned is null is an error has occured"
	identityList (country: String, region: String, city: String, services: [ServiceQ_input!], searchedOrSupplied: SearchedOrSupplied! = SUPPLIED_SERVICES): [Identity!]
	
	"version returns the version of the skills server"
	version: String!

}

type Mutation {

	"'insertLocation' creates 'country' if it doesn't exist and inserts 'region' into 'country if it doesn't exist; it returns the inserted 'Location1Short'; if region already exists in 'country', an error is raised"
	insertLocation (country: String!, region: String!): Location1Short!
	
	"'insertLocations' removes all locations in the base and inserts into it the new ones in 'locations'; all locations in all identities must be included in the set of new ones, otherwise an error is raised"
	insertLocations (locations: [Location1Short_input!]!): Boolean!
	
	"'insertService' creates 'category' if it doesn't exist and inserts 'speciality' into 'category'; it returns the inserted 'Service1'; 'category' and 'speciality' must not be null; if 'speciality' already exists in 'category', an error is raised"
	insertService (category: String!, speciality: String!): Service1!
	
	"'insertServices' removes all services in the base and inserts into it the new ones in 'services'; all services in all identities must be included in the set of new ones, otherwise an error is raised"
	insertServices (services: [Service1_input!]!): Boolean!
	
	"'insertIdentity' inserts 'id' into the base and returns its identity; if the 'nickname' of 'id' already exists in the base, an error is raised and 'insertIdentity' returns null; in 'id.locations', all 'country' / 'region' pairs must already exist in the base, and the same for 'category' / 'speciality' pairs in 'id.supplied_services' and 'id.searched_services'"
	insertIdentity (id: Identity_input!): Identity
	
	"'insertIdentities' inserts all identities in 'ids' into the base and returns themselves; if the 'nickname' of any identity in 'ids' already exists in the base, an error is raised and 'insertIdentities' returns null; in 'id.locations', all 'country' / 'region' pairs must already exist in the base, and the same for 'category' / 'speciality' pairs in 'id.supplied_services' and 'id.searched_services'"
	insertIdentities (ids: [Identity_input!]!): [Identity!]
	
	"'removeLocation' removes all locations if 'country' and 'region' are null, or else removes all regions out of 'country' if 'region' is null, or else removes 'region' out of 'country'(if 'country' becomes empty, it is removed too); there must be no identity in the removed locations, otherwise an error is raised; 'removeLocation' raises also an error if 'country' is null and 'region' is not null, or if 'country' doesn't exist or if 'region' doesn't exist in 'country', or if there are identities in the locations removed, and returns false in this case, otherwise, it returns true"
	removeLocation (country: String, region: String): Boolean!
	
	"'removeService' removes all services if 'category' and 'speciality' are null, or else removes all specialities out of 'category' if 'speciality' is null, or else removes 'speciality' out of 'category' (if 'category' becomes empty, it is removed too); there must be no identity providing or searching the removed services, otherwise an error is raised; 'removeService' raises also an error if 'category' is null and 'speciality' is not null, or if if 'category' doesn't exist or if 'speciality' doesn't exist in 'category', and returns false in this case, otherwise, it returns true"
	removeService (category: String, speciality: String): Boolean!
	
	"'renameLocation' renames 'country', 'region' or 'city' to 'newName', according to the value of 'location' and only for the values of 'country', 'region' and 'city' which are not null nor void; it returns false and raises an error if 'country' is void or doesn't exist, or if 'region' doesn't exist in 'country', or if 'region' is null or void and 'city' is not null nor void"
	renameLocation (location: Locations!, country: String!, region:String, city: String, newName: String!): Boolean!
	
	"'renameService' renames 'category' or 'speciality' to 'newName', according to the value of 'service' and only for the values of 'category' and 'speciality' which are not null nor void; it returns false and raises an error if 'category' is void or doesn't exist, or if 'speciality' doesn't exist in 'category'"
	renameService (service: Services!, category: String!, speciality: String, newName: String!): Boolean!
	
	"'removeIdentity' removes the identity whose nickname is 'nickname'; if this identity doesn't exist, it raises an error and returns null, otherwise, it returns the identity; If this identity was the last one in a city, this city is removed too"
	removeIdentity (nickname: String!): Identity
	
	"'changeIdentity' replaces the identity whose nickname is 'nickname' by 'newId'; if this identity doesn't exist or if the nickname of 'newId' is different from 'nickname' and already exists in the base, it raises an error and returns null, otherwise, it returns the old identtity"
	changeIdentity (nickname: String!, newId: Identity_input!): Identity
	
	"For every identities 'id' in 'ids', 'changeIdentities' replaces, in the base, the identity whose nickname is 'id.nickname' with 'id' and returns the list of previous versions; if the 'nickname' of any identity in 'ids' doesn't exist in the base, an error is raised and 'changeIdentities' returns null; in 'id.locations', all 'country' / 'region' pairs must already exist in the base, and the same for 'category' / 'speciality' pairs in 'id.supplied_services' and 'id.searched_services'"
	changeIdentities (ids: [Identity_input!]!): [Identity!]

}

"Avatar of String"
scalar Pubkey @specifiedBy(url: "https://github.com/nodes/common/doc/blob/master/rfc/0009_Duniter_Blockchain_Protocol_V11.md#public-key")
`

)
