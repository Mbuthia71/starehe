package models

import "time"

type Profile struct {
	UserID            string     `json:"user_id" db:"user_id"`
	FullName          string     `json:"full_name" db:"full_name"`
	Bio              *string    `json:"bio,omitempty" db:"bio"`
	AvatarURL        *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CoverURL         *string    `json:"cover_url,omitempty" db:"cover_url"`
	ClassYear       *int       `json:"class_year,omitempty" db:"class_year"`
	House            *House     `json:"house,omitempty" db:"house"`
	Career           *string    `json:"career,omitempty" db:"career"`
	Location         *string    `json:"location,omitempty" db:"location"`
	ProfileVisibility string    `json:"profile_visibility" db:"profile_visibility"`
	ContactVisibility string    `json:"contact_visibility" db:"contact_visibility"`
	CareerVisibility  string    `json:"career_visibility" db:"career_visibility"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type House string

const (
	HouseGikubu   House = "Gikubu"
	HouseNgala    House = "Ngala"
	HouseGeturo   House = "Geturo"
	HouseShaw     House = "Shaw"
	HouseHorsten  House = "Horsten"
	HouseMboya    House = "Mboya"
	HouseShell    House = "Shell"
	HouseChaka    House = "Chaka"
	HouseNjonjo   House = "Njonjo"
	HouseKirkley  House = "Kirkley"
	HouseMuriuki  House = "Muriuki"
	HouseKibaki   House = "Kibaki"
)

var HouseDescriptions = map[House]string{
	HouseGikubu:  "Named after co-founder Joseph Gikubu",
	HouseNgala:   "Named after the late Ronald Ngala, a cabinet minister in the Jomo Kenyatta government",
	HouseGeturo:  "Named after co-founder Geoffrey Gatama Geturo",
	HouseShaw:    "Named after the late Patrick David Shaw, a senior police officer",
	HouseHorsten: "Named after a late Danish Ambassador and benefactor of the School",
	HouseMboya:   "Named after the late Thomas Joseph Mboya, a patron of the School and cabinet minister",
	HouseShell:   "Named after the main supporter of the Centre since its inception, the Shell-BP petroleum company",
	HouseChaka:   "Named after Shaka Zulu the Zulu king and warrior",
	HouseNjonjo:  "Named after former Attorney-General Charles Njonjo",
	HouseKirkley: "Named after a friend of the Centre, Sir Lesley Kirkley",
	HouseMuriuki: "Named after longtime Board member Nick Muriuki Mugwandia",
	HouseKibaki:  "Named after the Patron of the School since 1969, third president Mwai Kibaki",
}

type Visibility string

const (
	VisibilityPublic      Visibility = "public"
	VisibilityConnections Visibility = "connections"
	VisibilityPrivate     Visibility = "private"
)
