package recovery

const (
	RepairProjectionCode  = "rebuild_projection"
	RepairJournalTailCode = "truncate_incomplete_journal_tail"
	RepairLeaseCode       = "release_abandoned_lease"
)
