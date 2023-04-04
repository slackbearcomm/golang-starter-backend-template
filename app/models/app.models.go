package models

type SearchFilter struct {
	SortBy     string
	SortDir    string
	Offset     int
	Limit      int
	IsFinal    *bool
	IsAccepted *bool
	IsApproved *bool
	IsArchived *bool
}
