package monica

func ContactToContactInput(contact Contact) ContactInput {
	return ContactInput{
		FirstName:              contact.FirstName,
		LastName:               contact.LastName,
		Nickname:               contact.Nickname,
		BirthdateDay:           contact.BirthdateDay,
		BirthdateMonth:         contact.BirthdateMonth,
		BirthdateYear:          contact.BirthdateYear,
		IsBirthdateKnown:       contact.IsBirthdateKnown,
		BirthdateIsAgeBased:    contact.BirthdateIsAgeBased,
		BirthdateAge:           contact.BirthdateAge,
		IsPartial:              contact.IsPartial,
		IsDeceased:             contact.IsDeceased,
		DeceasedDateDay:        contact.DeceasedDateDay,
		DeceasedDateMonth:      contact.DeceasedDateMonth,
		DeceasedDateYear:       contact.DeceasedDateYear,
		DeceasedDateIsAgeBased: contact.DeceasedDateIsAgeBased,
		IsDeceasedDateKnown:    contact.IsDeceasedDateKnown,
	}
}
