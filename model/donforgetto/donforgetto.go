package donforgetto

import (
	"database/sql"

	"github.com/phassans/banana/shared"
	"github.com/rs/zerolog"
)

type (
	donforgettoEngine struct {
		sql    *sql.DB
		logger zerolog.Logger
	}

	DonforgettoEngine interface {
		FetchUserEducation(userID string) ([]shared.Education, error)
		FetchUserCompanies(userID string) ([]shared.Company, error)
		FetchUserChannels(userID string) (shared.Channels, error)
	}
)

// NewDonforgettoEngine returns an instance of donforgettoEngine
func NewDonforgettoEngine(psql *sql.DB, logger zerolog.Logger) DonforgettoEngine {
	return &donforgettoEngine{psql, logger}
}

func (d *donforgettoEngine) FetchUserEducation(userID string) ([]shared.Education, error) {
	return d.addEducation(), nil
}

func (d *donforgettoEngine) FetchUserCompanies(userID string) ([]shared.Company, error) {
	return d.addCompanies(), nil
}

func (d *donforgettoEngine) FetchUserChannels(userID string) (shared.Channels, error) {
	return shared.Channels{EducationChannels: d.addEducationChannels(), CompanyChannels: d.addCompanyChannels()}, nil
}

func (d *donforgettoEngine) addEducation() []shared.Education {
	var educations []shared.Education
	educations = append(educations, shared.Education{Name: "stanford", Degree: "bachelors", Stream: "arts", Year: 2009})
	educations = append(educations, shared.Education{Name: "berkley", Degree: "masters", Stream: "science", Year: 2013})

	return educations
}

func (d *donforgettoEngine) addCompanies() []shared.Company {
	var companies []shared.Company
	companies = append(companies, shared.Company{Name: "McAfee", Location: "Denver"})
	companies = append(companies, shared.Company{Name: "PayPal", Location: "San Jose"})
	companies = append(companies, shared.Company{Name: "Atlassian", Location: "MountainView"})

	return companies
}

func (d *donforgettoEngine) addEducationChannels() []shared.Channel {
	var channels []shared.Channel
	channels = append(channels, shared.Channel{Name: "stanford"})
	channels = append(channels, shared.Channel{Name: "stanford-bachelors-arts"})
	channels = append(channels, shared.Channel{Name: "stanford-bachelors-arts-2009"})

	channels = append(channels, shared.Channel{Name: "berkley"})
	channels = append(channels, shared.Channel{Name: "berkley-masters-science"})
	channels = append(channels, shared.Channel{Name: "berkley-masters-science-2013"})

	return channels
}

func (d *donforgettoEngine) addCompanyChannels() []shared.Channel {
	var channels []shared.Channel
	channels = append(channels, shared.Channel{Name: "McAfee"})
	channels = append(channels, shared.Channel{Name: "McAfee-Denver"})

	channels = append(channels, shared.Channel{Name: "PayPal"})
	channels = append(channels, shared.Channel{Name: "PayPal-SanJose"})

	channels = append(channels, shared.Channel{Name: "Atlassian"})
	channels = append(channels, shared.Channel{Name: "Atlassian-MountainView"})

	return channels
}
