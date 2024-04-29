package snapshot

type Summary struct {
	Accepted []Snapshot
	Rejected []Snapshot
	Skipped  []Snapshot
}

func (s *Summary) AddAccepted(snapshot Snapshot) {
	s.Accepted = append(s.Accepted, snapshot)
}

func (s *Summary) AddRejected(snapshot Snapshot) {
	s.Rejected = append(s.Rejected, snapshot)
}

func (s *Summary) AddSkipped(snapshot Snapshot) {
	s.Skipped = append(s.Skipped, snapshot)
}
