package main

type Record struct {
	CreatedAt,
	ServiceName,
	SourceID,
	SourceURL,
	UserID,
	UserName,
	Text,
	SharedBy,
	CSVWriteAt string
}

func (r *Record) String() []string {
	var s []string
	s = append(s, r.CreatedAt)
	s = append(s, r.ServiceName)
	s = append(s, r.SourceID)
	s = append(s, r.SourceURL)
	s = append(s, r.UserID)
	s = append(s, r.UserName)
	s = append(s, r.Text)
	s = append(s, r.SharedBy)
	s = append(s, r.CSVWriteAt)
	return s
}
